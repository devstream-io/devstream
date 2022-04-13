package githubactions

import (
	"fmt"
)

func GetLanguage(l *Language) string {
	return fmt.Sprintf("%s-%s", l.Name, l.Version)
}

func BuildState(owner, org, repo string) map[string]interface{} {
	res := make(map[string]interface{})
	if owner != "" {
		res["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", owner, repo)
	} else {
		res["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", org, repo)
	}
	return res
}

func BuildReadState(path string) map[string]interface{} {
	res := make(map[string]interface{})
	res["workflowDir"] = path
	return res
}
