FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go mod verify
RUN GOOS=linux go build -ldflags="-s -w" ./cmd/web

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/web /go/bin/web
COPY --from=build /go/src/app/markdown /usr/bin/markdown
EXPOSE 4000
ENTRYPOINT ["/go/bin/web"]