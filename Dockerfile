FROM node:alpine AS build-fe
WORKDIR /build
COPY web .
RUN yarn --immutable
RUN yarn build

FROM golang:alpine AS build-be
WORKDIR /build
COPY cmd cmd
COPY pkg pkg
COPY internal internal
COPY go.mod .
COPY go.sum .
COPY migrations internal/embedded/migrations
COPY --from=build-fe /build/dist pkg/webserver/_webdist
RUN rm -f pkg/webserver/_webdist/.keep
RUN go build -o bin/yuri cmd/yuri/main.go

FROM alpine:latest
COPY --from=build-be /build/bin/yuri /var/opt/yuri
RUN apk add ffmpeg
EXPOSE 80
EXPOSE 6969
ENTRYPOINT ["/var/opt/yuri"]
