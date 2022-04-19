FROM golang:1.18-alpine as build
WORKDIR /src
COPY go.sum go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/server ./cmd/server
RUN CGO_ENABLED=0 go build -o /bin/cli ./cmd/cli

FROM scratch
COPY --from=build /bin/server /bin/server
COPY --from=build /bin/cli /bin/cli
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /bin
CMD ["./server"]
