package main

import (
	"alert/internal"
	"alert/pkg"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := pkg.Load()
	discordClient := internal.NewDiscord(config.WebhookURL, config.RPCTimeout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleShutdownSignals(cancel)

	for _, node := range config.Nodes {
		node := node
		go internal.MonitorNode(ctx, node, config, discordClient)
	}

	<-ctx.Done()
	log.Println("end")
}

func handleShutdownSignals(cancel context.CancelFunc) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	<-channel
	log.Println("stop signal")
	cancel()
}
