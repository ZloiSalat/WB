package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
)

type Stream interface {
	NewSubscription() error
	handleNATSMessage(*stan.Msg) []byte
	//StanMsgHandler(*stan.Msg) []byte
	GetMsgChannel() chan []byte
}

type NatsStreaming struct {
	ns    stan.Conn
	Data  string `json:"data"`
	msg   *stan.Msg
	msgCh chan []byte
}

func NewNatsConnection(msgCh chan []byte) (*NatsStreaming, error) {
	clientID := "your-client-id-"
	clusterID := "my-cluster"
	natsURL := "nats://localhost:4222" // Update with your NATS Streaming server URL

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}

	//defer sc.Close()
	return &NatsStreaming{
		ns:    sc,
		msgCh: msgCh,
	}, nil
}

func (sc *NatsStreaming) NewSubscription() error {
	subject := "your-subject"
	durableName := "durable-name"
	_, err := sc.ns.Subscribe(subject, func(msg *stan.Msg) {
		data := sc.handleNATSMessage(msg)
		sc.msgCh <- data // Send the data to the channel
	}, stan.DurableName(durableName))

	if err != nil {
		log.Fatalf("Error subscribing to NATS Streaming: %v", err)
	}
	//defer subscription.Close()
	return nil
}

func (sc *NatsStreaming) handleNATSMessage(msg *stan.Msg) []byte {

	if msg == nil {
		log.Println("No message received.")
		return nil
	}

	createNewStreamMsg := new(NatsStreaming)

	err := json.Unmarshal(msg.Data, createNewStreamMsg)
	if err != nil {
		log.Printf("Error unmarshaling NATS message: %v", err)
	}

	fmt.Printf("Received message: %s\n", createNewStreamMsg.Data)

	return []byte(createNewStreamMsg.Data)

}

func (sc *NatsStreaming) GetMsgChannel() chan []byte {
	return sc.msgCh
}
