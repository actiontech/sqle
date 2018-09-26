package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSectionMatch(t *testing.T) {
	ok, match := SectionMatch("a/b/c/d", "a/b/c/*")
	assert.True(t, ok)
	assert.Equal(t, "d", match[0])

	ok, match = SectionMatch("a/b/c/d", "a/b/*/d")
	assert.True(t, ok)
	assert.Equal(t, "c", match[0])

	ok, match = SectionMatch("a/b/c/d", "*/b/c/d")
	assert.True(t, ok)
	assert.Equal(t, "a", match[0])

	ok, match = SectionMatch("a/b/c/d", "a/b/c")
	assert.False(t, ok)

	ok, match = SectionMatch("a/b/c/d", "a/b/c/**")
	assert.True(t, ok)
	assert.Equal(t, "d", match[0])

	ok, match = SectionMatch("a/b/c/d", "a/**")
	assert.True(t, ok)
	assert.Equal(t, "b/c/d", match[0])

	ok, match = SectionMatch("a/b/c/d", "a/*/c/**")
	assert.True(t, ok)
	assert.Equal(t, "b", match[0])
	assert.Equal(t, "d", match[1])

	ok, match = SectionMatch("a/b/d", "a/b|c/d")
	assert.True(t, ok)
	assert.Equal(t, "b", match[0])

	ok, match = SectionMatch("a/b/d", "a/c|d/d")
	assert.False(t, ok)
	assert.Empty(t, match)

	ok, match = SectionMatch("a/b/d", "a/*/c|d")
	assert.True(t, ok)
	assert.Equal(t, "b", match[0])
	assert.Equal(t, "d", match[1])

	ok, match = SectionMatch("a/*/d", "a/\\*/c|d")
	assert.True(t, ok)
	assert.Equal(t, "d", match[0])

	ok, match = SectionMatch("a/b/c/d", "a/\\*/d")
	assert.False(t, ok)
	assert.Empty(t, match)
}
