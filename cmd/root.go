package cmd

import (
	"fmt"
	"os"

	"github.com/shirokurostone/hosts-manager/manager"
	"github.com/spf13/cobra"
)

var appname = "hosts-manager"
var version = "0.1.0"
var config *manager.Config

var rootCmd = &cobra.Command{
	Use:   "hosts-manager",
	Short: "hosts-manager",
	Long:  "hosts-manager",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config, _ = manager.InitializeConfig(configFilePath)
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		manager.SaveConfig(config)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", appname, version)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all hostsgroups",
	Long:  "List all hostsgroups",
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.ListGroup(config)
	},
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new hostgroup",
	Long:  "Create a new hostgroup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.NewGroup(config, args[0])
	},
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a specific hostgroup",
	Long:  "Edit a specific hostgroup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.EditGroup(config, args[0])
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove specific hostgroups",
	Long:  "Remove specific hostgroups",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.RemoveGroup(config, args)
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Print active hostgroups",
	Long:  "Print active hostgroups",
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.ShowGroup(config, args)
	},
}

var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate hostgroups",
	Long:  "Activate hostgroups",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.ActivateGroup(config, args)
	},
}

var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate hostgroups",
	Long:  "Deactivate hostgroups",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.DeactivateGroup(config, args)
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply hostgroups to the hosts file",
	Long:  "Apply hostgroups to the hosts file",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var hostsfile string
		if len(args) == 0 {
			hostsfile = "/etc/hosts"
		} else {
			hostsfile = args[0]
		}

		return manager.ApplyGroup(config, hostsfile, dryRunFlag)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var dryRunFlag bool
var configFilePath string

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(activateCmd)
	rootCmd.AddCommand(deactivateCmd)
	rootCmd.AddCommand(applyCmd)

	rootCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "", "config file path")
	applyCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "n", false, "run dry-run mode")
}
