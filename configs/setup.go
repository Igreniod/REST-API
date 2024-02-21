package configs

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func ConnectDB() *mongo.Client {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGOURI")))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

//		//ping the database
//		err = client.Ping(ctx, nil)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println("Connected to MongoDB")
//		return client
//	}
func ConnectDB() *mongo.Client {
	Env()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGOURI")))
	if err != nil {
		log.Fatal("Error : ", err)
	}
	defer cancel()

	// ctx, can = context.WithTimeout(context.Background(), 10*time.Second)
	// if err = client.Ping(ctx, readpref.Primary()); err != nil {
	// 	fmt.Printf("could not ping to mongo db service: %v\n", err)
	// 	return nil
	// }

	log.Println("DB Connected : ", os.Getenv("MONGOURI"))
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, dataBaseName string, collectionName string) *mongo.Collection {
	collection := client.Database(dataBaseName).Collection(collectionName)
	return collection
}
