build:
	go build -o api

run:
	go run *.go

dev:
	goreleaser --snapshot --skip-publish --rm-dist

dev-publish:
	goreleaser --snapshot --rm-dist
