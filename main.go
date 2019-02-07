package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Config 用来载入设置json文件的数据结构
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	Webport string `json:"webport"`
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

func indexpage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
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

func main() {
	config := loadconfig("config.json")
	fmt.Println(config)
	fmt.Println("hello, world")
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", indexpage)
	router.GET("/testjson", testjson)
	router.Run(":" + config.Webport)
}
