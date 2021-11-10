build:
	go build -o openstream ./cmd/openstream/

copy:
	mkdir -p plugins
	cp ../openstream-plugin-githubactions/githubactions_0.0.1.so plugins/

test:
	./openstream install repo-scaffolding-go
	./openstream reinstall repo-scaffolding-go
	./openstream uninstall repo-scaffolding-go
