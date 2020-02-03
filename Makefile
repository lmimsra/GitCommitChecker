# simpleなmakefile

build:
	GOOS=linux GOARCH=amd64 go build
	@echo "complete build!"

package: build
	zip -r deploy.zip git-commit-checker .env
	@echo "complete packaging"
	@echo "please deploy \"deploy.zip\""
