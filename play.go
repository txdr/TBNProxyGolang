package main

import (
	"errors"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
	"sync"
)

var proxyRunningMu sync.Mutex
var proxyRunning = false
var proxyFailReasonMu sync.Mutex
var proxyFailReason int
var serverConnection *minecraft.Conn
var localConnection *minecraft.Conn

func startProxy(remoteAddress string, port int, token *oauth2.Token) {
	proxyFailReason = -1
	if proxyRunning {
		killCurrentProxy(false)
	}

	var src oauth2.TokenSource
	if token != nil {
		src = auth.RefreshTokenSource(token)
	}

	p, err := minecraft.NewForeignStatusProvider(remoteAddress)
	if err != nil {
		panic(err)
	}
	listener, err := minecraft.ListenConfig{
		StatusProvider: p,
	}.Listen("raknet", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(c.(*minecraft.Conn), listener, src, remoteAddress)
	}
}

func handleConn(conn *minecraft.Conn, listener *minecraft.Listener, src oauth2.TokenSource, remoteAddress string) {
	serverConn, err := minecraft.Dialer{
		TokenSource: src,
		ClientData:  conn.ClientData(),
	}.Dial("raknet", remoteAddress)
	if err != nil {
		proxyFailReasonMu.Lock()
		proxyFailReason = 0
		proxyFailReasonMu.Unlock()
	}
	proxyRunning = true
	var g sync.WaitGroup
	g.Add(2)
	go func() {
		if err := conn.StartGame(serverConn.GameData()); err != nil {
			panic(err)
		}
		g.Done()
	}()
	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			panic(err)
		}
		g.Done()
	}()
	g.Wait()

	localConnection = conn
	initializeLocalPlayer()

	serverConnection = serverConn
	proxyRunningMu.Lock()
	proxyRunning = true
	proxyRunningMu.Unlock()

	localPlayer.SendMessage("Hello this is a test.")

	go func() {
		defer listener.Disconnect(conn, "connection lost")
		defer serverConn.Close()
		defer func() {
			proxyRunningMu.Lock()
			proxyRunning = false
			proxyRunningMu.Unlock()
		}()
		for {
			pk, err := conn.ReadPacket()
			if err != nil {
				return
			}

			if err := serverConn.WritePacket(pk); err != nil {
				var disc minecraft.DisconnectError
				if ok := errors.As(err, &disc); ok {
					_ = listener.Disconnect(conn, disc.Error())
				}
				return
			}
		}
	}()
	go func() {
		defer serverConn.Close()
		defer listener.Disconnect(conn, "connection lost")
		for {
			pk, err := serverConn.ReadPacket()
			if err != nil {
				var disc minecraft.DisconnectError
				if ok := errors.As(err, &disc); ok {
					_ = listener.Disconnect(conn, disc.Error())
				}
				return
			}
			if err := conn.WritePacket(pk); err != nil {
				return
			}
		}
	}()
}

func killCurrentProxy(transferring bool) {
	localConnection.Close()
	proxyRunningMu.Lock()
	proxyRunning = false
	proxyRunningMu.Unlock()
}
