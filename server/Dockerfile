# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine
RUN apk add font-terminus font-inconsolata font-dejavu font-noto font-noto-cjk font-awesome font-noto-extra
RUN apk add ffmpeg
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 3000
CMD [ "/app/main" ]
