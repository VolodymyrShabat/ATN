package business

import (
	"github.com/VolodymyrShabat/Test_ATN/src/storage"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
)

type BookService struct {
	bookDAO storage.BookDAO
}

func NewBookService(bookDAO storage.BookDAO) *BookService {
	return &BookService{
		bookDAO: bookDAO,
	}
}

func (bs *BookService) CreateBook(book models.Book) error {
	err := bs.bookDAO.CreateBook(book)
	return err
}

func (bs *BookService) UpdateBook(book models.Book) (*models.Book, error) {
	res, err := bs.bookDAO.UpdateBook(book)
	return res, err
}

func (bs *BookService) DeleteBook(id int) error {
	err := bs.bookDAO.DeleteBook(id)
	return err
}

func (bs *BookService) GetBookById(id int) (*models.Book, error) {
	book, err := bs.bookDAO.GetBookById(id)
	return book, err
}
