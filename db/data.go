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

func AddRecentSearchForUser(id string, location string) (dbId string, err error) {
	ctx, client, cancel, err := DbSetup()
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		// Create a new user and add location to recent and return nil
		newId, err := CreateNewUserForLocation(location)
		return newId, err
	}

	filter := bson.M{"_id": objId}
	update := bson.D{
		{Key: "$pull", Value: bson.D{{Key: "recents", Value: location}}},
	}
	resp, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	if resp.MatchedCount < 1 {
		newId, err := CreateNewUserForLocation(location)
		return newId, err
	}

	update = bson.D{
		{Key: "$push", Value: bson.D{{Key: "recents", Value: location}}},
	}
	resp, err = usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return id, nil
}

func CreateNewUserForLocation(location string) (string, error) {
	ctx, client, cancel, err := DbSetup()
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")

	userResult, err := usersCollection.InsertOne(ctx, bson.D{{Key: "recents", Value: bson.A{location}}, {Key: "favourites", Value: bson.A{}}})
	if err != nil {
		return "", err
	}
	return userResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func GetRecentsAndFavouritesForUser(id string) ([]string, []string, error) {
	ctx, client, cancel, err := DbSetup()
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return []string{}, []string{}, err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	dbId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return []string{}, []string{}, nil
	}

	resp := usersCollection.FindOne(ctx, bson.M{"_id": dbId})

	var decodedResponse bson.M
	err = resp.Decode(&decodedResponse)
	if err != nil {
		return []string{}, []string{}, err
	}

	recentLocations := []string{}
	for i := len(decodedResponse["recents"].(bson.A)) - 1; i >= 0; i-- {
		v := decodedResponse["recents"].(bson.A)[i]
		recentLocations = append(recentLocations, v.(string))
	}

	favouriteLocations := []string{}
	for i := len(decodedResponse["favourites"].(bson.A)) - 1; i >= 0; i-- {
		v := decodedResponse["favourites"].(bson.A)[i]
		favouriteLocations = append(favouriteLocations, v.(string))
	}

	return recentLocations, favouriteLocations, nil
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
	if err != nil {
		return []string{}, nil
	}

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
		return false, nil
	}

	return true, nil
}
