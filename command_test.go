package command_test

import (
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/find"
)

func TestFind_Current(t *testing.T) {
	result := run.Quick(command.Find("."))
	assertion.NoError(t, result.Err)
	// Should find at least current directory
}

func TestFind_Name(t *testing.T) {
	result := run.Quick(command.Find(".", command.Name("*.go")))
	assertion.NoError(t, result.Err)
}

func TestFind_Type(t *testing.T) {
	result := run.Quick(command.Find(".", command.FileType))
	assertion.NoError(t, result.Err)
}

func TestFind_MaxDepth(t *testing.T) {
	result := run.Quick(command.Find(".", command.MaxDepth(1)))
	assertion.NoError(t, result.Err)
}

func TestFind_FollowSymlinks(t *testing.T) {
	result := run.Quick(command.Find(".", command.FollowSymlinks))
	assertion.NoError(t, result.Err)
}

