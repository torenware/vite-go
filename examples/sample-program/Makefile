INSTALL_ONCE=src public package.json
CONFIG_FILE := vite.config.ts

ifneq (, $(wildcard ./vite.config.js))
  # if this is present, use it for config.
  CONFIG_FILE := vite.config.js
endif

$(INSTALL_ONCE): run-install-check
	:

run-install-check:
	./install-vue.sh

dist: $(CONFIG_FILE) src public package.json
	@echo "run build..."
	@npm run build

dist/manifest.json dist/assets: dist

go.mod:
	@go mod init vitemodtest
	@go get github.com/torenware/vite-go@latest

go.sum: go.mod

run: go.sum dist/manifest.json dist/assets test-template.tmpl
	@go run .


