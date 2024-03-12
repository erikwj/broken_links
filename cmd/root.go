/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "brokenlinks",
	Short: "A cli to validate a markdown tree for broken links",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		directory := dir
		f := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".md" {
				if debug {
					fmt.Printf("Validating %s \n", path)
				}

				if err := ValidateLinks(path); err != nil {
					fmt.Printf("Error validating links in file %s: %v\n", path, err)
				}
			}
			return nil
		}
		err := filepath.Walk(directory, f)

		if err != nil {
			fmt.Printf("Error walking the path %s: %v\n", directory, err)
			os.Exit(1)
		}
	},
}

var (
	dir   string
	debug bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&dir, "dir", "", "directory to be checked")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "expose detailed info on execution")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
