package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andresPirona/sqlmock-ap/domain/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func DbMockInit(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqldb,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func TestFindPerson_Found(t *testing.T) {
	sqlDB, db, mock := DbMockInit(t)
	defer sqlDB.Close()

	// alice
	person := entity.Person{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "Alice",
	}

	personImplement := NewPersonImplementation(db)
	persons := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(person.ID, person.Name)

	expectedSQL := "SELECT (.+) FROM `people`"
	mock.ExpectQuery(expectedSQL).WillReturnRows(persons)
	_, errF := personImplement.GetByID(person.ID)

	assert.Nil(t, errF)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindUser_NotFound(t *testing.T) {
	sqlDB, db, mock := DbMockInit(t)
	defer sqlDB.Close()

	personImplement := NewPersonImplementation(db)
	person := sqlmock.NewRows([]string{"id", "name"})

	expectedSQL := "SELECT (.+) FROM `people`"
	mock.ExpectQuery(expectedSQL).WillReturnRows(person)
	_, res := personImplement.GetByID(1)
	assert.True(t, errors.Is(res, gorm.ErrRecordNotFound))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreatePerson_shouldSuccess(t *testing.T) {
	// Inicializa la base de datos simulada y el mock
	sqlDB, db, mock := DbMockInit(t)
	defer sqlDB.Close()

	// Crea una instancia de la implementación de la persona
	personImplement := NewPersonImplementation(db)

	// Define un objeto persona para insertar
	person := entity.Person{
		Name: "Alice",
	}

	// Define la expectativa para la inserción
	expectedSQL := "INSERT INTO `people`"

	mock.ExpectBegin()

	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// Llama a la función para guardar a la persona
	personImplement.Save(person)

	// Verifica que se cumplan todas las expectativas
	assert.Nil(t, mock.ExpectationsWereMet())
}
