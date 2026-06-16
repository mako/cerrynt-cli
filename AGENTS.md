# Cerrynt CLI — Agent Instructions

This project is a terminal-based RSS reader built in Go using the Bubble Tea framework.

It is part of a larger ecosystem:
- Cerrynt Web (Rails backend)
- Cerrynt API (Rails-based API)
- Future mobile client

The CLI is the primary real-time reading interface.

---

## Core Goals

- Provide a fast, keyboard-driven RSS reader experience in the terminal
- Support multiple devices via shared backend state
- Store read/unread state centrally (via API)
- Be configurable via XDG-compliant config files
- Respect terminal color schemes by default, but allow overrides

---

## Tech Stack

- Language: Go
- UI: Bubble Tea + Lip Gloss
- HTTP: standard library or net/http + optional client wrapper
- Config: XDG Base Directory spec
- Data: consumed via Cerrynt API (Rails backend)

---

## Architecture Principles

- CLI is a thin client — no persistent business logic
- Backend (Rails API) is source of truth for:
  - feeds
  - articles
  - read/unread state
  - authentication
- CLI should cache minimally and safely (optional)

---

## UI Guidelines

- Vim-like keybindings where appropriate
- Fully keyboard-driven navigation
- Support customizable keymaps
- Respect terminal theme by default:
  - Do NOT hardcode colors unless explicitly configured
  - Prefer terminal capability detection
- Use Bubble Tea + Lip Gloss idiomatically

---

## Configuration

- Config path:
  $XDG_CONFIG_HOME/cerrynt/config.yaml

- Must support:
  - API base URL
  - auth token
  - keybindings
  - theme overrides (optional)

---

## Authentication

- Token-based auth against Rails backend
- Token stored locally in config or secure storage (TBD)
- CLI must never store passwords

---

## Development Rules

- Prefer small, composable packages
- Avoid over-engineering early abstractions
- Keep domain logic minimal in CLI
- Any new feature should first consider backend responsibility

---

## AI Usage Guidelines

This project is AI-assisted (Gemini CLI), but:

- All generated code must be reviewed
- No direct acceptance of AI output without understanding
- Prefer iterative refinement over full rewrites
- Always maintain consistency with architecture rules

---

## Future Extensions

- Offline cache mode
- Full-text search
- Mobile client integration
- Plugin system for custom feeds
