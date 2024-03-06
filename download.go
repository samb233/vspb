package vspb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetPackage(dir string, pkg *Package) error {
	repoDir := dir + "/repos"
	path := repoDir + "/" + pkg.Name

	if _, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("mkdir '%s' error: %w", path, err)
		}
	}

	if len(pkg.Repo) == 0 {
		return nil
	}

	if !isGit(pkg.Repo) {
		return grabDownload(pkg.Repo, path)
	}

	// TODO:
	// if already exists, check the version of repo
	return gitClone(pkg.Repo, pkg.Version, path, pkg.VersionIsBranch)
}

func grabDownload(addr, repoDir string) error {
	_, err := grab.Get(repoDir, addr)
	return err
}

func gitClone(addr, version, path string, isBranch bool) error {
	var needCheckout bool = false

	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	cloneOpts := &git.CloneOptions{
		URL:               addr,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if isGitHash(version) {
		needCheckout = true
	} else {
		cloneOpts.Depth = 1
		if len(version) > 0 {
			if isBranch {
				cloneOpts.ReferenceName = plumbing.NewBranchReferenceName(version)
			} else {
				cloneOpts.ReferenceName = plumbing.NewTagReferenceName(version)
			}
		}
	}

	r, err := git.PlainClone(abs, false, cloneOpts)
	if err != nil || !needCheckout {
		return err
	}

	tree, err := r.Worktree()
	if err != nil {
		return err
	}

	return tree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(version),
	})
}

func gitPull(addr, version, path string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	pullOpts := &git.PullOptions{
		RemoteURL:         addr,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	tree, err := r.Worktree()
	if err != nil {
		return err
	}

	err = tree.Pull(pullOpts)
	if err != nil || len(version) == 0 {
		return err
	}

	checkoutOpts := &git.CheckoutOptions{}
	if isGitHash(version) {
		checkoutOpts.Hash = plumbing.NewHash(version)
	} else {
		checkoutOpts.Branch = plumbing.NewTagReferenceName(version)
	}

	return tree.Checkout(checkoutOpts)
}

// TODO:
// does there has better way to check these two?

func isGit(addr string) bool {
	return strings.HasSuffix(addr, ".git")
}

func isGitHash(version string) bool {
	return len(version) == 40
}
