package command

type Name string
type Size string
type MaxDepth int

type fileType string

const (
	BlockSpecialType     fileType = "b"
	CharacterSpecialType fileType = "c"
	DirectoryType        fileType = "d"
	FileType             fileType = "f"
	LinkType             fileType = "l"
	FIFOType             fileType = "p"
	SocketType           fileType = "s"
)

type FollowSymlinksFlag bool

const (
	FollowSymlinks   FollowSymlinksFlag = true
	NoFollowSymlinks FollowSymlinksFlag = false
)

type flags struct {
	Name           Name
	Type           fileType
	Size           Size
	MaxDepth       MaxDepth
	FollowSymlinks FollowSymlinksFlag
}

func (n Name) Configure(flags *flags)               { flags.Name = n }
func (t fileType) Configure(flags *flags)           { flags.Type = t }
func (s Size) Configure(flags *flags)               { flags.Size = s }
func (m MaxDepth) Configure(flags *flags)           { flags.MaxDepth = m }
func (f FollowSymlinksFlag) Configure(flags *flags) { flags.FollowSymlinks = f }
