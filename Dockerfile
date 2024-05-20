FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . ./
RUN go build -o ./cmd/ cmd/server.go

FROM alpine as runner

# Устанавливаем bash и копируем wait-for-it.sh
RUN apk add --no-cache bash
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

COPY --from=builder /usr/local/src/cmd/ /

# Используем wait-for-it.sh для ожидания базы данных
CMD ["wait-for-it.sh", "ps-psql:5432", "--", "/server"]
