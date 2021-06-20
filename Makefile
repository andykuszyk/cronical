test:
	go test ./... -v

build:
	docker build -t andykuszyk/cronical:$(GITHUB_REF)

publish: test build
	docker push andykuszyk/cronical:$(GITHUB_REF)
