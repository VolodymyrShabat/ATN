package storage

import (
	"context"
	"fmt"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"time"
)

type UserDAO struct{}

func init() {
	_, err := GetConnection()
	if err != nil {
		log.Fatal(fmt.Errorf("error during getting database connection: %v", err))
	}
}

func (ud *UserDAO) CreateUser(user models.User) error {
	_, err := UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("error during inserting user: %v", err)
	}
	return nil
}

func (ud *UserDAO) GetUserById(id int) (*models.User, error) {
	var result primitive.M
	err := UserCollection.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&result)
	u := models.User{}
	if len(result) == 0 {
		return nil, nil
	}

	u.Id = ConvertInt(result["id"])
	u.Name = result["name"].(string)
	u.City = result["city"].(string)
	u.Login = result["login"].(string)
	u.Email = result["email"].(string)
	u.Age = ConvertInt(result["age"])

	return &u, err
}

func (ud *UserDAO) GetUserByLogin(login string) (*models.User, error) {
	var result primitive.M
	err := UserCollection.FindOne(context.TODO(), bson.D{{"login", login}}).Decode(&result)
	u := models.User{}
	if len(result) == 0 {
		return nil, nil
	}

	u.Id = ConvertInt(result["id"])
	u.Name = result["name"].(string)
	u.City = result["city"].(string)
	u.Login = result["login"].(string)
	u.Email = result["email"].(string)
	u.Password = result["password"].(string)
	u.Age = ConvertInt(result["age"])

	return &u, err
}

func (ud *UserDAO) GetUserByEmail(email string) (*models.User, error) {
	var result primitive.M
	err := UserCollection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&result)
	u := models.User{}
	if len(result) == 0 {
		return nil, nil
	}

	u.Id = ConvertInt(result["id"])
	u.Name = result["name"].(string)
	u.City = result["city"].(string)
	u.Login = result["login"].(string)
	u.Email = result["email"].(string)
	u.Password = result["password"].(string)
	u.Age = ConvertInt(result["age"])

	return &u, err
}

func (ud *UserDAO) UpdatePasswordResetToken(email, passwordResetToken string) error {
	query := bson.D{{Key: "email", Value: strings.ToLower(email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	result, err := UserCollection.UpdateOne(context.TODO(), query, update)
	if result.MatchedCount == 0 {
		return fmt.Errorf("no matches found for user email")
	}
	return err
}

func (ud *UserDAO) UpdatePasswordByResetToken(resetToken, hashedPassword string) error {
	query := bson.D{{Key: "passwordResetToken", Value: resetToken}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashedPassword}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}
	result, err := UserCollection.UpdateOne(context.TODO(), query, update)
	if result.MatchedCount == 0 {
		return fmt.Errorf("no matches found for this token")
	}
	return err
}
