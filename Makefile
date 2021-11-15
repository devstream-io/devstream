all: build copy

build:
	go build -trimpath -gcflags="all=-N -l" -o openstream ./cmd/openstream/

copy:
	mkdir -p plugins
	cp ../openstream-plugin-githubactions/githubactions_0.0.1.so plugins/
	cp ../openstream-plugin-argocd/argocd_0.0.1.so plugins/
	cp ../openstream-plugin-argocdapp/argocdapp_0.0.1.so plugins/
