INSTALL_ONCE=src public package.json

dist: vite.config.ts src public package.json
	@echo "run build..."
	@npm run build

dist/manifest.json dist/assets: dist

go.mod:
	@go mod init vitemodtest
	@go get github.com/torenware/vite-go@latest

go.sum: go.mod

run: go.sum dist/manifest.json dist/assets test-template.tmpl
	@go run .


