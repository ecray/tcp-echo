FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o tcp-echo .

FROM scratch
COPY --from=builder /app/tcp-echo /tcp-echo
EXPOSE 9095
ENTRYPOINT ["/tcp-echo"]
