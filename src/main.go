package main

import (
	"team01/database"
)

func main() {
	db := database.New()
	println(db)
}
