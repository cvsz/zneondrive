FROM golang:1.27.1-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/zneondrive-api ./cmd/server

FROM alpine:3.22
RUN addgroup -S zneon && adduser -S -G zneon zneon
USER zneon
COPY --from=build /out/zneondrive-api /usr/local/bin/zneondrive-api
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/zneondrive-api"]
