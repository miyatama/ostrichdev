package ostrich

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type CommandExecutorInterface interface {
	ExecCommand(command string, args []string) ([]string, error)
}

type CommandExecutor struct {
}

func (c *CommandExecutor) ExecCommand(command string, args []string) ([]string, error) {
	c.outputDebug(fmt.Sprintf("ExecCommand(): command: %s, args: %s", command, strings.Join(args, " ")))
	cmd := exec.Command(
		command,
		args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		c.outputDebug("ExecCommand catch error -----")
		c.outputDebug(fmt.Sprintf("stdout in error: %s", out))
		c.outputDebug(fmt.Sprintf("error description: %s", err.Error()))
		c.outputDebug(fmt.Sprintf("error: %#v", err))
		c.outputDebug("-----------------------------")
		return []string{}, 
			fmt.Errorf(
				"error: %s.command: %s, args: %s", 
				err.Error(), 
				command, 
				strings.Join(args, " "))
	}
	result := strings.Split(string(out), "\n")
	return result, nil
}

func (c *CommandExecutor) outputDebug(message string) {
	log.Printf("[DEBUG]: %s", message)
}
