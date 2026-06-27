FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o recommender .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/recommender .
COPY --from=builder /app/data ./data

ENV NUMBER_OF_LATENT_FACTORS=20
ENV LEARNING_RATE=0.01
ENV REGULARIZATION_PARAMETER=0.02
ENV NUMBER_OF_EPOCHS=20
ENV INITIALIZATION_MIN=-0.1
ENV INITIALIZATION_MAX=0.1

CMD ["./recommender"]