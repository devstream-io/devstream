package devlake

const (
	devLakeInstallYAMLDownloadURL = "https://raw.githubusercontent.com/merico-dev/lake/main/k8s-deploy.yaml"
	/// TODO(ironcore864): now the namespace is hard-coded instead of parsed from the YAML file
	defaultNamespace = "devlake"
)

// according to devLakeInstallYAMLDownloadURL
// a successful DevLake installation should have the following deployments
// (and corresponding services as well)
var devLakeDeployments = []string{
	"mysql",
	"grafana",
	"config-ui",
	"devlake",
}
