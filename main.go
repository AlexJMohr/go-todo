package main

import (
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed assets/*
var assetsFS embed.FS

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

func setupRouter() *gin.Engine {
	r := gin.Default()

	tmpl := template.Must(template.New("").ParseFS(templatesFS, "templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	sub, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		log.Fatal("failed to sub assets FS:", err)
	}
	r.StaticFS("/assets", http.FS(sub))

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
			return
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}

		task.Completed = !task.Completed
		result = db.Save(&task)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}

		c.HTML(http.StatusOK, "task.html", task)
	})

	return r
}

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		initDB()
		r := setupRouter()
		adapter := ginadapter.NewV2(r)
		lambda.Start(adapter.ProxyWithContext)
	} else {
		godotenv.Load(".env")
		initDB()
		r := setupRouter()
		if err := r.Run(); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}
}
