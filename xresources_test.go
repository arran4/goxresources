package xresources_test

import (
	"embed"
	"encoding/json"
	"io/fs"
	"path"
	"strings"
	"testing"

	"xresources"

	"golang.org/x/tools/txtar"
)

//go:embed testdata/txtar/*.txtar
var testdataFS embed.FS

type Options struct {
	Filter string `json:"filter"`
}

func TestTxtarCases(t *testing.T) {
	entries, err := fs.Glob(testdataFS, "testdata/txtar/*.txtar")
	if err != nil {
		t.Fatalf("glob fixtures: %v", err)
	}

	for _, fixture := range entries {
		fixture := fixture
		t.Run(strings.TrimSuffix(path.Base(fixture), ".txtar"), func(t *testing.T) {
			raw, err := testdataFS.ReadFile(fixture)
			if err != nil {
				t.Fatalf("read fixture %s: %v", fixture, err)
			}
			ar := txtar.Parse(raw)

			var input, expected string
			var opts Options

			for _, f := range ar.Files {
				name := strings.TrimSpace(f.Name)
				switch name {
				case "input.txt":
					input = string(f.Data)
				case "expected.txt":
					expected = string(f.Data)
				case "options.json":
					if len(f.Data) > 0 {
						if err := json.Unmarshal(f.Data, &opts); err != nil {
							t.Fatalf("unmarshal options.json: %v", err)
						}
					}
				}
			}

			// Parse the input
			doc, err := xresources.ParseString(input)
			if err != nil {
				t.Fatalf("failed to parse input: %v", err)
			}

			// Apply filter if specified
			if opts.Filter != "" {
				doc = doc.Filter(opts.Filter)
			}

			// Check circularity / expected string
			got := doc.String()
			
			// For robust comparison, normalize leading/trailing newlines
			gotTrimmed := strings.TrimSpace(got)
			expectedTrimmed := strings.TrimSpace(expected)

			if gotTrimmed != expectedTrimmed {
				t.Errorf("mismatch!\nWant:\n%s\n\nGot:\n%s", expectedTrimmed, gotTrimmed)
			}
		})
	}
}
