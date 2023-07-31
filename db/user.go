package db

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIdForUser() string {
	var dbUrl string
	dbUrl = "mongodb+srv://" + os.Getenv("MONGO_USERNAME") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster1.5jqwhvz.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")

	userResult, err := usersCollection.InsertOne(ctx, bson.D{{Key: "recents", Value: bson.A{}}, {Key: "favourites", Value: bson.A{}}})
	return userResult.InsertedID.(primitive.ObjectID).Hex()

}

func CheckUserId(id string) (bool, error) {
	var dbUrl string
	dbUrl = "mongodb+srv://" + os.Getenv("MONGO_USERNAME") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster1.5jqwhvz.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	dbId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, errors.New("Invalid Id")
	}

	resp := usersCollection.FindOne(ctx, bson.M{"_id": dbId})

	var decodedResponse bson.M
	err = resp.Decode(&decodedResponse)

	if err != nil {
		return false, err
	}

	return true, nil
}
