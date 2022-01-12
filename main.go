package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type User struct {
	Id   string `json:"id"`
	Name string `json:"firstname"`
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

	// router init
	router := gin.Default()

	// routes
	router.GET("/users", getUsers)

	// run server
	router.Run(os.Getenv("SERVER_IP") + ":" + os.Getenv("SERVER_PORT"))

}

// func getUsers(c *gin.Context) {
// 	// sql := "SELECT * FROM userstable"
// 	sql := "select * from usersdb;"
// 	res, err := db.Exec(sql)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.JSON(http.StatusOK, res)
// }

func getUsers(c *gin.Context) {
	res, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	users := []User{}
	for res.Next() {

		var user User
		err := res.Scan(&user.Id, &user.Name)

		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
		// fmt.Printf("%v\n", user)

	}
	c.JSON(http.StatusOK, users)
}

// CREATE TABLE IF NOT EXISTS users(
// 	id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
// 	firstname VARCHAR(30) NOT NULL)

// INSERT INTO table_name (id, firstname )VALUES ("aydin");
