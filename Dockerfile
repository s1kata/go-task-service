FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app ./main  

FROM alpine:3.19 
WORKDIR /app
COPY --from=builder /app/app .
RUN adduser -D appuser
USER appuser
EXPOSE 8080
ENTRYPOINT [ "./app" ]