package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ppreeper/fdup/finder"
	"github.com/urfave/cli/v2"
)

func main() {
	var delCount int
	app := &cli.App{
		Name:                   "fdup",
		Usage:                  "Find duplicate files by content hash",
		UseShortOptionHandling: true,
		Description: `fdup recursively walks one or more directories, groups files by
content hash, and reports groups of duplicates.

Delete mode:
  -d     Print all files as commented-out rm commands (#rm -vf ...)
  -dd    Keep the first file (commented-out), emit active rm commands for the rest`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "emit rm commands for duplicates (-d comments all, -dd keeps first)",
				Value:   false,
				Count:   &delCount,
			},
			&cli.IntFlag{
				Name:    "matchcount",
				Aliases: []string{"m"},
				Usage:   "minimum number of files in a duplicate group",
				Value:   2,
			},
			&cli.BoolFlag{
				Name:  "skip-empty",
				Usage: "skip zero-length files",
				Value: false,
			},
		},
		Action: func(cCtx *cli.Context) error {
			matchcount := cCtx.Int("matchcount")
			if matchcount < 2 {
				return fmt.Errorf("matchcount must be at least 2")
			}

			dirs := cCtx.Args().Slice()
			if len(dirs) == 0 {
				dirs = []string{"."}
			}

			opts := &finder.Options{
				SkipEmpty: cCtx.Bool("skip-empty"),
				WarnFunc: func(msg string) {
					fmt.Fprintf(os.Stderr, "warning: %s\n", msg)
				},
			}

			groups, err := finder.Find(dirs, matchcount, opts)
			if err != nil {
				return err
			}

			for _, group := range groups {
				for i, path := range group.Paths {
					if delCount >= 1 {
						prefix := "rm -vf"
						if delCount == 1 || i == 0 {
							prefix = "#rm -vf"
						}
						fmt.Printf("%s %s\n", prefix, shellQuote(path))
					} else {
						fmt.Println(path)
					}
				}
				fmt.Println()
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// shellQuote returns a POSIX shell-safe single-quoted string.
// Single quotes protect all special characters; embedded single quotes
// are handled by ending the quote, inserting an escaped single quote,
// and restarting the quote.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
