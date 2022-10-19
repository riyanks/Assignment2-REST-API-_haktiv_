package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

// Product is a representation of a product
type Item struct {
	Item_id     int             `form:"itemId" json:"itemId"`
	Item_code   string          `form:"itemCode" json:"itemCode"`
	Description string          `form:"description" json:"description"`
	Quantity    decimal.Decimal `form:"quantity" json:"quantity" sql:"type:decimal(16,2);"`
	Order_id    string          `from:"orderId" json:"orderId"`
}

// Result is an array of product
type Ordered struct {
	Ordered_at    string      `json:"orderedAt"`
	Customer_name string      `json:"customerName"`
	Items         interface{} `json:"items"`
}

// Main
func main() {
	db, err = gorm.Open("mysql", "root:@/assignment2?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&Item{})
	handleRequests()
}

func handleRequests() {
	log.Println("Start the development server at http://127.0.0.1:8888")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := Ordered{Ordered_at: "none", Customer_name: " not found"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := Ordered{Ordered_at: "none", Customer_name: " not allowed"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/item", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/item", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/item/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/item/{id}", updateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/item/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8888", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := io.ReadAll(r.Body)

	var barang Ordered
	json.Unmarshal(payloads, &barang)

	db.Create(&barang)

	result, err := json.Marshal(&barang)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get products")

	barang := []Ordered{}
	db.Find(&barang)

	results, err := json.Marshal(&barang)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Ordered

	db.First(&product, productID)

	// res := Result{Code: 200, Data: product, Message: "Success get product"}
	result, err := json.Marshal(&product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	payloads, _ := io.ReadAll(r.Body)

	var productUpdates Ordered
	json.Unmarshal(payloads, &productUpdates)

	var product Ordered
	db.First(&product, productID)
	db.Model(&product).Updates(productUpdates)

	// res := Result{Code: 200, Data: product, Message: "Success update product"}
	result, err := json.Marshal(&product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Ordered

	db.First(&product, productID)
	db.Delete(&product)

	result, err := json.Marshal(&product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
