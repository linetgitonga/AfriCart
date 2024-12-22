package database

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// GetNextSequenceValue generates the next custom_id value for a given sequence
func GetNextSequenceValue(sequenceName string) (int, error) {
    collection := Client.Database("store").Collection("counters")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"_id": sequenceName}
    update := bson.M{"$inc": bson.M{"sequence_value": 1}}
    opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

    var result struct {
        SequenceValue int `bson:"sequence_value"`
    }

    err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
    if err == mongo.ErrNoDocuments {
        // Initialize the counter if it doesn't exist
        _, err := collection.InsertOne(ctx, bson.M{
            "_id":            sequenceName,
            "sequence_value": 1,
        })
        if err != nil {
            return 0, err
        }
        return 1, nil
    } else if err != nil {
        return 0, err
    }

    return result.SequenceValue, nil
}
