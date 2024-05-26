package main

import (
	"team01/pkg/database"
)

func main() {
	db := database.New()
	println(db)
}
