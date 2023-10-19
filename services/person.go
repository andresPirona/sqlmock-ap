package services

import (
	"github.com/andresPirona/sqlmock-ap/domain/entity"
	"github.com/andresPirona/sqlmock-ap/domain/repository"
	"gorm.io/gorm"
)

type implementationPerson struct {
	db *gorm.DB
}

func (i implementationPerson) Save(object entity.Person) error {
	result := i.db.Debug().Create(&object)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (i implementationPerson) GetByID(id uint) (*entity.Person, error) {
	var user *entity.Person
	err := i.db.Debug().Where("id = ?", id).First(&user).Error
	return user, err
}

func NewPersonImplementation(db *gorm.DB) repository.PersonRepository {
	return &implementationPerson{
		db: db,
	}
}
