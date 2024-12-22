package models

import "time"
import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Subject   string             `bson:"subject" json:"subject"`
    Message   string             `bson:"message" json:"message"`
    Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
