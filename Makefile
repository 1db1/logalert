BINARY_NAME=logalert

UNAME := $(shell uname)

build:
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin ./
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux ./

install: build
ifeq ($(UNAME), Linux)
	mkdir -p /var/lib/${BINARY_NAME}
	cp ./bin/${BINARY_NAME}-linux /usr/sbin/${BINARY_NAME}
endif
ifeq ($(UNAME), Darwin)
	mkdir -p /usr/local/var/lib/${BINARY_NAME}
	cp ./bin/${BINARY_NAME}-darwin /usr/local/sbin/${BINARY_NAME}
endif

run:
ifeq ($(UNAME), Linux)
	/usr/sbin/${BINARY_NAME}
endif
ifeq ($(UNAME), Darwin)
	/usr/local/sbin/${BINARY_NAME}
endif

uninstall:
ifeq ($(UNAME), Linux)
	rm /usr/sbin/${BINARY_NAME}
	rm -rf /var/lib/${BINARY_NAME}
endif
ifeq ($(UNAME), Darwin)
	rm /usr/local/sbin/${BINARY_NAME}
	rm -rf /usr/local/var/lib/${BINARY_NAME}
endif

clean:
	go clean
	rm -f ./bin/*
