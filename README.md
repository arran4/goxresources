# xresources

A robust Go library for parsing, formatting, and filtering `.Xresources` files.

## When Should You Use It?

The `xresources` library is perfect for you when:
- You need to build a configuration tool that manipulates `.Xresources` files programmatically.
- You want to extract all settings related to a specific application (e.g., `XTerm` or `URxvt`) and modify them, without dropping the user's surrounding comments or breaking other applications.
- You require **circular testing**, meaning any file you parse can be written back exactly as it was, maintaining line continuations, spacing, and comments.

## Features
- Full support for comments (`!`), preprocessor macros (`#`), and blank lines.
- AST-based parsing that allows programmatic manipulation of `Key: Value` configurations.
- Handles multi-line values continuing with a backslash `\`.
- Application-specific prefix filtering (`Filter("AppPrefix")`) which preserves comments tightly coupled to an application's resource section.

## How to Use It

### Installation

```bash
go get github.com/your-org/xresources
```

### 1. Parsing and Formatting

You can read an `.Xresources` file and write it back symmetrically:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/your-org/xresources"
)

func main() {
	doc, err := xresources.ParseString(`
! XTerm settings
XTerm*faceName: Monospace
XTerm*faceSize: 10
	`)
	if err != nil {
		log.Fatal(err)
	}
	
	// Print it out exactly as it was read
	fmt.Print(doc.String())
}
```

### 2. Filtering by Application

Often you just want to grab the settings for a single app:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/your-org/xresources"
)

func main() {
    config := `
! General
*color0: black

! XTerm settings
XTerm*faceName: Monospace
XTerm*faceSize: 10

! URxvt settings
URxvt.font: xft:Monospace:size=10
`
	doc, err := xresources.ParseString(config)
	if err != nil {
		log.Fatal(err)
	}
	
	// Extract just the XTerm lines and associated comments
	xtermDoc := doc.Filter("XTerm")
	fmt.Print(xtermDoc.String())
    // Output will contain the XTerm settings with its leading comments.
}
```

## Testing Methodology

This library uses the `txtar` testing approach outlined by [Arran's Technical Blog](https://arran4.github.io/blog/post/2026/004-txtar-patterns-for-agents/).
Every scenario is written in a `testdata/txtar/*.txtar` fixture, representing isolated end-to-end setups where `options.json` acts as an operation configurator, and files map cleanly from `input.txt` -> operations -> `expected.txt` for highly reproducible agent and deterministic checks.
