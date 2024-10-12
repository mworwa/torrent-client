package peers

import (
	"encoding/binary"
	"fmt"
	"net"
)

type messageId int8

const (
	MessageChoke      messageId = 0
	MessageUnchoke    messageId = 1
	MessageInterested messageId = 2
	MessageBitfield   messageId = 5
	MessageRequest    messageId = 6
	MessagePiece      messageId = 7
)

type PeerMessage struct {
	Id      messageId
	Payload []byte
}

func SendMessage(conn net.Conn, message PeerMessage) error {

	lenght := uint32(len(message.Payload) + 1)
	buf := make([]byte, 4+lenght)

	binary.BigEndian.PutUint32(buf[0:4], lenght)
	buf[4] = byte(message.Id)
	copy(buf[5:], message.Payload)

	_, err := conn.Write(buf)
	if err != nil {
		return err
	}
	fmt.Printf("Message with id: %d sent\n", message.Id)

	return nil
}

func SendInterested(conn net.Conn) error {
	return SendMessage(conn, PeerMessage{Id: MessageInterested})
}

func SendUnchoke(conn net.Conn) error {
	return SendMessage(conn, PeerMessage{Id: MessageUnchoke})
}

func SendRequest(conn net.Conn, index, begin, lenght uint32) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], index)
	binary.BigEndian.PutUint32(payload[4:8], begin)
	binary.BigEndian.PutUint32(payload[8:12], lenght)

	return SendMessage(conn, PeerMessage{Id: MessageRequest, Payload: payload})
}

func ReadMessage(conn net.Conn) (*PeerMessage, error) {
	lenghtBuf := make([]byte, 4)
	_, err := conn.Read(lenghtBuf)
	if err != nil {
		return nil, err
	}

	lenght := binary.BigEndian.Uint32(lenghtBuf)
	// keep-alive message
	if lenght == 0 {
		return nil, nil
	}
	payloadBuf := make([]byte, lenght)
	_, err = conn.Read(payloadBuf)
	if err != nil {
		return nil, err
	}
	peerMessage := PeerMessage{Id: messageId(payloadBuf[0]), Payload: payloadBuf}
	fmt.Printf("Message recived id: %d\n", peerMessage.Id)

	return &peerMessage, nil
}
