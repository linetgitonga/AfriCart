package models

type Product struct {
    ID          string   `json:"id" bson:"_id,omitempty"`
    CustomID    int      `json:"custom_id" bson:"custom_id"`
    Name        string   `json:"name" bson:"name"`
    Brand       string   `json:"brand" bson:"brand"`
    Category    string   `json:"category" bson:"category"`
    Price       float64  `json:"price" bson:"price"`
    Sizes       []string `json:"sizes" bson:"sizes"`
    Description string   `json:"description" bson:"description"`
    Reviews     []Review `json:"reviews" bson:"reviews"`
    Images      []string `json:"images" bson:"images"`
}

type Review struct {
    UserEmail string `json:"user_email" bson:"user_email"`
    Rating    int    `json:"rating" bson:"rating"`
    Comment   string `json:"comment" bson:"comment"`
}
