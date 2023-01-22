package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"ashwin.com/go-auth/database"
	"ashwin.com/go-auth/helper"
	"ashwin.com/go-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(loginPassword string, databasePassword string) (bool, string) {
	check := true
	msg := ""
	err := bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(loginPassword))

	if err != nil {
		check = false
		msg = "invalid password"
	}

	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser *models.Users
		c.BindJSON(&newUser)

		validationErr := validator.New().Struct(newUser)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		documentCount, err := userCollection.CountDocuments(ctx, bson.M{"email": newUser.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		}

		if documentCount > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
			return
		}

		newUser.ID = primitive.NewObjectID()
		password := HashPassword(newUser.Password)
		newUser.Password = password

		InsertOne, err := userCollection.InsertOne(ctx, newUser)
		defer cancel()
		if err != nil {
			log.Panic(err)
		}
		c.JSON(http.StatusOK, InsertOne)
	}
}

func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginCredentials models.Users
		var foundCredentials models.Users
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

		err := c.BindJSON(&loginCredentials)
		if err != nil {
			log.Panic(err)

		}

		err = userCollection.FindOne(ctx, bson.M{"email": loginCredentials.Email}).Decode(&foundCredentials)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(loginCredentials.Password, foundCredentials.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, err := helper.GenerateToken(foundCredentials.Name, foundCredentials.Email)
		if err != nil {
			log.Panic(err)
		}

		c.SetCookie("jwt", token, 5*60, "/", "localhost", false, true)
		c.SetCookie("name", foundCredentials.Name, 3600, "/", "localhost", false, true)
		c.JSON(http.StatusOK, "Login Successfull")
	}
}

func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := c.Request.Cookie("name")
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, "Welcome "+username.Value)

	}
}
