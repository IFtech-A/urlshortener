FROM golang:1.16-alpine AS go-package

WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

FROM go-package AS go-builder

WORKDIR /app
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build -o shortener ./cmd/shortener/main.go

FROM node:14.16-alpine AS npm-builder

WORKDIR /app

COPY ./frontend/* ./

RUN npm ci
RUN npm run build

FROM alpine

WORKDIR /app

COPY --from=go-builder /app/shortener ./
COPY --from=npm-builder /app/build ./front

CMD [ "/app/shortener" ]
