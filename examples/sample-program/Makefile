JS_DIR=frontend
INSTALL_ONCE=$(JS_DIR)
CONFIG_FILE := vite.config.ts
# VITE_PORT=3000
GO_APP_PORT=4000
VITE_PID=/tmp/vite-script.pid
GO_PID=/tmp/vite-go.pid

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

dev:  stop_dev go.sum $(JS_DIR)
	@echo starting dev server
	@ VITE_PID=$(VITE_PID) ./start-vite.sh $(JS_DIR)
	@ go run . -pid $(GO_PID) &


build: clean go.sum dist/manifest.json dist/assets test-template.tmpl
	@echo building go binary...
	@go build -o test_program .

preview: build
	@echo run test_program
	@ ./test_program -env production -assets dist

stop_dev:
ifneq (,$(wildcard $(VITE_PID)))
	@ echo Stopping vitemodtest
	@! ps -p $$(cat $(VITE_PID)) &>/dev/null || kill $$(cat $(VITE_PID) 2>/dev/null) &> /dev/null
else
	@echo Vite already stopped
endif
ifneq (,$(wildcard $(GO_PID)))
	@echo Stopping go run
	@! ps -p $$(cat $(GO_PID)) &>/dev/null || kill $$(cat $(GO_PID) 2>/dev/null) > /dev/null
else
	@echo Go already stopped
endif
