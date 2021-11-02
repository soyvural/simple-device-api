# Build 
FROM golang:alpine AS build
ENV GO111MODULE=on

ARG BUILD_VERSION=development

WORKDIR /go/src/github.com/soyvural/simple-device-api
ADD . .

WORKDIR /go/src/github.com/soyvural/simple-device-api/cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/simpledeviceapi

# Deploy, after executable file created we don't need go anymore.
FROM scratch
ENV GIN_MODE=release
COPY --from=build /go/bin/simpledeviceapi /simpledeviceapi
CMD ["/simpledeviceapi", "--port", "8080"]