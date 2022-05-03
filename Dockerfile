FROM golang:1.17-alpine AS builder
ENV GO111MODULE=on
WORKDIR /BusinessWallet
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/server

FROM scratch
WORKDIR /bin/BusinessWallet/
COPY --from=builder /BusinessWallet/config.json ./config.json
COPY --from=builder /bin/server ./server
EXPOSE 4242
CMD ["./server"]