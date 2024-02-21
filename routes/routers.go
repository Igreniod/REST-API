package routes

import (
	"fmt"
	"latihandatabasegolang/configs"
	"latihandatabasegolang/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func SetRoutes() {
	configs.Env()
	router := mux.NewRouter()

	router.HandleFunc("/users", controllers.GetAllUser).Methods("GET")
	router.HandleFunc("/user/{userId}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/addnewuser", controllers.AddNewUser).Methods("POST")
	router.HandleFunc("/edituser/{userId}", controllers.EditUser).Methods("POST")
	router.HandleFunc("/deleteuser/{userId}", controllers.DeleteUser).Methods("DELETE")

	fmt.Printf("Server is running on port : %s\n", os.Getenv("SERVICE_ADDRESS"))
	http.ListenAndServe(":"+os.Getenv("SERVICE_ADDRESS"), router)
}
