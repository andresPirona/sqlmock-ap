package repository

import "github.com/andresPirona/sqlmock-ap/domain/entity"

type PersonRepository interface {
	Save(object entity.Person) error
	GetByID(id uint) (*entity.Person, error)
}
