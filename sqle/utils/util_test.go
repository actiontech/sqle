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

func TestRound(t *testing.T) {
	assert.Equal(t, float64(1), Round(1.11, 0))
	assert.Equal(t, float64(0), Round(1.111117, -2))
	assert.Equal(t, 1.1, Round(1.11, 1))
	assert.Equal(t, 1.11112, Round(1.111117, 5))
}

func TestSupplementalQuotationMarks(t *testing.T) {
	assert.Equal(t, "'asdf'", SupplementalQuotationMarks("'asdf'"))
	assert.Equal(t, "\"asdf\"", SupplementalQuotationMarks("\"asdf\""))
	assert.Equal(t, "`asdf`", SupplementalQuotationMarks("`asdf`"))
	assert.Equal(t, "``", SupplementalQuotationMarks(""))
	assert.Equal(t, "`asdf`", SupplementalQuotationMarks("asdf"))
	assert.Equal(t, "`\"asdf`", SupplementalQuotationMarks("\"asdf"))
	assert.Equal(t, "`asdf\"`", SupplementalQuotationMarks("asdf\""))
	assert.Equal(t, "`'asdf`", SupplementalQuotationMarks("'asdf"))
	assert.Equal(t, "`asdf'`", SupplementalQuotationMarks("asdf'"))
	assert.Equal(t, "``asdf`", SupplementalQuotationMarks("`asdf"))
	assert.Equal(t, "`asdf``", SupplementalQuotationMarks("asdf`"))
	assert.Equal(t, "`\"asdf'`", SupplementalQuotationMarks("\"asdf'"))
	assert.Equal(t, "`\"asdf``", SupplementalQuotationMarks("\"asdf`"))
	assert.Equal(t, "`'asdf\"`", SupplementalQuotationMarks("'asdf\""))
	assert.Equal(t, "`'asdf``", SupplementalQuotationMarks("'asdf`"))
	assert.Equal(t, "``asdf\"`", SupplementalQuotationMarks("`asdf\""))
	assert.Equal(t, "``asdf'`", SupplementalQuotationMarks("`asdf'"))
}
