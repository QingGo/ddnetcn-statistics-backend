package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Config 用来载入设置json文件的数据结构
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Schema   string `json:"schema"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	Webport string `json:"webport"`
}

func loadconfig(filename string) (config Config) {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&config)
	return config
}

func initalConnect(sqlConfigString string) *sql.DB {
	db, err := sql.Open("mysql",
		"ddnet:ddnet@tcp(127.0.0.1:3306)/ddnet")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func indexpage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func testjson(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(404, gin.H{
			"message": "请输入用户名",
		})
		return
	}
	fmt.Println(username)
	c.JSON(200, gin.H{
		"message": "Hello " + username,
	})
}

func querystat(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			c.JSON(404, gin.H{
				"message": "请输入用户名",
			})
			return
		}
		fmt.Println(username)
		var (
			count      int
			totalTime  float64
			totalPoint int
		)
		qurey1 := `
			SELECT
				COUNT( record_race.time ),
				sum( record_race.time ),
				sum( record_maps.Points ) 
			FROM
				record_race
				INNER JOIN record_maps ON record_race.Map = record_maps.Map 
			WHERE
				record_race.TIMESTAMP BETWEEN "2018-01-01" 
				AND "2019-01-01" 
			GROUP BY
				record_race.NAME 
			HAVING
				record_race.NAME = ?;`
		err := db.QueryRow(qurey1, username).Scan(&count, &totalTime, &totalPoint)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(200, gin.H{
			"message":    "Hello " + username,
			"username":   username,
			"count":      count,
			"totalTime":  totalTime,
			"totalPoint": totalPoint,
		})
	}
	return gin.HandlerFunc(fn)
}

func main() {

	config := loadconfig("config.json")
	sqlConfigString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Schema)
	fmt.Println(sqlConfigString)
	db := initalConnect(sqlConfigString)

	fmt.Println(config)
	fmt.Println("hello, world")
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("templates/*")
	router.GET("/", indexpage)
	router.GET("/testjson", testjson)
	router.GET("/querystat", querystat(db))
	router.Run(":" + config.Webport)

	db.Close()
}
