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
// Read state is tracked here in readIDs and applied in two ways:
//  1. Immediately: articleList.MarkRead() updates the live list so the user
//     sees the change as soon as they press ESC from an article.
//  2. On re-entry: when the user opens the same feed again, readIDs are
//     applied to the fresh article slice before constructing articleList.
package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mako/cerrynt-cli/internal/data"
	"github.com/mako/cerrynt-cli/internal/keymap"
	"github.com/mako/cerrynt-cli/internal/screens/articlelist"
	"github.com/mako/cerrynt-cli/internal/screens/articleview"
	"github.com/mako/cerrynt-cli/internal/screens/feedlist"
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
	readIDs     map[string]bool // in-memory read state; source of truth for MVP
	feedList    feedlist.Model
	articleList articlelist.Model
	articleView articleview.Model
}

// New constructs the root model with default terminal dimensions.
// The real dimensions arrive via tea.WindowSizeMsg shortly after startup.
func New() Model {
	return Model{
		active:   screenFeedList,
		width:    80,
		height:   24,
		readIDs:  make(map[string]bool),
		feedList: feedlist.New(data.Feeds()),
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
		// Build article list, applying in-memory read state to any articles
		// the user has already read in this session.
		articles := data.Articles(msg.Feed.ID)
		for i, a := range articles {
			if m.readIDs[a.ID] {
				articles[i].IsRead = true
			}
		}
		m.articleList = articlelist.New(msg.Feed, articles)
		m.active = screenArticleList
		return m, nil

	case articlelist.ArticleSelectedMsg:
		// Record read state before opening the article.
		m.readIDs[msg.Article.ID] = true
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
