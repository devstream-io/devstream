package golang

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	rs "github.com/devstream-io/devstream/internal/pkg/plugin/common/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	DefaultWorkPath   = ".github-repo-scaffolding-golang"
	TransitBranch     = "init-with-devstream"
	DefaultMainBranch = "main"
)

type Config struct {
	AppName   string
	ImageRepo string
	Repo      Repo
}

type Repo struct {
	Name  string
	Owner string
}

func pushToRemote(repoPath string, opts *rs.Options) error {
	ghOptions := &github.Option{
		Owner:    opts.Owner,
		Org:      opts.Org,
		Repo:     opts.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return err
	}

	var retErr error
	// It's ok to give the opts.Org to CreateRepo() when create a repository for a authenticated user.
	if err := ghClient.CreateRepo(opts.Org); err != nil {
		log.Errorf("Failed to create repo: %s.", err)
		return err
	}
	log.Infof("The repo %s has been created.", opts.Repo)

	defer func() {
		if retErr == nil {
			return
		}
		// need to clean the repo created when retErr != nil
		if err := ghClient.DeleteRepo(); err != nil {
			log.Errorf("Failed to delete the repo %s: %s.", opts.Repo, err)
		}
	}()

	if retErr = walkLocalRepoPath(repoPath, opts, ghClient); retErr != nil {
		log.Debugf("Failed to walk local repo-path: %s.", retErr)
		return retErr
	}

	mainBranch := getMainBranchName(opts)
	if retErr = mergeCommits(ghClient, mainBranch); retErr != nil {
		log.Debugf("Failed to merge commits: %s.", retErr)
		return retErr
	}

	return nil
}

func walkLocalRepoPath(repoPath string, opts *rs.Options, ghClient *github.Client) error {
	mainBranch := getMainBranchName(opts)

	if err := filepath.Walk(repoPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debugf("Walk error: %s.", err)
			return err
		}

		if info.IsDir() {
			log.Debugf("Found dir: %s.", path)
			return nil
		}

		log.Debugf("Found file: %s.", path)

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		githubPath := strings.Join(strings.Split(path, "/")[2:], "/")
		// the main branch needs a initial commit
		if strings.Contains(path, "gitignore") {
			err := ghClient.CreateFile(content, githubPath, mainBranch)
			if err != nil {
				log.Debugf("Failed to add the .gitignore file: %s.", err)
				return err
			}
			log.Debugf("Added the .gitignore file.")
			return ghClient.NewBranch(mainBranch, TransitBranch)
		}
		return ghClient.CreateFile(content, githubPath, TransitBranch)
	}); err != nil {
		return err
	}

	return nil
}

func mergeCommits(ghClient *github.Client, mainBranch string) error {
	number, err := ghClient.NewPullRequest(TransitBranch, mainBranch)
	if err != nil {
		return err
	}

	return ghClient.MergePullRequest(number, github.MergeMethodSquash)
}

func getMainBranchName(opts *rs.Options) string {
	if opts.Branch == "" {
		return DefaultMainBranch
	}
	return opts.Branch
}

func buildState(opts *rs.Options) map[string]interface{} {
	res := make(map[string]interface{})
	res["owner"] = opts.Owner
	res["org"] = opts.Org
	res["repoName"] = opts.Repo

	outputs := make(map[string]interface{})
	outputs["owner"] = opts.Owner
	outputs["org"] = opts.Org
	outputs["repo"] = opts.Repo
	if opts.Owner != "" {
		outputs["repoURL"] = fmt.Sprintf("https://github.com/%s/%s.git", opts.Owner, opts.Repo)
	} else {
		outputs["repoURL"] = fmt.Sprintf("https://github.com/%s/%s.git", opts.Org, opts.Repo)
	}
	res["outputs"] = outputs

	return res
}
