package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasPrefix(t *testing.T) {
	type args struct {
		s             string
		prefix        string
		caseSensitive bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{s: "hello, world", prefix: "hel", caseSensitive: true}, true},
		{"", args{s: "hello, world", prefix: "HEL", caseSensitive: true}, false},
		{"", args{s: "hello, world", prefix: "hel", caseSensitive: false}, true},
		{"", args{s: "hello, world", prefix: "HEL", caseSensitive: false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPrefix(tt.args.s, tt.args.prefix, tt.args.caseSensitive); got != tt.want {
				t.Errorf("HasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
	type args struct {
		s             string
		suffix        string
		caseSensitive bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{s: "hello, world", suffix: "rld", caseSensitive: true}, true},
		{"", args{s: "hello, world", suffix: "RLD", caseSensitive: true}, false},
		{"", args{s: "hello, world", suffix: "rld", caseSensitive: false}, true},
		{"", args{s: "hello, world", suffix: "RLD", caseSensitive: false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasSuffix(tt.args.s, tt.args.suffix, tt.args.caseSensitive); got != tt.want {
				t.Errorf("HasSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDuplicate(t *testing.T) {
	assert.Equal(t, []string{}, GetDuplicate([]string{"1", "2", "3"}))
	assert.Equal(t, []string{"2"}, GetDuplicate([]string{"1", "2", "2"}))
	assert.Equal(t, []string{"2", "3"}, GetDuplicate([]string{"1", "2", "2", "3", "3", "3"}))
}

func TestRemoveDuplicate(t *testing.T) {
	assert.Equal(t, []string{"1", "2", "3"}, RemoveDuplicate([]string{"1", "2", "3"}))
	assert.Equal(t, []string{"1", "2", "3"}, RemoveDuplicate([]string{"1", "2", "2", "3"}))
	assert.Equal(t, []string{"1", "2", "3"}, RemoveDuplicate([]string{"1", "2", "2", "3", "3", "3"}))
}
