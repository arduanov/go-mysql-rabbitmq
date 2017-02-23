package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/siddontang/go-mysql/canal"
	"github.com/streadway/amqp"
	"encoding/json"
)

var host = flag.String("host", "127.0.0.1", "MySQL host")
var port = flag.Int("port", 3306, "MySQL port")
var user = flag.String("user", "root", "MySQL user, must have replication privilege")
var password = flag.String("password", "", "MySQL password")

var flavor = flag.String("flavor", "mysql", "Flavor: mysql or mariadb")
var dataDir = flag.String("data-dir", "./var", "Path to store data, like master.info")
var serverID = flag.Int("server-id", 101, "Unique Server ID")
var mysqldump = flag.String("mysqldump", "", "mysqldump execution path")

var rabbitmq_host = flag.String("rabbitmq_host", "127.0.0.1", "RabbitMQ host")
var rabbitmq_port = flag.Int("rabbitmq_port", 5672, "RabbitMQ port")
var rabbitmq_user = flag.String("rabbitmq_user", "guest", "RabbitMQ user")
var rabbitmq_password = flag.String("rabbitmq_password", "guest", "RabbitMQ password")
var rabbitmq_exchange = flag.String("rabbitmq_exchange", "mysql", "RabbitMQ exchange name")

func main() {
	flag.Parse()

	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", *host, *port)
	cfg.User = *user
	cfg.Password = *password
	cfg.Flavor = *flavor
	cfg.DataDir = *dataDir

	cfg.ServerID = uint32(*serverID)
	cfg.Dump.ExecutionPath = *mysqldump
	cfg.Dump.DiscardErr = false

	c, err := canal.NewCanal(cfg)
	if err != nil {
		fmt.Printf("create canal err %v", err)
		os.Exit(1)
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", *rabbitmq_user, *rabbitmq_password, *rabbitmq_host, *rabbitmq_port)

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	err = ch.ExchangeDeclare(
		fmt.Sprintf("%s", *rabbitmq_exchange),
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	c.RegRowsEventHandler(&handler{ch})

	err = c.Start()
	if err != nil {
		fmt.Printf("start canal err %V", err)
		os.Exit(1)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sc

	c.Close()
}

type handler struct {
	channel *amqp.Channel
}

func (h *handler) Do(e *canal.RowsEvent) error {
	body, _ := json.Marshal(e)
	key := e.Table.Schema + "." + e.Table.Name + "." + e.Action
	return h.channel.Publish(
		fmt.Sprintf("%s", *rabbitmq_exchange),
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (h *handler) String() string {
	return "TestHandler"
}

func failOnError(err error, msg string) {
	if err != nil {
		//log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
