package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dutroctu/go-webservice/database"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("loading products...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("%d products loaded...", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}
	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}
	return prodMap, nil
}

func getProduct(productID int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products
	WHERE productId = ?`, productID)

	var product Product
	err := row.Scan(&product.ProductID, &product.Manufacturer,
		&product.Sku, &product.Upc, &product.PricePerUnit,
		&product.QuantityOnHand, &product.ProductName)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &product, nil
}

func removeProduct(productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.QueryContext(ctx, `DELETE FROM products where productId = ?`, productID)
	if err != nil {
		return err
	}
	return nil
}

func getProductList() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := database.DbConn.QueryContext(ctx, `SELECT productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products`)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	products := make([]Product, 0)
	for result.Next() {
		var product Product
		result.Scan(&product.ProductID, &product.Manufacturer,
			&product.Sku, &product.Upc, &product.PricePerUnit,
			&product.QuantityOnHand, &product.ProductName)
		products = append(products, product)
	}
	return products, nil
}

func insertProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO products
	(manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName) VALUES (?, ?, ?, ?, ?, ?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName)
	if err != nil {
		return 0, nil
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(insertID), nil
}

func updateProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `UPDATE products SET
	manufacturer=?,
	sku=?,
	upc=?,
	pricePerUnit=CAST(? AS DECIMAL(13,2)),
	quantityOnHand=?,
	productName=?
	WHERE productId=?`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductID)
	if err != nil {
		return err
	}
	return nil
}

func GetTopTenProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(ctx, `SELECT productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products ORDER BY quantityOnHand DESC LIMIT 10`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID, &product.Manufacturer,
			&product.Sku, &product.Upc, &product.PricePerUnit,
			&product.QuantityOnHand, &product.ProductName)
		products = append(products, product)
	}
	return products, nil
}
