FROM golang:1.20-buster as builder

WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum

# download if above files changed
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags '-w -s -buildid=' -a -o /bin/server ./cmd


CMD ["/bin/server"]