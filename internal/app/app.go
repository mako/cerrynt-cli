// Package app contains the root Bubble Tea model.
// It owns navigation state (which screen is active) and delegates Update()
// and View() to whichever screen is currently showing.
//
// Navigation flow:
//
//	FeedList → [Enter] → ArticleList → [Enter] → ArticleView
//	               ↑ [ESC]                  ↑ [ESC]
//
// Async data flow:
//
//	Init() fires fetchFeedsCmd → FeedsLoadedMsg / FetchErrorMsg
//	FeedSelectedMsg fires fetchArticlesCmd → ArticlesLoadedMsg / FetchErrorMsg
//
// All async results are handled here and pushed into screen models via
// SetFeeds / SetArticles / SetError — screens themselves never call the store.
package app

import (
	"context"

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

// ── Async result messages ────────────────────────────────────────────────────
// These are defined here because they are produced and consumed entirely within
// this package. Screens have no knowledge of them.

// FeedsLoadedMsg carries the result of an async feed fetch.
type FeedsLoadedMsg struct {
	Feeds []domain.Feed
}

// ArticlesLoadedMsg carries the result of an async article fetch.
type ArticlesLoadedMsg struct {
	Articles []domain.Article
}

// FetchErrorMsg carries an error from any failed async fetch.
type FetchErrorMsg struct {
	Err error
	// target identifies which screen should display the error.
	target screen
}

// ── Screen enum ──────────────────────────────────────────────────────────────

type screen int

const (
	screenFeedList screen = iota
	screenArticleList
	screenArticleView
)

// ── Root model ───────────────────────────────────────────────────────────────

// Model is the root Bubble Tea model.
// It implements the tea.Model interface and is passed to tea.NewProgram.
type Model struct {
	store       data.Store
	active      screen
	width       int
	height      int
	st          *state.State // single source of truth for local application state
	statePath   string       // path to persist state; empty disables persistence
	feedList    feedlist.Model
	articleList articlelist.Model
	articleView articleview.Model
}

// New constructs the root model. The feed list starts in the loading state;
// Init() fires the first async fetch.
func New(store data.Store, st *state.State, statePath string) Model {
	return Model{
		store:     store,
		active:    screenFeedList,
		width:     80,
		height:    24,
		st:        st,
		statePath: statePath,
		feedList:  feedlist.New(),
	}
}

// Init fires the initial feed fetch and starts the feedlist spinner.
//
// tea.Batch runs both commands concurrently: the spinner tick keeps the
// animation going while the fetch is in flight.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.feedList.Init(), // spinner tick
		m.fetchFeedsCmd(),
	)
}

// fetchFeedsCmd returns a tea.Cmd that calls the store asynchronously and
// produces either FeedsLoadedMsg or FetchErrorMsg.
func (m Model) fetchFeedsCmd() tea.Cmd {
	return func() tea.Msg {
		feeds, err := m.store.Feeds(context.Background())
		if err != nil {
			return FetchErrorMsg{Err: err, target: screenFeedList}
		}
		return FeedsLoadedMsg{Feeds: feeds}
	}
}

// fetchArticlesCmd returns a tea.Cmd that fetches articles for feedID
// asynchronously and produces either ArticlesLoadedMsg or FetchErrorMsg.
func (m Model) fetchArticlesCmd(feedID string) tea.Cmd {
	return func() tea.Msg {
		articles, err := m.store.Articles(context.Background(), feedID)
		if err != nil {
			return FetchErrorMsg{Err: err, target: screenArticleList}
		}
		return ArticlesLoadedMsg{Articles: articles}
	}
}

// ── Update ───────────────────────────────────────────────────────────────────

// Update is the central message handler.
//
// Order of handling:
//  1. Terminal resize — update dimensions before anything else.
//  2. Global quit — works from any screen.
//  3. Async result messages — FeedsLoadedMsg, ArticlesLoadedMsg, FetchErrorMsg.
//  4. Navigation messages — emitted by screens.
//  5. Delegation — all other messages go to the active screen.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 1. Track terminal dimensions. Do not return — let the message continue
	//    so the active screen can also handle it (e.g. articleview resize).
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

	// 3. Async result messages.
	switch msg := msg.(type) {

	case FeedsLoadedMsg:
		m.feedList = m.feedList.SetFeeds(msg.Feeds)
		return m, nil

	case ArticlesLoadedMsg:
		// Apply persisted read state before handing articles to the screen.
		articles := msg.Articles
		for i, a := range articles {
			if m.st.Read[a.ID] {
				articles[i].IsRead = true
			}
		}
		m.articleList = m.articleList.SetArticles(articles)
		return m, nil

	case FetchErrorMsg:
		switch msg.target {
		case screenFeedList:
			m.feedList = m.feedList.SetError(msg.Err)
		case screenArticleList:
			m.articleList = m.articleList.SetError(msg.Err)
		}
		return m, nil
	}

	// 4. Navigation messages emitted by screens.
	switch msg := msg.(type) {

	case feedlist.FeedSelectedMsg:
		// Switch to article list in loading state, then kick off async fetch.
		m.articleList = articlelist.New(msg.Feed)
		m.active = screenArticleList
		return m, tea.Batch(
			m.articleList.Init(), // spinner tick
			m.fetchArticlesCmd(msg.Feed.ID),
		)

	case articlelist.ArticleSelectedMsg:
		// Mark as read in local state and persist.
		m.st.Read[msg.Article.ID] = true
		if m.statePath != "" {
			_ = state.Save(m.statePath, m.st)
		}
		// Update the live list so ESC shows up-to-date read state.
		m.articleList = m.articleList.MarkRead(msg.Article.ID)
		m.articleView = articleview.New(msg.Article, m.width, m.height)
		m.active = screenArticleView
		return m, nil

	case articlelist.BackMsg:
		m.active = screenFeedList
		return m, nil

	case articleview.BackMsg:
		m.active = screenArticleList
		return m, nil
	}

	// 5. Delegate to the active screen.
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

// ── View ─────────────────────────────────────────────────────────────────────

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
