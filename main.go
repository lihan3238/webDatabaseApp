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
	//db, err := sql.Open("mysql", "root:123456@tcp(192.168.50.128:3306)/movie")
	//本地mysql
	db, err := sql.Open("mysql", "root:qcwh2018@tcp(localhost:3306)/movie")
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

	router.POST("/task", func(c *gin.Context) {
		// 获取前端发送的数据
		userId := c.PostForm("userId")
		keyword := c.PostForm("keyword")
		year := c.PostForm("year")
		selectedTask := c.PostForm("task")

		// 执行查询，获取结果
		results, err := executeQuery(db, userId, keyword, year, selectedTask)
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

type ResultB struct {
	Movie  string `json:"movie"`
	Rating string `json:"rating"`
}

type ResultC struct {
	Movie string `json:"movie"`
}

type ResultD struct {
	Movie string `json:"movie"`
}

type ResultE struct {
	Gender    string `json:"gender"`
	RatingNum string `json:"rating_num"`
}

// 修改 executeQuery 函数返回结果
// 修改 executeQuery 函数
func executeQuery(db *sql.DB, searchUserID string, searchKeyword string, searchYear string, task string) (interface{}, error) {
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

		//log.Println("Query results:", results)

		return results, nil
	case "task_b":
		query := `
			SELECT m.title, AVG(r.rating) AS avg_rating
			FROM movies m
			JOIN ratings r ON m.movieId = r.movieId
			WHERE m.title LIKE ?
			GROUP BY m.title
			ORDER BY avg_rating DESC;
		`
		args = []interface{}{"%" + searchYear + "%"}
		rows, err := db.Query(query, args...)
		if err != nil {
			log.Println("Error executing query:", err)
			return nil, err
		}
		defer rows.Close()

		var results []ResultB

		for rows.Next() {
			var result ResultB

			err := rows.Scan(&result.Movie, &result.Rating)
			if err != nil {
				return nil, err
			}

			results = append(results, result)
		}

		log.Println("Query results:", results)

		return results, nil

		// Add logic for task_b query
	case "task_c":
		query := `
		SELECT m.title
		FROM movies m
		JOIN (
			SELECT r.movieId, AVG(r.rating) AS avg_rating
			FROM ratings r
			WHERE r.movieId IN (
				SELECT m.movieId
				FROM movies m
				WHERE m.genres LIKE ?
			)
			GROUP BY r.movieId
			ORDER BY avg_rating DESC
			LIMIT 20
		) subquery ON m.movieId = subquery.movieId
		ORDER BY subquery.avg_rating DESC;
	`
		args = []interface{}{"%" + searchKeyword + "%"}

		log.Println(args)

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Println("Error executing query:", err)
			return nil, err
		}
		defer rows.Close()

		var results []ResultC

		for rows.Next() {
			var result ResultC

			err := rows.Scan(&result.Movie)
			if err != nil {
				return nil, err
			}

			results = append(results, result)
		}

		//log.Println("Query results:", results)

		return results, nil
		// Add logic for task_c query
	case "task_d":
		query := `
		SELECT m.title
		FROM movies m
		JOIN (
			SELECT r.movieId
			FROM ratings r
			JOIN users u ON r.userId = u.userId
			WHERE u.gender = ?
			GROUP BY r.movieId
			ORDER BY AVG(r.rating) DESC
			LIMIT 20
		) subquery ON m.movieId = subquery.movieId
		ORDER BY (SELECT AVG(rating) FROM ratings WHERE movieId = m.movieId) DESC;
		
	`
		args = []interface{}{searchKeyword}

		log.Println(args)

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Println("Error executing query:", err)
			return nil, err
		}
		defer rows.Close()

		var results []ResultD

		for rows.Next() {
			var result ResultD

			err := rows.Scan(&result.Movie)
			if err != nil {
				return nil, err
			}

			results = append(results, result)
		}

		//log.Println("Query results:", results)

		return results, nil
		// Add logic for task_d query
	case "task_e":
		query := `
		SELECT
		u.gender,
		COUNT(DISTINCT u.userId) AS user_count
	FROM
		users u
	JOIN
		ratings r ON u.userId = r.userId
	WHERE
		r.rating > ?
	GROUP BY
		u.gender;
	
	`
		args = []interface{}{searchKeyword}

		log.Println(args)

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Println("Error executing query:", err)
			return nil, err
		}
		defer rows.Close()

		var results []ResultE

		for rows.Next() {
			var result ResultE

			err := rows.Scan(&result.Gender, &result.RatingNum)
			if err != nil {
				return nil, err
			}

			results = append(results, result)
		}

		//log.Println("Query results:", results)

		return results, nil
		// Add logic for task_e query
	default:
		return nil, fmt.Errorf("Invalid task specified" + searchUserID)
	}
}

func ginHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
