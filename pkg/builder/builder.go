package builder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	client "github.com/mattmoor/bindings/pkg/github"
	"github.com/mattmoor/knobots/pkg/comment"
)

const (
	statusNotFound = "Not Found"
)

// TODO: support signoff.

type Builder struct {
	DryRun          bool
	Verbose         bool
	SkipLockChanges bool

	// Git options.
	Workspace    string
	Owner        string
	Repo         string
	Branch       string
	CommitBranch *string

	// PR options.
	Title     string
	Body      string
	Token     *string
	Signature *object.Signature
	Signoff   bool

	// Binding
	Username string
	Password string
}

// Commit will commit changes to the local branch, and then optionally branch a new branch if set.
func (b *Builder) Commit(r *git.Repository) (plumbing.ReferenceName, error) {
	// First, build the worktree.
	wt, err := r.Worktree()
	if err != nil {
		return "", fmt.Errorf("Error fetching worktree: %v", err)
	}

	// Check the status of the worktree, and if there aren't any changes
	// bail out we're done.
	st, err := wt.Status()
	if err != nil {
		return "", fmt.Errorf("Error fetching worktree status: %v", err)
	}
	changes := true
	if len(st) == 0 {
		log.Println("No changes")
		changes = false
	}

	nonGopkgCount := 0
	for p := range st {
		if path.Base(p) != "Gopkg.lock" {
			nonGopkgCount++
		}
		_, err = wt.Add(p)
		if err != nil {
			return "", fmt.Errorf("Error staging %q: %v", p, err)
		}
		// Display any changed we do find: `git status --porcelain`
		log.Printf("Added - %v", p)
	}
	// If we want to ignore lock changes, just return.
	if b.SkipLockChanges && nonGopkgCount == 0 {
		log.Println("Only Gopkg.lock files changed (skipping PR).")
		changes = false
	}

	if changes {
		commitMessage := b.Title + "\n\n" + b.Body
		// Commit the staged changes to the repo.
		if _, err := wt.Commit(commitMessage, &git.CommitOptions{Author: b.Signature}); err != nil {
			return "", fmt.Errorf("Error committing changes: %v", err)
		}
	}

	var rn plumbing.ReferenceName

	// Commit branch is optional.
	if b.CommitBranch != nil {
		headRef, err := r.Head()
		if err != nil {
			return "", fmt.Errorf("Error fetching workspace HEAD: %v", err)
		}
		// Create and checkout a new branch from the commit of the HEAD reference.
		// This should be roughly equivalent to `git checkout -b {new-branch}`
		newBranchName := plumbing.NewBranchReferenceName(*b.CommitBranch)
		log.Println("new branch", newBranchName)
		if err := wt.Checkout(&git.CheckoutOptions{
			Hash:   headRef.Hash(),
			Branch: newBranchName,
			Create: true,
			Force:  true,
		}); err != nil {
			return "", fmt.Errorf("Error checkout out new branch: %v", err)
		}
		rn = newBranchName
	} else {
		ref, err := r.Head()
		if err != nil {
			return "", err
		}
		rn = ref.Name()
	}

	return rn, nil
}

func (b *Builder) Push(r *git.Repository, rs config.RefSpec) error {
	// Push the branch to a remote to which we have write access.
	// TODO(mattmoor): What if the fork doesn't exist, or has another name?
	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name: b.Username,
		URLs: []string{fmt.Sprintf("https://github.com/%s/%s.git", b.Username, b.Repo)},
	})
	if err != nil {
		if err.Error() == "remote already exists" {
			// Skip.
		} else {
			return fmt.Errorf("Error creating new remote: %v", err)
		}
	}
	remote, err = r.Remote(b.Username)
	if err != nil {
		return fmt.Errorf("Error creating new remote: %v", err)
	}

	// Publish all local branches to the remote.
	err = remote.Push(&git.PushOptions{
		RemoteName: b.Username,
		RefSpecs:   []config.RefSpec{rs},
		Auth: &http.BasicAuth{
			Username: b.Username, // This can be anything.
			Password: b.Password,
		},
	})
	if err != nil {
		return fmt.Errorf("Error pushing to remote: %v", err)
	}

	return nil
}

func (b *Builder) PR(branch string) error {
	// Head has the form source-owner:branch, per the Github API docs.
	head := fmt.Sprintf("%s:%s", b.Username, branch)

	// Inject a signature into the body that will help us clean up matching older PRs.
	bodyWithSignature := comment.WithSignature(b.Title, b.Body)
	// Inject the token (if specified) into the body of the PR, so
	// that we can identify it's provenance.
	if b.Token != nil {
		bodyWithSignature = comment.WithSignature(*b.Token, *bodyWithSignature)
	}

	ctx := context.Background()
	ghc, err := client.New(ctx)
	if err != nil {
		return fmt.Errorf("Error creating github client: %v", err)
	}

	pr, _, err := ghc.PullRequests.Create(ctx, b.Owner, b.Repo, &github.NewPullRequest{
		Title: &b.Title,
		Body:  bodyWithSignature,
		Head:  &head,
		Base:  &b.Branch,
	})
	if err != nil {
		return fmt.Errorf("Error creating PR: %v", err)
	}

	log.Printf("Created PR: #%d", pr.GetNumber())
	return nil
}

func (b *Builder) Do() error {

	log.Printf("Debug: %#v\n\n", b)

	// Clean up older PRs as the first thing we do so that if the latest batch of
	// changes needs nothing we don't leave old PRs around.
	if err := b.cleanupOlderPRs(); err != nil {
		return fmt.Errorf("Error cleaning up PRs: %v", err)
	}

	r, err := git.PlainOpen(b.Workspace)
	if err != nil {
		return fmt.Errorf("Error opening /workspace: %v", err)
	}

	rn, err := b.Commit(r)
	if err != nil {
		return err
	}

	// Always match local name to remote name.
	rs := config.RefSpec(fmt.Sprintf("%s:%s", rn, rn))

	log.Printf("Pushing to refspec: %s", rs.String())

	if err := b.Push(r, rs); err != nil {
		return err
	}

	if err := b.PR(rn.String()); err != nil {
		return err
	}
	return nil
}

func (b *Builder) cleanupOlderPRs() error {
	ctx := context.Background()
	ghc, err := client.New(ctx)
	if err != nil {
		return err
	}

	closed := "closed"
	lopt := &github.PullRequestListOptions{
		Base: b.Branch,
	}
	for {
		prs, resp, err := ghc.PullRequests.List(ctx, b.Owner, b.Repo, lopt)
		if err != nil {
			ghe := &github.ErrorResponse{}
			if errors.As(err, &ghe) && ghe.Message == statusNotFound {
				return nil
			}
			return err
		}
		for _, pr := range prs {
			if comment.HasSignature(b.Signature.Name, pr.GetBody()) {
				_, _, err := ghc.PullRequests.Edit(ctx, b.Owner, b.Repo, pr.GetNumber(), &github.PullRequest{
					State: &closed,
				})
				if err != nil {
					return err
				}
			}
		}
		if resp.NextPage == 0 {
			break
		}
		lopt.Page = resp.NextPage
	}

	return nil
}
