
vue-install:
	@npm create vite@latest vue-temp -- --template vue
	@cd vue-temp; cp -r src public package.json ..
	@rm -r vue-temp
	@npm install

dist: vue-install vite.config.ts src public package.json
	@npm run build

dist/manifest.json: dist

go.mod:
	@go mod init vitemodtest
	@go get github.com/torenware/vite-go@latest

go.sum: go.mod

run: go.sum dist/manifest.json dist/assets test-template.tmpl
	@go run .

