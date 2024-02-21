package main

import (
	"latihandatabasegolang/configs"
	"latihandatabasegolang/routes"
)

func main() {
	configs.Env()
	// fmt.Println(os.Getenv("MONGOURI"))

	// router := mux.NewRouter()

	//run database
	configs.ConnectDB()

	// //routes
	routes.SetRoutes()

	// log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVICE_ADDRESS"), router))

}

// func init() {
// }
