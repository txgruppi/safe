VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || date '+%Y%m%d-dirty')

build: clean
	gox -tags release \
		-ldflags '-X github.com/txgruppi/safe/buildinfo.Version=$(VERSION)' \
		-output 'bin/{{.Dir}}_{{.OS}}_{{.Arch}}' \
		-osarch 'linux/amd64' \
		-osarch 'linux/arm64' \
		-osarch 'darwin/amd64' \
		-osarch 'windows/amd64'

compress: build
	upx -9 ./bin/*

clean:
	rm -rf ./bin