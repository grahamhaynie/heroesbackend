# build
FROM docker.io/golang:1.19 as build
RUN mkdir /gorestapi
COPY . /gorestapi
WORKDIR /gorestapi
RUN mkdir bin
RUN go build -tags netgo -o bin ./...

# run 
FROM docker.io/busybox
COPY --from=build /gorestapi/bin/* /go/bin/