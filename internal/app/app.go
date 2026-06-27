// Package app contains the root Bubble Tea model.
// It owns navigation state (which screen is active) and delegates Update()
// and View() to whichever screen is currently showing.
//
// Navigation flow:
//
//	FeedList → [Enter] → ArticleList → [Enter] → ArticleView
//	               ↑ [ESC]                  ↑ [ESC]
//
// Navigation is driven by typed messages: screens emit them, app receives
// them and switches the active screen. Screens have no knowledge of each other.
//
// Local application state (e.g. which articles have been read) is held in a
// *state.State and applied in two ways:
//  1. Immediately: articleList.MarkRead() updates the live list so the user
//     sees the change as soon as they press ESC from an article.
//  2. On re-entry: when the user opens the same feed again, state.Read is
//     applied to the fresh article slice before constructing articleList.
package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mako/cerrynt-cli/internal/data"
	"github.com/mako/cerrynt-cli/internal/domain"
	"github.com/mako/cerrynt-cli/internal/keymap"
	"github.com/mako/cerrynt-cli/internal/screens/articlelist"
	"github.com/mako/cerrynt-cli/internal/screens/articleview"
	"github.com/mako/cerrynt-cli/internal/screens/feedlist"
	"github.com/mako/cerrynt-cli/internal/state"
)

// screen is an enum of all possible active screens.
type screen int

const (
	screenFeedList screen = iota
	screenArticleList
	screenArticleView
)

// Model is the root Bubble Tea model.
// It implements the tea.Model interface and is passed to tea.NewProgram.
type Model struct {
	active      screen
	width       int
	height      int
	st          *state.State // single source of truth for local application state
	statePath   string       // path to persist state; empty disables persistence
	feedList    feedlist.Model
	articleList articlelist.Model
	articleView articleview.Model
}

// New constructs the root model with the provided feeds and local state.
// statePath is the file path used by Save; pass an empty string to disable
// persistence (useful in tests or when no writable location is available).
// The real terminal dimensions arrive via tea.WindowSizeMsg shortly after startup.
func New(feeds []domain.Feed, st *state.State, statePath string) Model {
	return Model{
		active:    screenFeedList,
		width:     80,
		height:    24,
		st:        st,
		statePath: statePath,
		feedList:  feedlist.New(feeds),
	}
}

// Init satisfies tea.Model. No startup commands needed.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is the central message handler.
//
// Order of handling:
//  1. Terminal resize — update dimensions before anything else so that
//     newly constructed screens get the correct size.
//  2. Global quit — works from any screen.
//  3. Navigation messages — emitted by screens, handled here to switch screens.
//  4. Delegation — all other messages go to the active screen's Update().
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 1. Track terminal dimensions. We do not return here — the message
	//    continues to be processed so the active screen can also resize itself.
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = sizeMsg.Width
		m.height = sizeMsg.Height
	}

	// 2. Global quit.
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, keymap.Default.Quit) {
			return m, tea.Quit
		}
	}

	// 3. Navigation messages.
	switch msg := msg.(type) {

	case feedlist.FeedSelectedMsg:
		// Build article list, applying persisted read state to each article.
		articles := data.Articles(msg.Feed.ID)
		for i, a := range articles {
			if m.st.Read[a.ID] {
				articles[i].IsRead = true
			}
		}
		m.articleList = articlelist.New(msg.Feed, articles)
		m.active = screenArticleList
		return m, nil

	case articlelist.ArticleSelectedMsg:
		// Mark the article as read in the shared state and persist immediately.
		m.st.Read[msg.Article.ID] = true
		if m.statePath != "" {
			// Errors are intentionally ignored: the in-memory state is correct
			// and the TUI should not crash due to a transient disk failure.
			// A future version may surface this via a status bar notification.
			_ = state.Save(m.statePath, m.st)
		}
		// Update the live list immediately so ESC returns to an up-to-date view.
		m.articleList = m.articleList.MarkRead(msg.Article.ID)
		// Pass current dimensions so word-wrap and scrolling are correct.
		m.articleView = articleview.New(msg.Article, m.width, m.height)
		m.active = screenArticleView
		return m, nil

	case articlelist.BackMsg:
		// Return to the feed list without reconstructing it (preserves cursor).
		m.active = screenFeedList
		return m, nil

	case articleview.BackMsg:
		// Return to the article list without reconstructing it (preserves cursor
		// and the read state applied by MarkRead above).
		m.active = screenArticleList
		return m, nil
	}

	// 4. Delegate to the active screen.
	switch m.active {
	case screenFeedList:
		var cmd tea.Cmd
		m.feedList, cmd = m.feedList.Update(msg)
		return m, cmd

	case screenArticleList:
		var cmd tea.Cmd
		m.articleList, cmd = m.articleList.Update(msg)
		return m, cmd

	case screenArticleView:
		var cmd tea.Cmd
		m.articleView, cmd = m.articleView.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the currently active screen.
func (m Model) View() string {
	switch m.active {
	case screenFeedList:
		return m.feedList.View()
	case screenArticleList:
		return m.articleList.View()
	case screenArticleView:
		return m.articleView.View()
	}
	return ""
}
