package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "", SupplementalQuotationMarks(""))
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
	assert.Equal(t, "`s`", SupplementalQuotationMarks("s"))
}

func TestLowerCaseMapAdd(t *testing.T) {

	cases := []struct {
		rawKey       string
		lowerCaseKey string
		expected     bool
	}{
		{"idx_1", "idx_1", true},
		{"idx_1", "IDX_1", false},
		{"IDX_1", "IDX_1", false},
		{"IDX_1", "idx_1", true},
		{"IDX_1", "idx_2", false},
	}

	for i := range cases {
		c := cases[i]
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := LowerCaseMap{}
			m.Add(c.rawKey)
			_, exist := m[c.lowerCaseKey]
			assert.Equal(t, exist, c.expected)
		})
	}

}

func TestLowerCaseMapExist(t *testing.T) {
	cases := []struct {
		rawKey   string
		paramKey string
		expected bool
	}{
		{"IDX_1", "idx_1", true},
		{"idx_1", "idx_1", true},
		{"idx_1", "IDX_1", true},
		{"IDX_1", "IDX_1", true},
		{"IDX_1", "idx_2", false},
		{"idx_1", "idx_2", false},
		{"idx_1", "IDX_2", false},
		{"IDX_1", "IDX_2", false},
	}

	for i := range cases {
		c := cases[i]
		t.Run("", func(t *testing.T) {
			m := LowerCaseMap{}
			m.Add(c.rawKey)
			assert.Equal(t, m.Exist(c.paramKey), c.expected)
		})
	}

}

func TestLowerCaseMapDelete(t *testing.T) {
	cases := []struct {
		rawKey              string
		paramKey            string
		deletedSuccessfully bool
	}{
		{"idx_1", "idx_1", true},
		{"IDX_1", "idx_1", true},
		{"IDX_1", "IDX_1", true},
		{"idx_1", "IDX_1", true},
	}

	for i := range cases {
		c := cases[i]
		t.Run("", func(t *testing.T) {
			m := LowerCaseMap{}
			m.Add(c.rawKey)
			m.Delete(c.paramKey)
			_, exist := m[c.rawKey]
			assert.Equal(t, !exist, c.deletedSuccessfully)
		})
	}
}

func Test_IsClosed(t *testing.T) {
	c1 := make(chan struct{})
	if IsClosed(c1) {
		t.Error("channel should not be closed")
	}
	close(c1)
	if !IsClosed(c1) {
		t.Error("channel should be closed")
	}
	if !IsClosed(nil) {
		t.Error("nil channel should be deemed as closed")
	}
	c2 := make(chan struct{}, 1)
	c2 <- struct{}{}
	if IsClosed(c2) {
		t.Error("c2 is not closed")
	}
	close(c2)
	if !IsClosed(c2) {
		t.Error("c2 is closed")
	}
}

func TestIncrementalAverageFloat64(t *testing.T) {
	give := []float64{1, 2, 3, 3, 2, 1}
	want := []float64{1, 1.5, 2, 2.25, 2.2, 2}
	var average float64 = 0
	var count int = 0
	for index := range give {
		average = IncrementalAverageFloat64(average, give[index], count, 1)
		assert.Equal(t, average, want[index])
		count++
	}
}

func TestMaxFloat64(t *testing.T) {
	var lessThan [2]float64 = [2]float64{1.111, 2.222}
	var moreThan [2]float64 = [2]float64{2.222, 1.111}
	var equal [2]float64 = [2]float64{2.222, 2.222}
	assert.Equal(t, float64(2.222), MaxFloat64(lessThan[0], lessThan[1]))
	assert.Equal(t, float64(2.222), MaxFloat64(moreThan[0], moreThan[1]))
	assert.Equal(t, float64(2.222), MaxFloat64(equal[0], equal[1]))
}

func TestIsGitHttpURL(t *testing.T) {

	trueCases := []string{
		"https://github.com/golang/go.git",
		"http://github.com/user/repo.git",
	}

	falseCases := []string{
		"https://github.com/user/repo",
		"ftp://github.com/user/repo.git",
		"git@github.com:user/repo.git",
	}

	for _, tc := range trueCases {
		assert.True(t, IsGitHttpURL(tc), "Expected %q to be a valid Git Http URL", tc)
	}

	for _, tc := range falseCases {
		assert.False(t, IsGitHttpURL(tc), "Expected %q to be an invalid Git Http URL", tc)
	}
}

func TestFullFuzzySearchRegexp(t *testing.T) {
	testCases := []struct {
		input       string
		wantMatch   []string
		wantNoMatch []string
	}{
		{
			"Hello",
			[]string{"heyHelloCode", "HElLO", "Sun_hello", "HelLo_Jack"},
			[]string{"GoLang is awesome", "I love GOLANG", "GoLangGOLANGGolang"},
		},
		{
			"Golang",
			[]string{"GoLang is awesome", "I love GOLANG", "GoLangGOLANGGolang"},
			[]string{"language", "hi", "heyHelloCode", "HElLO", "Sun_hello", "HelLo_Jack"},
		}, {
			".*(?i)",
			[]string{"GoLang .*(?i) awesome", "I love GO^.*(?i)SING", "GoLangGO.*(?i)Golang"},
			[]string{"language", "hi", "heyHelloCode", "HElLO", "Sun_hello", "HelLo_Jack"},
		},
	}

	for _, tc := range testCases {
		reg := FullFuzzySearchRegexp(tc.input)

		// Positive cases
		for _, s := range tc.wantMatch {
			if !reg.MatchString(s) {
				t.Errorf("Expected %q to match %v", s, reg)
			}
		}

		// Negative cases
		for _, s := range tc.wantNoMatch {
			if reg.MatchString(s) {
				t.Errorf("Expected %q NOT to match %v", s, reg)
			}
		}
	}
}
