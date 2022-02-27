package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-playground/validator"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Login    string             `json:"login,omitempty" bson:"login,omitempty" validate:"min=3,max=50"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,email"`
	Password string             `json:"password,omitempty" bson:"password,omitempty" validate:"min=6,max=50"`
}

type Result struct {
	ResultMessage string `json:"resultmessage,omitempty" bson:"resultmessage,omitempty"`
}

var validate *validator.Validate
var user User
var ctx, cancel = context.WithTimeout(context.Background(), 25*time.Second)
var client, erring = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
var collection = client.Database("theregistrdeveloper").Collection("user")

func JSONregistration(c *gin.Context) {
	c.Writer.Header().Add("content-type", "application/json")
	json.NewDecoder(c.Request.Body).Decode(&user)
	ctx, _ := context.WithTimeout(context.Background(), 25*time.Second)
	validate = validator.New()
	if err := validate.Struct(user); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(c.Writer).Encode(Result{ResultMessage: err.Error()})
	} else {
		collection.InsertOne(ctx, user)
		json.NewEncoder(c.Writer).Encode(http.StatusOK)
	}
}

func JSONLogin(c *gin.Context) {
	c.Writer.Header().Add("content-type", "application/json")
	json.NewDecoder(c.Request.Body).Decode(&user)
	ctx, _ := context.WithTimeout(context.Background(), 25*time.Second)
	err := collection.FindOne(ctx, user).Decode(&user)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(c.Writer).Encode(Result{ResultMessage: err.Error()})
		return
	} else {
		json.NewEncoder(c.Writer).Encode(http.StatusOK)
	}
}

func main() {
	defer cancel()
	if erring != nil {
		return
	}
	r := gin.Default()
	r.POST("/registration", JSONregistration)
	r.GET("/login", JSONLogin)
	r.Run()
}
