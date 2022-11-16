package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

func consumer() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		log.Printf("failed to start consumer, error: %v", err)
		return
	}
	partitionList, err := consumer.Partitions("test")
	if err != nil {
		log.Printf("failed to get list of partition, error: %v", err)
		return
	}
	log.Println(partitionList)

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("test", int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("failed to start consumer for partition %d, error: %v", partition, err)
			return
		}
		defer pc.AsyncClose()
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				log.Printf("Partition: %d Offset: %d Key: %v Value: %v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
}
