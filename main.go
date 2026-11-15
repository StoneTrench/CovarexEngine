package main

import (
	"log/slog"

	"github.com/StoneTrench/CovarexEngine/src/client"
)

func main() {
	err := client.Init()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
