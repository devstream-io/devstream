build:
	go get ./...
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/githubactions_0.0.1.so ./cmd/githubactions/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/argocd_0.0.1.so ./cmd/argocd/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/argocdapp_0.0.1.so ./cmd/argocdapp/
	go build -trimpath -gcflags="all=-N -l" -o dsm ./cmd/devstream/
