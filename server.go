package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)
/* HANDLING OPERATIONS */

// types

// HandlerFunction represents a fiber handler
type HandlerFunction func(c *fiber.Ctx) error
// HandlerType is an enum to represent a specific handler type instead of using string names to reach.
type HandlerType int
// Handler represents a handle operation by keeping the Handling Route and Handling Function as pair.
type Handler struct{
	Route string
	Function HandlerFunction
}
// GuestLedgerBookHandlers is a map structure which provides accessing current related Handler by its HandlerType.
type GuestLedgerBookHandlers map[HandlerType]Handler

// constants

const(
	GetAllBook HandlerType = iota
	GetGuestLedger
	AddGuestLedger
)

//init

func(handlers GuestLedgerBookHandlers) InitHandlers(){
	handlers[GetAllBook] = Handler{
		Route:    "/book",
		Function: GetAllBookHandler,
	}
	handlers[GetGuestLedger] = Handler{
		Route:    "/book/:email",
		Function: GetGuestLedgerHandler,
	}
	handlers[AddGuestLedger] = Handler{
		Route:    "/add",
		Function: AddGuestLedgerHandler,
	}
}
//Handling Functions

func GetAllBookHandler(c *fiber.Ctx) error{
	return c.JSON(GuestLedgerBook)
}
func GetGuestLedgerHandler(c *fiber.Ctx) error{
	email := c.Params("email")
	if len(email) < 0{
		return errors.New("email is empty")
	}
	guestLedger,err := findGuestLedgerByEmail(email)
	if err!= nil {
		return err
	}
	return c.JSON(guestLedger)
}
func AddGuestLedgerHandler(c *fiber.Ctx) error{
	var guestLedger GuestLedger
	if err := c.BodyParser(&guestLedger); err != nil {
		log.Fatal(err)
	}
	if len(guestLedger.Email) > 0 {
		GuestLedgerBook = append(GuestLedgerBook,guestLedger)
	}
	return c.JSON(GuestLedgerBook)
}

/* SERVER INITIALIZE */

// server variables
var serverInitalized = false // to control is server prepared
var Handlers = GuestLedgerBookHandlers{}
var app = fiber.New()

//server functions

func InitServer() {
	app.Use(cors.New())
	Handlers.InitHandlers()
	// Getting All Books
	app.Get(Handlers[GetAllBook].Route, Handlers[GetAllBook].Function)
	// Getting A Guest Ledger By Email
	app.Get(Handlers[GetGuestLedger].Route,Handlers[GetGuestLedger].Function)
	// Posting A Guest Ledger
	app.Post(Handlers[AddGuestLedger].Route,Handlers[AddGuestLedger].Function)
	serverInitalized = true
}

func StartServer(port int)error{
	if !serverInitalized {
		InitServer()
	}
	err := app.Listen(fmt.Sprintf(":%d",port))
	return err
}


