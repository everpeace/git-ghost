package cmd

import (
	"errors"
	"fmt"
	"git-ghost/pkg/ghost/git"
	"os"

	"github.com/spf13/cobra"
)

type globalFlags struct {
	srcDir      string
	ghostPrefix string
	ghostRepo   string
	baseCommit  string
}

var (
	Version  string
	Revision string
)

var RootCmd = &cobra.Command{
	Use:   "git-ghost",
	Short: "git-ghost",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Use == "version" {
			return nil
		}
		err := validateEnvironment()
		if err != nil {
			return err
		}
		err = globalOpts.Validate()
		if err != nil {
			return err
		}
		return nil
	},
}

var globalOpts globalFlags

func init() {
	cobra.OnInitialize()
	currentDir := os.Getenv("PWD")
	RootCmd.PersistentFlags().StringVar(&globalOpts.srcDir, "src-dir", currentDir, "source directory which you create ghost from")
	ghostPrefixEnv := os.Getenv("GHOST_PREFIX")
	if ghostPrefixEnv == "" {
		ghostPrefixEnv = "ghost"
	}
	RootCmd.PersistentFlags().StringVar(&globalOpts.ghostPrefix, "ghost-prefix", ghostPrefixEnv, "prefix of ghost branch name")
	ghostRepoEnv := os.Getenv("GHOST_REPO")
	RootCmd.PersistentFlags().StringVar(&globalOpts.ghostRepo, "ghost-repo", ghostRepoEnv, "git refspec for ghost commits repository")
	RootCmd.PersistentFlags().StringVar(&globalOpts.baseCommit, "base-commit", "HEAD", "base commit hash for generating ghost commit.")
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of git-ghost",
	Long:  `Print the version number of git-ghost`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("git-ghost %s (revision: %s)", Version, Revision)
	},
}

func validateEnvironment() error {
	err := git.ValidateGit()
	if err != nil {
		return errors.New("git is required")
	}
	return nil
}

func (flags *globalFlags) Validate() error {
	if flags.srcDir == "" {
		return errors.New("src-dir must be specified")
	}
	if flags.ghostPrefix == "" {
		return errors.New("ghost-prefix must be specified")
	}
	if flags.ghostRepo == "" {
		return errors.New("ghost-repo must be specified")
	}
	if flags.baseCommit != "" {
		err := git.ValidateRefspec(".", flags.baseCommit)
		if err != nil {
			return errors.New("base-commit is not a valid object")
		}
	}
	return nil
}
