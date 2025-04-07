package main

import (
	"todo-backend/cassandra"
	"todo-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func main() {
	cassandra.Connect()
	defer cassandra.Session.Close()

	r := gin.Default()


	r.GET("/tasks", routes.GetTasks)
	r.GET("/tasks/completed", routes.GetCompletedTasks)
	r.GET("/tasks/pending", routes.GetPendingTasks)
	r.POST("/tasks", routes.AddTask)
	r.PUT("/tasks/:id", routes.UpdateTask)
	r.DELETE("/tasks/:id", routes.DeleteTask)

	initTasks()

	r.Run(":8080")
}

func initTasks() {
	tasks := []string{
		"Buy groceries", "Complete Go project", "Call mom", "Read a book", "Water the plants",
		"Clean the kitchen", "Do the laundry", "Work out", "Write a blog post", "Fix the bug",
		"Learn Cassandra", "Refactor code", "Pay bills", "Check emails", "Backup files",
		"Plan the trip", "Take a walk", "Meditate", "Cook dinner", "Organize desk",
	}
	for _, title := range tasks {
		id := cassandra.Session.Query("INSERT INTO tasks (id, title, completed) VALUES (?, ?, ?)",
			gocql.TimeUUID(), title, false)
		id.Exec()
	}
}
