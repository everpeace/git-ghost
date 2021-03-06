// Copyright 2019 Preferred Networks, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/pfnet-research/git-ghost/pkg/util/errors"

	log "github.com/sirupsen/logrus"
)

func JustOutputCmd(cmd *exec.Cmd) ([]byte, errors.GitGhostError) {
	wd, _ := os.Getwd()
	log.WithFields(log.Fields{
		"pwd":     wd,
		"command": strings.Join(cmd.Args, " "),
	}).Debug("exec")
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	bytes, err := cmd.Output()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return []byte{}, errors.New(s)
		}
		return []byte{}, errors.WithStack(err)
	}
	return bytes, nil
}

func JustRunCmd(cmd *exec.Cmd) errors.GitGhostError {
	wd, _ := os.Getwd()
	log.WithFields(log.Fields{
		"pwd":     wd,
		"command": strings.Join(cmd.Args, " "),
	}).Debug("exec")
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return errors.WithStack(err)
	}
	return nil
}

func GetExitCode(err error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return -1
}
