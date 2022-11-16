package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

func SendMessage() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// test msg
	msg := &sarama.ProducerMessage{}
	msg.Topic = "test"
	msg.Value = sarama.StringEncoder("this is a test log")
	// connect kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Printf("producer close, error: %v", err)
		return
	}
	defer client.Close()

	// send message
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		log.Printf("send message failed, error: %v", err)
		return
	}
	log.Printf("pid: %v offset: %v\n", pid, offset)
}
