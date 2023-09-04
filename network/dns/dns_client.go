package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// Set the DNS server address
	dnsServer := "8.8.8.8:53"

	// Domain name you want to resolve
	domain := os.Args[1]
	fmt.Println("hostName:", domain)
	// Create a UDP connection to the DNS server
	conn, err := net.Dial("udp", dnsServer)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		return
	}
	defer conn.Close()

	// Prepare a DNS query message
	dnsQuery := makeDNSQuery(domain)

	// Send the DNS query to the server
	_, err = conn.Write(dnsQuery)
	if err != nil {
		fmt.Println("Error sending DNS query:", err)
		return
	}
	respBuffer := make([]byte, 4096)
	n, err := conn.Read(respBuffer)
	if err != nil {
		fmt.Println("Error reading data from connection", err)
	}
	parseDNSResponse(respBuffer[:n], domain)
}

func parseDNSResponse(respBuffer []byte, hostName string) {
	// dns response header:
	if len(respBuffer) < 13 {
		log.Fatalf("Malformed response, len should be atleast 12")
	}
	txnId := respBuffer[0:2]
	if !bytes.Equal(txnId, []byte{0x00, 0x01}) {
		log.Fatalf("response txnId is not same as reqId")
	}
	// skip 2 bytes of flags and 2 bytes of questions
	answerRRs := respBuffer[6:8]
	if !bytes.Equal(answerRRs, []byte{0x00, 0x01}) {
		log.Fatalf("answer record should be 1")
	}
	// skip 2 bytes of Auth RRs and Additional RRs
	// skip query -- this depends on your host
	queryLen := 0

	for _, part := range strings.Split(hostName, ".") {
		queryLen++ // length octet
		queryLen += len(part)
	}
	queryLen++ // for end of hostname
	index := 12 + queryLen
	if len(respBuffer) < index {
		log.Fatalf("Malformed response, len should be more than %d", index)
	}
	index += 2 + 2 // type and class

	// Resource record format
	/*
		+ 2 // domain name offset
		+ 2 // type
		+ 2 // class
		+ 4 // ttl
		+ 2 // rdlength
		// answer record starts here
	*/
	answerIndex := index + 2 + 2 + 2 + 4 + 2
	if len(respBuffer) < answerIndex {
		log.Fatalf("Malformed response, len should be more than %d", answerIndex)
	}
	answer := respBuffer[answerIndex : answerIndex+4]

	fmt.Println("hostIP:", net.IP(answer))
}
func makeDNSQuery(domain string) []byte {
	// Prepare a basic DNS query for A records (IPv4 addresses)
	// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
	dnsQuery := []byte{
		0x00, 0x01, // ID (16 bits)
		0x01, 0x00, // Flags (QR=0, Opcode=0, AA=0, TC=0, RD=1, RA=0, Z=0, RCODE=0) (16 bits)
		0x00, 0x01, // Questions (QDCOUNT) (16 bits)
		0x00, 0x00, // Answer RRs (ANCOUNT) (16 bits)
		0x00, 0x00, // Authority RRs (NSCOUNT) (16 bits)
		0x00, 0x00, // Additional RRs (ARCOUNT) (16 bits)
	}

	// Prepare QNAME
	for _, part := range strings.Split(domain, ".") {
		// length octet
		dnsQuery = append(dnsQuery, byte(len(part)))
		dnsQuery = append(dnsQuery, []byte(part)...)
	}

	// zero length octet -- end of QNAME
	dnsQuery = append(dnsQuery, 0x00)
	// Add the QTYPE (A record type for IPv4) to the query
	dnsQuery = append(dnsQuery, 0x00, 0x01)
	// Add the QCLASS (IN for internet) to the query
	dnsQuery = append(dnsQuery, 0x00, 0x01)

	return dnsQuery
}
