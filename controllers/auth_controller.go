package controllers

import (
	"context"
	"fmt"
	"github.com/dastardlyjockey/ecommerce/database"
	tokens "github.com/dastardlyjockey/ecommerce/helpers"
	"github.com/dastardlyjockey/ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var UserCollection = database.Collection(database.Client, "user")

var validate = validator.New()

func hashedPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func comparePassword(user, verifiedUser string) error {
	err := bcrypt.CompareHashAndPassword([]byte(verifiedUser), []byte(user))
	if err != nil {
		return err
	}
	return nil
}

func SignUp(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to bind JSON: %v", err)})
		return
	}

	err = validate.Struct(user)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid user data: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access emails in the documents"})
		return
	}

	if count > 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "The email is already registered"})
		return
	}

	count, err = UserCollection.CountDocuments(ctx, bson.M{"phone_number": user.PhoneNumber})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access phone numbers in the documents"})
		return
	}

	if count > 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "The phone number is already registered"})
		return
	}

	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()

	hashedPassword, err := hashedPassword(*user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to encrypt password"})
		log.Println(err)
		return
	}

	user.Password = &hashedPassword

	token, refreshToken, err := tokens.GenerateToken(*user.FirstName, *user.LastName, *user.Email, *user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate the token"})
		log.Println(err)
		return
	}

	user.Token = &token
	user.RefreshToken = &refreshToken

	user.AddressDetails = make([]models.Address, 0)
	user.OrderStatus = make([]models.Order, 0)

	user.UserCart = make([]models.Product, 0)

	user.CreatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed in registering the time"})
		log.Println(err)
		return
	}

	_, err = UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add the user to the database"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, "registered successfully")
}

func SignIn(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access the user"})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var verifiedUser models.User
	err = UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&verifiedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Email not registered"})
		log.Println(err.Error())
		return
	}

	err = comparePassword(*user.Password, *verifiedUser.Password)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid password"})
		log.Println(err.Error())
		return
	}

	token, refreshToken, err := tokens.GenerateToken(*verifiedUser.FirstName, *verifiedUser.LastName, *verifiedUser.Email, *verifiedUser.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate the token"})
		log.Println(err)
		return
	}

	err = tokens.UpdateToken(ctx, token, refreshToken, verifiedUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the token"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, verifiedUser)

}
