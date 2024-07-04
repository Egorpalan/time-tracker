package handlers

import (
	"log"
	"net/http"
	"time"
	"fmt"

	"github.com/Egorpalan/time-tracker/internal/models"
	"github.com/Egorpalan/time-tracker/internal/repository"
	"github.com/gin-gonic/gin"
)

// GetUsers returns a list of users with optional filtering and pagination
func GetUsers(c *gin.Context) {
	var users []models.User
	query := repository.DB

	if passportNumber := c.Query("passportNumber"); passportNumber != "" {
		query = query.Where("passport_number = ?", passportNumber)
	}
	if surname := c.Query("surname"); surname != "" {
		query = query.Where("surname LIKE ?", "%"+surname+"%")
	}
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if address := c.Query("address"); address != "" {
		query = query.Where("address LIKE ?", "%"+address+"%")
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	var pageNum, limitNum int
	fmt.Sscanf(page, "%d", &pageNum)
	fmt.Sscanf(limit, "%d", &limitNum)
	offset := (pageNum - 1) * limitNum

	query = query.Offset(offset).Limit(limitNum).Find(&users)

	c.JSON(http.StatusOK, users)
	log.Printf("GetUsers: Retrieved %d users", len(users))
}

// AddUser adds a new user
func AddUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("AddUser: Error binding JSON: %v", err)
		return
	}

	if err := repository.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("AddUser: Error creating user in DB: %v", err)
		return
	}

	c.JSON(http.StatusOK, user)
	log.Printf("AddUser: User created successfully: %+v", user)
}

// UpdateUser updates an existing user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := repository.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Printf("UpdateUser: User not found for id %s", id)
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("UpdateUser: Error binding JSON: %v", err)
		return
	}

	if err := repository.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("UpdateUser: Error saving user in DB: %v", err)
		return
	}

	c.JSON(http.StatusOK, user)
	log.Printf("UpdateUser: User updated successfully: %+v", user)
}

// DeleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := repository.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Printf("DeleteUser: User not found for id %s", id)
		return
	}

	if err := repository.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("DeleteUser: Error deleting user from DB: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	log.Printf("DeleteUser: User deleted successfully with id %s", id)
}

// GetUserTasks returns the time spent by a user on tasks within a specified period
func GetUserTasks(c *gin.Context) {
	userID := c.Param("id")
	startDate := c.Query("start")
	endDate := c.Query("end")

	var tasks []models.Task
	query := repository.DB.Where("user_id = ?", userID)

	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
			log.Printf("GetUserTasks: Error parsing start date: %v", err)
			return
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
			log.Printf("GetUserTasks: Error parsing end date: %v", err)
			return
		}
		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	}

	query.Order("end_time - start_time desc").Find(&tasks)

	c.JSON(http.StatusOK, tasks)
	log.Printf("GetUserTasks: Retrieved %d tasks for user id %s", len(tasks), userID)
}

// StartTask starts a new task for a user
func StartTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("StartTask: Error binding JSON: %v", err)
		return
	}

	task.StartTime = time.Now()
	if err := repository.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("StartTask: Error creating task in DB: %v", err)
		return
	}

	c.JSON(http.StatusOK, task)
	log.Printf("StartTask: Task started successfully: %+v", task)
}

// EndTask ends a task for a user
func EndTask(c *gin.Context) {
	id := c.Param("task_id")
	var task models.Task

	if err := repository.DB.Where("id = ?", id).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		log.Printf("EndTask: Task not found for id %s", id)
		return
	}

	task.EndTime = time.Now()
	if err := repository.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("EndTask: Error saving task in DB: %v", err)
		return
	}

	c.JSON(http.StatusOK, task)
	log.Printf("EndTask: Task ended successfully: %+v", task)
}
