package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	msgCh := make(chan []byte)

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	memoryCache, err := NewCache()
	if err != nil {
		log.Fatal(err)
	}
	natsSt, err := NewNatsConnection(msgCh)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	go natsSt.NewSubscription()
	go store.CreateUserFromNATS(natsSt.GetMsgChannel())

	server := NewAPIServer(":3000", store, memoryCache, natsSt)
	server.Run()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

}
