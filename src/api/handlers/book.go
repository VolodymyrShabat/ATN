package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/VolodymyrShabat/Test_ATN/src/business"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type BookService interface {
	CreateBook(book models.Book) error
	UpdateBook(book models.Book) (*models.Book, error)
	GetBookById(id int) (*models.Book, error)
	DeleteBook(id int) error
}

type BookHandler struct {
	BookService BookService
	UserService UserAccountService
	validator   *validator.Validate
}

func NewBookHandler(BookService *business.BookService, UserAccountService *business.UserService, validator *validator.Validate) *BookHandler {
	return &BookHandler{
		BookService: BookService,
		UserService: UserAccountService,
		validator:   validator,
	}
}

func (b *BookHandler) CreateBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var book models.Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during unmarshalling request: %v", err))
			return
		}

		if validationErr := b.validator.Struct(&book); validationErr != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during validation: %v", validationErr))
			return
		}

		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			SendResponse(rw, http.StatusBadRequest, "token is empty")
			return
		}
		var authorizationToken = strings.Split(bearerToken, " ")[1]
		if authorizationToken == "" {
			SendResponse(rw, http.StatusBadRequest, "token is empty")
			return
		}
		token, err := b.UserService.VerifyToken(authorizationToken)
		if err != nil {
			SendResponse(rw, http.StatusForbidden, fmt.Sprintf("error during verifying token: %v", err))
			return
		}

		u, err := b.UserService.GetUserById(token.UserId)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during getting user: %v", err))
			return
		}

		book.Creator = u.Id
		id := uuid.New()
		book.Id = int(id.ID())
		err = b.BookService.CreateBook(book)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during creating book: %v", err))
			return
		}

		SendResponse(rw, http.StatusCreated, fmt.Sprintf("book successfully created with id: %v", int(id.ID())))
	}
}

func (b *BookHandler) UpdateBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strId, ok := vars["id"]
		if !ok {
			SendResponse(rw, http.StatusBadRequest, "id in url is missing")
			return
		}

		id, err := strconv.Atoi(strId)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during converting url id: %v", err))
			return
		}

		updatedBook, err := b.BookService.GetBookById(id)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during getting book by id: %v", err))
			return
		}
		if updatedBook == nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("book with id :%v not found", id))
			return
		}

		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			SendResponse(rw, http.StatusBadRequest, "token is empty")
			return
		}

		var authorizationToken = strings.Split(bearerToken, " ")[1]
		if authorizationToken == "" {
			SendResponse(rw, http.StatusBadRequest, "token is empty")
			return
		}

		token, err := b.UserService.VerifyToken(authorizationToken)
		if err != nil {
			SendResponse(rw, http.StatusForbidden, fmt.Sprintf("error during verifying token: %v", err))
			return
		}

		if token.UserId != updatedBook.Creator {
			SendResponse(rw, http.StatusForbidden, "user not allowed")
			return
		}

		var newBook models.Book
		err = json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during unmarshalling request: %v", err))
			return
		}

		if validationErr := b.validator.Struct(&newBook); validationErr != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during validation: %v", validationErr))
		}

		updatedBook.About = newBook.About
		updatedBook.Name = newBook.Name
		_, err = b.BookService.UpdateBook(*updatedBook)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during updating book: %v", err))
			return
		}

		SendResponse(rw, http.StatusOK, "book successfully updated")
	}
}

func (b *BookHandler) DeleteBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strId, ok := vars["id"]
		if !ok {
			SendResponse(rw, http.StatusBadRequest, "id in url is missing")
			return
		}

		id, err := strconv.Atoi(strId)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during converting url id: %v", err))
			return
		}

		deletedBook, err := b.BookService.GetBookById(id)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during getting book by id: %v", err))
			return
		}

		if deletedBook == nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("book with id :%v not found", id))
			return
		}

		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			SendResponse(rw, http.StatusBadRequest, "token is empty")
			return
		}

		var authorizationToken = strings.Split(bearerToken, " ")[1]
		if authorizationToken == "" {
			SendResponse(rw, http.StatusBadRequest, "token is empty")
			return
		}

		token, err := b.UserService.VerifyToken(authorizationToken)
		if err != nil {
			SendResponse(rw, http.StatusForbidden, fmt.Sprintf("error during verifying token: %v", err))
			return
		}

		if token.UserId != deletedBook.Creator {
			SendResponse(rw, http.StatusForbidden, "user not allowed")
			return
		}

		err = b.BookService.DeleteBook(id)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during deleting book: %v", err))
			return
		}

		SendResponse(rw, http.StatusOK, "book successfully deleted")
	}
}

func (b *BookHandler) GetByBookId() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strId, ok := vars["id"]
		if !ok {
			SendResponse(rw, http.StatusBadRequest, "id in url is missing")
			return
		}

		id, err := strconv.Atoi(strId)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during converting url id: %v", err))
			return
		}

		b, err := b.BookService.GetBookById(id)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during getting book by id: %v", err))
			return
		}

		if b == nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("book with id :%v not found", id))
			return
		}

		SendResponse(rw, http.StatusOK, fmt.Sprintf("%v", b))
	}
}
