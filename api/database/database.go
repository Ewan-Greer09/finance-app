package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Ewan-Greer09/finance-app/api/config"
	"github.com/Ewan-Greer09/finance-app/api/models"
)

type SQLite struct {
	DB *gorm.DB
}

type Database interface {
	AddExpense(link models.Expense) error
	AddIncome(link models.Income) error
	GetExpenses() ([]models.Expense, error)
	GetIncomes() ([]models.Income, error)
	DeleteExpense(id int) error
	DeleteIncome(id int) error

	GetUser(username string) (models.User, error)
	CreateUser(user models.User) error

	Close() error
}

func NewDatabase(cf config.Config) *SQLite {
	db, err := gorm.Open(sqlite.Open(cf.API.DatabaseName), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	err = db.AutoMigrate(&models.Expense{}, &models.Income{}, &models.User{})
	if err != nil {
		log.Panic(err)
	}

	return &SQLite{
		DB: db,
	}
}

// Adds an Expense to the database
func (d *SQLite) AddExpense(expense models.Expense) error {
	tx := d.DB.Model(models.Expense{}).Create(&expense)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Gets 10 Expenses from the database
func (d *SQLite) GetExpenses() ([]models.Expense, error) {
	var expenses []models.Expense
	tx := d.DB.Model(models.Expense{}).Order("created_at desc").Limit(10).Find(&expenses)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return expenses, nil
}

func (d *SQLite) DeleteExpense(id int) error {
	tx := d.DB.Model(models.Expense{}).Delete(&models.Expense{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Adds an Income to the database
func (d *SQLite) AddIncome(link models.Income) error {
	tx := d.DB.Model(models.Income{}).Create(&link)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Gets 10 Incomes from the database
func (d *SQLite) GetIncomes() ([]models.Income, error) {
	var incomes []models.Income
	tx := d.DB.Model(models.Income{}).Order("created_at desc").Limit(10).Find(&incomes)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return incomes, nil
}

func (d *SQLite) DeleteIncome(id int) error {
	tx := d.DB.Model(models.Income{}).Delete(&models.Income{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *SQLite) GetUser(username string) (models.User, error) {
	var user models.User
	tx := d.DB.Model(models.User{}).Where("username = ?", username).First(&user)
	if tx.Error != nil {
		return models.User{}, tx.Error
	}
	return user, nil
}

func (d *SQLite) CreateUser(user models.User) error {
	tx := d.DB.Model(models.User{}).Create(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *SQLite) Close() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
