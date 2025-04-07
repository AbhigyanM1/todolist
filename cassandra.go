package cassandra

import (
	"log"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func Connect() {
	cluster := gocql.NewCluster("127.0.0.1") // or "localhost"
	cluster.Port = 9042
	cluster.Keyspace = "todoapp"
	cluster.Consistency = gocql.Quorum

	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
}
