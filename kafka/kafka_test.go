package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestNewConsumerGroup(t *testing.T) {
	var brokers = []string{"192.168.2.150:9092"}
	ctx, cancle := context.WithCancel(context.Background())
	k, err := NewKafkaStub(ctx, brokers)
	if err != nil {
		log.Fatalln(err)
		return
	}

	k.NewConsumerGroup(ctx, "ccs-gozero-group", []string{"ccs-associatapi-broadcast", "ccs-strategyapi-broadcast", "ccs-associatapi-channel", "ccs-strategyapi-channel"}, func(b []byte) {
		fmt.Println("msg", string(b))
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ctx.Done():
		log.Println("terminating by context canceled")
	case <-sig:
		log.Println("terminating via signal")
	}
	cancle()
}
