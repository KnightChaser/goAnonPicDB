FROM golang:1.21.5-alpine3.18 AS gobuilder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o server main.go

FROM scratch
COPY --from=gobuilder /build .
EXPOSE 8080

CMD ["./server"]
