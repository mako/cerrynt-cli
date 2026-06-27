// Package state manages local application state that is persisted between
// sessions. It is concerned only with client-side state (e.g. read article
// IDs) and is independent of configuration, domain logic, and the UI.
package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// State holds all local application state that is written to disk.
type State struct {
	// Read maps article IDs to true for every article the user has read.
	// Absent keys are treated as unread; false values should not be written.
	Read map[string]bool `json:"read"`
}

// DefaultPath returns the default path for the cerrynt state file, following
// the XDG Base Directory specification:
//
//   - If $XDG_DATA_HOME is set to an absolute path, use $XDG_DATA_HOME/cerrynt/state.json.
//   - Otherwise fall back to $HOME/.local/share/cerrynt/state.json.
//
// DefaultPath only resolves the path; it does not create directories or files.
func DefaultPath() (string, error) {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" || !filepath.IsAbs(base) {
		// Non-absolute values are ignored per XDG spec.
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("state: resolve home directory: %w", err)
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "cerrynt", "state.json"), nil
}

// Load reads and parses the JSON state file at the given path.
// It returns an error for any problem, including when the file does not exist
// (os.IsNotExist will be true in that case). The caller is responsible for
// deciding what to do when the file is missing — for example, starting with
// an empty State rather than treating it as a fatal error.
func Load(path string) (*State, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("state: open %q: %w", path, err)
	}
	defer f.Close()

	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("state: parse %q: %w", path, err)
	}
	return &s, nil
}

// Save writes s to the given path atomically.
//
// Atomicity is achieved by writing to a temporary file in the same directory
// as path and then renaming it into place. Because the temporary file and the
// destination are on the same filesystem, the rename is atomic on Linux: a
// concurrent reader will always see either the old file or the new file,
// never a partial write.
//
// Save creates the parent directory (and any missing ancestors) if needed.
// It uses mode 0700 for directories and 0600 for the state file.
func Save(path string, s *State) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("state: create directory %q: %w", dir, err)
	}

	// Write to a temp file in the same directory as path so that the
	// subsequent rename stays on the same filesystem.
	tmp, err := os.CreateTemp(dir, ".state-*.json.tmp")
	if err != nil {
		return fmt.Errorf("state: create temp file in %q: %w", dir, err)
	}
	tmpPath := tmp.Name()

	// On any failure, remove the temp file so we don't leave debris.
	committed := false
	defer func() {
		if !committed {
			_ = os.Remove(tmpPath)
		}
	}()

	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ") // human-readable on disk for easier debugging
	if err := enc.Encode(s); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("state: encode to %q: %w", tmpPath, err)
	}

	// Sync before close to ensure data reaches the OS page cache.
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("state: sync %q: %w", tmpPath, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("state: close %q: %w", tmpPath, err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("state: rename %q to %q: %w", tmpPath, path, err)
	}
	committed = true

	// Restrict permissions on the state file (may contain article IDs that
	// reveal reading habits). The temp file inherits the umask; we chmod
	// after rename so the final file always has the correct mode.
	if err := os.Chmod(path, 0600); err != nil {
		return fmt.Errorf("state: chmod %q: %w", path, err)
	}

	return nil
}
