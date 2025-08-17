package opt

// Custom types for parameters
type Name string
type Type string
type Size string
type MaxDepth int

// Boolean flag types with constants
type FollowSymlinksFlag bool
const (
	FollowSymlinks   FollowSymlinksFlag = true
	NoFollowSymlinks FollowSymlinksFlag = false
)

// Flags represents the configuration options for the find command
type Flags struct {
	Name           Name               // Name pattern to match
	Type           Type               // File type (f=file, d=directory)
	Size           Size               // Size constraint
	MaxDepth       MaxDepth           // Maximum depth to search
	FollowSymlinks FollowSymlinksFlag // Follow symbolic links
}

// Configure methods for the opt system
func (n Name) Configure(flags *Flags) { flags.Name = n }
func (t Type) Configure(flags *Flags) { flags.Type = t }
func (s Size) Configure(flags *Flags) { flags.Size = s }
func (m MaxDepth) Configure(flags *Flags) { flags.MaxDepth = m }
func (f FollowSymlinksFlag) Configure(flags *Flags) { flags.FollowSymlinks = f }
