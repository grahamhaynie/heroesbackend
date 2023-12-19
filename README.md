# Second app backend
HTTP server for second app

## TODO 
- fix pipeline
- fix photo uploading logic. kinda goofy. also fix case of duplicate filename
- add database
- openapi
- correct return codes
- kube
- fix windows debug popup

## Usage
```
curl localhost:8080/api/heroes
curl localhost:8080/api/heroes/12
curl -XPOST localhost:8080/api/heroes -H "Content-Type: application/json" --data '{"Id":69,"Name":"bob","Power":"none","AlterEgo":"nobody"}'
curl -XPUT localhost:8080/api/heroes -H "Content-Type: application/json" --data '{"Id":19,"Name":"notbob","Power":"making rocks","AlterEgo":"somebody"}'
curl -XDELETE localhost:8080/api/heroes/19
curl localhost:8080/api/heroes?name=torn
```