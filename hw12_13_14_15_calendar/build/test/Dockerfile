FROM golang:1.20

ENV CODE_DIR /test_integration/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

#COPY . ${CODE_DIR}

#COPY . ${CODE_DIR}

#CMD go test -tags integration ./...