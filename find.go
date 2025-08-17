package find

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	localopt "github.com/yupsh/find/opt"
	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
)

// Flags represents the configuration options for the find command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Find creates a new find command with the given parameters
func Find(parameters ...any) yup.Command {
	return command(opt.Args[string, Flags](parameters...))
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	// Default to current directory if no paths specified
	paths := c.Positional
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, path := range paths {
		if err := c.walkPath(ctx, path, stdout, stderr, 0); err != nil {
			fmt.Fprintf(stderr, "find: %s: %v\n", path, err)
		}
	}

	return nil
}

func (c command) walkPath(ctx context.Context, path string, output, stderr io.Writer, depth int) error {
	// Check max depth
	if c.Flags.MaxDepth > 0 && depth > int(c.Flags.MaxDepth) {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Check if current path matches criteria
	if c.matchesCriteria(path, info) {
		fmt.Fprintln(output, path)
	}

	// Recurse into directories
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			// Check for cancellation before each directory entry
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return err
			}

			subPath := filepath.Join(path, entry.Name())

			// Handle symlinks based on flag
			if entry.Type()&os.ModeSymlink != 0 && !bool(c.Flags.FollowSymlinks) {
				continue
			}

			if err := c.walkPath(ctx, subPath, output, stderr, depth+1); err != nil {
				fmt.Fprintf(stderr, "find: %s: %v\n", subPath, err)
			}
		}
	}

	return nil
}

func (c command) matchesCriteria(path string, info os.FileInfo) bool {
	// Check name pattern
	if c.Flags.Name != "" {
		matched, _ := filepath.Match(string(c.Flags.Name), filepath.Base(path))
		if !matched {
			return false
		}
	}

	// Check type
	if c.Flags.Type != "" {
		switch string(c.Flags.Type) {
		case "f":
			if !info.Mode().IsRegular() {
				return false
			}
		case "d":
			if !info.IsDir() {
				return false
			}
		case "l":
			if info.Mode()&os.ModeSymlink == 0 {
				return false
			}
		}
	}

	// Additional size, time checks would go here
	// This is a simplified implementation

	return true
}
