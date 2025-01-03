package main

import (
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	"github.com/rahulgit24/goprojects/go-fiber-crm-basic/database"
	"github.com/rahulgit24/goprojects/go-fiber-crm-basic/lead"
)

func setUpRoutes(app *fiber.App){
	app.Get("/api/v1/lead",lead.GetLeads)
	app.Get("/api/v1/lead/:id",lead.GetLead)
	app.Post("/api/v1/lead",lead.NewLead)
	app.Delete("/api/v1/lead/:id",lead.DeleteLead)
}

func initDb(){
	var err error
	database.DBConn, err = gorm.Open("sqlite3","leads.db")
	if err!=nil{
		panic("Failed to connect database")
	}
	fmt.Printf("Connection opened to database")
	database.DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("Database Migrated")
}

func main(){
	app := fiber.New()
	initDb()
	setUpRoutes(app)
	app.Listen(3000)
	defer database.DBConn.Close()
}