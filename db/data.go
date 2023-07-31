package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbSetup() (context.Context, *mongo.Client, context.CancelFunc, error) {
	var dbUrl string
	dbUrl = "mongodb+srv://" + os.Getenv("MONGO_USERNAME") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster1.5jqwhvz.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		return nil, nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	return ctx, client, cancel, nil
}

func AddRecentSearchForUser(id string, location string) error {
	ctx, client, cancel, err := DbSetup()

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	dbId, err := primitive.ObjectIDFromHex(id)

	// Id is already checked

	// Remove location from Recents if present
	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": dbId},
		bson.D{
			{Key: "$pull", Value: bson.D{{Key: "recents", Value: location}}},
		})
	if err != nil {
		return err
	}

	// Add location to end of Recents
	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": dbId},
		bson.D{
			{Key: "$push", Value: bson.D{{Key: "recents", Value: location}}},
		})
	if err != nil {
		return err
	}

	return nil
}

func GetRecentsForUser(id string) ([]string, error) {
	ctx, client, cancel, err := DbSetup()

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return []string{}, err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	dbId, err := primitive.ObjectIDFromHex(id)

	resp := usersCollection.FindOne(ctx, bson.M{"_id": dbId})

	var decodedResponse bson.M
	err = resp.Decode(&decodedResponse)
	if err != nil {
		return []string{}, err
	}

	locations := []string{}
	for i := len(decodedResponse["recents"].(bson.A)) - 1; i >= 0; i-- {
		v := decodedResponse["recents"].(bson.A)[i]
		locations = append(locations, v.(string))
	}

	return locations, nil
}

func GetFavouritesForUser(id string) ([]string, error) {
	ctx, client, cancel, err := DbSetup()

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return []string{}, err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	dbId, err := primitive.ObjectIDFromHex(id)

	resp := usersCollection.FindOne(ctx, bson.M{"_id": dbId})

	var decodedResponse bson.M
	err = resp.Decode(&decodedResponse)
	if err != nil {
		return []string{}, err
	}

	locations := []string{}
	for i := len(decodedResponse["favourites"].(bson.A)) - 1; i >= 0; i-- {
		v := decodedResponse["favourites"].(bson.A)[i]
		locations = append(locations, v.(string))
	}

	return locations, nil
}

func HandleFavouriteForUser(id string, location string) (string, error) {
	ctx, client, cancel, err := DbSetup()

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	dbId, err := primitive.ObjectIDFromHex(id)

	resp, err := usersCollection.UpdateOne(ctx, bson.M{"_id": dbId},
		bson.D{
			{Key: "$pull", Value: bson.D{{Key: "favourites", Value: location}}},
		})
	if err != nil {
		return "", err
	}

	if resp.ModifiedCount == 1 {
		return "Removed from favourites", nil
	}

	resp, err = usersCollection.UpdateOne(ctx, bson.M{"_id": dbId},
		bson.D{
			{Key: "$push", Value: bson.D{{Key: "favourites", Value: location}}},
		})

	if err != nil {
		return "", err
	}

	return "Added to favourites", nil
}

func IsFavourite(id string, location string) (bool, error) {
	ctx, client, cancel, err := DbSetup()

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return false, err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	dbId, err := primitive.ObjectIDFromHex(id)

	resp := usersCollection.FindOne(ctx, bson.M{"_id": dbId, "favourites": bson.M{"$elemMatch": bson.M{"$eq": location}}})

	var decodedResponse bson.M
	err = resp.Decode(&decodedResponse)

	if err != nil {
		return false, err
	}

	return true, nil
}
