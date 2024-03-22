/*
Copyright © 2024 @erikwj
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/erikwj/brokenlinks/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "brokenlinks",
	Short: "A cli to validate a markdown tree for broken links",
	Long: `A cli to validate a markdown tree for broken links

	Currently support for:
	- image links in png, svg, or gif format
	- web links [manually for now]
	- file links in same directory
	- internal references to [other] markdown files headers
	`,
	// Execution
	Run: func(cmd *cobra.Command, args []string) {
		directory := dir
		extension := ext

		// validate that directory is not empty
		if directory == "" {
			fmt.Println("Error: directory is required")
			// print usage
			_ = cmd.Usage()
			os.Exit(1)
		}

		f := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == extension {
				if verbose {
					fmt.Fprintf(cmd.OutOrStdout(), "# Validating %s \n", path)
				}

				if err := internal.ValidateLinks(path, extension, errors_only); err != nil {
					fmt.Printf("# Error validating links in file %s: %v\n", path, err)
				}
			}
			return nil
		}
		err := filepath.Walk(directory, f)

		if err != nil {
			fmt.Printf("# Error walking the path %s: %v\n", directory, err)
			os.Exit(1)
		}
	},
}

var (
	dir         string
	ext         string
	verbose     bool
	errors_only bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	RootCmd.PersistentFlags().StringVar(&ext, "ext", ".md", "File extension to be filtered on")
	RootCmd.PersistentFlags().StringVar(&dir, "dir", "", "Required: directory to be checked")
	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Optional: print file names that are being checked; default: false")
	RootCmd.PersistentFlags().BoolVar(&errors_only, "errors_only", false, "Optional: print only errors, no weblinks; default: false")

}
