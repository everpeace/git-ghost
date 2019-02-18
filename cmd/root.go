package cmd

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util/errors"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

type globalFlags struct {
	srcDir       string
	ghostWorkDir string
	ghostPrefix  string
	ghostRepo    string
	verbose      bool
}

func (gf globalFlags) WorkingEnvSpec() types.WorkingEnvSpec {
	return types.WorkingEnvSpec{
		SrcDir:          gf.srcDir,
		GhostWorkingDir: gf.ghostWorkDir,
		GhostRepo:       gf.ghostRepo,
	}
}

var (
	Version  string
	Revision string
)

var RootCmd = &cobra.Command{
	Use:           "git-ghost",
	Short:         "git-ghost",
	SilenceErrors: false,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Use == "version" {
			return nil
		}
		err := validateEnvironment()
		if err != nil {
			return err
		}
		err = globalOpts.SetDefaults()
		if err != nil {
			return err
		}
		err = globalOpts.Validate()
		if err != nil {
			return err
		}
		if globalOpts.verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.ErrorLevel)
		}
		return nil
	},
}

var globalOpts globalFlags

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().StringVar(&globalOpts.srcDir, "src-dir", "", "source directory which you create ghost from (default to PWD env)")
	RootCmd.PersistentFlags().StringVar(&globalOpts.ghostWorkDir, "ghost-working-dir", "", "local root directory for git-ghost interacting with ghost repository (default to a temporary directory)")
	RootCmd.PersistentFlags().StringVar(&globalOpts.ghostPrefix, "ghost-prefix", "", "prefix of ghost branch name (default to GIT_GHOST_PREFIX env, or ghost)")
	RootCmd.PersistentFlags().StringVar(&globalOpts.ghostRepo, "ghost-repo", "", "git remote url for ghosts repository (default to GIT_GHOST_REPO env)")
	RootCmd.PersistentFlags().BoolVarP(&globalOpts.verbose, "verbose", "v", false, "verbose mode")
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

func validateEnvironment() errors.GitGhostError {
	err := git.ValidateGit()
	if err != nil {
		return errors.New("git is required")
	}
	return nil
}

func (flags *globalFlags) SetDefaults() errors.GitGhostError {
	if globalOpts.srcDir == "" {
		globalOpts.srcDir = os.Getenv("PWD")
	}
	if globalOpts.ghostWorkDir == "" {
		globalOpts.ghostWorkDir = os.TempDir()
	}
	if globalOpts.ghostPrefix == "" {
		ghostPrefixEnv := os.Getenv("GIT_GHOST_PREFIX")
		if ghostPrefixEnv == "" {
			ghostPrefixEnv = "ghost"
		}
		globalOpts.ghostPrefix = ghostPrefixEnv
	}
	if globalOpts.ghostRepo == "" {
		globalOpts.ghostRepo = os.Getenv("GIT_GHOST_REPO")
	}
	return nil
}

func (flags *globalFlags) Validate() errors.GitGhostError {
	if flags.srcDir == "" {
		return errors.New("src-dir must be specified")
	}
	_, err := os.Stat(flags.ghostWorkDir)
	if err != nil {
		return errors.Errorf("ghost-working-dir is not found (value: %v)", flags.ghostWorkDir)
	}
	if flags.ghostPrefix == "" {
		return errors.New("ghost-prefix must be specified")
	}
	if flags.ghostRepo == "" {
		return errors.New("ghost-repo must be specified")
	}
	return nil
}
