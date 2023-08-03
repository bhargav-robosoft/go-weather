package db

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbSetup() (context.Context, *mongo.Client, context.CancelFunc, error) {
	dbUrl := "mongodb+srv://" + os.Getenv("MONGO_USERNAME") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster1.5jqwhvz.mongodb.net/?retryWrites=true&w=majority"
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
		newId, err := CreateNewUserForLocation(location, true)
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
		newId, err := CreateNewUserForLocation(location, true)
		return newId, err
	}

	update = bson.D{
		{Key: "$push", Value: bson.D{{Key: "recents", Value: location}}},
	}
	_, err = usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return id, nil
}

func CreateNewUserForLocation(location string, isRecent bool) (dbId string, err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return "", err
	}
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")

	var userResult *mongo.InsertOneResult
	if isRecent {
		userResult, err = usersCollection.InsertOne(ctx, bson.D{{Key: "recents", Value: bson.A{location}}, {Key: "favourites", Value: bson.A{}}})
		if err != nil {
			return "", err
		}
	} else {
		userResult, err = usersCollection.InsertOne(ctx, bson.D{{Key: "recents", Value: bson.A{}}, {Key: "favourites", Value: bson.A{location}}})
		if err != nil {
			return "", err
		}
	}

	return userResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func GetRecentsAndFavouritesForUser(id string) (recentLocations []string, favouriteLocations []string, err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return []string{}, []string{}, err
	}
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

	recentLocations = []string{}
	for i := len(decodedResponse["recents"].(bson.A)) - 1; i >= 0; i-- {
		v := decodedResponse["recents"].(bson.A)[i]
		recentLocations = append(recentLocations, v.(string))
	}

	favouriteLocations = []string{}
	for i := len(decodedResponse["favourites"].(bson.A)) - 1; i >= 0; i-- {
		v := decodedResponse["favourites"].(bson.A)[i]
		favouriteLocations = append(favouriteLocations, v.(string))
	}

	return recentLocations, favouriteLocations, nil
}

func GetFavouritesForUser(id string) (favouriteLocations []string, err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return []string{}, err
	}
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

func HandleFavouriteForUser(id string, location string) (dbId string, response string, err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return "", "", err
	}
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return "", "", err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		// Create a new user and add location to recent and return nil
		newId, err := CreateNewUserForLocation(location, false)
		return newId, "Added to favourites", err
	}

	resp, err := usersCollection.UpdateOne(ctx, bson.M{"_id": objId},
		bson.D{
			{Key: "$pull", Value: bson.D{{Key: "favourites", Value: location}}},
		})
	if err != nil {
		return "", "", err
	}

	if resp.MatchedCount < 1 {
		newId, err := CreateNewUserForLocation(location, false)
		return newId, "Added to favourites", err
	}

	if resp.ModifiedCount == 1 {
		return id, "Removed from favourites", nil
	}

	// Modified count 0
	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": objId},
		bson.D{
			{Key: "$push", Value: bson.D{{Key: "favourites", Value: location}}},
		})

	if err != nil {
		return "", "", err
	}

	return id, "Added to favourites", nil
}

func ClearRecents(id string) (err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return err
	}
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	dbId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}

	resp, err := usersCollection.UpdateOne(ctx, bson.M{"_id": dbId},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "recents", Value: bson.A{}}}},
		})
	if err != nil {
		return err
	}

	if resp.MatchedCount == 0 {
		return errors.New("Id not found")
	}

	return nil
}

func ClearFavourites(id string) (err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return err
	}
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	goWeatherDatabase := client.Database("go-weather")
	usersCollection := goWeatherDatabase.Collection("users")
	dbId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}

	resp, err := usersCollection.UpdateOne(ctx, bson.M{"_id": dbId},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "favourites", Value: bson.A{}}}},
		})
	if err != nil {
		return err
	}

	if resp.MatchedCount == 0 {
		return errors.New("Id not found")
	}

	return nil
}

func IsFavourite(id string, location string) (isFavourite bool, err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return false, err
	}
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
