# Second app backend
HTTP server for second app

## TODO 
- unit tests for mongo
- add contexts https://go.dev/blog/context
- fix photo uploading logic. kinda goofy. also fix case of duplicate filename. use channels.
- maybe cleanup main into more files.
- finish todos
- cleanup git 
- openapi
- check return codes
- kube

## Mongodb
Before running with mongodb as the database, start a mongodb docker container (on windows, start docker desktop first)
```
docker run --name mongo -d -p 27017:27017 mongo
```
To run app with mongodb, specify the -u flag with the URI of the mongdb. For usage with above docker container, provide `-u mongodb://localhost:27017`

## Building
```
mkdir bin
go build -tags netgo -o bin ./...
```

## Usage
See above mongodb section.
```
cd bin
./heroes -h
```

To curl the webapp:
```
curl localhost:8080/api/heroes
curl localhost:8080/api/heroes/12
curl -XPOST localhost:8080/api/heroes -H "Content-Type: application/json" --data '{"Id":69,"Name":"bob","Power":"none","AlterEgo":"nobody"}'
curl -XPUT localhost:8080/api/heroes -H "Content-Type: application/json" --data '{"Id":19,"Name":"notbob","Power":"making rocks","AlterEgo":"somebody"}'
curl -XDELETE localhost:8080/api/heroes/19
curl localhost:8080/api/heroes?name=torn
```

Alternatively, use the angular frontend. link TODO