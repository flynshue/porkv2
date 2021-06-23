package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type GHRepo struct {
	Dir     string
	owner   string
	project string
	repo    *git.Repository
}

func NewGHRepo(repo string) (*GHRepo, error) {
	values := strings.Split(repo, "/")
	if len(values) != 2 {
		return nil, fmt.Errorf("must supply repository in owner/project format")
	}
	return &GHRepo{owner: values[0], project: values[1]}, nil
}

func (g *GHRepo) RepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s/%s.git", g.owner, g.project)
}

func (g *GHRepo) Clone() error {
	fullPath := filepath.Join(viper.GetString("location"), g.owner, g.project)
	opts := &git.CloneOptions{
		URL:  g.RepositoryURL(),
		Auth: &http.BasicAuth{Username: "anyuser", Password: viper.GetString("token")},
	}
	repo, err := git.PlainClone(fullPath, false, opts)
	if err != nil {
		return err
	}
	g.repo = repo
	g.Dir = fullPath
	return nil
}

func (g *GHRepo) Checkout(branch string, create bool) error {
	opts := &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Create: create,
	}
	if create {
		head, err := g.repo.Head()
		if err != nil {
			return err
		}
		opts.Hash = head.Hash()
	}
	tree, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	return tree.Checkout(opts)
}
