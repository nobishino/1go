package c_test

import (
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/nobishino/1go/c"
)

func TestTokenize(t *testing.T) {
	testcases := [...]struct {
		title  string
		in     string
		expect []string
	}{
		{
			title:  "1 term",
			in:     "5",
			expect: []string{"5"},
		},
		{
			title:  "2 terms#1",
			in:     "5+3",
			expect: []string{"5", "+", "3"},
		},
		{
			title:  "2 terms#2",
			in:     "5-31",
			expect: []string{"5", "-", "31"},
		},
		{
			title:  "3 terms",
			in:     "5+13-4",
			expect: []string{"5", "+", "13", "-", "4"},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			got := c.Tokenize(tt.in)
			if diff := cmp.Diff(got, tt.expect); diff != "" {
				t.Errorf("differs: (-got +expect)\n%s", diff)
			}
		})
	}
}

func TestAddSub(t *testing.T) {
	testcases := [...]struct {
		in     string
		expect string
	}{
		{
			in:     "5+20-4",
			expect: "./testdata/1.s",
		},
	}
	for _, tt := range testcases {
		got := c.Compile(tt.in)
		expect := readTestFile(t, tt.expect)
		if got != expect {
			t.Errorf("\n[expect]\n%s\n[got]\n%s", expect, got)
		}
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	s, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(s)
}
