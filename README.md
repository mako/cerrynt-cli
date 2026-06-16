# Cerrynt

Cerrynt is a terminal-first RSS reader ecosystem written in Go using Bubble Tea and Lip Gloss.

It is designed as part of a larger system that will include:
- a Rails backend (data source of truth)
- a future web UI
- potential mobile clients

The CLI is intentionally a thin client focused on fast navigation and readability.

---

## Current Status

🚧 Early MVP (Phase 1)

Implemented:
- Feed list navigation
- Article list navigation
- Article reader view
- Vim-style keybindings (j/k/Enter/Esc)
- Mock data layer
- Screen-based architecture using Bubble Tea

---

## Architecture

- Go (CLI)
- Bubble Tea (TUI framework)
- Lip Gloss (styling)
- Rails backend (planned)

The application currently uses static mock data.
No networking or persistence is implemented yet by design.

---

## Design Goals

- Fast, keyboard-driven RSS reading experience
- Vim-like navigation
- Clean separation of UI and domain logic
- Backend-agnostic CLI (sync layer will come later)
- Extensible architecture for future multi-device usage

---

## Project Status

This project is under active development as a learning + portfolio project.

It is being built iteratively with AI-assisted development and strong architectural review.

---

## License

MIT (see LICENSE)
