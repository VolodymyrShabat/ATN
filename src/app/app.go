package app

import (
	"github.com/VolodymyrShabat/Test_ATN/src/business"
	"github.com/VolodymyrShabat/Test_ATN/src/storage"
)

type Services struct {
	User *business.UserService
	Book *business.BookService
}

func CreateServices(JwtSecretKey, Salt []byte) *Services {
	services := Services{
		User: business.NewUserService(storage.UserDAO{}, JwtSecretKey, Salt),
		Book: business.NewBookService(storage.BookDAO{}),
	}

	return &services
}
