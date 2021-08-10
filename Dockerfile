FROM golang:1.13-alpine as build

WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -o app

FROM scratch
COPY --from=build /build /
EXPOSE 9580
ENTRYPOINT [ "/app" ]