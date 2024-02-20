package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/urfave/cli/v2"
)

var files = make(map[[sha256.Size]byte][]string)

func checkDuplicate(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil
	}
	if info.IsDir() {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil
	}
	digest := sha256.Sum256(data)
	files[digest] = append(files[digest], path)

	return nil
}

func main() {
	var delCount int
	app := &cli.App{
		Name:  "fdup",
		Usage: "Find duplicate files",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Add Delete commands to duplicate files",
				Value:   false,
				Count:   &delCount,
			},
			&cli.IntFlag{
				Name:    "matchcount",
				Aliases: []string{"m"},
				Usage:   "Minimum count of duplicate files to show",
				Value:   1,
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.Int("matchcount") < 1 {
				return fmt.Errorf("matchcount must be greater than 0")
			}
			dirs := cCtx.Args().Slice()
			if len(dirs) == 0 {
				dirs = append(dirs, ".")
			}

			resfiles := make(map[[sha256.Size]byte][]string)

			for _, dir := range dirs {
				err := filepath.Walk(dir, checkDuplicate)
				if err != nil {
					return err
				}
				for digest, v := range files {
					if len(v) > cCtx.Int("matchcount") {
						resfiles[digest] = v
					}
				}
			}

			for _, filelist := range resfiles {
				slices.Sort(filelist)
				for k, filename := range filelist {
					if cCtx.Bool("delete") && delCount == 1 {
						fmt.Println("#rm -vf \"" + filename + "\"")
					} else if cCtx.Bool("delete") && delCount > 1 && k == 0 {
						fmt.Println("#rm -vf \"" + filename + "\"")
					} else if cCtx.Bool("delete") && delCount > 1 && k > 0 {
						fmt.Println("rm -vf \"" + filename + "\"")
					} else {
						fmt.Println(filename)
					}
				}
				fmt.Println()
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
