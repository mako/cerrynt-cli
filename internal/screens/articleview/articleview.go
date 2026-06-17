// Package articleview implements the article reading screen.
// It receives a single article and displays its content.
// Content is scrollable with j/k when it exceeds the terminal height.
package articleview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mako/cerrynt-cli/internal/domain"
	"github.com/mako/cerrynt-cli/internal/keymap"
	"github.com/mako/cerrynt-cli/internal/styles"
)

const statusRows = 2

// BackMsg is emitted when the user presses ESC to return to the article list.
type BackMsg struct{}

// Model is the Bubble Tea model for the article view screen.
type Model struct {
	article domain.Article
	width   int
	height  int
	offset  int // index of the first visible line
}

// New constructs an articleview Model. width and height should be the current
// terminal dimensions so word-wrap and scrolling are correct from the first render.
func New(article domain.Article, width, height int) Model {
	return Model{article: article, width: width, height: height}
}

// Init satisfies tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key input and window resize.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Re-clamp in case content is now fully visible after a resize.
		m.offset = clamp(m.offset, 0, maxOffset(m.buildLines(), m.height))
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Default.Down):
			lines := m.buildLines()
			m.offset = clamp(m.offset+1, 0, maxOffset(lines, m.height))

		case key.Matches(msg, keymap.Default.Up):
			m.offset = clamp(m.offset-1, 0, maxOffset(m.buildLines(), m.height))

		case key.Matches(msg, keymap.Default.Back):
			return m, func() tea.Msg { return BackMsg{} }
		}
	}

	return m, nil
}

// View renders the visible window of content and a status bar.
func (m Model) View() string {
	lines := m.buildLines()

	// The status bar occupies 2 rows: a blank separator line + the bar text.
	// This constant must match the literal "\n\n" prefix in the return below.
	visible := m.height - statusRows
	if visible < 1 {
		visible = 1
	}

	offset := clamp(m.offset, 0, maxOffset(lines, m.height))
	end := offset + visible
	if end > len(lines) {
		end = len(lines)
	}

	content := strings.Join(lines[offset:end], "\n")

	// Build status bar text.
	max := maxOffset(lines, m.height)
	var bar string
	if max > 0 {
		// Content is taller than the screen — show a scroll percentage.
		pct := offset * 100 / max
		bar = fmt.Sprintf("[%d%%]  j/k scroll  •  esc back  •  q quit", pct)
	} else {
		bar = "esc back  •  q quit"
	}

	// "\n\n" = blank separator line between content and status bar (= statusRows).
	return content + "\n\n" + styles.StatusBar.Render(bar)
}

// buildLines renders the full article into a flat slice of display lines.
// All layout decisions live here so View() is just a slice-and-join operation
// and scroll position is always relative to a stable line numbering.
func (m Model) buildLines() []string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	// Narrow by 2 so text does not press against the very edge of the terminal.
	wrap := w - 2
	if wrap < 20 {
		wrap = 20
	}

	var lines []string

	a := m.article

	// Title: each word-wrapped line is its own slice element so that the
	// scroll-window calculation in View() counts display rows correctly.
	// styles.Title has no margin, so per-line rendering is safe.
	for _, l := range wordWrap(a.Title, wrap) {
		lines = append(lines, styles.Title.Render(l))
	}
	lines = append(lines, "") // blank line after title

	// Metadata
	lines = append(lines,
		styles.Faint.Render("Published  "+a.Published.Format("2 Jan 2006  15:04")),
		styles.Faint.Render("URL        "+a.URL),
		"", // blank line before body
	)

	// Body
	if a.Summary != "" {
		lines = append(lines, wordWrap(a.Summary, wrap)...)
	} else {
		lines = append(lines, styles.Faint.Render("No summary available."))
	}

	return lines
}

// maxOffset returns the maximum valid scroll offset for a given line count
// and terminal height. Returns 0 if all content fits on screen.
func maxOffset(lines []string, height int) int {
	visible := height - statusRows
	if visible < 1 {
		visible = 1
	}
	if max := len(lines) - visible; max > 0 {
		return max
	}
	return 0
}

// clamp returns v clamped to [lo, hi].
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// wordWrap breaks text into lines of at most width characters on word boundaries.
func wordWrap(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	current := words[0]

	for _, word := range words[1:] {
		if len(current)+1+len(word) > width {
			lines = append(lines, current)
			current = word
		} else {
			current += " " + word
		}
	}

	return append(lines, current)
}
