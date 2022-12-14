/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"darktool/settings"
	"darktool/tools/validator"
	"fmt"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate freelancer config files for being correct and if not, try to fix",
	Long: `Freelancer folder is automatically discovered in any child folders
or you can set its location with ENV variable DARKTOOL_FREELANCER_FOLDER`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validate called")
		validator.Run()
		fmt.Println("OK")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.PersistentFlags().BoolVarP(&settings.DryRun, "dry", "d", false, "enable dry for checks without writing to file / good for CI")
	validateCmd.PersistentFlags().StringVarP(&settings.FreelancerFreelancerLocation, "search", "s", settings.FreelancerFreelancerLocation, "Freelancer location to search for validate running")
}
