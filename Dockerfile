FROM golang:alpine AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download


FROM golang:alpine AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd


FROM scratch AS runner
COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
CMD [ "/app" ]