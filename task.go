
package models

import "github.com/gocql/gocql"

type Task struct {
	ID        gocql.UUID `json:"id"`
	Title     string     `json:"title"`
	Completed bool       `json:"completed"`
}
