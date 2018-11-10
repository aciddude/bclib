package network

import (
	"git.posc.in/cw/watchers/serial"
	"math/rand"
	"bytes"
	"encoding/binary"
	"time"
	"fmt"
	"net"
	log "github.com/sirupsen/logrus"
)

// New initializes peer structure
func (p *Peer) New() {
	p.ip = net.ParseIP("37.59.38.74")
	p.port = "8333"
}

// New initializes network structure
func (n *Network) New() {
	n.networkMagic = 0xD9B4BEF9 // Maybe LE
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333
	n.nPeers = 0
}

func (n *Network) AddPeer(p Peer) {
	n.peers = append(n.peers, p)
	n.nPeers++
}

func (n *Network) netAddr() {
}

// NetworkVersion sends the protocol version to the selected peer
func (n *Network) NetworkVersion(id uint32) ([]byte, error) {
	if id >= n.nPeers {
		return nil, fmt.Errorf("NetworkVersion: (id %d) >= (nPeers %d)", id, n.nPeers)
	}
	peer := n.peers[id]

	b := bytes.NewBuffer([]byte{})

	binary.Write(b, binary.LittleEndian, uint32(n.version)) // Protocol version, 70015
	binary.Write(b, binary.LittleEndian, uint64(n.services)) // Network services
	binary.Write(b, binary.LittleEndian, uint64(time.Now().Unix())) // Timestamp

	// Network address of receiver (26)
	// b.Write(c.PeerAddr.NetAddr.Bytes())
	b.Write(peer.ip) // Network address of receiver
	b.Write([]byte(peer.port)) // Network port of receiver

	// Network address of emitter (26)
	b.Write(bytes.Repeat([]byte{0}, 26))

	binary.Write(b, binary.LittleEndian, uint64(rand.Intn(2^64))) // nonce, 8 bytes

	binary.Write(b, binary.LittleEndian, uint64(len(n.userAgent)))
	b.Write([]byte(n.userAgent))
	// b.Write([]byte{0})

	binary.Write(b, binary.LittleEndian, uint32(0)) // Last blockheight received
	b.WriteByte(1)  // don't notify me about txs (BIP37)

	log.Info(fmt.Sprintf("version %x", b.Bytes()))
	return b.Bytes(), nil
}

// SendRawMsg sends command and payload
func (n *Network) NetworkMsg(cmd string, pl []byte) ([]byte, error) {
	var sbuf [24]byte

	// fmt.Println(c.ConnID, "sent", cmd, len(pl))

	binary.LittleEndian.PutUint32(sbuf[0:4], n.networkMagic)
	copy(sbuf[4:16], cmd) // version
	binary.LittleEndian.PutUint32(sbuf[16:20], uint32(len(pl)))

	chksum := serial.DoubleSha256(pl[:])
	copy(sbuf[20:24], chksum[:4])

	// c.append_to_send_buffer(sbuf[:])
	// c.append_to_send_buffer(pl) // payload
	msg := append(sbuf[:], pl...)

	log.Info(fmt.Sprintf("raw msg %x", msg))

	return msg, nil
}

// this function assumes that there is enough room inside sendBuf
/*
func append_to_send_buffer(d []byte) {
 room_left := SendBufSize - c.SendBufProd
 if room_left>=len(d) {
  copy(c.sendBuf[c.SendBufProd:], d)
  room_left = c.SendBufProd+len(d)
 } else {
  copy(c.sendBuf[c.SendBufProd:], d[:room_left])
  copy(c.sendBuf[:], d[room_left:])
 }
 c.SendBufProd = (c.SendBufProd + len(d)) & SendBufMask
}
*/
