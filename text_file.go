package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your_module_name/api"
	"github.com/your_module_name/cassandra"
	"github.com/your_module_name/kafka"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStartConsumerAPI(t *testing.T) {
	router := gin.Default()
	router.GET("/consume", api.StartConsumerAPI)

	req, _ := http.NewRequest("GET", "/consume", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Consumer started in background")
}

func TestInsertSentenceAPI(t *testing.T) {
	router := gin.Default()
	router.POST("/insert", api.InsertSentenceAPI)

	body := {"sentence": "Test sentence"}
	req, _ := http.NewRequest("POST", "/insert", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Sentence inserted successfully")
}

func TestPublishSentencesAPI(t *testing.T) {
	router := gin.Default()
	router.POST("/publish", api.PublishSentencesAPI)

	req, _ := http.NewRequest("POST", "/publish", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestConnectToCassandra(t *testing.T) {
	session, err := cassandra.ConnectToCassandra()
	assert.NoError(t, err)
	assert.NotNil(t, session)
	defer session.Close()
}

func TestInsertSentence(t *testing.T) {
	session, err := cassandra.ConnectToCassandra()
	assert.NoError(t, err)
	defer session.Close()

	err = cassandra.InsertSentence(session, "Unit test sentence")
	assert.NoError(t, err)
}

func TestNewProducer(t *testing.T) {
	producer, err := kafka.NewProducer()
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestPublishMessage(t *testing.T) {
	producer, err := kafka.NewProducer()
	assert.NoError(t, err)
	defer producer.Close()

	err = kafka.PublishMessage(producer, "test-topic", "Test message")
	assert.NoError(t, err)
}

func TestNewConsumer(t *testing.T) {
	consumer, err := kafka.NewConsumer()
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
}

func TestConsumeMessages(t *testing.T) {
	consumer, err := kafka.NewConsumer()
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	go kafka.ConsumeMessages(consumer)
	// No assertion — test passes if no panic
}