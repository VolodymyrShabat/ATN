package models

import "fmt"

type Book struct {
	Id      int    `json:"id"`
	Name    string `json:"name,omitempty" validate:"required"`
	About   string `json:"about,omitempty" validate:"required"`
	Creator int    `json:"creator"`
}

func (b *Book) String() string {
	return fmt.Sprintf("Book id: %d\n name: %s\n about: %s\n creator_id %d\n", b.Id, b.Name, b.About, b.Creator)
}
