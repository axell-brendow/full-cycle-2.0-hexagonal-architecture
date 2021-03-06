package db_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/full-cycle-2.0-hexagonal-architecture/adapters/db"
	"github.com/full-cycle-2.0-hexagonal-architecture/application"
	"github.com/stretchr/testify/require"
)

var Db *sql.DB

func setUp() {
	var err error

	Db, err = sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal(err.Error())
	}

	db.CreateTable(Db)
	createProduct(Db)
}

func createProduct(db *sql.DB) {
	query := `insert into products values("abc", "Product Test", 0, "disabled");`

	stmt, err := db.Prepare(query)

	if err != nil {
		log.Fatal(err.Error())
	}

	stmt.Exec()
}

func TestProductDb_Get(t *testing.T) {
	setUp()
	defer Db.Close()

	productDb := db.NewProductDb(Db)

	product, err := productDb.Get("abc")

	require.Nil(t, err)
	require.Equal(t, "Product Test", product.GetName())
	require.Equal(t, 0.0, product.GetPrice())
	require.Equal(t, "disabled", product.GetStatus())
}

func TestProductDb_Save(t *testing.T) {
	setUp()
	defer Db.Close()
	productDb := db.NewProductDb(Db)

	product := application.NewProduct()
	product.Name = "Product Test"
	product.Price = 25

	productResult, err := productDb.Save(product)
	require.Nil(t, err)
	require.Equal(t, product.Name, productResult.GetName())
	require.Equal(t, product.Price, productResult.GetPrice())
	require.Equal(t, product.Status, productResult.GetStatus())

	product.Status = "enabled"

	productResult, err = productDb.Save(product)
	require.Nil(t, err)
	require.Equal(t, product.Name, productResult.GetName())
	require.Equal(t, product.Price, productResult.GetPrice())
	require.Equal(t, product.Status, productResult.GetStatus())
}
