package kafka

import (
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/GophKeeper/internal/storage/bbolt"
)

type Controller struct {
	lg           *zap.SugaredLogger
	kafkaHost    string
	bboltService *bbolt.ServiceImpl
}

func NewController(lg *zap.SugaredLogger, kafkaHost string, bboltService *bbolt.ServiceImpl) *Controller {
	return &Controller{
		lg:           lg,
		kafkaHost:    kafkaHost,
		bboltService: bboltService,
	}
}

func (c *Controller) InitKafkaTopics() error {
	conn, err := kafka.Dial("tcp", c.kafkaHost)
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfig := getTopics()

	err = controllerConn.CreateTopics(topicConfig...)
	if err != nil {
		return fmt.Errorf("create topics: %w", err)
	}
	return nil
}
