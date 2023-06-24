FROM golang:1.20 as builder

WORKDIR /go/src/app

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

# download if above files changed
RUN go mod download

COPY . .

RUN go build -o /bin/server ./cmd

EXPOSE 8080


CMD ["/bin/server"]