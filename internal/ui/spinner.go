package ui

import (
	"fmt"
	"sync"
	"time"
)

var (
	currSpinner *spinner
	spinnerMu   sync.Mutex
)

type spinner struct {
	message string
	stop    chan struct{}
	done    chan struct{}
}

func (s *spinner) start() {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	// 5 dots animation pattern
	frames := []string{
		C(White, ".") + C(Dim, "...."),
		C(Dim, ".") + C(White, ".") + C(Dim, "..."),
		C(Dim, "..") + C(White, ".") + C(Dim, ".."),
		C(Dim, "...") + C(White, ".") + C(Dim, "."),
		C(Dim, "....") + C(White, "."),
	}
	i := 0

	// Initial print
	fmt.Printf("\r  "+C(Bold+BrightCyan, "❯")+"  "+C(Bold+White, s.message)+" %s", frames[0])

	for {
		select {
		case <-ticker.C:
			i++
			fmt.Printf("\r  "+C(Bold+BrightCyan, "❯")+"  "+C(Bold+White, s.message)+" %s", frames[i%len(frames)])
		case <-s.stop:
			fmt.Print("\r\033[K")
			close(s.done)
			return
		}
	}
}

func startDoing(msg string) {
	spinnerMu.Lock()
	defer spinnerMu.Unlock()

	// Stop previous if any
	if currSpinner != nil {
		close(currSpinner.stop)
		<-currSpinner.done
	}

	currSpinner = &spinner{
		message: msg,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
	go currSpinner.start()
}

func stopDoing() {
	spinnerMu.Lock()
	defer spinnerMu.Unlock()

	if currSpinner != nil {
		close(currSpinner.stop)
		<-currSpinner.done
		currSpinner = nil
	}
}
