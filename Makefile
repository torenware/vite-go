
BROWSERFY := $(shell command -v browserify 2> /dev/null )

clean:
	@echo clean up preable files...
	@- rm -rf package-lock.js node_modules react/preamble.js
	@echo cleaned up.

node_modules: package.json
	@echo installing node dependencies
	@npm install

node_modules/react-refresh/runtime.js: node_modules

node_modules/react-refresh: node_modules/react-refresh/runtime.js
	@echo installing react-refresh
	@npm install react-refresh@latest

react/preamble.js: node_modules/react-refresh/runtime.js
ifndef BROWSERFY
	@echo browserfy not found so installing
	@npm install -g browserify
endif
	@echo installing preamble
	@browserify react/refresh-loader.js -o react/preamble.js

react: react/preamble.js

test:
	@echo running tests...
	@go test -v .

