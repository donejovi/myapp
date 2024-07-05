package main

import (
	"myapp/database"
	"myapp/routers"
)

func main() {
	database.Connect()
	database.Migrate()

	r := routers.SetupRouter()
	r.Run(":8080")

}
