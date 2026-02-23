FROM golang:latest as builder

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY . .

RUN CGO_ENABLED=0 go build -o simple-api-gateway ./cmd/simple-api-gateway/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/simple-api-gateway ./
COPY --from=builder /app/.env ./
COPY --from=builder /app/frontend ./ 

EXPOSE 4000
EXPOSE 50051
EXPOSE 50052
EXPOSE 50053
EXPOSE 50054
EXPOSE 50055
EXPOSE 27017
EXPOSE 5432

ENTRYPOINT [ "./simple-api-gateway" ]