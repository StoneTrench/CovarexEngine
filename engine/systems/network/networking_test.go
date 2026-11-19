package network_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/StoneTrench/CovarexEngine/engine/core/log"
	"github.com/StoneTrench/CovarexEngine/engine/platform"
	"github.com/StoneTrench/CovarexEngine/engine/systems/network"
)

func Test(t *testing.T) {
	log.Init(new(platform.VFS).Init(nil, true), "", true)

	slog.Info("Server started")
	server, err := network.StartServer(context.Background(), 25565)
	if err != nil {
		t.Error(err)
		return
	}
	eh1 := server.PacketSource.Subscribe(func(e network.NetworkEvent) error {
		slog.Info(fmt.Sprintf("Server event: %v", e))
		return nil
	})
	defer server.PacketSource.Unsubscribe(eh1)
	defer server.Close()

	var client *network.Client
	for i := range 10 {
		slog.Info(fmt.Sprintf("Client started: %d", i))
		client, err = network.StartClient(context.Background(), "127.0.0.1:25565")
		if err != nil {
			t.Error(err)
			return
		}
		eh2 := client.PacketSource.Subscribe(func(e network.NetworkEvent) error {
			slog.Info(fmt.Sprintf("Client event: %v", e))
			return nil
		})
		defer client.PacketSource.Unsubscribe(eh2)
		defer client.Close()
	}

	slog.Info("Started sleeping for 5s")
	time.Sleep(time.Second * 10)
}
