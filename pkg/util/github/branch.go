package github

import (
	"fmt"

	"github.com/google/go-github/v42/github"

	"github.com/merico-dev/stream/internal/pkg/log"
)

func (c *Client) NewBranch(baseBranch, newBranch string) error {
	refStr := fmt.Sprintf("heads/%s", baseBranch)
	ref, _, err := c.Git.GetRef(c.Context, c.Owner, c.Repo, refStr)
	if err != nil {
		log.Debugf("Failed to get the ref for %s: %s.", refStr, err)
		return err
	}
	log.Debugf("Got the ref: Ref %s, URL %s, nodeId %s, Obj: %s.",
		ref.GetRef(), ref.GetURL(), ref.GetNodeID(), ref.GetObject().String())

	newRef := fmt.Sprintf("heads/%s", newBranch)
	_, _, err = c.Git.CreateRef(c.Context, c.Owner, c.Repo, &github.Reference{
		Ref: &newRef,
		Object: &github.GitObject{
			Type: nil,
			SHA:  ref.GetObject().SHA,
			URL:  nil,
		},
	})
	return err
}
