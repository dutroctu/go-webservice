package main

import (
	"net/http"

	"github.com/dutroctu/go-webservice/product"
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
	product.SetupRoutes(apiBasePath)

	http.ListenAndServe(":5000", nil)
}
