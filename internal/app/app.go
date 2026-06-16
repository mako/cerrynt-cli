// Package app contains the root Bubble Tea model.
// It owns navigation state (which screen is active) and delegates Update()
// and View() to whichever screen is currently showing.
//
// Navigation flow:
//   FeedList → [Enter] → ArticleList → [Enter] → ArticleView
//                              ↑ [ESC]                ↑ [ESC]
//
// Navigation is driven by typed messages: screens emit them, app receives
// them and switches the active screen. Screens have no knowledge of each other.
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
	feedList    feedlist.Model
	articleList articlelist.Model
	articleView articleview.Model
}

// New constructs the root model, loading feeds from the data layer.
func New() Model {
	return Model{
		active:   screenFeedList,
		feedList: feedlist.New(data.Feeds()),
	}
}

// Init satisfies tea.Model. No startup commands needed.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is the central message handler.
//
// The order of handling matters:
//  1. Global quit is checked first so it works from any screen.
//  2. Navigation messages (emitted by screens as tea.Cmd results) switch
//     the active screen and construct the next screen's model.
//  3. All other messages are delegated to the active screen's Update().
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 1. Global quit — handled before any screen sees the message.
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, keymap.Default.Quit) {
			return m, tea.Quit
		}
	}

	// 2. Navigation messages emitted by screens.
	//
	// Each case constructs the target screen fresh from the relevant data,
	// then switches active. Screens are replaced on navigation; scroll
	// position is not preserved (accepted trade-off for MVP simplicity).
	switch msg := msg.(type) {
	case feedlist.FeedSelectedMsg:
		articles := data.Articles(msg.Feed.ID)
		m.articleList = articlelist.New(msg.Feed, articles)
		m.active = screenArticleList
		return m, nil

	case articlelist.ArticleSelectedMsg:
		m.articleView = articleview.New(msg.Article)
		m.active = screenArticleView
		return m, nil

	case articlelist.BackMsg:
		m.active = screenFeedList
		return m, nil

	case articleview.BackMsg:
		m.active = screenArticleList
		return m, nil
	}

	// 3. Delegate all other messages to the active screen.
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
