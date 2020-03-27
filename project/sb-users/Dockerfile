FROM golang

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["./main", "user:password@tcp(db:3306)/sports_betting_db"]