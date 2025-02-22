package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors" // Thêm thư viện CORS
	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var tasks []Task
var idCounter = 1

func main() {
	router := gin.Default()
	// 💡 Thêm middleware CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Cho phép frontend truy cập
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))
	// Tạo công việc mới
	router.POST("/tasks", func(c *gin.Context) {
		var newTask Task
		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTask.ID = idCounter
		idCounter++
		tasks = append(tasks, newTask)
		c.JSON(http.StatusCreated, newTask)
	})

	// Lấy danh sách công việc
	router.GET("/tasks", func(c *gin.Context) {
		c.JSON(http.StatusOK, tasks)
	})

	// Cập nhật trạng thái hoàn thành
	router.PUT("/tasks/:id", func(c *gin.Context) {
		var updateData struct {
			Title     string `json:"title"`
			Completed bool   `json:"completed"`
		}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		idParam := c.Param("id")
		for i, task := range tasks {
			if fmt.Sprintf("%d", task.ID) == idParam {
				// tasks[i].Title = updateData.Title
				tasks[i].Completed = updateData.Completed
				c.JSON(http.StatusOK, tasks[i])
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Công việc không tồn tại"})
	})

	// Xóa công việc
	router.DELETE("/tasks/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		for i, task := range tasks {
			if fmt.Sprintf("%d", task.ID) == idParam {
				tasks = append(tasks[:i], tasks[i+1:]...) // Xóa task khỏi slice
				c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Công việc không tồn tại"})
	})

	router.Run(":8080")
}
