package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/yunduansing/gocommon/gen"
	"github.com/yunduansing/gocommon/logger"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Reader struct {
	ready chan bool
	data  chan KafkaMsg
}

type KafkaMsg struct {
	data *sarama.ConsumerMessage
	ack  func()
}

type Message struct {
	Id        uint64 `json:"id"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Reader) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Reader) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Reader) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Printf("Message claimed: key = %s, value = %s, timestamp = %v, topic = %s", gen.ByteToString(message.Key), gen.ByteToString(message.Value), message.Timestamp, message.Topic)
		consumer.data <- KafkaMsg{
			data: message,
			ack: func() {
				session.MarkMessage(message, "")
			},
		}

	}

	return nil
}

type KafkaStub struct {
	ctx    context.Context
	client sarama.Client
	group  map[string]*sarama.ConsumerGroup
	sender sarama.AsyncProducer
	reader Reader
}

func NewKafkaStub(ctx context.Context, brokers []string) (k *KafkaStub, err error) {
	k = new(KafkaStub)
	k.ctx = ctx
	k.client, err = newSaramaClient(brokers)
	if err != nil {
		return
	}
	k.sender, err = sarama.NewAsyncProducerFromClient(k.client)
	if err != nil {
		return
	}
	k.group = make(map[string]*sarama.ConsumerGroup)
	k.reader = Reader{
		ready: make(chan bool),
		data:  make(chan KafkaMsg),
	}
	return
}

func StartReceive(r Reader, f func([]byte)) {
	for {
		select {
		case d, ok := <-r.data:
			if !ok {
				return
			}
			f(d.data.Key)
			d.ack()
		}
	}
}

func (k *KafkaStub) NewConsumerGroup(ctx context.Context, groupName string, topics []string, f func([]byte)) {
	if _, exists := k.group[groupName]; exists {
		return
	} else {
		consumerGroup, err := sarama.NewConsumerGroupFromClient(groupName, k.client)
		if err != nil {
			log.Panicf("Error creating consumer group client: %v", err)
		}
		k.group[groupName] = &consumerGroup
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := (*k.group[groupName]).Consume(ctx, topics, &k.reader); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			k.reader.ready = make(chan bool)
		}
	}()
	go StartReceive(k.reader, f)
	<-k.reader.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	defer k.sender.AsyncClose()
	if err := k.client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func (k *KafkaStub) SendMessageAsync(ctx context.Context, topic, content string) {
	go func() {
		// [!important] 异步生产者发送后必须把返回值从 Errors 或者 Successes 中读出来 不然会阻塞 sarama 内部处理逻辑 导致只能发出去一条消息
		select {
		case _ = <-k.sender.Successes():
		case e := <-k.sender.Errors():
			if e != nil {
				logger.Error(fmt.Sprintf("[Producer] err:%v msg:%+v \n", e.Msg, e.Err))
			}
		}
	}()

	msg := &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(content)}
	// 异步发送只是写入内存了就返回了，并没有真正发送出去
	// sarama 库中用的是一个 channel 来接收，后台 goroutine 异步从该 channel 中取出消息并真正发送
	k.sender.Input() <- msg
}

func newSaramaClient(brokers []string) (sarama.Client, error) {
	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.MaxMessageBytes = 1000000 * 5
	config.Net.MaxOpenRequests = 1
	return sarama.NewClient(brokers, config)
}
