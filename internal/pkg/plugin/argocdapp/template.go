package argocdapp

var argoCDAppTemplate = `---
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{.App.Name}}
  namespace: {{.App.Namespace}}
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  destination:
    namespace: {{.Destination.Namespace}}
    server: {{.Destination.Server}}
  project: default
  source:
    helm:
      valueFiles:
      - {{.Source.Valuefile}}
    path: {{.Source.Path}}
    repoURL: {{.Source.RepoURL}}
    targetRevision: HEAD
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
`
