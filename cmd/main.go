package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/xth0m/misen/internal/parser"
	"github.com/xth0m/misen/internal/recipe"
	"github.com/xth0m/misen/internal/util"
)

func loadRecipe(filename string) *recipe.Recipe {
	p := parser.NewParser(filename)
	r, err := p.Parse()

	if err != nil {
		util.Panic(2, err.Error())
	} else if len(p.Errors) > 0 {
		for _, e := range p.Errors {
			fmt.Fprintln(os.Stderr, e)
		}

		os.Exit(1)
	}

	return r
}

func viewCmd(ctx context.Context, c *cli.Command) error {

	r := loadRecipe(c.StringArg("filename"))
	r.Print(os.Stdout)

	return nil
}

func main() {
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
					&cli.StringArg{Name: "filename"},
				},

				Action: viewCmd,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
