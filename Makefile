build:
go get ./...
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/githubactions_0.0.1.so ./plugins/githubactions/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/argocd_0.0.1.so ./plugins/argocd/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/argocdapp_0.0.1.so ./plugins/argocdapp/
	go build -trimpath -gcflags="all=-N -l" -o openstream ./cmd/openstream/
