BIN_DIR := ./bin/

certs:
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -key server.key -out server.pem -days 3650

docs:
	swag init --parseDependency --parseInternal --output docs

build:
	mkdir ${BIN_DIR}
	cp ./config.yaml ${BIN_DIR}
	go build -trimpath -ldflags="-w -s" -o ${BIN_DIR}server

run:
	go run .

clean:
	rm -r ${BIN_DIR}

test:
	go test -v

build-docker: build
	docker build . -t api-rest

run-docker: build-docker
	docker run -p 9090:9090 api-rest