package github

import (
	"fmt"
	"net/http"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
)

// Workflow is the struct for a GitHub Actions Workflow.
type Workflow struct {
	CommitMessage    string
	WorkflowFileName string
	WorkflowContent  string
}

func (c *Client) AddWorkflow(workflow *Workflow, branch string) error {
	sha, err := c.getFileSHA(workflow.WorkflowFileName)
	if err != nil {
		return err
	}
	if sha != "" {
		log.Infof("GitHub Actions workflow %s already exists.", workflow.WorkflowFileName)
		return nil
	}

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(workflow.CommitMessage),
		Content: []byte(workflow.WorkflowContent),
		Branch:  github.String(branch),
	}

	log.Infof("Creating GitHub Actions workflow %s ...", workflow.WorkflowFileName)
	_, _, err = c.Client.Repositories.CreateFile(
		c.Context,
		c.GetRepoOwner(),
		c.Repo,
		generateGitHubWorkflowFileByName(workflow.WorkflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Successf("Github Actions workflow %s created.", workflow.WorkflowFileName)
	return nil
}
func (c *Client) DeleteWorkflow(workflow *Workflow, branch string) error {
	sha, err := c.getFileSHA(workflow.WorkflowFileName)
	if err != nil {
		return err
	}
	if sha == "" {
		log.Successf("Github Actions workflow %s already removed.", workflow.WorkflowFileName)
		return nil
	}

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(workflow.CommitMessage),
		SHA:     github.String(sha),
		Branch:  github.String(branch),
	}

	log.Infof("Deleting GitHub Actions workflow %s ...", workflow.WorkflowFileName)
	_, _, err = c.Client.Repositories.DeleteFile(
		c.Context,
		c.GetRepoOwner(),
		c.Repo,
		generateGitHubWorkflowFileByName(workflow.WorkflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Successf("GitHub Actions workflow %s removed.", workflow.WorkflowFileName)
	return nil
}

// VerifyWorkflows get the workflows with names "wf1.yml", "wf2.yml", then:
// If all workflows is ok => return ({"wf1.yml":nil, "wf2.yml:nil}, nil)
// If some error occurred => return (nil, error)
// If wf1.yml is not found => return ({"wf1.yml":error("not found"), "wf2.yml:nil},nil)
func (c *Client) VerifyWorkflows(workflows []*Workflow) (map[string]error, error) {
	wsFiles := make([]string, 0, len(workflows))
	for _, w := range workflows {
		wsFiles = append(wsFiles, w.WorkflowFileName)
	}
	fmt.Printf("Workflow files: %v", wsFiles)

	_, dirContent, resp, err := c.Client.Repositories.GetContents(
		c.Context,
		c.GetRepoOwner(),
		c.Repo,
		".github/workflows",
		&github.RepositoryContentGetOptions{},
	)

	// StatusCode == 404
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		log.Errorf("GetContents return with status code 404.")
		retMap := mapz.FillMapWithStrAndError(wsFiles, fmt.Errorf("not found"))
		return retMap, nil
	}

	// error reason is not 404
	if err != nil {
		log.Errorf("GetContents failed with error: %s.", err)
		return nil, err
	}

	// StatusCode == 200
	log.Success("GetContents return with status code 200.")
	var filesInRemoteDir = make([]string, 0)
	for _, f := range dirContent {
		log.Infof("Found remote file: %s.", f.GetName())
		filesInRemoteDir = append(filesInRemoteDir, f.GetName())
	}

	lostFiles := mapset.NewSet(wsFiles...).Difference(mapset.NewSet(filesInRemoteDir...)).ToSlice()
	// all files exist
	if len(lostFiles) == 0 {
		log.Info("All files exist.")
		retMap := mapz.FillMapWithStrAndError(wsFiles, nil)
		return retMap, nil
	}
	// some files lost
	log.Warn("Some files lost.")
	retMap := mapz.FillMapWithStrAndError(wsFiles, nil)
	for _, f := range lostFiles {
		log.Infof("Lost file: %s.", f)
		retMap[f] = fmt.Errorf("not found")
	}
	return retMap, nil
}

func (c *Client) GetWorkflowPath() (string, error) {
	_, _, resp, err := c.Client.Repositories.GetContents(
		c.Context,
		c.GetRepoOwner(),
		c.Repo,
		".github/workflows",
		&github.RepositoryContentGetOptions{},
	)

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return resp.Request.URL.Path, nil
}

func (c *Client) FetchRemoteContent(wsFiles []string) ([]string, map[string]error, error) {
	var filesInRemoteDir = make([]string, 0)
	_, dirContent, resp, err := c.Repositories.GetContents(
		c.Context,
		c.GetRepoOwner(),
		c.Repo,
		".github/workflows",
		&github.RepositoryContentGetOptions{},
	)

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		log.Error("GetContents returned with status code 404.")
		retMap := mapz.FillMapWithStrAndError(wsFiles, fmt.Errorf("not found"))
		return nil, retMap, nil
	}

	// error reason is not 404
	if err != nil {
		log.Errorf("GetContents failed with error: %s.", err)
		return nil, nil, err
	}

	log.Info("GetContents successfully.")
	for _, f := range dirContent {
		log.Infof("Found remote file: %s.", f.GetName())
		filesInRemoteDir = append(filesInRemoteDir, f.GetName())
	}
	return filesInRemoteDir, nil, nil
}
