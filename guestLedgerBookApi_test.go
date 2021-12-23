package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/assert" // add Testify package
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Test struct{
	Description string
	Route string
	ExpectedCode int
	Method string `default:"GET"`
	Body	io.Reader
}

func TestGetGuestLedgerBook(t *testing.T){

	test := Test{
			Description: "get http status 200,when successfully get all guest ledger book",
			Route: "/book",
			ExpectedCode: http.StatusOK,
			Method: http.MethodGet,
	}
	app:= fiber.New()
	app.Use(cors.New())
	GuestLedgerBook = append(GuestLedgerBook,GuestLedger{Email: "akl@gmail.com",Message: "Hi There!"})
	app.Get("/book",func(c *fiber.Ctx) error{
		return c.JSON(GuestLedgerBook)
	})

	req := httptest.NewRequest(test.Method,test.Route, nil)
	resp,_ := app.Test(req,1)

	currentLedgerBookAsByte ,_:= json.Marshal(GuestLedgerBook)//Marshal converts given interface to json represents as []byte. And the string(returned_[]byte) will give as json represents. Unmarshal to reverse
	sendedGuestLedgerBookAsByte,_ := ioutil.ReadAll(resp.Body)
	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.JSONEq(t, string(currentLedgerBookAsByte),string(sendedGuestLedgerBookAsByte),"Current book and sended book are totally equal")
}

func TestGetGuestLedgerByEmail(t *testing.T){
	emailToSearch := "akl@gmail.com" // test parameter
	test:= Test{
		Description:  "get http status 200, when successfully get a guest ledger by email",
		Route:        fmt.Sprintf("/book/%s",emailToSearch) ,
		ExpectedCode: http.StatusOK,
		Method: http.MethodGet,
	}
	app:= fiber.New()
	app.Use(cors.New())
	GuestLedgerBook = append(GuestLedgerBook,GuestLedger{Email: emailToSearch,Message: "Hi There!"}) // add a mock data to look it goes correctly

	app.Get("/book/:email", func(c *fiber.Ctx) error {
		email := c.Params("email")
		if len(email) < 0{
			return errors.New("email is empty")
		}
		guestLedger,err := findGuestLedgerByEmail(email)
		if err!= nil {
			return err
		}
		return c.JSON(guestLedger)
	})

	req := httptest.NewRequest(test.Method,test.Route, nil)
	resp,_ := app.Test(req,1)
	currentGuestLedger,_ := findGuestLedgerByEmail(emailToSearch)
	currentGuestLedgerAsByte,_ := json.Marshal(currentGuestLedger)
	sendedGuestLedgerAsByte,_ := ioutil.ReadAll(resp.Body)

	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.JSONEq(t, string(currentGuestLedgerAsByte),string(sendedGuestLedgerAsByte),"Sended guest ledger and current guest ledger ( has the email) are equal")
}

func TestPostNewGuestLedger(t *testing.T){
	test:= Test{
		Description:  "get http status 200, when post a guest ledger is succeed",
		Route:        "/add",
		ExpectedCode: http.StatusOK,
		Method: http.MethodPost,
		Body:  strings.NewReader(`{"email": "axd@gmail.com", "message": "Hello"}`),
	}
	app:= fiber.New()
	app.Use(cors.New())
	GuestLedgerBook = append(GuestLedgerBook,GuestLedger{Email: "akl@gmail.com",Message: "Hello"})

	app.Post("/add", func(c *fiber.Ctx) error {
		var guestLedger GuestLedger
		if err := c.BodyParser(&guestLedger); err != nil {
			log.Fatal(err)
		}
		if len(guestLedger.Email) > 0 {
			GuestLedgerBook = append(GuestLedgerBook,guestLedger)
		}
		return c.JSON(GuestLedgerBook)
	})

	req := httptest.NewRequest(test.Method,test.Route, test.Body)

	//if you do not set this, the post operation body will come empty like {email:"",message:""} and will give 422 Unprocessible Entity HTTP Status Code
	req.Header.Set("Content-Type","application/json")

	resp,_ := app.Test(req,1)
	currentLedgerBook,_:= json.Marshal(GuestLedgerBook)//Marshal converts given interface to json represents as []byte. And the string(returned_[]byte) will give as json represents. Unmarshal to reverse
	sendedLedgerBook,_ := ioutil.ReadAll(resp.Body) // Read response body as byte to compare with current ledger book.
	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)// First compare to status code to be expected
	assert.JSONEq(t,string(currentLedgerBook),string(sendedLedgerBook),"Both post response and current ledger book are same") //Compare for data came successfully with post operation
}
