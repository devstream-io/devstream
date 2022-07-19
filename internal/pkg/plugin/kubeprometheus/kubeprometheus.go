package kubeprometheus

var defaultDeploymentTplList = []string{
	// ${release-name}-grafana
	"%s-grafana",
	// ${release-name}-kube-prometheus-stack-operator
	"%s-kube-prometheus-stack-operator",
	// ${release-name}-kube-state-metrics
	"%s-kube-state-metrics",
}

var defaultDaemonsetTplList = []string{
	// ${release-name}-prometheus-node-exporter
	"%s-prometheus-node-exporter",
}

var defaultStatefulsetTplList = []string{
	// alertmanager-${release-name}-kube-prometheus-stack-alertmanager
	"alertmanager-%s-kube-prometheus-stack-alertmanager",
	// prometheus-${release-name}-kube-prometheus-stack-prometheus
	"prometheus-%s-kube-prometheus-stack-prometheus",
}
