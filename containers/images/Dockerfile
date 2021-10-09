FROM golang:1.16-alpine
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/app main.go

FROM alpine
WORKDIR /src

COPY --from=0 /bin/app /bin/app

ENTRYPOINT ["/bin/app"]
