package cmd

import (
	"errors"
	"fmt"
	"git-ghost/pkg/ghost"
	"git-ghost/pkg/ghost/git"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

type pushFlags struct {
	localBase string
}

func init() {
	RootCmd.AddCommand(NewPushCommand())
}

func NewPushCommand() *cobra.Command {
	var (
		pushOpts pushFlags
	)
	command := &cobra.Command{
		Use:   "push",
		Short: "generate and push a ghost commit to remote repository",
		Long:  "generate and push a ghost commit to remote repository",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := ghost.Push(ghost.PushOptions{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostPrefix:     globalOpts.ghostPrefix,
				GhostRepo:       globalOpts.ghostRepo,
				RemoteBase:      globalOpts.baseCommit,
				LocalBase:       pushOpts.localBase,
			})
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			if resp.LocalModBranch != nil {
				fmt.Println(resp.LocalModBranch.LocalModHash)
			}
		},
	}
	command.PersistentFlags().StringVar(&pushOpts.localBase, "local-base", "HEAD", "git refspec used to create a local modification patch from")
	return command
}

func (flags pushFlags) Validate() error {
	if flags.localBase == "" {
		return errors.New("local-base must be specified")
	}
	err := git.ValidateRefspec(".", flags.localBase)
	if err != nil {
		return errors.New("local-base is not a valid object")
	}
	return nil
}
