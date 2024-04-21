package helpers

import (
	"context"
	"github.com/dastardlyjockey/ecommerce/controllers"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"time"
)

var secretKey = os.Getenv("SECRET_KEY")

type signedDetails struct {
	firstName   string
	lastName    string
	phoneNumber string
	email       string
	jwt.RegisteredClaims
}

func GenerateToken(firstName, lastName, email, phoneNumber string) (token, refreshToken string, err error) {
	claims := &signedDetails{
		firstName:   firstName,
		lastName:    lastName,
		email:       email,
		phoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(30) * time.Minute)),
		},
	}

	refreshClaims := &signedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(24) * time.Hour)),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		log.Println("Error signing the token")
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS384, refreshClaims).SignedString(secretKey)
	if err != nil {
		log.Println("Error signing the refresh token")
	}

	return token, refreshToken, err
}

func UpdateToken(ctx context.Context, signedToken, signedRefreshToken, userID string) error {
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	updatedAt, err := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}

	updateObj = append(updateObj, bson.E{"updated_at", updatedAt})

	filter := bson.M{"user_id": userID}

	_, err = controllers.UserCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}})
	if err != nil {
		return err
	}

	return nil
}
