package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"context"
	"log"
	"time"
)

type Post struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}

type PostRepository struct {
	Client         *mongo.Client
	DBName         string
	CollectionName string
	Collection     *mongo.Collection
}

func NewPostRepository(client *mongo.Client, dbName, collectionName string) *PostRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &PostRepository{
		Client:         client,
		DBName:         dbName,
		CollectionName: collectionName,
		Collection:     collection,
	}

}

func (pr *PostRepository) EnsureCollectionExists() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tempPost := Post{Title: "Temp Title", Content: "Temp Content"}
	_, err := pr.Collection.InsertOne(ctx, tempPost)
	if err != nil {
		log.Println("Error ensuring collection exists:", err)
		return err
	}

	_, err = pr.Collection.DeleteOne(ctx, bson.M{"title": "Temp Title"})
	if err != nil {
		log.Println("Error deleting temporary post:", err)
		return err
	}

	return nil
}

func (pr *PostRepository) CreatePost(post Post) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := pr.Collection.InsertOne(ctx, post)
	if err != nil {
		log.Println("Error creating post:", err)
		return nil, err
	}
	return result, nil
}

func (pr *PostRepository) GetPost(id string) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ID format:", err)
		return nil, err
	}

	var post Post
	err = pr.Collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No document found with that ID")
			return nil, nil
		}
		log.Println("Error fetching post:", err)
		return nil, err
	}
	return &post, nil
}

func (pr *PostRepository) GetAllPost() ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := pr.Collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error finding posts:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []Post
	for cursor.Next(ctx) {
		var post Post
		err := cursor.Decode(&post)
		if err != nil {
			log.Println("Error decoding post:", err)
			continue
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return posts, nil
}

func (pr *PostRepository) DeletePost(id string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ID format:", err)
		return nil, err
	}

	result, err := pr.Collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Println("Error deleting post:", err)
		return nil, err
	}
	return result, nil
}
