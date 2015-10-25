# Building Location and Trip Planner Service Project

```
cd LocationRestWebServices
export GOPATH=$PWD # Set GO PATH
go get -d -v ./... # Get Dependencies
go get gopkg.in/tomb.v2 #Get dependencies
go build -v ./...  # Build the project
go build src/server/Server.go # Build Server 
```

# Run the server using
```
go run src/server/Server.go
```
or 
```
./Server 
```

# API calls
## Get list of all locations

```
curl -X GET http://127.0.0.1:4000/locations
```

Sample response is : 
HTTP Response CODE : 200 
```
{"id":"561b2eba5cb8b322702fbe82","name":"Renu Batheja","address":"228 Park Ave","city":"New York","state":"NY","zipcode":"10003","coordinate":{"lat":40.737705,"lng":-73.9887481}}
 {"id":"561b2f365cb8b322702fbe83","name":"Kapil Bhalla","address":"978 Henderson Ave","city":"New York","state":"NY","zipcode":"94086","coordinate":{"lat":40.6365687,"lng":-74.1206583}}
 {"id":"561b304d5cb8b322702fbe84","name":"Tanvi Bhalla","address":"978 Henderson Ave","city":"Sunnyvale","state":"CA","zipcode":"94086","coordinate":{"lat":37.3553733,"lng":-122.0046902}}
```

## Get list of locations by ID

```
curl -X GET http://127.0.0.1:4000/locations/561b2eba5cb8b322702fbe82
```

Sample response is : 
HTTP Response CODE : 200 
```
{"id":"561b2eba5cb8b322702fbe82","name":"Renu Batheja","address":"228 Park Ave","city":"New York","state":"NY","zipcode":"10003","coordinate":{"lat":40.737705,"lng":-73.9887481}}
```

## Create a new location 

```
curl -X POST http://127.0.0.1:4000/locations -H 'Content-Type:application/json'  -d '{ "name" : "Prince Smith", "address" : "1230 Henderson Aveunue", "city" : "Sunnyvale", "state" : "CA", "zip" : "94086" }'
```

Sample response is : 
HTTP Response CODE : 201
```
{"id":"562c18ff5cb8b32170cd1a72","name":"Prince Smith","address":"1230 Henderson Aveunue","city":"Sunnyvale","state":"CA","zipcode":"","coordinate":{"lat":37.3541509,"lng":-122.0047329}}
``` 


## Update existing location

```
curl -X PUT http://127.0.0.1:4000/locations/562c18ff5cb8b32170cd1a72 -H 'Content-Type:application/json'  -d '{ "address" : "228 Park Ave", "city" : "New York", "state" : "NY", "zip" : "10003" }'
```

Sample response is : 
HTTP Response CODE : 201 
```
{"id":"562c18ff5cb8b32170cd1a72","name":"Prince Smith","address":"1901 Halford Aveunue","city":"Santa Clara","state":"CA","zipcode":"","coordinate":{"lat":37.3580054,"lng":-121.9971017}}
```

## Delete existing location
```
curl -X DELETE http://127.0.0.1:4000/locations/562c18ff5cb8b32170cd1a72
```
HTTP Response CODE : 200 
```
```

Verify that above record is deleted

```
curl -X GET http://127.0.0.1:4000/locations/562c18ff5cb8b32170cd1a72
```

above doesn't return any data