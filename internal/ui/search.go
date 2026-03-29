package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/term"
	"github.com/yousefbustamiii/package-mate/internal/components"
)

const (
	modeSearch = 0
	modeAction = 1 // result-navigation mode (j/k between results)
)

// bannerWidth is the visual column width of the PACKAGE ASCII art (including 2-space left margin).
// The box inner = bannerWidth - 4 (2 margin + ┌ + ┐).
const bannerWidth = 59

// ActionFunc is called when the user confirms a selection from the action menu.
// choice: 1=Install/Update  2=Uninstall  3=Information & Versions
type ActionFunc func(itemName string, choice int)

// SearchResult is a matched entry with its sort priority.
type SearchResult struct {
	Entry       SectionEntry
	SectionName string
	priority    int // 1=name starts-with  2=name contains  3=desc contains  0=no query
}

type searchState struct {
	query      []rune
	results    []SearchResult
	page       int
	cursor     int
	groups     []SectionGroup
	pageSize   int // always 10
	totalLines int // lines printed last render (for erase-then-redraw)
	mode       int // modeSearch or modeAction
}

// ShowSearch is the public entry point.
func ShowSearch(groups []SectionGroup, onSelect ActionFunc) {
	// ── Static header (printed once, scrolls up naturally) ────────────────────
	totalCount := 0
	installedCount := 0
	for _, g := range groups {
		totalCount += len(g.Entries)
		for _, e := range g.Entries {
			if e.Status != components.StatusNotInstalled {
				installedCount++
			}
		}
	}

	Banner()
	fmt.Println("  " + C(Bold+White, "usage: ") + C(Cyan, "mate <tool> ") + C(Dim, "(interactive menu)"))
	Blank()
	fmt.Printf("  "+C(Bold+White, "Installed: ")+C(BrightGreen+Bold, "%d")+" / "+C(Bold+White, "%d")+"\n", installedCount, totalCount)
	Blank()
	fmt.Printf("  "+C(BrightGreen+Bold, "✓")+"  "+C(Dim, "---> Installed & Managed (up to date)")+"\n")
	Blank()
	fmt.Printf("  "+C(Yellow+Bold, "↻")+"  "+C(Dim, "---> Installed & Managed, update available")+"\n")
	Blank()
	fmt.Printf("  "+C(BrightCyan+Bold, "⚙")+"  "+C(Dim, "---> Installed & Unmanaged")+"\n")
	Blank()
	fmt.Printf("  "+C(Bold+White, "?")+"  "+C(Dim, "---> Multiple Installations")+"\n")
	Blank()

	// ── Initial state ─────────────────────────────────────────────────────────
	state := &searchState{
		groups:   groups,
		pageSize: 10,
		mode:     modeSearch,
	}
	refilter(state)

	// ── Raw mode ──────────────────────────────────────────────────────────────
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\033[?25h")
	}()
	fmt.Print("\033[?25l")

	// ── Event loop ────────────────────────────────────────────────────────────
	for {
		render(state)

		var b [4]byte
		n, _ := os.Stdin.Read(b[:])

		exit, selectChosen := handleKey(b[:], n, state)
		if exit {
			if state.totalLines > 0 {
				fmt.Printf("\033[%dA\033[J", state.totalLines)
			}
			return
		}
		if selectChosen {
			page := currentPageResults(state)
			if state.cursor >= len(page) {
				continue
			}
			chosen := page[state.cursor]

			// Restore terminal for the numbered menu
			_ = term.Restore(int(os.Stdin.Fd()), oldState)
			fmt.Print("\033[?25h")

			// Erase dynamic region
			if state.totalLines > 0 {
				fmt.Printf("\033[%dA\033[J", state.totalLines)
			}
			state.totalLines = 0

			// Show the original numbered action menu and execute — then done
			choice := PromptActionMenu(chosen.Entry.Name)
			if choice > 0 {
				onSelect(chosen.Entry.Name, choice)
			}
			return
		}
	}
}

// ── Filter ────────────────────────────────────────────────────────────────────

func filterResults(groups []SectionGroup, query []rune) []SearchResult {
	q := strings.ToLower(strings.TrimSpace(string(query)))
	var out []SearchResult
	for _, g := range groups {
		for _, e := range g.Entries {
			if q == "" {
				out = append(out, SearchResult{e, g.Label, 0})
				continue
			}
			nl := strings.ToLower(e.Name)
			dl := strings.ToLower(e.Desc)
			if strings.HasPrefix(nl, q) {
				out = append(out, SearchResult{e, g.Label, 1})
			} else if strings.Contains(nl, q) {
				out = append(out, SearchResult{e, g.Label, 2})
			} else if strings.Contains(dl, q) {
				out = append(out, SearchResult{e, g.Label, 3})
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].priority < out[j].priority
	})
	return out
}

func refilter(state *searchState) {
	state.results = filterResults(state.groups, state.query)
	state.page = 0
	state.cursor = 0
}

// ── Pagination ────────────────────────────────────────────────────────────────

func currentPageResults(s *searchState) []SearchResult {
	start := s.page * s.pageSize
	if start >= len(s.results) {
		return nil
	}
	end := start + s.pageSize
	if end > len(s.results) {
		end = len(s.results)
	}
	return s.results[start:end]
}

func totalPages(s *searchState) int {
	if len(s.results) == 0 {
		return 1
	}
	return (len(s.results) + s.pageSize - 1) / s.pageSize
}

// ── Render ────────────────────────────────────────────────────────────────────

func render(state *searchState) {
	if state.totalLines > 0 {
		fmt.Printf("\033[%dA\033[J", state.totalLines)
	}

	lines := 0
	w := termWidth()

	// 1. Input box (fixed width matching banner)
	lines += printInputBox(state.query, state.mode)

	// 2. Blank
	fmt.Print("\r\n")
	lines++

	// 3. Results block
	pageResults := currentPageResults(state)
	if len(pageResults) == 0 {
		if len(state.query) > 0 {
			fmt.Printf("\r\033[K  %s\r\n", C(Dim, "No results found."))
		} else {
			fmt.Printf("\r\033[K  %s\r\n", C(Dim, "Type to search..."))
		}
		fmt.Print("\r\n")
		fmt.Print("\r\n")
		lines += 3 // message + 2 blanks
	} else {
		// Name column width for this page
		nameColWidth := 10
		for _, r := range pageResults {
			if len(r.Entry.Name) > nameColWidth {
				nameColWidth = len(r.Entry.Name)
			}
		}
		nameColWidth += 2

		for i, r := range pageResults {
			printResultRow(r, i == state.cursor, nameColWidth, w)
			fmt.Print("\r\n") // blank 1
			fmt.Print("\r\n") // blank 2
			lines += 3        // result line + 2 blanks
		}
	}

	// 4. Nav line
	lines += printNavLine(state, len(pageResults))

	state.totalLines = lines
}

// printInputBox renders the search input box.
// Width is fixed to match the PACKAGE banner width.
// Returns lines used (always 3).
func printInputBox(query []rune, mode int) int {
	const boxInner = bannerWidth - 4 // 59 - 4 = 55

	// Top border — bold white
	fmt.Printf("\r\033[K  %s\r\n", C(Bold+White, "┌"+strings.Repeat("─", boxInner)+"┐"))

	// Content layout: │  ⌕  [text/cursor]  [padding]  │
	const iconStr = "  ⌕  " // 5 visible chars
	const iconLen = 5
	const rightPad = 2
	contentArea := boxInner - iconLen - rightPad

	var textStr string
	var textLen int

	if len(query) == 0 {
		ph := "Search packages..."
		if len(ph) > contentArea {
			ph = ph[:contentArea]
		}
		textStr = C(Dim, ph)
		textLen = len(ph)
	} else {
		qs := string(query)
		maxDisplay := contentArea - 1 // leave 1 for cursor block
		if maxDisplay < 0 {
			maxDisplay = 0
		}
		if len(qs) > maxDisplay {
			qs = qs[len(qs)-maxDisplay:]
		}
		// In action mode, dim the query slightly to indicate we're navigating
		if mode == modeAction {
			textStr = C(Dim, qs) + C(Dim, "│")
		} else {
			textStr = C(White+Bold, qs) + C(BrightCyan, "█")
		}
		textLen = len(qs) + 1
	}

	pad := contentArea - textLen
	if pad < 0 {
		pad = 0
	}

	fmt.Printf("\r\033[K  %s%s%s%s%s\r\n",
		C(Bold+White, "│"),
		C(BrightCyan, iconStr),
		textStr,
		strings.Repeat(" ", pad+rightPad),
		C(Bold+White, "│"),
	)

	// Bottom border — bold white
	fmt.Printf("\r\033[K  %s\r\n", C(Bold+White, "└"+strings.Repeat("─", boxInner)+"┘"))

	return 3
}

// printResultRow renders one result row.
func printResultRow(r SearchResult, highlighted bool, nameColWidth, w int) {
	const statusLen = 4 // "[✓] "

	// Status prefix — exact same logic as ShowAllTools
	statusPrefix := C(Dim, "[ ] ")
	switch r.Entry.Status {
	case components.StatusInstalled:
		statusPrefix = C(Dim, "[") + C(BrightGreen+Bold, "✓") + C(Dim, "] ")
	case components.StatusOutdated:
		statusPrefix = C(Dim, "[") + C(Yellow+Bold, "↻") + C(Dim, "] ")
	case components.StatusUnmanaged:
		statusPrefix = C(Dim, "[") + C(BrightCyan+Bold, "⚙") + C(Dim, "] ")
	}

	// Name — same as ShowAllTools
	nameStr := C(r.Entry.Color+Bold, r.Entry.Name)
	if r.Entry.HasMultiple {
		nameStr += " " + C(Bold+White, "?")
	}

	selector := "  "
	if highlighted {
		selector = C(Cyan, "❯ ")
	}

	nameVisLen := len(r.Entry.Name)
	if r.Entry.HasMultiple {
		nameVisLen += 2
	}
	namePad := nameColWidth - nameVisLen
	if namePad < 1 {
		namePad = 1
	}

	// Description — truncated to fit remaining terminal width
	usedFixed := 2 + 2 + statusLen + nameColWidth + namePad + 2
	descAvail := w - usedFixed
	desc := r.Entry.Desc
	if descAvail < 4 {
		desc = ""
	} else if len(desc) > descAvail {
		desc = desc[:descAvail-1] + "…"
	}

	fmt.Printf("\r\033[K  %s%s%s%s%s\r\n",
		selector,
		statusPrefix,
		nameStr,
		strings.Repeat(" ", namePad),
		C(Dim, desc),
	)
}

// printNavLine renders the bottom nav/hint line. Returns lines used (always 1).
func printNavLine(state *searchState, pageLen int) int {
	tp := totalPages(state)
	total := len(state.results)

	var parts []string

	if state.mode == modeSearch {
		// Search mode: Enter moves to result navigation
		if pageLen > 0 {
			parts = append(parts, C(White, "[↵]")+" "+C(Dim, "navigate results"))
		}
		parts = append(parts, C(Dim, fmt.Sprintf("page %d/%d", state.page+1, tp)))
		parts = append(parts, C(Dim, fmt.Sprintf("(%d results)", total)))
		parts = append(parts, C(White, "[esc]")+" "+C(Dim, "exit"))
	} else {
		// Action mode: Enter selects, s goes back, n/p paginate
		parts = append(parts, C(White, "[↵]")+" "+C(Dim, "select"))
		parts = append(parts, C(White, "[s]")+" "+C(Dim, "back to search"))
		if state.page < tp-1 {
			parts = append(parts, C(White, "[n]")+" "+C(Dim, "next page"))
		}
		if state.page > 0 {
			parts = append(parts, C(White, "[p]")+" "+C(Dim, "prev page"))
		}
		parts = append(parts, C(Dim, fmt.Sprintf("page %d/%d", state.page+1, tp)))
		parts = append(parts, C(White, "[esc]")+" "+C(Dim, "exit"))
	}

	fmt.Printf("\r\033[K  %s\r\n", strings.Join(parts, C(Dim, "   ")))
	return 1
}

// ── Key handling ──────────────────────────────────────────────────────────────

// handleKey returns (exit, selectChosen).
// selectChosen=true means Enter was pressed in modeAction to confirm a result.
func handleKey(b []byte, n int, state *searchState) (bool, bool) {
	// Escape sequences (arrow keys)
	if n == 3 && b[0] == 27 && b[1] == '[' {
		switch b[2] {
		case 'A': // Arrow Up
			moveCursorUp(state)
		case 'B': // Arrow Down
			moveCursorDown(state)
		}
		return false, false
	}

	if n == 1 {
		switch b[0] {
		case 3, 27: // Ctrl+C, Esc — exit from any mode
			return true, false

		case 13: // Enter
			if state.mode == modeSearch {
				// Switch to result-navigation mode if results exist
				if len(currentPageResults(state)) > 0 {
					state.mode = modeAction
				}
			} else {
				// Confirm the highlighted result
				if len(currentPageResults(state)) > 0 {
					return false, true
				}
			}

		case 127, 8: // Backspace
			if state.mode == modeSearch && len(state.query) > 0 {
				state.query = state.query[:len(state.query)-1]
				refilter(state)
			} else if state.mode == modeAction {
				// Backspace in action mode goes back to search
				state.mode = modeSearch
			}

		default:
			if state.mode == modeSearch {
				// All printable chars go to the query
				if b[0] >= 32 && b[0] < 127 {
					state.query = append(state.query, rune(b[0]))
					refilter(state)
					state.cursor = 0
				}
			} else {
				// Action (result-navigation) mode controls
				switch b[0] {
				case 's', 'S':
					state.mode = modeSearch
				case 'j', 'J':
					moveCursorDown(state)
				case 'k', 'K':
					moveCursorUp(state)
				case 'n', 'N':
					nextPage(state)
				case 'p', 'P':
					prevPage(state)
				}
			}
		}
	}

	return false, false
}

// ── Cursor & pagination helpers ───────────────────────────────────────────────

func moveCursorDown(state *searchState) {
	page := currentPageResults(state)
	if len(page) == 0 {
		return
	}
	state.cursor = (state.cursor + 1) % len(page)
}

func moveCursorUp(state *searchState) {
	page := currentPageResults(state)
	if len(page) == 0 {
		return
	}
	state.cursor = (state.cursor - 1 + len(page)) % len(page)
}

func nextPage(state *searchState) {
	if state.page < totalPages(state)-1 {
		state.page++
		state.cursor = 0
	}
}

func prevPage(state *searchState) {
	if state.page > 0 {
		state.page--
		state.cursor = 0
	}
}
