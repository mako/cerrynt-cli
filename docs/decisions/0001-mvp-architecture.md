# ADR 0001 — MVP Architecture

**Status:** Accepted  
**Date:** 2026-06-16  
**Authors:** @mako, Antigravity (AI pair)

---

## Context

Cerrynt CLI is a terminal-first RSS reader written in Go using
[Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss).

It is the first component of a broader ecosystem that will eventually include:

- **cerrynt-web** — a Ruby on Rails backend and web UI
- **cerrynt-mobile** — a possible future mobile client

The Rails backend will become the source of truth for feeds, articles,
read/unread state, user accounts, and cross-device synchronisation.

The CLI is intentionally a **thin client**: it delegates business logic and
persistence to the backend and keeps its own footprint minimal.

### Project constraints

- The primary author is a senior Ruby/JavaScript developer learning Go.
- Go idioms should be followed, but code should remain understandable.
- Complexity must be introduced only when there is a demonstrated need.
- The MVP exists to learn Bubble Tea and establish a clean architecture,
  not to ship a production-ready feature set.

### MVP scope (hard boundaries)

The MVP must not include:

- RSS feed fetching or parsing
- HTTP networking
- Backend or API integration
- Authentication
- Any form of persistence or caching

The MVP will use **static mock data** only. This constraint is intentional.

---

## Decision

### 1. Repository structure

```
cerrynt-cli/
├── cmd/
│   └── cerrynt/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── domain/
│   │   └── domain.go
│   ├── data/
│   │   └── mock.go
│   ├── screens/
│   │   ├── feedlist/
│   │   │   └── feedlist.go
│   │   ├── articlelist/
│   │   │   └── articlelist.go
│   │   └── articleview/
│   │       └── articleview.go
│   ├── styles/
│   │   └── styles.go
│   ├── keymap/
│   │   └── keymap.go
│   └── config/
│       └── config.go
├── docs/
│   └── decisions/
│       └── 0001-mvp-architecture.md
├── AGENTS.md
├── go.mod
├── go.sum
└── README.md
```

### 2. Package responsibilities

| Package | Responsibility |
|---|---|
| `cmd/cerrynt` | Entry point. Initialises config and starts the Bubble Tea program. |
| `internal/app` | Root Bubble Tea `Model`. Owns navigation state (active screen enum). Delegates `Update()` and `View()` to the active screen. |
| `internal/domain` | Plain Go structs: `Feed`, `Article`. No Bubble Tea, no logic, no tags. |
| `internal/data` | Returns `[]domain.Feed` and `[]domain.Article`. Mock implementation for MVP; interface will accommodate API and cache sources later. |
| `internal/screens/feedlist` | Self-contained Bubble Tea component for the feed list. |
| `internal/screens/articlelist` | Self-contained Bubble Tea component for the article list. |
| `internal/screens/articleview` | Self-contained Bubble Tea component for the article reader. |
| `internal/styles` | Centralised Lip Gloss style definitions. All screens source styles from here. |
| `internal/keymap` | Keybinding definitions using `charmbracelet/bubbles/key`. Shared across all screens for MVP. |
| `internal/config` | Loads `$XDG_CONFIG_HOME/cerrynt/config.yaml`. Returns a `Config` struct. Stubbed for MVP. |

### 3. Bubble Tea architecture

Cerrynt uses the **Elm Architecture** as implemented by Bubble Tea:

- Each screen is a self-contained model with `Init()`, `Update()`, and `View()`.
- Screens communicate upward by returning **typed message structs** via `tea.Cmd`,
  not by calling methods on the parent.
- The root `app.Model` inspects messages to drive navigation transitions.

**Navigation flow:**

```
FeedList
  → [Enter] → FeedSelectedMsg → app constructs ArticleList → switches screen
      → [Enter] → ArticleSelectedMsg → app constructs ArticleView → switches screen
          → [ESC] → BackMsg → app constructs ArticleList → switches screen
      → [ESC] → BackMsg → app switches back to FeedList
  → [q] → quit from any screen
```

### 4. Domain model

```go
// internal/domain/domain.go

type Feed struct {
    ID          string
    Title       string
    URL         string
    UnreadCount int
}

type Article struct {
    ID        string
    FeedID    string
    Title     string
    URL       string
    Summary   string
    Published time.Time
    IsRead    bool
}
```

These are intentionally simple value types. No methods, no database tags,
no JSON tags in the MVP. They will be extended when API integration begins.

### 5. Screen lifecycle

Screen models are **constructed (or replaced) on navigation events** rather
than being kept alive permanently. This is the simpler approach for MVP.

Scroll position is not preserved when navigating back. This is a deliberate
trade-off: the added complexity of retaining sub-model state is not justified
until there is a demonstrated need from real usage.

If scroll position preservation becomes important, the `app.Model` can be
updated to retain sub-models instead of replacing them — this change is
localised to `internal/app` and does not affect the screen packages.

### 6. Styles

All Lip Gloss style definitions live in `internal/styles`. Screens must not
define their own ad-hoc styles. This prevents style drift and provides a single
place to implement theme overrides later.

For MVP, styles will use terminal defaults where possible (no hardcoded colours)
with only structural styling (padding, borders, width) defined explicitly.

### 7. Keybindings

A single shared `KeyMap` struct is used for MVP. Per-screen keymaps are
considered premature optimisation at this stage. The keymap will be loaded
from `Config` in Phase 2 to support user customisation.

Default bindings follow vim conventions: `j`/`k` to move, `Enter` to open,
`ESC` to go back, `q` to quit.

---

## Future architectural considerations

### Article rendering (`internal/render` or `internal/article`)

Article content rendering is expected to become a dedicated concern as the
project grows. Potential responsibilities include:

- Markdown or HTML rendering within the terminal
- Readability extraction from fetched HTML
- Syntax highlighting for code blocks in articles
- Wrapping and reflowing text for variable terminal widths

**This is not implemented in the MVP.** It is noted here so that rendering
logic is not casually embedded in `internal/screens/articleview` in a way
that becomes difficult to extract later.

When the need arises, a dedicated package (`internal/render` or
`internal/article`) should be introduced rather than growing `articleview`
beyond its remit as a display component.

### API client (`internal/api`)

When backend integration begins, an `internal/api` package will be introduced.
It will consume the same `domain.Feed` and `domain.Article` types, replacing
the mock implementation in `internal/data`. The screens should require no
changes for this transition.

### Configuration (`internal/config`)

Config loading is stubbed for MVP. Phase 2 will implement full XDG-compliant
loading from `$XDG_CONFIG_HOME/cerrynt/config.yaml`, covering API URL,
auth token, keybindings, and theme overrides.

---

## Consequences

### Positive

- Package naming accurately reflects intent (`domain`, `screens`, `data`, `styles`).
- Screen lifecycle is simple: construct on navigate, discard on back.
- Style centralisation prevents drift and enables future theming.
- The domain types are stable across MVP → API transition.
- Bubble Tea message passing keeps screens fully decoupled and unit-testable.
- Complexity is deferred to phases where it is actually needed.

### Negative / trade-offs

- Scroll position is lost on back navigation in MVP. Accepted consciously.
- A single `KeyMap` means all screens share the same binding definitions,
  even if some keys are only meaningful on certain screens. Acceptable for MVP.
- `internal/data` currently contains only a mock implementation. The package
  currently contains mock data only, but is expected to become the entry point
  for future data sources.

### Risks

- Bubble Tea's component model is unfamiliar. The early phases will involve
  learning the framework as much as building features. Keeping MVP small
  directly mitigates this risk.
- Future screen-specific keybindings (e.g., article scrolling vs. list
  navigation) may require refactoring `internal/keymap`. This is low risk
  given the keymap package is small and self-contained.
