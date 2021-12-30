FROM golang:1.17 AS build
WORKDIR /go/src
COPY go ./go
COPY go.mod .
COPY go.sum .
COPY main.go .

ENV CGO_ENABLED=0
#RUN go get -d -v ./... 

RUN go build -a -installsuffix cgo -o gigs .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY ./api/openapi.yaml ./api/openapi.yaml
COPY --from=build /go/src/gigs .
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]
