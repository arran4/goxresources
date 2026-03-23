package xresources

import (
	"os"
	"path/filepath"
)

// ParseFile parses an Xresources document from a given file path.
func ParseFile(filepath string) (*Document, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file)
}

// LoadConfig holds the configuration for loading .Xresources files.
type LoadConfig struct {
	UseXDG      bool
	UseHomeDir  bool
	MergeSystem bool
	CustomPaths []string
}

// UseXDG represents an option to load from XDG paths.
type UseXDG bool

// UseHomeDir represents an option to load from the home directory.
type UseHomeDir bool

// MergeSystem represents an option to load from system paths.
type MergeSystem bool

// CustomPaths represents an option to load from specific file paths.
type CustomPaths []string

func defaultLoadConfig() LoadConfig {
	return LoadConfig{
		UseXDG:      false,
		UseHomeDir:  false,
		MergeSystem: false,
		CustomPaths: nil,
	}
}

// Load dynamically loads and merges `.Xresources` files based on the provided options.
// Options are applied using a type-switched variadic system.
func Load(opts ...any) (*Document, error) {
	cfg := defaultLoadConfig()

	for _, opt := range opts {
		switch o := opt.(type) {
		case UseXDG:
			cfg.UseXDG = bool(o)
		case UseHomeDir:
			cfg.UseHomeDir = bool(o)
		case MergeSystem:
			cfg.MergeSystem = bool(o)
		case CustomPaths:
			cfg.CustomPaths = o
		}
	}

	mergedDoc := &Document{}
	var pathsToLoad []string

	if cfg.MergeSystem {
		pathsToLoad = append(pathsToLoad, "/etc/X11/Xresources")
	}

	if cfg.UseHomeDir {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			pathsToLoad = append(pathsToLoad, filepath.Join(homeDir, ".Xresources"))
		}
	}

	if cfg.UseXDG {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				xdgConfigHome = filepath.Join(homeDir, ".config")
			}
		}
		if xdgConfigHome != "" {
			pathsToLoad = append(pathsToLoad, filepath.Join(xdgConfigHome, "X11", "Xresources"))
		}
	}

	pathsToLoad = append(pathsToLoad, cfg.CustomPaths...)

	for _, path := range pathsToLoad {
		doc, err := ParseFile(path)
		if err == nil && doc != nil {
			mergedDoc.Nodes = append(mergedDoc.Nodes, doc.Nodes...)
		}
	}

	return mergedDoc, nil
}
