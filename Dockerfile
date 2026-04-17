FROM golang:1.25.4-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd/server


FROM alpine:3.21

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=build /app/app .


CMD [ "./app" ]
