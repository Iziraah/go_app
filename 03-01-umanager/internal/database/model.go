package database

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Link struct {
	ID        primitive.ObjectID `bson:"id"`
	Title     string             `bson:"title,omitempty"`
	URL       string             `bson:"url"`
	Images    []string           `bson:"images"`
	Tags      []string           `bson:"tags"`
	UserID    string             `bson:"userID"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
