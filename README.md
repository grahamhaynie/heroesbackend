# Second app backend
HTTP server for second app

## TODO 
- unit tests for database
- comment code
- add env variables database
- add contexts https://go.dev/blog/context
- fix photo uploading logic. kinda goofy. also fix case of duplicate filename. use channels.
- finish todos
- openapi
- check return codes
- kube
- fix windows debug popup

## Mongodb
Before running, start a mongodb docker container (on windows, start docker desktop first...)
```

```

## Usage
Start mongo docker container
```
curl localhost:8080/api/heroes
curl localhost:8080/api/heroes/12
curl -XPOST localhost:8080/api/heroes -H "Content-Type: application/json" --data '{"Id":69,"Name":"bob","Power":"none","AlterEgo":"nobody"}'
curl -XPUT localhost:8080/api/heroes -H "Content-Type: application/json" --data '{"Id":19,"Name":"notbob","Power":"making rocks","AlterEgo":"somebody"}'
curl -XDELETE localhost:8080/api/heroes/19
curl localhost:8080/api/heroes?name=torn
```