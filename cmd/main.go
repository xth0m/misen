package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/xth0m/misen/internal/parser"
)

func main() {
	var recipeFilename string

	cmd := &cli.Command{
		Name:    "misen",
		Version: "v0.0.1",
		Usage:   "A recipe markup language",

		Suggest:               true,
		EnableShellCompletion: true,

		Commands: []*cli.Command{
			{
				Name:  "view",
				Usage: "view a misen recipe in a formatted way",

				Arguments: []cli.Argument{
					&cli.StringArg{Name: "filename", Destination: &recipeFilename},
				},

				Action: func(ctx context.Context, c *cli.Command) error {

					parser := parser.NewParser(recipeFilename)
					r, err := parser.Parse()

					if err != nil {
					}

					if len(parser.Errors) > 0 {
						for _, e := range parser.Errors {
							fmt.Println(e)
						}

						return nil
					}

					r.Print(os.Stdout)

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
