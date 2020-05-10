package ostrich

import (
	"fmt"
)
type GitCommand struct {
	executor CommandExecutorInterface
}

func (g *GitCommand) Clone(repository string) error {
	_, err := g.executor.ExecCommand("git", []string{"clone", repository})
	return err
}

func (g *GitCommand) Checkout(branch string) error {
	_, err := g.executor.ExecCommand("git", []string{"checkout", "-b", branch})
	return err
}

func (g *GitCommand) Pull(branch string) error {
	_, err := g.executor.ExecCommand("git", []string{"pull", "origin", branch})
	return err
}

func (g *GitCommand) Branch() ([]string, error) {
	return g.executor.ExecCommand("git", []string{"branch"})
}

func (g *GitCommand) Show(commitId string) ([]string, error) {
	return g.executor.ExecCommand("git", []string{"show", commitId})
}

func (g *GitCommand) Commit(message string) error {
	_, err := g.executor.ExecCommand("git", []string{"commit", "-m", message})
	return err
}
func (g *GitCommand) Push(branch string) error {
	_, err := g.executor.ExecCommand("git", []string{"push", "-f", "origin", branch})
	return err
}
func (g *GitCommand) Version() ([]string, error) {
	return g.executor.ExecCommand("git", []string{"--version"})
}

func (g *GitCommand) Add(filepath string) error {
	_, err := g.executor.ExecCommand("git", []string{"add", filepath})
	return err
}

func (g *GitCommand) Rm(filepath string) error {
	_, err := g.executor.ExecCommand("git", []string{"rm", filepath})
	return err
}

func (g *GitCommand) Reset(branch string) error {
	_, err := g.executor.ExecCommand("git", []string{"reset", "--hard", fmt.Sprintf("origin/%s",branch)})
	return err
}

func (g *GitCommand) Fetch() error {
	_, err := g.executor.ExecCommand("git", []string{"fetch"})
	return err
}
