package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
type aroduct struct {
	ID    int             `form:"id" json:"id"`
	Code  string          `form:"code" json:"code"`
	Name  string          `form:"name" json:"name"`
	Price decimal.Decimal `form:"price" json:"price" sql:"type:decimal(16,2);"`
}

// Result is an array of product
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Main
func main() {
	db, err = gorm.Open("mysql", "root:root@/recordings?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&aroduct{})
	handleRequests()

}

func handleRequests() {
	fmt.Println("Start the development server at http://127.0.0.1:9999")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := Result{Code: 404, Message: "Method not found"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := Result{Code: 403, Message: "Method not allowed"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/wkwk/ea", getCertainNamesWithoutResult).Methods("get")
	myRouter.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")
	myRouter.HandleFunc("/api/products/NoResultStruct", createProductWithoutResultStruct).Methods("POST")
	myRouter.HandleFunc("/api/products/range_price/{minPrice}/{maxPrice}", getRangePrice).Methods("get")

	log.Fatal(http.ListenAndServe(":9999", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var product aroduct
	json.Unmarshal(payloads, &product)

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Success create product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func createProductWithoutResultStruct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)
	var product aroduct
	json.Unmarshal(payloads, &product)

	db.Create(&product)

	res, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func getCertainNamesWithoutResult(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: getCertainNamesWithoutResult")
	productq := []aroduct{}
	// db.Find(&productq)
	// var result Result
	db.Raw("SELECT * FROM aroducts limit 5").Scan(&productq)
	// db.Exec("select * from aroducts")
	res, err := json.Marshal(productq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get products")

	productz := []aroduct{}
	db.Find(&productz)

	res := Result{Code: 200, Data: productz, Message: "Success get products"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getRangePrice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get range price")
	vars := mux.Vars(r)
	minPrice := vars["minPrice"]
	maxPrice := vars["maxPrice"]

	var productz []aroduct
	db.Raw("select distinct * from aroducts where price between " + minPrice + " and " + maxPrice).Scan(&productz)
	res, err := json.Marshal(productz)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product aroduct

	// db.First(&product, productID)
	db.Raw("select * from aroducts where id = " + productID).Scan(&product)

	res := Result{Code: 200, Data: product, Message: "Success get product"}
	result, err := json.Marshal(res)

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

	payloads, _ := ioutil.ReadAll(r.Body)

	var productUpdates aroduct
	json.Unmarshal(payloads, &productUpdates)

	var product aroduct
	db.First(&product, productID)
	db.Model(&product).Updates(productUpdates)

	res := Result{Code: 200, Data: product, Message: "Success update product"}
	result, err := json.Marshal(res)

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

	var product aroduct

	db.First(&product, productID)
	db.Delete(&product)

	res := Result{Code: 200, Message: "Success delete product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
