package lib_test

import (
	"testing"

	"go-git-finder/lib"
)

func TestGithubRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"classic PAT", "x-access-token:ghp_1234567890abcdef1234567890abcdef123456", 1},
		{"no token", "no token here", 0},
		{"fine-grained PAT", "token: github_pat_11ABCDEF0_xyz1234567890123456789012345678", 1},
		{"oauth token", "gho_abcdefghijklmnopqrstuvwxyz1234", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := lib.GithubRegex.FindAllString(tt.input, -1)
			if len(matches) != tt.expected {
				t.Errorf("expected %d matches, got %d for input %q", tt.expected, len(matches), tt.input)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"example.com", "http://example.com/.git/config"},
		{"https://example.com/", "https://example.com/.git/config"},
		{"  http://test.com  ", "http://test.com/.git/config"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := lib.ParseURL(tt.input)
			if got != tt.expected {
				t.Errorf("ParseURL(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
