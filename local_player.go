package main

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var localPlayer *LocalPlayer

type LocalPlayer struct {
	EntityRuntimeID uint64
	Position        mgl32.Vec3
}

func initializeLocalPlayer() {
	localPlayer = &LocalPlayer{
		EntityRuntimeID: localConnection.GameData().EntityRuntimeID,
		Position:        localConnection.GameData().PlayerPosition,
	}
}

func (p *LocalPlayer) SendMessage(message string) {
	err := localConnection.WritePacket(&packet.Text{
		TextType: packet.TextTypeChat,
		Message:  message,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to send message to client \"%s\".", message))
	}
}

func (p *LocalPlayer) SendMessageToServer(message string) {
	err := serverConnection.WritePacket(&packet.Text{
		TextType: packet.TextTypeChat,
		Message:  message,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to send message to server \"%s\".", message))
	}
}
