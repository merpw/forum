package server

import (
	"testing"
)

func TestShortenContent(t *testing.T) {
	longcontent := `lorem ipsum dolor sit amet consectetur adipiscing
	 elit sed do eiusmod tempor incididunt ut labore et dolore magna
	  aliqua ut enim ad minim veniam quis nostrud exercitation ullamco 
	  laboris nisi ut aliquip ex ea commodo`
	shortened := `lorem ipsum dolor sit amet consectetur adipiscing
	 elit sed do eiusmod tempor incididunt ut labore et dolore magna
	  aliqua ut enim ad minim veniam quis nostrud exercitation ullamco 
	  laboris nisi...`
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Short content",
			content:  "This is a short content",
			expected: "This is a short content",
		},
		{
			name:     "Long content",
			content:  longcontent,
			expected: shortened,
		},
		{
			name:     "Empty content",
			content:  "",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := shortenContent(test.content)
			if result != test.expected {
				t.Errorf("Expected \"%s\", but got \"%s\"", test.expected, result)
			}
		})
	}
}
