# Second app backend
HTTP server for second app

## TODO 
- unit tests for mongo
- finish todos
- cleanup git - cleanup history, add release, tag
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