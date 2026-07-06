FROM golang:1.25

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

EXPOSE 8080

CMD ["go", "run", "./cmd/api"]
