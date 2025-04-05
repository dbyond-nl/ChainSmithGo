FROM golang:alpine AS build-stage

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=v0.0.0-devel

WORKDIR /go/src/app
COPY . .

RUN go get -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -v -a -ldflags="-X 'github.com/dbyond-nl/chainsmithgo/internal/version/appVersion=${VERSION}'" -o chainsmith ./cmd/chainsmith

FROM alpine AS export-stage
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=build-stage /go/bin/chainsmith /usr/bin/
CMD /usr/bin/chainsmith
