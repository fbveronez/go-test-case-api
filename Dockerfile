FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache git

# definir GOPATH e PATH
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# instalar ferramentas
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# gerar swagger
RUN swag init -g cmd/api/main.go -o docs

EXPOSE 8080

# usar air para hot reload
CMD ["air"]