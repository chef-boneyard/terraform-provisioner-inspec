.PHONY: prep/plugins build/linux build/darwin build/install build

PROVISIONER_BINARY_NAME=terraform-provisioner-inspec
PLUGINS_DIR=~/.terraform.d/plugins

prep/plugins:
	mkdir -p ${PLUGINS_DIR}

build/darwin: prep/plugins
	CGO_ENABLED=0 GOOS=darwin installsuffix=cgo go build -o ./${PROVISIONER_BINARY_NAME}
	
build/linux: prep/plugins
	CGO_ENABLED=0 GOOS=linux installsuffix=cgo go build -o ./${PROVISIONER_BINARY_NAME}

install:
	cp ./${PROVISIONER_BINARY_NAME} ${PLUGINS_DIR}/${PROVISIONER_BINARY_NAME}

test:
	go test -v