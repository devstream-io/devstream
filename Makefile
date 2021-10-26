build:
	go build -o openstream ./cmd/openstream/

copy:
	mkdir -p plugins
	cp ../openstream-plugin-repo-scaffolding-go/repo-scaffolding-go.so plugins

test:
	./openstream install repo-scaffolding-go
	./openstream reinstall repo-scaffolding-go
	./openstream uninstall repo-scaffolding-go
