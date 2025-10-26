package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"ordersystem/model"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseHandler struct {
	dbConn *gorm.DB
}

func NewDatabaseHandler() (*DatabaseHandler, error) {
	slog.Info("Connecting to database")
	// connect to db
	dsn, err := getDsn()
	if err != nil {
		return nil, err
	}
	dbConn, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// create tables and migrate
	err = dbConn.AutoMigrate(&model.Drink{}, &model.Order{})
	if err != nil {
		return nil, err
	}
	// add test data to db
	err = prepopulate(dbConn)
	if err != nil {
		return nil, err
	}
	return &DatabaseHandler{dbConn: dbConn}, nil
}

func getDsn() (string, error) {
	dbUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return "", errors.New("environment variable 'POSTGRES_USER' is not set")
	}
	dbPw, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return "", errors.New("environment variable 'POSTGRES_PASSWORD' is not set")
	}
	dbName, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return "", errors.New("environment variable 'POSTGRES_DB' is not set")
	}
	dbPort, ok := os.LookupEnv("POSTGRES_TCP_PORT")
	if !ok {
		return "", errors.New("environment variable 'POSTGRES_TCP_PORT' is not set")
	}
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return "", errors.New("environment variable 'DB_HOST' is not set")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbHost, dbUser, dbPw, dbName, dbPort)
	return dsn, nil
}

func prepopulate(dbConn *gorm.DB) error {
	// check if drinks table already has data
	var count int64
	err := dbConn.Model(&model.Drink{}).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		// already populated
		return nil
	}

	drinks := []model.Drink{
		{Name: "Sprite Zero", Price: 3.25, Description: "The best thing ever created."},
		{Name: "Water", Price: 0.0, Description: "The second best thing ever created."},
	}
	if err := dbConn.Create(&drinks).Error; err != nil {
		return err
	}

	orders := []model.Order{
		{DrinkID: 1, Amount: 2}, // 2x Sprite Zero
		{DrinkID: 2, Amount: 1}, // 1x Water
	}
	if err := dbConn.Create(&orders).Error; err != nil {
		return err
	}

	return nil
}

func (db *DatabaseHandler) GetDrinks() (drinks []model.Drink, err error) {
	err = db.dbConn.Find(&drinks).Error
	if err != nil {
		return nil, err
	}
	return drinks, nil
}

func (db *DatabaseHandler) GetOrders() (orders []model.Order, err error) {
	err = db.dbConn.Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

const totalledStmt = `SELECT drink_id, SUM(amount) AS total_amount_ordered FROM orders WHERE deleted_at IS NULL GROUP BY drink_id ORDER BY drink_id;`

func (db *DatabaseHandler) GetTotalledOrders() (totals []model.DrinkOrderTotal, err error) {
	err = db.dbConn.Raw(totalledStmt).Scan(&totals).Error
	if err != nil {
		return nil, err
	}
	return totals, nil
}

func (db *DatabaseHandler) AddOrder(order *model.Order) error {
	err := db.dbConn.Create(order).Error
	if err != nil {
		return err
	}
	return nil
}
