FROM node:20-alpine AS frontend-builder

WORKDIR /app

COPY frontend-react-js/package*.json ./frontend-react-js/

WORKDIR /app/frontend-react-js

RUN npm install

WORKDIR /app

COPY frontend-react-js ./frontend-react-js
WORKDIR /app/frontend-react-js
RUN npm run build


FROM golang:latest AS go-builder

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY . .

RUN mkdir -p ./frontend-dist
COPY --from=frontend-builder /app/frontend-react-js/dist ./cmd/simple-api-gateway/frontend-dist

RUN go mod download

RUN CGO_ENABLED=0 go build -o simple-api-gateway ./cmd/simple-api-gateway/main.go


FROM scratch

WORKDIR /app

COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=go-builder /app/simple-api-gateway ./
COPY --from=go-builder /app/.env ./
COPY --from=go-builder /app/frontend ./ 
COPY --from=go-builder /app/certs ./certs

EXPOSE 4000

ENTRYPOINT [ "./simple-api-gateway" ]