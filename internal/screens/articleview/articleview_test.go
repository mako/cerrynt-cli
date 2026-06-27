package articleview

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// clamp
// ---------------------------------------------------------------------------

func TestClamp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    int
		lo   int
		hi   int
		want int
	}{
		{"within range", 5, 0, 10, 5},
		{"at lower bound", 0, 0, 10, 0},
		{"at upper bound", 10, 0, 10, 10},
		{"below lower bound", -5, 0, 10, 0},
		{"above upper bound", 15, 0, 10, 10},
		{"lo equals hi", 7, 5, 5, 5},
		{"all zeros", 0, 0, 0, 0},
		{"negative range", -3, -10, -1, -3},
		{"v below negative range", -15, -10, -1, -10},
		{"v above negative range", 5, -10, -1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := clamp(tt.v, tt.lo, tt.hi)
			if got != tt.want {
				t.Errorf("clamp(%d, %d, %d) = %d, want %d", tt.v, tt.lo, tt.hi, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// maxOffset
// ---------------------------------------------------------------------------

func TestMaxOffset(t *testing.T) {
	t.Parallel()

	// Helper: produce a slice of n strings.
	makeLines := func(n int) []string {
		lines := make([]string, n)
		for i := range lines {
			lines[i] = "line"
		}
		return lines
	}

	// statusRows = 2, so visible = height - 2.
	// maxOffset = max(0, len(lines) - visible)
	tests := []struct {
		name   string
		nLines int
		height int
		want   int
	}{
		{"content fits exactly", 8, 10, 0},    // visible=8, 8 lines → 0
		{"content fits with room", 5, 10, 0},  // visible=8, 5 lines → 0
		{"content overflows by 1", 9, 10, 1},  // visible=8, 9 lines → 1
		{"content overflows by 5", 13, 10, 5}, // visible=8, 13 lines → 5
		{"zero lines", 0, 10, 0},
		{"height smaller than statusRows", 3, 1, 2}, // visible clamped to 1, 3-1=2
		{"height equals statusRows", 3, 2, 2},       // visible clamped to 1
		{"empty height", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := maxOffset(makeLines(tt.nLines), tt.height)
			if got != tt.want {
				t.Errorf("maxOffset(%d lines, height=%d) = %d, want %d",
					tt.nLines, tt.height, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// wordWrap
// ---------------------------------------------------------------------------

func TestWordWrap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		text  string
		width int
		// wantLines is the expected output as a single string with lines
		// separated by "|" for readability. Each segment is one slice element.
		wantLines []string
	}{
		{
			name:      "empty string",
			text:      "",
			width:     20,
			wantLines: nil,
		},
		{
			name:      "whitespace-only",
			text:      "   ",
			width:     20,
			wantLines: nil,
		},
		{
			name:      "single short word",
			text:      "hello",
			width:     20,
			wantLines: []string{"hello"},
		},
		{
			name:      "two words that fit on one line",
			text:      "hello world",
			width:     20,
			wantLines: []string{"hello world"},
		},
		{
			name:      "line break at word boundary",
			text:      "hello world",
			width:     5,
			wantLines: []string{"hello", "world"},
		},
		{
			name:      "exactly at width boundary",
			text:      "hi there",
			width:     8, // "hi there" is 8 chars — fits
			wantLines: []string{"hi there"},
		},
		{
			name:      "one char over width forces break",
			text:      "hi there!",
			width:     8, // "hi there!" = 9 chars — "there!" goes to next line
			wantLines: []string{"hi", "there!"},
		},
		{
			name:  "multiple lines",
			text:  "one two three four five",
			width: 9,
			// "one two" = 7, "three" = 5, "four" = 4, "five" = 4
			// "one two" fits (7 ≤ 9); adding "three" → 11 > 9 → break
			// "three four" = 10 > 9 → break after "three"
			// "four five" = 9 ≤ 9 → fits
			wantLines: []string{"one two", "three", "four five"},
		},
		{
			name:      "single word longer than width",
			text:      "superlongword",
			width:     5,
			wantLines: []string{"superlongword"}, // not split mid-word
		},
		{
			name:      "leading and trailing spaces are ignored",
			text:      "  hello world  ",
			width:     20,
			wantLines: []string{"hello world"},
		},
		{
			name:      "multiple internal spaces collapsed",
			text:      "hello   world",
			width:     20,
			wantLines: []string{"hello world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := wordWrap(tt.text, tt.width)

			if tt.wantLines == nil {
				if len(got) != 0 {
					t.Errorf("wordWrap(%q, %d) = %v, want nil/empty", tt.text, tt.width, got)
				}
				return
			}

			if len(got) != len(tt.wantLines) {
				t.Errorf("wordWrap(%q, %d): got %d lines, want %d\ngot:  %v\nwant: %v",
					tt.text, tt.width, len(got), len(tt.wantLines), got, tt.wantLines)
				return
			}

			for i, line := range got {
				if line != tt.wantLines[i] {
					t.Errorf("wordWrap(%q, %d) line %d = %q, want %q",
						tt.text, tt.width, i, line, tt.wantLines[i])
				}
			}
		})
	}
}

// TestWordWrapNoMidWordSplit verifies that no line ever exceeds the width
// unless a single word is itself longer than the width (which cannot be split).
func TestWordWrapNoMidWordSplit(t *testing.T) {
	t.Parallel()

	text := "The quick brown fox jumps over the lazy dog near the riverbank"
	width := 15

	lines := wordWrap(text, width)
	for _, line := range lines {
		// Reconstruct whether this was a forced overrun (single long word).
		isSingleWord := !strings.Contains(line, " ")
		if len(line) > width && !isSingleWord {
			t.Errorf("line %q exceeds width %d and is not a single word", line, width)
		}
	}

	// Reassemble and verify no words were lost.
	joined := strings.Join(lines, " ")
	if joined != text {
		t.Errorf("wordWrap lost or altered words\ngot:  %q\nwant: %q", joined, text)
	}
}
