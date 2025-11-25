package database

import (
	"database/sql"
	"fmt"
	"testing"

	"cleanarch/internal/entity"

	"github.com/stretchr/testify/suite"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	fmt.Println("Setting up suite")
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	fmt.Println("tearing down test")
	suite.Db.Exec("DELETE FROM orders")
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	fmt.Println("tearing down suite")
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestOrderRepository_ListAll() {
	stmt, err := suite.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	suite.NoError(err)
	defer stmt.Close()
	ordersPrepared := []entity.Order{
		{ID: "1", Price: 10.0, Tax: 2.0, FinalPrice: 12.0},
		{ID: "2", Price: 20.0, Tax: 4.0, FinalPrice: 24.0},
		{ID: "3", Price: 30.0, Tax: 6.0, FinalPrice: 36.0},
	}
	for _, order := range ordersPrepared {
		_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
		suite.NoError(err)
	}

	repo := NewOrderRepository(suite.Db)
	orders, err := repo.FindAll()
	suite.NoError(err)
	suite.Len(orders, len(ordersPrepared))
	suite.Equal(ordersPrepared, orders)
}
