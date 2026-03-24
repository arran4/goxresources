package xresources

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Document represents a complete Xresources file.
type Document struct {
	Nodes []Node
}

// Node represents a single element in an Xresources file.
type Node interface {
	isNode()
}

// Comment represents a comment line (starts with !).
type Comment struct {
	Text string
}

func (c Comment) isNode() {}

// PreprocessorDirective represents a preprocessor line (starts with #).
type PreprocessorDirective struct {
	Text string
}

func (p PreprocessorDirective) isNode() {}

// Resource represents a key-value pair.
type Resource struct {
	Key   string
	Value string
}

func (r Resource) isNode() {}

// EmptyLine represents a blank line.
type EmptyLine struct{}

func (e EmptyLine) isNode() {}

// Raw represents an unrecognized or malformed line, kept for circularity.
type Raw struct {
	Text string
}

func (r Raw) isNode() {}

// Parse parses an Xresources document from an io.Reader.
func Parse(r io.Reader) (*Document, error) {
	doc := &Document{}
	scanner := bufio.NewScanner(r)
	
	// Increase buffer size to handle long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var currentLine strings.Builder
	inContinuation := false

	for scanner.Scan() {
		line := scanner.Text()

		if inContinuation {
			if strings.HasSuffix(line, "\\") {
				currentLine.WriteString("\n")
				currentLine.WriteString(line[:len(line)-1])
				inContinuation = true
			} else {
				currentLine.WriteString("\n")
				currentLine.WriteString(line)
				inContinuation = false
				processLine(doc, currentLine.String())
				currentLine.Reset()
			}
			continue
		}

		if strings.HasSuffix(line, "\\") {
			currentLine.WriteString(line[:len(line)-1])
			inContinuation = true
			continue
		}

		processLine(doc, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if inContinuation {
		processLine(doc, currentLine.String())
	}

	return doc, nil
}

func processLine(doc *Document, line string) {
	trimmed := strings.TrimSpace(line)

	if trimmed == "" {
		doc.Nodes = append(doc.Nodes, EmptyLine{})
		return
	}

	if strings.HasPrefix(trimmed, "!") {
		doc.Nodes = append(doc.Nodes, Comment{Text: line})
		return
	}

	if strings.HasPrefix(trimmed, "#") {
		doc.Nodes = append(doc.Nodes, PreprocessorDirective{Text: line})
		return
	}

	colonIdx := strings.Index(line, ":")
	if colonIdx != -1 {
		// Key cannot have spaces before colon in a valid resource, but let's be lenient
		// or at least trim the key. 
		key := strings.TrimSpace(line[:colonIdx])
		
		// Value is everything after the colon, typically trimmed of leading spaces.
		value := line[colonIdx+1:]
		value = strings.TrimLeft(value, " \t")
		
		doc.Nodes = append(doc.Nodes, Resource{
			Key:   key,
			Value: value,
		})
		return
	}

	doc.Nodes = append(doc.Nodes, Raw{Text: line})
}

// ParseBytes is a helper to parse from a byte slice.
func ParseBytes(b []byte) (*Document, error) {
	return Parse(bytes.NewReader(b))
}

// ParseString is a helper to parse from a string.
func ParseString(s string) (*Document, error) {
	return Parse(strings.NewReader(s))
}

// WriteTo writes the document back to an io.Writer.
func (d *Document) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64
	
	writeStr := func(s string) error {
		n, err := io.WriteString(w, s)
		totalWritten += int64(n)
		return err
	}

	for _, node := range d.Nodes {
		switch n := node.(type) {
		case EmptyLine:
			if err := writeStr("\n"); err != nil {
				return totalWritten, err
			}
		case Comment:
			if err := writeStr(n.Text + "\n"); err != nil {
				return totalWritten, err
			}
		case PreprocessorDirective:
			if err := writeStr(n.Text + "\n"); err != nil {
				return totalWritten, err
			}
		case Raw:
			if err := writeStr(n.Text + "\n"); err != nil {
				return totalWritten, err
			}
		case Resource:
			// Resources with newlines need to be formatted back with escaping
			val := n.Value
			if strings.Contains(val, "\n") {
				lines := strings.Split(val, "\n")
				// join back with "\\\n"
				val = strings.Join(lines, "\\\n")
			}
			
			line := n.Key + ":\t" + val + "\n"
			if err := writeStr(line); err != nil {
				return totalWritten, err
			}
		}
	}
	
	return totalWritten, nil
}

// String returns the document as a string.
func (d *Document) String() string {
	var buf strings.Builder
	_, _ = d.WriteTo(&buf)
	return buf.String()
}

// Filter returns a new Document containing only the resources that match
// the given prefix (e.g., application name like "XTerm"). It preserves
// comments and empty lines that might be associated with the resources,
// but for a strict filter, it might only return matching resources.
// Here we return only the matching Resource nodes and any preceding comments.
func (d *Document) Filter(appPrefix string) *Document {
	newDoc := &Document{}
	
	// Ensure we match "App*" or "App."
	// We'll check if the resource key starts with the prefix followed by '*' or '.'
	matchPrefix := appPrefix
	
	// A basic block collection to keep comments above the matching resources
	var pendingComments []Node
	
	for _, node := range d.Nodes {
		switch n := node.(type) {
		case Comment, EmptyLine, PreprocessorDirective, Raw:
			pendingComments = append(pendingComments, n)
		case Resource:
			// Check if it matches
			if strings.HasPrefix(n.Key, matchPrefix+".") || strings.HasPrefix(n.Key, matchPrefix+"*") || n.Key == matchPrefix {
				// Matched, append pending comments and this resource
				newDoc.Nodes = append(newDoc.Nodes, pendingComments...)
				newDoc.Nodes = append(newDoc.Nodes, n)
				pendingComments = nil // reset
			} else {
				// Didn't match, discard pending comments
				pendingComments = nil
			}
		}
	}
	
	return newDoc
}
