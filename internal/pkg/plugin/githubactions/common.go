package githubactions

import "fmt"

func GetLanguage(l *Language) string {
	return fmt.Sprintf("%s-%s", l.Name, l.Version)
}
