package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/xth0m/misen/internal/parser"
	"github.com/xth0m/misen/internal/recipe"
	"github.com/xth0m/misen/internal/util"

	"github.com/logrusorgru/aurora/v4"
	"github.com/urfave/cli/v3"
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
	if c.StringArg("filename") == "" {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	if c.Float("scale") <= 0 {
		util.Panic(1, "scale value must be > 0")
	}

	r := loadRecipe(c.StringArg("filename"))
	r.Scale = c.Float("scale")

	recipe.Au = aurora.New(aurora.WithColors(!c.Bool("no-colors")))
	r.Print(os.Stdout)

	return nil
}

func doctorCmd(ctx context.Context, c *cli.Command) error {
	if c.StringArg("filename") == "" {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	loadRecipe(c.StringArg("filename"))

	return nil
}

func main() {
	cmd := &cli.Command{
		Name:    "misen",
		Version: "v0.1.0",
		Usage:   "A recipe markup language",

		Suggest:               true,
		EnableShellCompletion: true,

		Commands: []*cli.Command{
			{
				Name:  "view",
				Usage: "view a misen recipe in a formatted way",

				Arguments: []cli.Argument{
					&cli.StringArg{ Name: "filename" },
				},

				Flags: []cli.Flag{
					&cli.BoolFlag{ Name: "no-colors", Value: false },
					&cli.FloatFlag{
						Name: "scale",
						Value: 1,
					},
				},

				Action: viewCmd,
			},
			{
				Name: "doctor",
				Usage: "check if a misen recipe is correct",

				Arguments: []cli.Argument{
					&cli.StringArg{ Name: "filename" },
				},

				Action: doctorCmd,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
