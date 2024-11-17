package main

import (
	"fmt"
	"net"
)

// Header section structure
type DnsHeader struct {
	// ID represents a unique identifier for a DNS query or response.
	// It is a 16-bit unsigned integer (2 bytes).
	ID uint16
	// QR represents the Query/Response flag in a DNS message.
	// It is a single bit field where a value of 0 indicates a query,
	// and a value of 1 indicates a response.
	QR uint8
	// OPCODE represents the operation code for a DNS query, which indicates the kind of query being made.
	// It is an 8-bit unsigned integer where the most common values are:
	// 0 - Standard query (QUERY)
	// 1 - Inverse query (IQUERY)
	// 2 - Server status request (STATUS)
	OPCODE uint8
	// AA represents the Authoritative Answer flag in the DNS header.
	// It indicates that the responding DNS server is authoritative for the domain name in the query.
	AA uint8
	// TC represents the Truncation flag in the DNS header. It is a single bit field
	// that indicates whether the message was truncated due to length greater than
	// that permitted on the transmission channel.
	TC uint8
	// RD represents the Recursion Desired flag in the DNS header.
	// It is a single bit field that indicates whether the client
	// wants the DNS server to perform recursive query resolution.
	RD uint8
	// RA represents the Recursion Available flag in the DNS header.
	// It indicates whether the DNS server supports recursive queries.
	RA uint8
	// Z represents a placeholder for an 8-bit unsigned integer value.
	Z uint8
	// RCODE represents the response code in a DNS message, indicating the status of the response.
	// It is an 8-bit unsigned integer where different values correspond to different response statuses.
	// For example, 0 indicates no error, 1 indicates a format error, 2 indicates a server failure, etc.
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
	data[2] = byte(h.QR) // data[2] = 1 << 7 / uint8(packet[2]) << 7
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
