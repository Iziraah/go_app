package links

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

const collection = "links"

func New(db *mongo.Database, timeout time.Duration) *Repository {
	return &Repository{db: db, timeout: timeout}
}

type Repository struct {
	db      *mongo.Database
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateReq) (database.Link, error) {
    var l database.Link

    ctx, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()

    collection := r.db.Collection(collection)
    _, err := collection.InsertOne(ctx, req)
    if err != nil {
        return l, err
    }

    return req, nil
}

func (r *Repository) FindByUserAndURL(ctx context.Context, link, userID string) (database.Link, error) {
    var l database.Link

    ctx, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()

    collection := r.db.Collection(collection)
    err := collection.FindOne(ctx, bson.M{"link": link, "userID": userID}).Decode(&l)
    if err != nil {
        return l, err
    }

    return l, nil
}

func (r *Repository) FindByCriteria(ctx context.Context, criteria Criteria) ([]database.Link, error) {
    var links []database.Link

    ctx, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()

    collection := r.db.Collection(collection)
    cursor, err := collection.Find(ctx, bson.M{"field": value})
    if err != nil {
        return links, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var l database.Link
        if err := cursor.Decode(&l); err != nil {
            return links, err
        }
        links = append(links, l)
    }

    return links, nil
}