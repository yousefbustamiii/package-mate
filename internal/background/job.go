package background

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
)

// JobStatus describes the current state of a background job.
type JobStatus string

const (
	StatusRunning JobStatus = "running"
	StatusDone    JobStatus = "done"
	StatusFailed  JobStatus = "failed"
	StatusAborted JobStatus = "aborted"
)

// Job represents a background installation task persisted to disk.
type Job struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Action        string    `json:"action"` // install | update | install-override | brew-override
	PID           int       `json:"pid"`
	Status        JobStatus `json:"status"`
	StartedAt     time.Time `json:"started_at"`
	FinishedAt    time.Time `json:"finished_at,omitempty"`
	Error         string    `json:"error,omitempty"`
	OldBinaryPath string    `json:"old_binary_path,omitempty"`
	OldFormula    string    `json:"old_formula,omitempty"`
	IsOldCask     bool      `json:"is_old_cask,omitempty"`
}

// Dir returns the directory where job state files live.
func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".package-mate", "jobs")
}

func ensureDir() error {
	return os.MkdirAll(Dir(), 0755)
}

// FilePath returns the path to this job's JSON state file.
func (j *Job) FilePath() string {
	return filepath.Join(Dir(), j.ID+".json")
}

// Save writes the job state to disk using file locking for atomicity.
func (j *Job) Save() error {
	data, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return err
	}

	path := j.FilePath()
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Acquire exclusive lock
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	// Truncate and write
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

// LoadJob reads a single job state file using a shared lock for safety.
func LoadJob(path string) (*Job, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Acquire shared lock
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH); err != nil {
		return nil, err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	var j Job
	if err := json.NewDecoder(f).Decode(&j); err != nil {
		return nil, err
	}
	return &j, nil
}

// IsAlive returns true if the job's process is still running.
func (j *Job) IsAlive() bool {
	if j.PID <= 0 {
		return false
	}
	err := syscall.Kill(j.PID, 0)
	return err == nil
}

// Elapsed returns a human-readable duration since the job started.
func (j *Job) Elapsed() string {
	if j.StartedAt.IsZero() {
		return "—"
	}
	d := time.Since(j.StartedAt).Round(time.Second)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
}

// LoadAll reads all job files from disk, newest first.
// It reconciles stale "running" states for dead processes.
func LoadAll() []*Job {
	entries, err := os.ReadDir(Dir())
	if err != nil {
		return nil
	}

	var jobs []*Job
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		j, err := LoadJob(filepath.Join(Dir(), e.Name()))
		if err != nil {
			continue
		}
		// Reconcile: if job says running but process is dead, mark failed.
		if j.Status == StatusRunning && !j.IsAlive() {
			j.Status = StatusFailed
			j.Error = "process exited unexpectedly"
			_ = j.Save()
		}
		jobs = append(jobs, j)
	}

	sort.Slice(jobs, func(i, k int) bool {
		return jobs[i].StartedAt.After(jobs[k].StartedAt)
	})
	return jobs
}

// Abort sends SIGTERM to the entire process group of the job and marks it aborted.
func (j *Job) Abort() error {
	if j.PID <= 0 {
		return fmt.Errorf("no PID for job %s", j.ID)
	}
	err := syscall.Kill(-j.PID, syscall.SIGTERM)
	if err != nil {
		_ = syscall.Kill(j.PID, syscall.SIGTERM)
	}
	j.Status = StatusAborted
	return j.Save()
}

// GetRunningJob returns a running job for the given itemName, if any.
func GetRunningJob(itemName string) *Job {
	jobs := LoadAll()
	for _, j := range jobs {
		if j.Name == itemName && j.Status == StatusRunning {
			return j
		}
	}
	return nil
}

// Enqueue creates a job file and launches the background installer process.
func Enqueue(itemName, action, oldBinaryPath, oldFormula string, isOldCask bool) (*Job, error) {
	if err := ensureDir(); err != nil {
		return nil, fmt.Errorf("cannot create jobs directory: %w", err)
	}

	// ── Atomic Queue Lock ───────────────────────────────────────────────────
	lockPath := filepath.Join(Dir(), ".enqueue.lock")
	lockFile, err := os.OpenFile(lockPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("cannot open queue lock: %w", err)
	}
	defer lockFile.Close()

	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX); err != nil {
		return nil, fmt.Errorf("cannot acquire queue lock: %w", err)
	}
	defer syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)

	// Double-check for existing running jobs for this item to prevent race conditions.
	if existing := GetRunningJob(itemName); existing != nil {
		return nil, fmt.Errorf("a background task for '%s' is already in progress (PID: %d)", itemName, existing.PID)
	}

	j := &Job{
		ID:            fmt.Sprintf("%d", time.Now().UnixNano()),
		Name:          itemName,
		Action:        action,
		Status:        StatusRunning,
		StartedAt:     time.Now(),
		OldBinaryPath: oldBinaryPath,
		OldFormula:    oldFormula,
		IsOldCask:     isOldCask,
	}

	// Write the job file first so the child can find it.
	if err := j.Save(); err != nil {
		return nil, fmt.Errorf("cannot save job: %w", err)
	}
	// ── End Atomic Block ────────────────────────────────────────────────────

	// Find our own executable.
	exe, err := os.Executable()
	if err != nil {
		_ = os.Remove(j.FilePath())
		return nil, fmt.Errorf("cannot locate executable: %w", err)
	}

	// Open /dev/null to silence child stdout/stderr.
	devNull, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err != nil {
		_ = os.Remove(j.FilePath())
		return nil, fmt.Errorf("cannot open /dev/null: %w", err)
	}

	// Fork a detached child process.
	cmd := exec.Command(exe, "_bg-exec", j.ID)
	cmd.Stdin = devNull
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	// Setsid detaches from the controlling terminal and creates a new session+process group.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		_ = devNull.Close()
		_ = os.Remove(j.FilePath())
		return nil, fmt.Errorf("cannot start background process: %w", err)
	}
	_ = devNull.Close()

	// Record the PID so mate bg can find and abort it.
	j.PID = cmd.Process.Pid
	if err := j.Save(); err != nil {
		return nil, err
	}

	return j, nil
}

// CleanOld removes finished jobs older than 24 hours.
func CleanOld() {
	jobs := LoadAll()
	cutoff := time.Now().Add(-24 * time.Hour)
	for _, j := range jobs {
		if j.Status != StatusRunning && j.StartedAt.Before(cutoff) {
			_ = os.Remove(j.FilePath())
		}
	}
}

// CleanAllFinished removes all jobs that are not currently running.
func CleanAllFinished() {
	jobs := LoadAll()
	for _, j := range jobs {
		if j.Status != StatusRunning {
			_ = os.Remove(j.FilePath())
		}
	}
}

// NotifyFinish sends a desktop notification for a completed job.
func (j *Job) NotifyFinish() {
	title := "Package Mate"
	msg := fmt.Sprintf("✓ %s finished %s!", j.Name, j.Action)
	if j.Status == StatusFailed {
		title = "Package Mate (Error)"
		msg = fmt.Sprintf("✗ %s %s failed: %s", j.Name, j.Action, j.Error)
	}

	home, _ := os.UserHomeDir()
	iconPath := filepath.Join(home, "Desktop", "package-mate", "icon.png")
	_ = beeep.Notify(title, msg, iconPath)
}
