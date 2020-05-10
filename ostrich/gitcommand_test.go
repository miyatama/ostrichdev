package ostrich

import (
	"errors"
	"testing"
)

type DummyExecutor struct {
	Command     string
	Args        []string
	ReturnError bool
	Result      []string
}

func (d *DummyExecutor) ExecCommand(command string, args []string) ([]string, error) {
	d.Command = command
	d.Args = args

	if d.ReturnError {
		return []string{}, errors.New("raise error")
	}
	return d.Result, nil
}

func TestGitClone(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter", func(t *testing.T) {
		executor.ReturnError = false
		repository := "http://github.com/x/y.git"
		err := git.Clone(repository)
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"clone",
			repository,
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		repository := "http://github.com/x/y.git"
		err := git.Clone(repository)
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitCheckout(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter", func(t *testing.T) {
		executor.ReturnError = false
		branch := "develop"
		err := git.Checkout(branch)
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"checkout",
			"-b",
			branch,
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		branch := "develop"
		err := git.Checkout(branch)
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitPull(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter", func(t *testing.T) {
		executor.ReturnError = false
		remoteBranch := "develop"
		err := git.Pull(remoteBranch)
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"pull",
			"origin",
			remoteBranch,
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		remoteBranch := "develop"
		err := git.Pull(remoteBranch)
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitBranch(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter and result", func(t *testing.T) {
		executor.ReturnError = false
		executor.Result = []string{
			"row01",
			"row02",
			"row03",
		}
		results, err := git.Branch()
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"branch",
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
		expectResults := []string{
			"row01",
			"row02",
			"row03",
		}

		for i, expectResult := range expectResults {
			if expectResult != results[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					expectResult,
					results[i])
			}

		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		_, err := git.Branch()
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitShow(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter and result", func(t *testing.T) {
		executor.ReturnError = false
		executor.Result = []string{
			"row01",
			"row02",
			"row03",
		}
		commitId := "ABCDEFG"
		results, err := git.Show(commitId)
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"show",
			commitId,
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
		expectResults := []string{
			"row01",
			"row02",
			"row03",
		}

		for i, expectResult := range expectResults {
			if expectResult != results[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					expectResult,
					results[i])
			}

		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		commitId := "ABCDEFG"
		_, err := git.Show(commitId)
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitCommit(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter and result", func(t *testing.T) {
		executor.ReturnError = false
		message := "ABCDEFG"
		err := git.Commit(message)
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"commit",
			"-m",
			message,
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		message := "ABCDEFG"
		err := git.Commit(message)
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitPush(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter and result", func(t *testing.T) {
		executor.ReturnError = false
		branch := "develop"
		err := git.Push(branch)
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"push",
			"-f",
			"origin",
			branch,
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		branch := "develop"
		err := git.Push(branch)
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitAdd(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter and result", func(t *testing.T) {
		executor.ReturnError = false
		filename := "fileA"
		err := git.Add(filename)
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"add",
			filename,
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		filename := "fileA"
		err := git.Add(filename)
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitRm(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter and result", func(t *testing.T) {
		executor.ReturnError = false
		filename := "fileA"
		err := git.Rm(filename)
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"rm",
			filename,
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		filename := "fileA"
		err := git.Rm(filename)
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitReset(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter and result", func(t *testing.T) {
		executor.ReturnError = false
		branch := "develop"
		err := git.Reset(branch)
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"reset",
			"--hard",
			"origin/develop",
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		branch := "develop"
		err := git.Reset(branch)
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}

func TestGitFetch(t *testing.T) {
	executor := &DummyExecutor{}
	git := GitCommand{
		executor: executor,
	}

	t.Run("execute command parameter and result", func(t *testing.T) {
		executor.ReturnError = false
		err := git.Fetch()
		if err != nil {
			t.Fatal("invalid return.")
		}
		expectCommand := "git"
		expectArgs := []string{
			"fetch",
		}
		if expectCommand != executor.Command {
			t.Fatalf(
				"invalid command.expect: %s, result: %s",
				expectCommand,
				executor.Command)
		}

		for i, arg := range expectArgs {
			if arg != executor.Args[i] {
				t.Fatalf(
					"invalid args %d.expect: %s, result: %s",
					i,
					arg,
					executor.Args[i])
			}
		}
	})
	t.Run("return error", func(t *testing.T) {
		executor.ReturnError = true
		err := git.Fetch()
		if err == nil {
			t.Fatal("invalid return.")
		}
	})
}
