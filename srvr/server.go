package srvr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Book (struct) : inventory structure of books in the shop
type Book struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Author string `json:"author,omitempty"`
	Count  int    `json:"count,omitempty"`
	//Author *Author `json:"address,omitempty"`
}

// User (struct) :  info structure of the users
type User struct {
	ID       int
	Name     string
	username string
	password string
}

var server http.Server
var router = mux.NewRouter()
var books []Book
var ids map[int]int
var count = 2

var users []User

var logger string

var v bool
var bypassLogin bool

func init() {

	books = append(books, Book{ID: "1", Name: "Pride and Prejudice", Author: "Jane Austen", Count: 5})
	books = append(books, Book{ID: "2", Name: "Things fall apart,", Author: "Chinua Achebe", Count: 9})
	ids = make(map[int]int)
	ids[1] = 1
	ids[2] = 1
	users = append(users, User{ID: 1, Name: "Tom", username: "tom95", password: "pass1"})
	users = append(users, User{ID: 2, Name: "Harry", username: "harry88", password: "pass2"})
	router.HandleFunc("/books", GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", GetBook).Methods("GET")
	router.HandleFunc("/books", CreateBook).Methods("POST")
	router.HandleFunc("/books/{id}", UpdateBook).Methods("UPDATE")
	router.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")

	logger = "Start:\n"

}

func isAuthorised(r *http.Request) bool {
	authorised := false
	//Authorization in Header has base and user's credentials encrypted in it
	rcvdEncrptdAuthArr := strings.Split(r.Header.Get("Authorization"), " ")
	if len(rcvdEncrptdAuthArr) != 2 {
		//log.Fatal("request body doesnt have proper authorization format")
		return false
	}

	byteCredStr, err := base64.StdEncoding.DecodeString(rcvdEncrptdAuthArr[1])
	if err != nil {
		//log.Fatal("couldn't decode credentials to string")
		return false
	}
	credential := string(byteCredStr)
	//log.Println("credential = ", credential)
	for _, user := range users {
		tempCred := strings.Join([]string{user.username, user.password}, ":")
		if tempCred == credential {
			logger = logger + "Authorised. \n"
			authorised = true
			break
		}

	}
	return authorised
}

// GetBooks : Display all books from the books variable
func GetBooks(w http.ResponseWriter, r *http.Request) {
	if !isAuthorised(r) && !bypassLogin {
		logger = logger + "Not Authorized. \n"
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(books)

}

// GetBook : get a single book by id
func GetBook(w http.ResponseWriter, r *http.Request) {

	if !isAuthorised(r) && !bypassLogin {
		logger = logger + "Not Authorized."
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	//If ID field is empty, return bad request
	_, err := strconv.Atoi(vars["id"])

	//id is not convertable, therefore alphabets exist
	if err != nil {
		logger = logger + "Bad ID Found \n"
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Find a match by id
	for _, item := range books {
		if item.ID == vars["id"] {
			//add it to existing list
			var b []Book
			b = append(b, item)
			json.NewEncoder(w).Encode(b)
			return
		}
	}
	//if no match is found,
	logger = logger + "No Match Found \n"
	w.WriteHeader(http.StatusNoContent)
}

// CreateBook : create a new book entry
func CreateBook(w http.ResponseWriter, r *http.Request) {

	if !isAuthorised(r) && !bypassLogin {
		logger = logger + "Not Authorized. "
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//add a new book entry in the specified index
	var book Book
	// decode json to get the struct equivalent
	//and save that to our book variable
	err := json.NewDecoder(r.Body).Decode(&book)

	//If json can not be decoded to struct, return bad request
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		logger = logger + "Json could not be decoded \n"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//If ID field is invalid, return bad request
	newKeyInt, err := strconv.Atoi(book.ID)
	//id is not convertable, therefore alphabets exist
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ids[newKeyInt] == 1 {
		//Duplicate Found
		w.WriteHeader(http.StatusConflict)
		return
	}

	//if no duplicate exists

	ids[newKeyInt] = 1

	//fmt.Println("Added ID = ", vars["id"])
	//add the new entry to our existing book entries
	books = append(books, book)
	json.NewEncoder(w).Encode(books)
}

// UpdateBook : create a new book entry
func UpdateBook(w http.ResponseWriter, r *http.Request) {

	if !isAuthorised(r) && !bypassLogin {
		logger = logger + "Not Authorized. \n"
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// extract parameters from URL
	vars := mux.Vars(r)
	newKey := vars["id"]
	newKeyInt, err := strconv.Atoi(newKey)

	//If ID field is invalid, return bad request
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//If requested key DOESNT exist, it will have the value 0
	if ids[newKeyInt] != 1 {
		//ID not found, Send badrequest status code
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var book Book
	// get the struct equivalent of the json
	//and save that to our book variable
	err = json.NewDecoder(r.Body).Decode(&book)
	//If json can not be decoded to struct, return bad request
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Find the index of the book and
	//update book entry in the correspondin index
	//Iterate over books to find the book by id
	for index, item := range books {
		if item.ID == vars["id"] {
			books[index].ID = book.ID
			if len(book.ID) == 0 {
				books[index].ID = vars["id"]
			}
			books[index].Name = book.Name
			books[index].Author = book.Author
			books[index].Count = book.Count
			break
		}
	}

	var b []Book
	b = append(b, book)
	json.NewEncoder(w).Encode(b)

}

// DeleteBook : Delete a book
func DeleteBook(w http.ResponseWriter, r *http.Request) {

	if !isAuthorised(r) && !bypassLogin {
		logger = logger + "Not Authorized. \n"
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// extract parameters from URL
	vars := mux.Vars(r)
	newKeyInt, err := strconv.Atoi(vars["id"])
	//If ID field is invalid, return bad request
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//initially, deleted = flase
	deleted := false
	//Iterate over books to find the book by id
	for index, item := range books {
		if item.ID == vars["id"] {
			//Delete the book with matching ID
			books = append(books[:index], books[index+1:]...)
			deleted = true
			break
		}
	}

	if !deleted {
		w.WriteHeader(http.StatusGone)
		return
	}
	//remove from ids array

	ids[newKeyInt] = 0

	//return the books
	json.NewEncoder(w).Encode(books)
}

//PostMain Former main function
func PostMain(port string, verbose bool, noLogin bool) {

	// f := flag.Int("f", 1234, "help message for flagname")
	// n := flag.String("name", "John Doe", "Help mesage for NAME")
	// //(name, shorthand,value,usage)
	//v = flag.BoolP("vFlag", "v", false, "help message")

	// fmt.Println("ip has value ", *f)
	// fmt.Println("name has value ", *n)
	//fmt.Println("V has value ", *v)
	bypassLogin = noLogin
	v = verbose
	if v {
		log.Println(v, "Running in verbose mode")
	} else {
		log.Println("Not running in verbose mode")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	server = http.Server{Addr: ":" + port, Handler: router}

	go startServer()
	<-stop
	stopServer()
	if v {
		logger = logger + " :End"
		fmt.Println(logger)
	}

}
func startServer() {
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
func stopServer() {
	log.Println("Shutting down the server within 5 seconds...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if cancel != nil {
		//log.Println("Server is not running.")
	}
	server.Shutdown(ctx)
	log.Printf("Server Closed")

}
