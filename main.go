package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `json:"id" binding:"uuid"`
	Task      string    `json:"task" form:"task" binding:"required"`
	Completed bool      `json:"completed" default:"false"`
}

var tasks = []Task{}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", tasks)
	})

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/api/tasks", func(c *gin.Context) {
		var newTask Task
		if err := c.ShouldBind(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTask.ID = uuid.New()
		tasks = append(tasks, newTask)
		c.HTML(http.StatusCreated, "task.html", newTask)
	})

	r.POST("/api/todos/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		taskID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		for i, task := range tasks {
			if task.ID == taskID {
				tasks[i].Completed = !tasks[i].Completed
				c.HTML(http.StatusOK, "task.html", tasks[i])
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	})

	err := r.Run()
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
