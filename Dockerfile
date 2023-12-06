# build
FROM docker.io/golang:1.19 as build
ARG VER
RUN mkdir /gorestapi
COPY . /gorestapi
WORKDIR /gorestapi
RUN go build -ldflags="-X 'main.Version=$VER'" -tags netgo -o bin ./...

# run 
FROM docker.io/busybox
COPY --from=build /gorestapi/bin/* /go/bin/