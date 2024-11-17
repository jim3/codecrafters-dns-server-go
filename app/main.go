package main

import (
	"fmt"
	"net"
)

// Header section structure
type DnsHeader struct {
	ID      uint16 // 2 bytes
	QR      uint8  // 1 bit
	OPCODE  uint8
	AA      uint8
	TC      uint8
	RD      uint8
	RA      uint8
	Z       uint8
	RCODE   uint8
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func (h *DnsHeader) PacketParser(packet []byte) []byte {
	data := make([]byte, 12)
	h.ID = uint16(packet[0])<<8 | uint16(packet[1])
	h.QR = uint8(packet[2]) << 7

	data[0] = byte(h.ID >> 8)
	data[1] = byte(h.ID)
	data[2] = byte(h.QR) // data[2] = 1 << 7
	return data
}

func main() {
	fmt.Println("Starting UDP server...")
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)
	for {
		// Copies datagram payload into `buf`
		size, remoteAddr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		header := &DnsHeader{}
		r := header.PacketParser(buf[:size])

		_, err = udpConn.WriteToUDP(r, remoteAddr)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
