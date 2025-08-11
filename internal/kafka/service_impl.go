package kafka

import (
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/GophKeeper/internal/storage/sqlite"
)

type KafkaServiceImpl struct {
	lg            *zap.SugaredLogger
	kafkaHost     string
	sqliteService sqlite.SqliteService
}

func NewKafkaService(lg *zap.SugaredLogger, kafkaHost string, sqliteService sqlite.SqliteService) KafkaService {
	return &KafkaServiceImpl{
		lg:            lg,
		kafkaHost:     kafkaHost,
		sqliteService: sqliteService,
	}
}

func InitKafkaTopics(kafkaHost string) error {
	conn, err := kafka.Dial("tcp", kafkaHost)
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
