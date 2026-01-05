package main

import (
	"fmt"
	"log"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/database"

	"github.com/labstack/echo/v4"
)

func startServer(port string) {
	e := echo.New()
	e.GET("/hey", func(c echo.Context) error {
		return c.String(200, "hey!")
	})
	e.Logger.Fatal(e.Start(":" + port))

}
func main() {
	// Initialize database
	dbPath := "./data/chat.db"
	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.RunMigrations("./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	PORT := "8080"
	fmt.Println("starting server at port", PORT)
	startServer(PORT)
}
