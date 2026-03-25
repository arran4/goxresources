# xresources

A robust Go library for parsing, formatting, and filtering `.Xresources` files.

## Why is this important?

`.Xresources` has historically been the standard way to configure core X11 applications. By supporting a robust parsing method, configuration tools can interact deeply with user setups without requiring destructive overhauls of their files. Preserving comments, spacing, and application-specific blocks keeps user settings intact while enabling programmatic configuration updates.

## When Should You Use It?

The `xresources` library is perfect for you when:
- You need to build a configuration tool that manipulates `.Xresources` files programmatically. This is extremely useful for a class of applications like terminal emulators (e.g., `XTerm`, `URxvt`) or window managers (e.g., `i3`, `dwm`, or `xmonad`) that rely on system-wide or user-level resource settings.
- You want to extract all settings related to a specific application and modify them, without dropping the user's surrounding comments or breaking other applications.
- (Note: On Wayland, there is no direct equivalent to a central `.Xresources` file, as configurations are mostly decentralized or managed via standard configuration files per application, often using formats like TOML, YAML, or INI located in `~/.config/`).

## Features
- Full support for comments (`!`), preprocessor macros (`#`), and blank lines.
- AST-based parsing that allows programmatic manipulation of `Key: Value` configurations.
- Handles multi-line values continuing with a backslash `\`.
- Application-specific prefix filtering (`Filter("AppPrefix")`) which preserves comments tightly coupled to an application's resource section.

## How to Use It

### Installation

#### Go get (Library)
```bash
go get github.com/arran4/goxresources
```

### 1. Parsing and Formatting

You can read an `.Xresources` file and write it back symmetrically:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/arran4/goxresources"
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
	
	"github.com/arran4/goxresources"
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

### 3. Loading and Merging Files Automatically

You can load and merge settings dynamically from different sources like XDG config paths and home directories using our flexible variadic loader:

```go
package main

import (
	"fmt"
	"log"

	"github.com/arran4/goxresources"
)

func main() {
	// Auto load the user's .Xresources from their home directory
	doc, err := xresources.Load(
		xresources.UseHomeDir(true),
		// You can also load an arbitrary user's .Xresources by specifying a custom path:
		// xresources.CustomPaths{"/home/otheruser/.Xresources"},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(doc.String())
}
```

## Common Properties

While applications can define any resource, some are widely standardized across X11 systems:

### Font & Display (Xft)
- `Xft.dpi`: Sets the DPI (Dots Per Inch) scaling (e.g., `96`, `192`).
- `Xft.antialias`: Enables font anti-aliasing (`1` or `true`, `0` or `false`).
- `Xft.hinting`: Enables font hinting (`1` or `true`, `0` or `false`).
- `Xft.hintstyle`: Sets hinting level (`hintnone`, `hintslight`, `hintmedium`, `hintfull`).
- `Xft.rgba`: Subpixel rendering type (`rgb`, `bgr`, `vrgb`, `vbgr`, or `none`).
- `Xft.lcdfilter`: LCD filter type (e.g., `lcddefault`).
- `Xft.autohint`: Enables autohinter (`1` or `true`, `0` or `false`).

### Mouse Cursors
- `Xcursor.theme`: The cursor theme name (e.g., `Adwaita`, `Vanilla-DMZ`).
- `Xcursor.size`: The size of the cursor in pixels (e.g., `16`, `32`, `48`).

### Common Terminal Emulators
#### XTerm
- `XTerm*faceName`: The TrueType font used (e.g., `xft:Monospace:size=10`).
- `XTerm*faceSize`: Size of the font.
- `XTerm*background`: Background color (e.g., `#282C34`).
- `XTerm*foreground`: Foreground text color.

#### URxvt (rxvt-unicode)
- `URxvt.font`: Main font string (e.g., `xft:Monospace:size=10`).
- `URxvt.scrollBar`: Toggles the scrollbar (`true`, `false`).
- `URxvt.geometry`: Sets initial window size (`80x24`).
- `URxvt.transparent`: Enables pseudo-transparency.
- `URxvt.depth`: Used for true transparency (e.g., `32`).

### Colors
Most terminal emulators share a common color palette scheme:
- `*color0` to `*color7`: Standard ANSI colors (Black, Red, Green, Yellow, Blue, Magenta, Cyan, White).
- `*color8` to `*color15`: Bright ANSI colors.
- `*background`: Default background.
- `*foreground`: Default foreground.

## More Applications

Many other X11 applications can be configured via `.Xresources`. For application-specific configuration details, refer to their official documentation:

- **XScreenSaver**: [XScreenSaver Documentation](https://www.jwz.org/xscreensaver/man.html#10)
- **xpdf**: [xpdf Manual](https://www.xpdfreader.com/xpdf-man.html)
- **rxvt-unicode (URxvt)**: [urxvt(1) Manual Page](http://pod.tst.eu/http://cvs.schmorp.de/rxvt-unicode/doc/rxvt.1.pod)
- **XTerm**: [xterm(1) Manual Page](https://invisible-island.net/xterm/manpage/xterm.html#h3-RESOURCES)
- **rofi**: [rofi(1) Manual Page (Legacy Xresources)](https://man.archlinux.org/man/rofi.1#XRESOURCES_AND_THEMING)

## Syntax of the Document

The `.Xresources` document syntax uses the following elements:
- **Comments**: Lines starting with `!` are ignored or treated as comments.
- **Preprocessor Directives**: Lines starting with `#` are used for `#define`, `#include`, etc.
- **Resources**: Key-value pairs defined as `Key: Value`, where `Key` identifies an application or resource path (often separated by `*` or `.`), and `Value` is the content.
- **Line Continuations**: Multi-line strings can be formed using a trailing `\` at the end of a line.
- **Empty Lines**: Ignored functionally, but preserved by this parser to ensure identical rewrites.
