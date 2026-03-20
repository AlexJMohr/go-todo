package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Task      string    `json:"task" form:"task" binding:"required" gorm:"not null"`
	Completed bool      `json:"completed" default:"false" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
}

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	db.AutoMigrate(&Task{})
}

func main() {
	godotenv.Load(".env")

	initDB()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	r.GET("/", func(ctx *gin.Context) {
		tasks := []Task{}
		result := db.Order("created_at desc").Find(&tasks)
		if result.Error != nil {
			ctx.String(http.StatusInternalServerError, "Failed to load tasks: %v", result.Error)
			return
		}
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
		result := db.Create(&newTask)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
			return
		}
		c.HTML(http.StatusCreated, "task.html", newTask)
	})

	r.PATCH("/api/tasks/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		taskID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		var task Task
		result := db.First(&task, "id = ?", taskID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		}

		task.Completed = !task.Completed
		result = db.Save(&task)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}

		c.HTML(http.StatusOK, "task.html", task)
	})

	err := r.Run()
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
