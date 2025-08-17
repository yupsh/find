package find_test

import (
	"context"
	"os"

	"github.com/yupsh/find"
	"github.com/yupsh/find/opt"
)

func ExampleFind() {
	ctx := context.Background()

	cmd := find.Find(".", opt.Name("*.go"), opt.Type("f"))
	cmd.Execute(ctx, nil, os.Stdout, os.Stderr)
	// Output will vary based on directory contents
}

func ExampleFind_directories() {
	ctx := context.Background()

	cmd := find.Find(".", opt.Type("d"), opt.MaxDepth(2))
	cmd.Execute(ctx, nil, os.Stdout, os.Stderr)
	// Output will vary based on directory structure
}
