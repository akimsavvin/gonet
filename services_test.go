package gonet

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

type IRepo interface {
	Insert(id string)
}

type Repo struct {
	id string
}

func NewRepo() *Repo {
	return &Repo{
		id: uuid.NewString(),
	}
}

func (r *Repo) Insert(id string) {
	fmt.Printf("Inserted %s\n", id)
}

type IService interface {
	GetID() string
}

type Service struct {
	id   string
	repo IRepo
}

func NewService(repo IRepo) *Service {
	return &Service{
		id:   uuid.NewString(),
		repo: repo,
	}
}

func (s *Service) GetID() string {
	s.repo.Insert(s.id)
	return s.id
}

func TestSingleton(t *testing.T) {
	AddService[IRepo](NewRepo, LifetimeSingleton)
	AddService[IService](NewService, LifetimeSingleton)

	id1 := GetService[IService]().GetID()
	id2 := GetService[IService]().GetID()

	if id1 != id2 {
		t.Errorf("Got different instances for singleton service: %s and %s", id1, id2)
	}
}

func TestTransient(t *testing.T) {
	AddService[IRepo](NewRepo, LifetimeTransient)
	AddService[IService](NewService, LifetimeTransient)
	id1 := GetService[IService]().GetID()
	id2 := GetService[IService]().GetID()

	if id1 == id2 {
		t.Errorf("Got same instances for transient service: %s", id1)
	}
}
