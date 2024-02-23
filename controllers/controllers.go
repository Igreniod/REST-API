package controllers

import (
	"context"
	"encoding/json"
	"latihandatabasegolang/configs"
	"latihandatabasegolang/models"
	"latihandatabasegolang/response"
	"latihandatabasegolang/tools"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, os.Getenv("DATABASE_NAME"), "Users")
var validate = validator.New()

func GetUser(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	params := mux.Vars(req)
	userId := params["userId"]
	var user models.User
	defer cancel()
	// objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"uid": userId}).Decode((&user))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		response := response.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(res).Encode(response)
		return
	}

	res.WriteHeader(http.StatusOK)
	response := response.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": user}}
	json.NewEncoder(res).Encode(response)

}

func GetAllUser(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		response := response.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(res).Encode(response)
		return
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			response := response.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(res).Encode(response)
		}

		users = append(users, singleUser)
	}

	res.WriteHeader(http.StatusOK)
	response := response.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": users}}
	json.NewEncoder(res).Encode(response)
}

func AddNewUser(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		response := response.UserResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(res).Encode(response)
		return
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		response := response.UserResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": validationErr.Error()}}
		json.NewEncoder(res).Encode(response)
		return
	}

	//Validator if user input symbol and number in input name
	if !tools.ValidateUserName(user.Nama) {
		res.WriteHeader(http.StatusBadRequest)
		response := response.UserResponse{Status: http.StatusBadRequest, Message: "Input Nama invalid"}
		json.NewEncoder(res).Encode(response)
		return
	}

	//Data yang akan di input ke database
	userId := tools.GenerateRandomString(10)
	newUser := models.User{
		Uid:    userId,
		Nama:   strings.TrimSpace(user.Nama),
		Email:  user.Email,
		Alamat: user.Alamat,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		response := response.UserResponse{Status: http.StatusInternalServerError, Message: "Gagal menginput data ke Database", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(res).Encode(response)
		return
	}

	res.WriteHeader(http.StatusCreated)
	response := response.UserResponse{Status: http.StatusCreated, Message: "Success", Data: map[string]interface{}{"data": result}}
	json.NewEncoder(res).Encode(response)
}

func EditUser(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	params := mux.Vars(req)
	userId := params["userId"]
	var user models.User
	defer cancel()

	//validate the request body
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		response := response.UserResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(res).Encode(response)
		return
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		response := response.UserResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": validationErr.Error()}}
		json.NewEncoder(res).Encode(response)
		return
	}

	update := bson.M{"nama": user.Nama, "email": user.Email, "alamat": user.Alamat}

	result, err := userCollection.UpdateOne(ctx, bson.M{"uid": userId}, bson.M{"$set": update})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		response := response.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"Data": err.Error()}}
		json.NewEncoder(res).Encode(response)
		return
	}

	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"uid": userId}).Decode(&updatedUser)

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			response := response.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"Data": err.Error()}}
			json.NewEncoder(res).Encode(response)
			return
		}
	}

	res.WriteHeader(http.StatusOK)
	response := response.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": updatedUser}}
	json.NewEncoder(res).Encode(response)
}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	params := mux.Vars(req)
	userId := params["userId"]
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"uid": userId})

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		response := response.UserResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(res).Encode(response)
		return
	}

	if result.DeletedCount < 1 {
		res.WriteHeader(http.StatusNotFound)
		response := response.UserResponse{Status: http.StatusNotFound, Message: "Error", Data: map[string]interface{}{"data": "User with specified ID not found!"}}
		json.NewEncoder(res).Encode(response)
		return
	}

	res.WriteHeader(http.StatusOK)
	response := response.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": "User deleted!!!"}}
	json.NewEncoder(res).Encode(response)
}
