package command

import (
	`context`
	`fmt`
	`io`
	`io/fs`
	`os`
	`path/filepath`
	`strings`

	yup "github.com/gloo-foo/framework"
)

// Dir represents a directory path
type Dir string

type command yup.Inputs[Dir, flags]

func Find(parameters ...any) yup.Command {
	// Initialize recognizes Dir and just parses (no file opening)
	return command(yup.Initialize[Dir, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Determine starting directories
		searchPaths := p.Positional
		if len(searchPaths) == 0 {
			searchPaths = []Dir{"."}
		}

		// Process each starting path
		for _, startPath := range searchPaths {
			startPathStr := string(startPath)
			// Clean the path
			startPathStr = filepath.Clean(startPathStr)

			// Check if path exists
			info, err := os.Stat(startPathStr)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "find: %s: %v\n", startPathStr, err)
				continue
			}

			// If it's a single file, check if it matches filters
			if !info.IsDir() {
				if p.matchesFilters(startPathStr, info) {
					_, _ = fmt.Fprintln(stdout, startPathStr)
				}
				continue
			}

			// Walk directory tree
			err = filepath.WalkDir(startPathStr, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					_, _ = fmt.Fprintf(stderr, "find: %s: %v\n", path, err)
					return nil // Continue walking
				}

				// Calculate depth relative to start path
				relPath, _ := filepath.Rel(startPathStr, path)
				depth := 0
				if relPath != "." {
					depth = strings.Count(relPath, string(os.PathSeparator)) + 1
				}

				// Check MaxDepth
				if p.Flags.MaxDepth > 0 && depth > int(p.Flags.MaxDepth) {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				// Get file info
				info, err := d.Info()
				if err != nil {
					_, _ = fmt.Fprintf(stderr, "find: %s: %v\n", path, err)
					return nil
				}

				// Apply filters
				if p.matchesFilters(path, info) {
					_, _ = fmt.Fprintln(stdout, path)
				}

				return nil
			})

			if err != nil {
				_, _ = fmt.Fprintf(stderr, "find: %s: %v\n", startPathStr, err)
			}
		}

		return nil
	}
}

// matchesFilters checks if a file matches all the specified filters
func (p command) matchesFilters(path string, info fs.FileInfo) bool {
	// Name filter (pattern matching)
	if p.Flags.Name != "" {
		matched, err := filepath.Match(string(p.Flags.Name), filepath.Base(path))
		if err != nil || !matched {
			return false
		}
	}

	// fileType filter
	if p.Flags.Type != "" {
		typeStr := string(p.Flags.Type)
		switch typeStr {
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

	// Size filter (simplified: +n, -n, n)
	if p.Flags.Size != "" {
		// TODO: Implement size filtering logic if needed
		// This would require parsing +n, -n, n and comparing with info.Size()
	}

	return true
}
