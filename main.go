package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func checkErr(err error) {
	if err != nil {
		panic(err) // throws exception on error
	}
}

type productType struct {
	Id              string `json:"id"`
	Img             string `json:"img"`
	Title           string `json:"title"`
	DiscountedPrice int    `json:"discountedPrice"`
	OriginalPrice   int    `json:"originalPrice"`
	Brand           string `json:"brand"`
}

func getItems(c *gin.Context) {
	db := Connect()
	rows, err := db.Query(`SELECT id, image, title, discountedprice, originalprice, brand FROM "Product" ORDER By id`)
	checkErr(err)
	var products []productType
	for rows.Next() {
		var p productType
		err = rows.Scan(&p.Id, &p.Img, &p.Title, &p.DiscountedPrice, &p.OriginalPrice, &p.Brand)
		checkErr(err)
		fmt.Println(p)
		products = append(products, p)
	}
	c.IndentedJSON(http.StatusOK, products)
}

type productDetailType struct {
	Id              int               `json:"id"`
	Sku_id          string            `json:"skuid"`
	Brand           string            `json:"brand"`
	Title           string            `json:"title"`
	Reviews         int               `json:"reviews"`
	Stars           float32           `json:"stars"`
	DiscountedPrice float32           `json:"discountedPrice"`
	OriginalPrice   float32           `json:"originalPrice"`
	Grades          []gradeDetailType `json:"grade"`
	Details         string            `json:"details"`
	Warranty        string            `json:"warranty"`
	Returns         string            `json:"returns"`
}

type gradeDetailType struct {
	Name    string          `json:"name"`
	BagSize []bagDetailType `json:"bagSize"`
}

type bagDetailType struct {
	Size        int    `json:"size"`
	Stock       int    `json:"stock"`
	MinOrder    int    `json:"minOrderQty"`
	MaxOrder    int    `json:"maxOrderQty"`
	Dimension   string `json:"dimension"`
	DateOfAvail string `json:"dateOfAvail"`
	Origin      string `json:"countryOfOrigin"`
}

func getProduct(c *gin.Context) {
	db := Connect()
	id := c.Param("id")
	rows0, err0 := db.Query(`select * from "Product" where id = $1`, id)
	checkErr(err0)
	var productDetails productDetailType
	for rows0.Next() {
		var image string
		err := rows0.Scan(&productDetails.Id, &productDetails.Sku_id, &productDetails.Brand, &productDetails.Title,
			&productDetails.Reviews, &productDetails.DiscountedPrice, &productDetails.OriginalPrice, &productDetails.Details,
			&productDetails.Warranty, &productDetails.Returns, &productDetails.Stars, &image)
		checkErr(err)
	}
	rows1, err := db.Query(`select distinct b."id", b."name" from "Stock" a, "Grade" b, "Product" d where d."id" = $1 and a."productId" = $1 and a."gradeId" = b."id"  order by b."id";`, id)
	checkErr(err)
	for rows1.Next() {
		var id int
		var name string
		err := rows1.Scan(&id, &name)
		productDetails.Grades = append(productDetails.Grades, gradeDetailType{name, []bagDetailType{}})
		checkErr(err)
	}
	rows2, err := db.Query(`select b."name", c."size", a."stock", a."minOrder", a."maxOrder",  c."dimension", a."dateOfAvail", a."origin"  from "Stock" a, "Grade" b, "BagSize" c, "Product" d where d."id" = $1 and a."productId" = $1 and a."gradeId" = b."id" and a."bagSizeId" = c."id";`, id)
	checkErr(err)
	for rows2.Next() {
		var bag bagDetailType
		var gradeName string
		err := rows2.Scan(&gradeName, &bag.Size, &bag.Stock, &bag.MinOrder,
			&bag.MaxOrder, &bag.Dimension, &bag.DateOfAvail, &bag.Origin)
		for i := 0; i < len(productDetails.Grades); i++ {
			if productDetails.Grades[i].Name == gradeName {
				productDetails.Grades[i].BagSize = append(productDetails.Grades[i].BagSize, bag)
				break
			}
		}
		checkErr(err)
	}
	c.IndentedJSON(http.StatusOK, productDetails)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/product", getItems)
	router.GET("/product/:id", getProduct)
	fmt.Printf("%v", router.Run("localhost:8080"))
}

const (
	DB_USER     = "rachitpratapsingh"
	DB_PASSWORD = "postgres"
	DB_NAME     = "oms"
)

func Connect() *sql.DB {
	// dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", "postgres://lkfjzdte:NKG4gaEOzSoYj1jG_1pXUXGBqvB9jtji@tiny.db.elephantsql.com/lkfjzdte")
	// db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	return db
}
