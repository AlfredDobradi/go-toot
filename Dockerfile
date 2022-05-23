FROM golang:1.18-buster AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -tags netgo -o toot ./cmd/...

FROM debian:buster
COPY --from=build /build/toot /usr/bin/toot
RUN chmod +x /usr/bin/toot
ENTRYPOINT ["/usr/bin/toot"]
