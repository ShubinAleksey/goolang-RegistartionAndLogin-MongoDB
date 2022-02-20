package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/validator.v2"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Login    string             `json:"login,omitempty" bson:"login,omitempty" validate:"min=3,max=50"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty" validate:"regexp=^([A-Za-z]|[0-9])+@"`
	Password string             `json:"password,omitempty" bson:"password,omitempty" validate:"min=6,max=50"`
}

var client *mongo.Client

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	r := gin.Default()
	r.POST("/registration", func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		var r *http.Request = c.Request
		w.Header().Add("content-type", "application/json")
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		collection := client.Database("theregistrdeveloper").Collection("user")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := validator.Validate(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Ошибка": err.Error(),
			})
		} else {
			collection.InsertOne(ctx, user)
			json.NewEncoder(w).Encode(http.StatusOK)
		}
	})
	r.GET("/allregistered", func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		var r *http.Request = c.Request
		w.Header().Add("content-type", "application/json")
		var users []User
		json.NewDecoder(r.Body).Decode(&users)
		collection := client.Database("theregistrdeveloper").Collection("user")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "Ошибка": "` + err.Error() + `"}`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var user User
			cursor.Decode(&user)
			users = append(users, user)
		}
		if err := cursor.Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "Ошибка": "` + err.Error() + `"}`))
			return
		}
		json.NewEncoder(w).Encode(users)
	})
	r.GET("/login", func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		var r *http.Request = c.Request
		w.Header().Add("content-type", "application/json")
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		collection := client.Database("theregistrdeveloper").Collection("user")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := collection.FindOne(ctx, user).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "Ошибка": "` + err.Error() + `"}`))
			return
		} else {
			json.NewEncoder(w).Encode(http.StatusOK)
		}
	})
	r.Run()
}
