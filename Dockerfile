### Builder
FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o gotrade .

### Main image
FROM alpine

RUN apk add libc6-compat
COPY --from=builder /app/gotrade /usr/local/bin/gotrade
RUN chmod +x /usr/local/bin/gotrade

CMD [ "gotrade" ]
