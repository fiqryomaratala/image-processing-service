FROM node:22-alpine

WORKDIR /app/frontend

COPY frontend/ ./

CMD ["sh", "-c", "echo \"Frontend placeholder for future phases\" && tail -f /dev/null"]
