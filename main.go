package main

import (
	"context"

	"github.com/quail-ink/quail-cli/cmd"
)

func main() {
	ctx := context.Background()
	cmd.ExecuteContext(ctx)
}
