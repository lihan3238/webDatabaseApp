// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 设置数据库连接信息
	//虚拟机mysqlcluster
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.50.128:3306)/movie")
	//本地mysql
	//db, err := sql.Open("mysql", "root:qc@tcp(localhost:3306)/movie")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")

	// 设置静态文件目录
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "static") // 静态文件路径

	// 设置路由规则
	router.GET("/", ginHtml)

	router.POST("/taska", func(c *gin.Context) {
		// 获取前端发送的数据
		userId := c.PostForm("userId")
		keyword := c.PostForm("keyword")
		year := c.PostForm("year")
		tag := c.PostForm("tag")
		selectedTask := c.PostForm("task")

		// 执行查询，获取结果
		results, err := executeQuery(db, userId, keyword, year, tag, selectedTask)
		if err != nil {
			log.Println("Error executing query:", err)
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}

		// 将查询结果以合适的格式返回给前端
		c.JSON(http.StatusOK, results)
	})

	router.Run(":8080")
}

type ResultA struct {
	Movie  string   `json:"movie"`
	Rating string   `json:"rating"`
	Tag    []string `json:"tag"`
}

// 修改 executeQuery 函数返回结果
// 修改 executeQuery 函数
func executeQuery(db *sql.DB, searchUserID string, searchKeyword string, searchYear string, searchTag string, task string) ([]ResultA, error) {
	var args []interface{}

	switch task {

	case "task_a":
		query := `
    SELECT
        m.title,
        r.rating,
        SUBSTRING_INDEX(GROUP_CONCAT(gt.tag ORDER BY g.relevance DESC), ',', 3) AS top_tags
    FROM
        ratings r
    JOIN
        movies m ON r.movieId = m.movieId
    LEFT JOIN
        genomescores g ON r.movieId = g.movieId
	LEFT JOIN
		genometags gt ON g.tagId = gt.tagId
    WHERE
        r.userId = ?
    GROUP BY
        r.movieId, r.rating, r.timestamp
    ORDER BY
        r.timestamp DESC
`
		searchUserID, err := strconv.Atoi(searchUserID)
		if err != nil {
			return nil, err
		}
		args = []interface{}{searchUserID}
		rows, err := db.Query(query, args...)
		if err != nil {
			log.Println("Error executing query:", err)
			return nil, err
		}
		defer rows.Close()

		var results []ResultA

		for rows.Next() {
			var result ResultA
			var topTags string

			err := rows.Scan(&result.Movie, &result.Rating, &topTags)
			if err != nil {
				return nil, err
			}

			// 将逗号分隔的标签字符串拆分为切片
			result.Tag = strings.Split(topTags, ",")

			results = append(results, result)
		}

		log.Println("Query results:", results)

		return results, nil
	case "task_b":
		query := `

`
		searchUserID, err := strconv.Atoi(searchUserID)
		if err != nil {
			return nil, err
		}
		args = []interface{}{searchUserID}
		rows, err := db.Query(query, args...)
		if err != nil {
			log.Println("Error executing query:", err)
			return nil, err
		}
		defer rows.Close()

		var results []ResultA

		for rows.Next() {
			var result ResultA
			var topTags string

			err := rows.Scan(&result.Movie, &result.Rating, &topTags)
			if err != nil {
				return nil, err
			}

			// 将逗号分隔的标签字符串拆分为切片
			result.Tag = strings.Split(topTags, ",")

			results = append(results, result)
		}

		log.Println("Query results:", results)

		return results, nil
		// Add logic for task_b query
	case "task_c":
		// Add logic for task_c query
	case "task_d":
		// Add logic for task_d query
	case "task_e":
		// Add logic for task_e query
	default:
		return nil, fmt.Errorf("Invalid task specified" + searchUserID)
	}
	return nil, nil
}

func ginHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
