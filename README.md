# Hotel Reservation

## Project Outline
- users -> book room from an hotel
- admins -> going to check reservation/bookings
- Authentication and Authorization -> JWT Tokens
- Hotels -> CURD API -> JSON
- Rooms -> CURD API -> JSON
- Scripts -> database management -> seeding migration


## Resoures
### Mongodb
Documentation
```
https://mongodb.com/docs/drivers/go/current/quick-start

```
Installing mongodb client
```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber
Documentation
```
https://gofiber.io
```
Installing gofiber
```
go get github.com/gofiber/fiber/v2
```
## Docker
### Installing mongodb as a docker container
```
docker run --name mongodb -d mongo:latest   -p 27017:27017
```

