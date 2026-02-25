# ---------- Frontend Build ----------
FROM node:18-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend .
RUN npm run build && npx next export

# ---------- Backend Build ----------
FROM golang:1.22-alpine AS backend
WORKDIR /app
COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

COPY backend .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# ---------- Final Stage ----------
FROM alpine:latest
WORKDIR /root/

COPY --from=backend /app/backend/app .
COPY --from=frontend /app/frontend/out ./frontend/out

EXPOSE 8080
CMD ["./app"]