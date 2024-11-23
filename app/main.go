package main

import (
	"bytes"
	"fmt"
	"net"
)

// Header
type DnsHeader struct {
	ID      uint16
	QR      uint8
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

// Question
type DnsQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

func (h *DnsHeader) PacketParser(packet []byte) []byte {
	data := make([]byte, 12)

	h.ID = uint16(packet[0])<<8 | uint16(packet[1])
	data[0] = byte(h.ID >> 8) // 0x04
	data[1] = byte(h.ID)      // 0xd2

	h.QR = uint8(packet[2]) << 7
	data[2] = byte(h.QR)

	// Set ODCOUNT to 1
	h.QDCOUNT = 1
	data[4] = 0x00 // data[4] = byte(h.QDCOUNT >> 8)
	data[5] = 0x01

	return data
}

func (q *DnsQuestion) QParser(packet []byte, header *DnsHeader) []byte {
	// Check minimum packet length
	if len(packet) < 13 {
		return []byte{}
	}
	// length of the domain
	domainLen := packet[12]
	if domainLen == 0 {
		return []byte{}
	}

	// append the domain length to the Name field
	q.Name = append(q.Name, domainLen)

	for i, v := range packet {
		if v == domainLen {
			start := i + 1
			domain := packet[start : start+int(domainLen)]
			// append the domain to the Name field
			q.Name = append(q.Name, domain...)
			// append the tld length
			tldLen := packet[start+int(domainLen)]
			q.Name = append(q.Name, tldLen)
			idx := bytes.Index(packet, []byte{tldLen})
			tldStart := idx + 1
			tld := packet[tldStart : tldStart+int(tldLen)]

			// append the tld to the Name field
			q.Name = append(q.Name, tld...)

			// append the null terminator
			zeroIndex := len(packet) - 5
			nullTerminator := packet[zeroIndex]
			q.Name = append(q.Name, nullTerminator)

			// append the type and class
			qtypeStart := zeroIndex + 1
			qtype := packet[qtypeStart : qtypeStart+2]
			q.Name = append(q.Name, qtype...)
			qclassStart := qtypeStart + 2
			qclass := packet[qclassStart : qclassStart+2]
			q.Name = append(q.Name, qclass...)
			break
		}
	}
	return q.Name
}

func combineResponse(header, question []byte) []byte {
	return append(header, question...)
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
		headerBytes := header.PacketParser(buf[:size])

		question := &DnsQuestion{}
		questionBytes := question.QParser(buf[:size], header)

		response := combineResponse(headerBytes, questionBytes)

		_, err = udpConn.WriteToUDP(response, remoteAddr)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
