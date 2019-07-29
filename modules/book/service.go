package book

import (
	"github.com/go-kit/kit/log"
	"math/rand"
	"nocai/gokit-demo/modules/client/auth"
)

var books = []Book{
	{Name: "Java", Author: "Java-Author"},
	{Name: "Golang", Author: "Golang-Author"},
}

type Service interface {
	Books() ([]Book, error)
	GetById(id int64) (*Book, error)
}

type service struct {
	l          log.Logger
	authClient auth.Service
}

func NewService(l log.Logger, authClient auth.Service) Service {
	return &service{l: l, authClient: authClient}
}

func (ser service) Books() ([]Book, error) {
	ser.l.Log("method", "Books")
	//return books, nil
	//panic("bbbbbbbb")
	//return nil, fmt.Errorf("aaaaaaaaaa")
	//return nil, returncodes.Fail("fail")
	//return books, returncodes.ErrBook

	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	wg.Done()
	//	panic("panic")
	//}()
	//wg.Wait()
	return books, nil
}

func (ser service) GetById(id int64) (*Book, error) {
	ser.l.Log("method", "GetById")
	return &books[rand.Intn(len(books))], nil
}
