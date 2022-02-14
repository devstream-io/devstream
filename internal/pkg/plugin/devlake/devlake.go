package devlake

const devLakeInstallYAMLDownloadURL = "https://raw.githubusercontent.com/merico-dev/lake/main/k8s-deploy.yaml"
const devLakeInstallYAMLFileName = "devlake-k8s-deploy.yaml"

func validateParams(param *Param) []error {
	// at the moment, devlake plugin doesn't have any parameters yet
	return make([]error, 0)
}
