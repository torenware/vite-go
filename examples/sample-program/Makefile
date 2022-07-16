JS_DIR=frontend
INSTALL_ONCE=$(JS_DIR)
CONFIG_FILE := vite.config.ts
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

# we don't even need a config file
# since Vanilla builds lack these.
$(JS_DIR)/dist:  $(JS_DIR) 
	@echo "run javascript build..."
	@cd $(JS_DIR); node_modules/.bin/vite build  --manifest manifest.json

$(JS_DIR)/node_modules:
	@cd $(JS_DIR); npm install

$(JS_DIR)/dist/manifest.json $(JS_DIR)/dist/assets: $(JS_DIR)/dist

go.mod:
	@go mod init vitemodtest
	@go get github.com/torenware/vite-go@latest

go.sum: go.mod

dev:  stop_dev go.sum $(JS_DIR)
	@echo starting dev server
	@ VITE_PID=$(VITE_PID) ./start-vite.sh $(JS_DIR)
	@ go run . -pid $(GO_PID) &

dev_go: stop_dev go.sum
	@echo starting go server only...
	@ go run . -pid $(GO_PID) &


build: clean go.sum $(JS_DIR)/dist/manifest.json $(JS_DIR)/dist/assets test-template.tmpl
	@echo building go binary...
	@go build -o test_program .

preview: stop_preview  $(JS_DIR)/node_modules build
	@echo run test_program
	@ ./test_program -env production -assets $(JS_DIR)/dist -pid /tmp/vite-go.pid

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

stop_preview:
ifneq (,$(wildcard $(GO_PID)))
	@echo Stopping test_program
	@! ps -p $$(cat $(GO_PID)) &>/dev/null || kill $$(cat $(GO_PID) 2>/dev/null) > /dev/null
else
	@echo Go already stopped
endif

