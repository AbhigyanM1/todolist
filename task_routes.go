package routes

import (
	"net/http"
	"strings"
	"todo-backend/cassandra"
	"todo-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func GetTasks(c *gin.Context) {
	status := strings.ToLower(c.Query("status"))
	search := strings.ToLower(c.Query("search"))

	var query string
	if status == "true" || status == "false" {
		query = "SELECT id, title, completed FROM tasks WHERE completed = " + status + " ALLOW FILTERING"
	} else {
		query = "SELECT id, title, completed FROM tasks"
	}

	var tasks []models.Task
	iter := cassandra.Session.Query(query).Iter()
	var task models.Task

	for iter.Scan(&task.ID, &task.Title, &task.Completed) {
		if search == "" || strings.Contains(strings.ToLower(task.Title), search) {
			tasks = append(tasks, task)
		}
	}
	if err := iter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func AddTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	task.ID = gocql.TimeUUID()
	if err := cassandra.Session.Query("INSERT INTO tasks (id, title, completed) VALUES (?, ?, ?)",
		task.ID, task.Title, task.Completed).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := gocql.ParseUUID(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var update struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = cassandra.Session.Query("UPDATE tasks SET title = ?, completed = ? WHERE id = ?",
		update.Title, update.Completed, id).Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := gocql.ParseUUID(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if err := cassandra.Session.Query("DELETE FROM tasks WHERE id = ?", id).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func GetCompletedTasks(c *gin.Context) {
	getFilteredTasks(c, true)
}

func GetPendingTasks(c *gin.Context) {
	getFilteredTasks(c, false)
}

func getFilteredTasks(c *gin.Context, completed bool) {
	var tasks []models.Task
	iter := cassandra.Session.Query("SELECT id, title, completed FROM tasks WHERE completed = ? ALLOW FILTERING", completed).Iter()
	var task models.Task

	for iter.Scan(&task.ID, &task.Title, &task.Completed) {
		tasks = append(tasks, task)
	}
	if err := iter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
