package storage

import (
	"context"
	"fmt"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type BookDAO struct{}

func init() {
	_, err := GetConnection()
	if err != nil {
		log.Fatal(fmt.Errorf("error during getting database connection: %v", err))
	}
}

func (ud *BookDAO) CreateBook(book models.Book) error {
	_, err := BookCollection.InsertOne(context.TODO(), book)
	if err != nil {
		return fmt.Errorf("error during inserting user: %v", err)
	}
	return nil
}

func (ud *BookDAO) GetBookById(id int) (*models.Book, error) {
	var result primitive.M
	err := BookCollection.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&result)
	b := models.Book{}
	if len(result) == 0 {
		return nil, nil
	}

	b.Id = ConvertInt(result["id"])
	b.Name = result["name"].(string)
	b.About = result["about"].(string)
	b.Creator = ConvertInt(result["creator"])

	return &b, err
}

func (ud *BookDAO) UpdateBook(book models.Book) (*models.Book, error) {
	filter := bson.D{{"id", book.Id}}
	after := options.After
	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	update := bson.D{{"$set", bson.D{{"about", book.About}, {"name", book.Name}}}}
	updateResult := BookCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	err := updateResult.Decode(&result)

	b := models.Book{}
	b.Id = ConvertInt(result["id"])
	b.Name = result["name"].(string)
	b.About = result["about"].(string)
	b.Creator = ConvertInt(result["creator"])

	return &b, err
}

func (ud *BookDAO) DeleteBook(id int) error {
	opts := options.Delete().SetCollation(&options.Collation{})
	_, err := BookCollection.DeleteOne(context.TODO(), bson.D{{"id", id}}, opts)
	return err
}

func ConvertInt(val interface{}) int {
	var i int

	switch t := val.(type) {
	case int:
		i = t
	case int8:
		i = int(t)
	case int16:
		i = int(t)
	case int32:
		i = int(t)
	case int64:
		i = int(t)
	}
	return i
}
