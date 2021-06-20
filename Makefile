test:
	go test ./... -v

publish: test
	./scripts/docker-publish.sh
