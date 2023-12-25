// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Movie 表示电影实体结构
type Movie struct {
	Title    string    `json:"title"`
	Rating   float64   `json:"rating"`
	Tags     []string  `json:"tags"`
	DateTime time.Time `json:"dateTime"`
}

// UserMovieInfo 表示用户电影信息实体结构
type UserMovieInfo struct {
	Movie   Movie    `json:"movie"`
	TopTags []string `json:"topTags"`
}

func main() {
	// 设置数据库连接信息
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
	router.GET("/", func(c *gin.Context) {
		ginHtml(c)
	})

	router.POST("/", func(c *gin.Context) {
		userID := c.PostForm("userId")
		keyword := c.PostForm("keyword")
		year := c.PostForm("year")
		tag := c.PostForm("tag")
		task := c.PostForm("task")

		var result interface{}

		switch task {
		case "task_a":
			result = searchTaskA(db, userID, keyword, year, tag)
		case "task_b":
			result = searchTaskB(db, year)
		case "task_c":
			result = searchTaskC(db, tag)
		case "task_d":
			result = searchTaskD(db, userID)
		case "task_e":
			minRating, err := strconv.ParseFloat(keyword, 64)
			if err != nil {
				log.Println("Error converting keyword to float64:", err)
				// 处理转换错误的情况，可以返回错误信息给前端
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating format"})
				return
			}
			result = searchTaskE(db, userID, minRating)
		}

		c.JSON(http.StatusOK, gin.H{"result": result})
	})

	router.Run(":8080")
}

func searchTaskA(db *sql.DB, userID, keyword, year, tag string) []UserMovieInfo {
	rows, err := db.Query(`
		SELECT m.title, r.rating, t.tag, FROM_UNIXTIME(r.timestamp) as dateTime
		FROM ratings r
		JOIN movies m ON r.movieId = m.movieId
		LEFT JOIN tags t ON r.userId = t.userId AND r.movieId = t.movieId
		WHERE r.userId = ?
		ORDER BY r.timestamp DESC
	`, userID)

	if err != nil {
		log.Println("Error executing query:", err)
		return nil
	}
	defer rows.Close()

	var userMovieInfos []UserMovieInfo

	for rows.Next() {
		var title, tag string
		var rating float64
		var dateTime time.Time

		err := rows.Scan(&title, &rating, &tag, &dateTime)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		userMovieInfo := UserMovieInfo{
			Movie: Movie{
				Title:    title,
				Rating:   rating,
				DateTime: dateTime,
			},
			TopTags: []string{tag},
		}

		userMovieInfos = append(userMovieInfos, userMovieInfo)
	}

	return userMovieInfos
}
func searchTaskB(db *sql.DB, year string) []Movie {
	// 查询不同年代的电影并按受欢迎程度排序的示例代码
	rows, err := db.Query(`
		SELECT m.title, AVG(r.rating) as avgRating
		FROM movies m
		JOIN ratings r ON m.movieId = r.movieId
		WHERE m.year = ?
		GROUP BY m.title
		ORDER BY avgRating DESC
		LIMIT 20
	`, year)

	if err != nil {
		log.Println("Error executing query:", err)
		return nil
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var title string
		var avgRating float64

		err := rows.Scan(&title, &avgRating)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		movie := Movie{
			Title:  title,
			Rating: avgRating,
		}

		movies = append(movies, movie)
	}

	return movies
}

func searchTaskC(db *sql.DB, genre string) []Movie {
	// 查询某一风格最受欢迎的20部电影的示例代码
	rows, err := db.Query(`
		SELECT m.title, AVG(r.rating) as avgRating
		FROM movies m
		JOIN ratings r ON m.movieId = r.movieId
		WHERE FIND_IN_SET(?, m.genres) > 0
		GROUP BY m.title
		ORDER BY avgRating DESC
		LIMIT 20
	`, genre)

	if err != nil {
		log.Println("Error executing query:", err)
		return nil
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var title string
		var avgRating float64

		err := rows.Scan(&title, &avgRating)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		movie := Movie{
			Title:  title,
			Rating: avgRating,
		}

		movies = append(movies, movie)
	}

	return movies
}

func searchTaskD(db *sql.DB, gender string) []Movie {
	// 根据用户性别推荐最受欢迎的电影20部电影的示例代码
	rows, err := db.Query(`
		SELECT m.title, AVG(r.rating) as avgRating
		FROM movies m
		JOIN ratings r ON m.movieId = r.movieId
		JOIN users u ON r.userId = u.userId
		WHERE u.gender = ?
		GROUP BY m.title
		ORDER BY avgRating DESC
		LIMIT 20
	`, gender)

	if err != nil {
		log.Println("Error executing query:", err)
		return nil
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var title string
		var avgRating float64

		err := rows.Scan(&title, &avgRating)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		movie := Movie{
			Title:  title,
			Rating: avgRating,
		}

		movies = append(movies, movie)
	}

	return movies
}

func searchTaskE(db *sql.DB, gender string, minRating float64) []Movie {
	// 区分性别，查询高于某个评分的打分情况的示例代码
	rows, err := db.Query(`
		SELECT m.title, r.rating
		FROM movies m
		JOIN ratings r ON m.movieId = r.movieId
		JOIN users u ON r.userId = u.userId
		WHERE u.gender = ? AND r.rating > ?
		ORDER BY r.rating DESC
	`, gender, minRating)

	if err != nil {
		log.Println("Error executing query:", err)
		return nil
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var title string
		var rating float64

		err := rows.Scan(&title, &rating)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		movie := Movie{
			Title:  title,
			Rating: rating,
		}

		movies = append(movies, movie)
	}

	return movies
}

func ginHtml(c *gin.Context) {
	// 传递给 HTML 模板的数据
	data := gin.H{"user_name": "lihan", "age": 32, "status": http.StatusOK, "data": gin.H{"id": 1, "name": "lihan"}}
	c.HTML(http.StatusOK, "index.html", data)
}
