package vcstool

import (
	"fmt"
	"strings"

	"github.com/cbegin/graven/internal/domain"
	"github.com/cbegin/graven/internal/util"
)

type GitVCSTool struct{}

type Validator func(stdout, stderr string) error

func (g *GitVCSTool) Tag(project *domain.Project, remote, tagName string) error {
	remoteName := "origin"

	if remote != "" {
		remoteName = remote
	}

	if sout, serr, err := util.RunCommand(project.ProjectPath(), nil, "git", "tag", tagName); err != nil {
		fmt.Printf("Tagging  %v\n%v\n", sout, serr)
		return err
	}

	if sout, serr, err := util.RunCommand(project.ProjectPath(), nil, "git", "push", "--tags", remoteName); err != nil {
		fmt.Printf("PushTags %v\n%v\n", sout, serr)
		return err
	}
	return nil
}

func (g *GitVCSTool) VerifyRepoState(project *domain.Project, remote, branch string) error {
	remoteName := "origin"
	branchName := "master"

	if remote != "" {
		remoteName = remote
	}

	if branch != "" {
		branchName = branch
	}

	fmt.Println("Check if on expected branch (e.g. master)")
	if err := verifyGitState(func(stdout, stderr string) error {
		actualBranch := strings.TrimSpace(stdout)
		if actualBranch != branchName {
			return fmt.Errorf("Expected to be on branch %v but found branch %v", branchName, actualBranch)
		}
		return nil
	}, project, "rev-parse", "--abbrev-ref", "HEAD"); err != nil {
		return err
	}

	fmt.Println("Check for uncommitted changes")
	if err := verifyGitState(func(stdout, stderr string) error {
		if strings.TrimSpace(stdout) != "" || strings.TrimSpace(stderr) != "" {
			return fmt.Errorf("Cannot release with uncommitted changes.")
		}
		return nil
	}, project, "status", "--porcelain"); err != nil {
		return err
	}

	fmt.Println("Check if changes exist on server")
	if err := verifyGitState(func(stdout, stderr string) error {
		lineCount := len(strings.Split(strings.TrimSpace(stderr), "\n"))
		if lineCount > 2 {
			return fmt.Errorf("Changes were detected on the remote %v for branch %v.", remoteName, branchName)
		}
		return nil
	}, project, "fetch", "--dry-run", remoteName, branchName); err != nil {
		return err
	}

	fmt.Println("Check if local changes are pushed")
	if err := verifyGitState(func(stdout, stderr string) error {
		parts := strings.Split(strings.TrimSpace(stdout), "\n")
		if strings.TrimSpace(parts[0]) != strings.TrimSpace(parts[1]) {
			return fmt.Errorf("Not all local changes for branch %v have been pushed to remote %v.", branchName, remoteName)
		}
		return nil
	}, project, "rev-parse", branchName, fmt.Sprintf("%v/%v", remoteName, branchName)); err != nil {
		return err
	}

	return nil
}

func verifyGitState(validator Validator, project *domain.Project, args ...string) error {
	sout, serr, err := util.RunCommand(project.ProjectPath(), nil, "git", args...)
	if err != nil {
		return fmt.Errorf("ERROR running Git command:\n%v\n%v\n", serr, err)
	}
	return validator(sout, serr)
}
