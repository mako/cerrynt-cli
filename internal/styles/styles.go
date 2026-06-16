// Package styles defines all Lip Gloss styles used across the application.
// Screens must not declare their own ad-hoc styles — all style definitions
// live here so that theming can be applied in one place later.
//
// For now, styles use only text attributes (bold, faint, underline) and no
// hardcoded colours. This naturally respects whatever the user's terminal
// theme is set to.
package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Title is used for screen headings.
	Title = lipgloss.NewStyle().Bold(true).MarginBottom(1)

	// Selected highlights the cursor row in a list.
	// Bold with a ">" prefix is used instead of colour to stay terminal-agnostic.
	Selected = lipgloss.NewStyle().Bold(true)

	// Normal is the default list item style.
	Normal = lipgloss.NewStyle()

	// Faint is used for read/seen items.
	Faint = lipgloss.NewStyle().Faint(true)

	// StatusBar is used for the help line at the bottom of each screen.
	StatusBar = lipgloss.NewStyle().Faint(true).MarginTop(1)
)
