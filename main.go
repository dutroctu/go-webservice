package main

import (
	"net/http"

	"github.com/dutroctu/go-webservice/database"
	"github.com/dutroctu/go-webservice/product"
	"github.com/dutroctu/go-webservice/receipt"
	_ "github.com/go-sql-driver/mysql"
)

// type URL struct {
// 	Scheme     string
// 	Opaque     string
// 	User       *UserInfo
// 	Host       string
// 	Path       string
// 	RawPath    string
// 	ForceQuery string
// 	RawQuery   string
// 	Fragment   string
// }

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)

	http.ListenAndServe(":5000", nil)
}
