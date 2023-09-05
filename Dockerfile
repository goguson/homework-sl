FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" .

FROM golang:alpine

WORKDIR /app

RUN apk update && apk add --no-cache curl

COPY --from=builder /app/homework-object-storage /usr/bin/

EXPOSE 8080

ENTRYPOINT ["homework-object-storage"]
