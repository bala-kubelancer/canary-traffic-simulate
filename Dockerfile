# Dockerfile
FROM golang:alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app ./main.go

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /app /usr/local/bin/unidevops
COPY static /static
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/unidevops"]

