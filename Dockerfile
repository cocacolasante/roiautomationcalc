FROM node:20-alpine AS embed-builder
WORKDIR /tool
COPY tool/package*.json ./
RUN npm install
COPY tool/ ./
RUN npm run build:embed

FROM node:20-alpine AS admin-builder
WORKDIR /admin
COPY admin/package*.json ./
RUN npm install
COPY admin/ ./
RUN npm run build

FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/templates ./templates
COPY --from=embed-builder /tool/dist ./embed
COPY --from=admin-builder /admin/dist ./admin-ui
EXPOSE 8093
CMD ["./server"]
