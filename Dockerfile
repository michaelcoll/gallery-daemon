# Start by building the application.
FROM golang:1.19 as build

ARG VERSION

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o /go/bin/gallery-daemon -ldflags="-s -w -X 'github.com/michaelcoll/gallery-daemon/cmd.version=$VERSION'"

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=build /go/bin/gallery-daemon /bin/gallery-daemon

ENV NAME=docker-daemon
ENV OWNER=no@name.com
ENV FOLDER=.

EXPOSE 9000

CMD ["gallery-daemon", "-n", ""]
