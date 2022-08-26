package gitea

import (
	"log"
	"code.gitea.io/sdk/gitea"
)


func (c *Client) InitRepo() error {
	opt:=gitea.CreateRepoOption{
		Name: c.Repo,
		Description:c.Description,
		IssueLabels: c.Labels,
		DefaultBranch: c.Branch,
	}
	switch c.Visibility {
	case "public":
		opt.Private=false
	case "private":
		opt.Private=true
	default:
		opt.Private=false
	}
	_,_,err:=c.CreateRepo(opt)
	if err!=nil{
		log.Errorf("Failed to create repo: %s.", err)
		
		
		return err
	}
	log.Infof("The repo %s has been created.", c.Repo)
	return nil
}
