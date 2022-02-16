package githubactions

import "fmt"

func GetLanguage(l *Language) string {
	return fmt.Sprintf("%s-%s", l.Name, l.Version)
}

func BuildState(owner, repo string) map[string]interface{} {
	res := make(map[string]interface{})
	res["workflowDir"] = fmt.Sprintf("repos/%s/%s/contents/.github/workflows", owner, repo)
	return res
}
