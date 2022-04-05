JS_DIR=frontend
INSTALL_ONCE=$(JS_DIR)
CONFIG_FILE := vite.config.ts
VITE_PORT=3000
GO_APP_PORT=4000

clean:
	@echo Cleaning up...
	@- rm -r test_program dist
	@echo cleaned.

ifneq (, $(wildcard ./vite.config.js))
  # if this is present, use it for config.
  CONFIG_FILE := vite.config.js
endif

$(INSTALL_ONCE): run-install-check
	:

run-install-check:
	./install-vue.sh

dist: $(CONFIG_FILE) $(JS_DIR)
	@echo "run javascript build..."
	@cd $(JS_DIR); npm run build

dist/manifest.json dist/assets: dist

go.mod:
	@go mod init vitemodtest
	@go get github.com/torenware/vite-go@latest

go.sum: go.mod

dev: go.sum $(JS_DIR)
	@echo starting dev server
	@-pkill -SIGTERM -f "npm run dev"
	@cd $(JS_DIR); npm run dev &
	@go run .

build: clean go.sum dist/manifest.json dist/assets test-template.tmpl
	@echo building go binary...
	@go build -o test_program .

preview: build
	@echo run test_program
	@ ./test_program -env production -assets dist



