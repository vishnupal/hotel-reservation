package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vishnupal/hotel-reservation/api"
	"github.com/vishnupal/hotel-reservation/api/middleware"
	"github.com/vishnupal/hotel-reservation/db"
)

// const collection = 'user'

func main() {
	config := fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.JSON(map[string]string{"error": err.Error()})
		},
	}

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	var (
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		hotelhandler = api.NewHotelHandler(store)
		app          = fiber.New(config)
		auth         = app.Group("api")
		apiv1        = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	// auth Handler

	auth.Post("/auth", authHandler.HandleAunthenticate)

	// Version Apis
	// user Hnadler
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandlerGetUsers)
	apiv1.Get("/user/:id", userHandler.HandlerGetUser)

	// Hotel Handler
	apiv1.Get("/hotel", hotelhandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelhandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelhandler.HandleGetRooms)
	app.Listen(*listenAddr)
}
