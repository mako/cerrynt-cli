package state_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mako/cerrynt-cli/internal/state"
)

// -- DefaultPath tests -------------------------------------------------------

func TestDefaultPath_XDGDataHomeSet(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "/custom/data")

	got, err := state.DefaultPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "/custom/data/cerrynt/state.json"
	if got != want {
		t.Errorf("DefaultPath() = %q, want %q", got, want)
	}
}

func TestDefaultPath_XDGDataHomeUnset(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "")

	got, err := state.DefaultPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasSuffix(got, filepath.Join("cerrynt", "state.json")) {
		t.Errorf("DefaultPath() = %q, expected suffix %q", got, filepath.Join("cerrynt", "state.json"))
	}
	if !strings.Contains(got, filepath.Join(".local", "share")) {
		t.Errorf("DefaultPath() = %q, expected to contain %q", got, filepath.Join(".local", "share"))
	}
}

func TestDefaultPath_XDGDataHomeRelativeIgnored(t *testing.T) {
	// Per XDG spec, relative paths must be ignored; fall back to ~/.local/share.
	t.Setenv("XDG_DATA_HOME", "relative/data")

	got, err := state.DefaultPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.HasPrefix(got, "relative/data") {
		t.Errorf("DefaultPath() = %q, should have ignored relative XDG_DATA_HOME", got)
	}
	if !strings.HasSuffix(got, filepath.Join("cerrynt", "state.json")) {
		t.Errorf("DefaultPath() = %q, expected suffix %q", got, filepath.Join("cerrynt", "state.json"))
	}
}

func TestDefaultPath_XDGDataHomeAbsolutePreserved(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "/srv/appdata")

	got, err := state.DefaultPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "/srv/appdata/cerrynt/state.json"
	if got != want {
		t.Errorf("DefaultPath() = %q, want %q", got, want)
	}
}

// -- Load tests --------------------------------------------------------------

func TestLoad_MissingFile(t *testing.T) {
	_, err := state.Load("/nonexistent/path/state.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist to be wrapped, got: %v", err)
	}
}

func TestLoad_ValidState(t *testing.T) {
	content := `{"read": {"101": true, "202": true}}`
	path := writeTempState(t, content)

	s, err := state.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s.Read) != 2 {
		t.Fatalf("len(Read) = %d, want 2", len(s.Read))
	}
	if !s.Read["101"] {
		t.Error("Read[\"101\"] should be true")
	}
	if !s.Read["202"] {
		t.Error("Read[\"202\"] should be true")
	}
}

func TestLoad_EmptyReadMap(t *testing.T) {
	content := `{"read": {}}`
	path := writeTempState(t, content)

	s, err := state.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Read) != 0 {
		t.Errorf("len(Read) = %d, want 0", len(s.Read))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeTempState(t, `{not valid json}`)

	_, err := state.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// -- Save tests --------------------------------------------------------------

func TestSave_CreatesParentDirectories(t *testing.T) {
	dir := t.TempDir()
	// Use a path whose parent subdirectory does not exist yet.
	path := filepath.Join(dir, "a", "b", "c", "state.json")

	s := &state.State{Read: map[string]bool{"42": true}}
	if err := state.Save(path, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("state file not found after Save: %v", err)
	}
}

func TestSave_WritesCorrectContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")

	s := &state.State{Read: map[string]bool{"101": true, "303": true}}
	if err := state.Save(path, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read saved file: %v", err)
	}

	var got state.State
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("saved file is not valid JSON: %v\ncontent: %s", err, data)
	}
	if !got.Read["101"] || !got.Read["303"] {
		t.Errorf("unexpected content: %+v", got)
	}
}

func TestSave_Roundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")

	original := &state.State{Read: map[string]bool{"aaa": true, "bbb": true}}
	if err := state.Save(path, original); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	loaded, err := state.Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if len(loaded.Read) != len(original.Read) {
		t.Fatalf("roundtrip len mismatch: got %d, want %d", len(loaded.Read), len(original.Read))
	}
	for id := range original.Read {
		if !loaded.Read[id] {
			t.Errorf("Read[%q] missing after roundtrip", id)
		}
	}
}

func TestSave_NoTempFileLeftOnSuccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	s := &state.State{Read: map[string]bool{"1": true}}
	if err := state.Save(path, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("temp file left behind after successful Save: %s", e.Name())
		}
	}
}

func TestSave_FilePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")

	s := &state.State{Read: map[string]bool{}}
	if err := state.Save(path, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	// Expect owner read/write only (0600), masking off file-type bits.
	got := info.Mode().Perm()
	if got != 0600 {
		t.Errorf("file permissions = %04o, want 0600", got)
	}
}

func TestSave_OverwritesExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")

	// First write.
	if err := state.Save(path, &state.State{Read: map[string]bool{"old": true}}); err != nil {
		t.Fatalf("first Save error: %v", err)
	}

	// Second write should replace it cleanly.
	if err := state.Save(path, &state.State{Read: map[string]bool{"new": true}}); err != nil {
		t.Fatalf("second Save error: %v", err)
	}

	loaded, err := state.Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if loaded.Read["old"] {
		t.Error("old key should not be present after overwrite")
	}
	if !loaded.Read["new"] {
		t.Error("new key should be present after overwrite")
	}
}

// -- helpers -----------------------------------------------------------------

// writeTempState writes JSON content to a temp file and returns its path.
func writeTempState(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempState: %v", err)
	}
	return path
}
