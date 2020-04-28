package connection

import (
	"fmt"
	"github.com/pion/webrtc/v2"
	"log"
)

type Connection struct {
	*webrtc.PeerConnection
	*webrtc.DataChannel

	connectionState webrtc.ICEConnectionState
	channelState    string

	messageCallback func(msg []byte)
}

var config = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

func (c *Connection) State() string {
	if c == nil {
		return ""
	}
	return c.connectionState.String()
}

func (c *Connection) Send(msg []byte) error {
	if c == nil {
		return nil
	}

	if c.channelState != "open" {
		return nil
	}

	return c.DataChannel.Send(msg)
}

func (c *Connection) ChannelState() string {
	if c == nil {
		return ""
	}
	return c.channelState
}

func (c *Connection) OnStateChange(state webrtc.ICEConnectionState) {
	c.connectionState = state
}

func (c *Connection) OnMessageReceived(msg webrtc.DataChannelMessage) {
	fmt.Printf("Message from DataChannel '%s'\n", string(msg.Data))
	c.messageCallback(msg.Data)
}

func CreateForOffer(messageCallback func(msg []byte)) (*Connection, error) {
	conn := &Connection{}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal("failed to create peer connection", err)
	}

	peerConnection.OnICEConnectionStateChange(conn.OnStateChange)

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		log.Fatal("failed to create data channel", err)
	}
	dataChannel.OnOpen(func() {
		conn.channelState = "open"
		conn.DataChannel = dataChannel
	})
	dataChannel.OnMessage(conn.OnMessageReceived)

	conn.PeerConnection = peerConnection
	conn.DataChannel = dataChannel

	conn.messageCallback = messageCallback
	return conn, nil
}

func CreateForAnswer(messageCallback func(msg []byte)) (*Connection, error) {
	conn := &Connection{}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal("failed to create peer connection", err)
	}

	peerConnection.OnICEConnectionStateChange(conn.OnStateChange)
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())
		d.OnOpen(func() {
			conn.channelState = "open"
			conn.DataChannel = d
		})
		d.OnMessage(conn.OnMessageReceived)

	})

	conn.PeerConnection = peerConnection
	conn.messageCallback = messageCallback
	return conn, nil
}
