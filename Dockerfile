FROM golang:1.17 as build

WORKDIR /app

COPY go.* /app/

RUN go mod download

COPY . .

RUN go build -o /app/superQueue

# Need glibc
FROM gcr.io/distroless/base
COPY --from=build /app/superQueue /app/

CMD ["/app/superQueue" ]
