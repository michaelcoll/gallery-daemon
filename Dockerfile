# Start by building the application.
FROM golang:1 AS build

ARG BUILDTIME
ARG VERSION
ARG REVISION

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o /go/bin/gallery-daemon -ldflags="-s -w -X 'github.com/michaelcoll/gallery-daemon/cmd.version=$VERSION'"

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=build /go/bin/gallery-daemon /bin/gallery-daemon

ENV DAEMON_NAME=docker-daemon
ENV OWNER=no@na.me

VOLUME /media

EXPOSE 9000

CMD ["gallery-daemon", "serve", "-f", "/media", "--local-db"]
