package internal

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	defaultShellCommand = "bash"
	// ShellCommand is the command which is called in the ShellContainer.
	ShellCommand []string
	// ShellContainer is the container used for shell command.
	ShellContainer string
)

// ShellCmd opens a shell in the environment's default application container.
func ShellCmd(cmd *cobra.Command, args []string) error {
	if CheckRegexInString("^pwa-studio", GetEnvType()) {
		SetShellContainer("node")
		SetDefaultShellCommand("sh")
	}

	if len(args) > 0 {
		ShellCommand = ExtractUnknownArgs(cmd.Flags(), args)
	} else {
		ShellCommand = ExtractUnknownArgs(cmd.Flags(), []string{defaultShellCommand})
	}

	log.Debugln("command:", ShellCommand)
	log.Debugln("container:", ShellContainer)

	var passedArgs []string
	passedArgs = append(passedArgs, "exec", ShellContainer)
	passedArgs = append(passedArgs, ShellCommand...)

	err := EnvRunDockerCompose(passedArgs, false)
	if err != nil {
		return err
	}

	return nil
}

// SetShellContainer changes the container used for the reward shell command.
func SetShellContainer(s string) {
	ShellContainer = s
}

// SetDefaultShellCommand changes the command invoked by reward shell command.
func SetDefaultShellCommand(s string) {
	defaultShellCommand = s
}
