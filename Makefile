default: build

test:
	go test ./...

build:
	go build -o bin/yatas-gcp

update:
	go get -u 
	go mod tidy

install: build
	mkdir -p ~/.yatas.d/plugins/github.com/padok-team/yatas-gcp/local/
	mv ./bin/yatas-gcp ~/.yatas.d/plugins/github.com/padok-team/yatas-gcp/local/

release: test
	standard-version
	git push --follow-tags origin main 
