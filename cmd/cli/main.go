package main

import (
	"os"

	"github.com/bluetuith-org/daemon/cmd/app"
)

func main() {
	app.New().Run(os.Args)
}
