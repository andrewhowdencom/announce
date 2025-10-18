# build stage
FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /announce .

# final stage
FROM scratch

COPY --from=builder /announce /announce

ENTRYPOINT ["/announce"]