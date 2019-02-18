package ghost

import (
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util"
	"git-ghost/pkg/util/errors"

	log "github.com/Sirupsen/logrus"
)

// PullOptions represents arg for Pull func
type PullOptions struct {
	types.WorkingEnvSpec
	*types.CommitsBranchSpec
	*types.PullableDiffBranchSpec
}

func pullAndApply(spec types.PullableGhostBranchSpec, we types.WorkingEnv) errors.GitGhostError {
	pulledBranch, err := spec.PullBranch(we)
	if err != nil {
		return errors.WithStack(err)
	}
	return pulledBranch.Apply(we)
}

// Pull pulls ghost branches and apply to workind directory
func Pull(options PullOptions) errors.GitGhostError {
	log.WithFields(util.ToFields(options)).Debug("pull command with")
	we, err := options.WorkingEnvSpec.Initialize()
	if err != nil {
		return errors.WithStack(err)
	}
	defer util.LogDeferredGitGhostError(we.Clean)

	if options.CommitsBranchSpec != nil {
		err := pullAndApply(*options.CommitsBranchSpec, *we)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if options.PullableDiffBranchSpec != nil {
		err := pullAndApply(*options.PullableDiffBranchSpec, *we)
		return errors.WithStack(err)
	}

	log.WithFields(util.ToFields(options)).Warn("pull command has nothing to do with")
	return nil
}
