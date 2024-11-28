package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Db struct {
	*sql.DB
}

type Employee struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Dept   string `json:"dept"`
	Salary int    `json:"salary"`
}

func (db *Db) getemployees(context *gin.Context) {

	query := `select * from employees;`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("error occured", err)
	}
	var employees []Employee
	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.Id, &employee.Name, &employee.Dept, &employee.Salary); err != nil {
			fmt.Println("error occured while fetching values from databse")
			context.JSON(http.StatusNotFound, gin.H{"message": err})
		}
		employees = append(employees, employee)
		// fmt.Println(employees)
	}
	context.JSON(http.StatusOK, map[string][]Employee{
		"data": employees,
	})
}
func (db *Db) getemployeesByname(c *gin.Context) {
	var response Employee
	queryParams := c.Params
	name, _ := queryParams.Get("name")
	// convertingId, _ := strconv.Atoi(id)
	query := fmt.Sprintf(`select * from employees where name ='%s'`, name)
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error while fetching data", err)
	}
	for rows.Next() {
		if err := rows.Scan(&response.Id, &response.Name, &response.Dept, &response.Salary); err != nil {
			fmt.Println("Error while scanning rows")
			c.JSON(http.StatusNotFound, gin.H{"message": err})
		}
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"record": response,
	})

}

func (db *Db) addEmployees(c *gin.Context) {
	var employee Employee
	if err := c.BindJSON(&employee); err != nil {
		fmt.Println("Error while binding employee", err)
	}
	query := `INSERT into employees(id,name,dept,salary) VALUES($1,$2,$3,$4)`
	_, errC := db.Query(query, employee.Id, employee.Name, employee.Dept, employee.Salary)
	if errC != nil {
		fmt.Println("Error while inserting data into tables", errC)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errC})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully inserted"})
}

func ConnectDatabase() (*Db, error) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("unable to get the directory")
	}

	err = godotenv.Load(dir + "/.env")
	if err != nil {
		fmt.Println("error while loading env vars")
	}

	// Extracting connection details from environment variables
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")

	// Creating a PostgreSQL connection string
	dbo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	fmt.Println("dbo: ", dbo)
	// Opening the database connection
	db, err := sql.Open("postgres", dbo)
	if err != nil {
		return nil, err
	}

	return &Db{db}, nil
}
func main() {
	db, err := ConnectDatabase()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/employees", db.getemployees)
	router.GET("/employees/:name", db.getemployeesByname)
	router.POST("/employees", db.addEmployees)
	router.Run("localhost:9090")

}
