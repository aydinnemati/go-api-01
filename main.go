package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type User struct {
	Id        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func main() {
	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// database connection
	dburi := os.Getenv("DB_USER_NAME") + ":" + os.Getenv("DB_USER_PASS") + "@tcp(" + os.Getenv("DB_IP") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")
	fmt.Println(dburi)
	db, err = sql.Open("mysql", dburi)
	if err != nil {
		panic(err)
	}

	// db init
	dbInit(os.Getenv("TBL_01"), os.Getenv("TBL_02"))

	// router init
	router := gin.Default()

	// pass db
	database := &database{db: db}

	// routes
	router.GET("/users", database.getUsers)
	router.POST("/users", database.addUser)

	// run server
	router.Run(os.Getenv("SERVER_IP") + ":" + os.Getenv("SERVER_PORT"))

}

func dbInit(table01 string, table02 string) {
	sql01 := "CREATE TABLE IF NOT EXISTS " + table01 + " (id INT(6) PRIMARY KEY, firstname VARCHAR(30) NOT NULL)"
	res, err := db.Exec(sql01)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	sql02 := "CREATE TABLE IF NOT EXISTS " + table02 + " (id INT(6) PRIMARY KEY, lastname VARCHAR(30) NOT NULL)"
	res02, err := db.Exec(sql02)
	if err != nil {
		panic(err)
	}
	fmt.Println(res02)
}

type database struct {
	db *sql.DB
}

func (database *database) getUsers(c *gin.Context) {
	getUsername(database.db)
	getUserlastname(database.db)
	// ##################################################################### parse and return per user
}

func getUsername(db *sql.DB) []User {
	res, err := db.Query("SELECT * FROM usernames")
	if err != nil {
		log.Fatal(err)
	}
	usersnames := []User{}
	for res.Next() {

		var user User
		err := res.Scan(&user.Id, &user.Firstname)

		if err != nil {
			log.Fatal(err)
		}
		usersnames = append(usersnames, user)
		// fmt.Printf("%v\n", user)
	}
	fmt.Println(usersnames)
	return usersnames
}

func getUserlastname(db *sql.DB) []User {
	res, err := db.Query("SELECT * FROM userslastname")
	if err != nil {
		log.Fatal(err)
	}
	usersnames := []User{}
	for res.Next() {

		var user User
		err := res.Scan(&user.Id, &user.Firstname)

		if err != nil {
			log.Fatal(err)
		}
		usersnames = append(usersnames, user)
		// fmt.Printf("%v\n", user)
	}
	return usersnames
}

func (database *database) addUser(c *gin.Context) {

	var user User
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Printf("error %s", err)
		c.JSON(501, gin.H{"error": err})
	}
	// fmt.Println(user.Firstname)
	var wg01 sync.WaitGroup
	chanel01 := make(chan bool)
	if checkIfexists(db, user.Id, user.Firstname) {

		go addUsername(database.db, user.Id, user.Firstname, wg01, chanel01)
		go addUserLastname(database.db, user.Id, user.Lastname, wg01)
		wg01.Wait()

	} else {
		c.Error(errors.New("user is existed"))
	}
	c.JSON(200, gin.H{"user": user})
}

func addUsername(db *sql.DB, userid int64, username string, wg sync.WaitGroup, chanel chan bool) {
	wg.Add(1)
	sqlcommand := fmt.Sprintf("INSERT INTO usernames (id, firstname) VALUES (%d, '%v' )", userid, username)
	_, err := db.Query(sqlcommand)
	if err != nil {
		fmt.Println(err)
	}
	wg.Done()
	// fmt.Println(res)
}

func addUserLastname(db *sql.DB, userid int64, userlastname string, wg sync.WaitGroup) {
	wg.Add(1)
	sqlcommand := fmt.Sprintf("INSERT INTO userslastname (id, lastname) VALUES (%d, '%v' )", userid, userlastname)
	_, err := db.Query(sqlcommand)
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()
	// fmt.Println(res)

}

func checkIfexists(db *sql.DB, userid int64, username string) bool {
	sqlcommand := fmt.Sprintf("SELECT * FROM usernames WHERE firstname='%v';", username)
	var id int64
	var name string
	row := db.QueryRow(sqlcommand)
	row.Scan(&id, &name)
	switch {
	case id != userid && name != username:
		return true
	default:
		return false
	}
}
