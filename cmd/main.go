package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/xth0m/misen/internal/parser"
	"github.com/xth0m/misen/internal/util"
)


func main() {
	var recipeFilename string

	cmd := &cli.Command{
		Name: "misen",
		Version: "v0.0.1",
		Usage: "A recipe markup language",

		Suggest: true,
		EnableShellCompletion: true,

		Commands: []*cli.Command{
			{
				Name: "view",
				Usage: "view a misen recipe in a formatted way",

				Arguments: []cli.Argument{
					&cli.StringArg{ Name: "filename", Destination: &recipeFilename },
				},

				Action: func(ctx context.Context, c *cli.Command) error {
					if err := parser.Parse(recipeFilename); err != nil {
						util.Panic(1, err.Error())
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

