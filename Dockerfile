FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o simple-api-gateway ./cmd/simple-api-gateway/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/simple-api-gateway ./
COPY --from=builder /app/.env ./
COPY --from=builder /app/frontend ./ 

EXPOSE 4000

ENTRYPOINT [ "./simple-api-gateway" ]