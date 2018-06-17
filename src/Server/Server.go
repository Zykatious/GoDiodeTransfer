package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	//Set flags
	dirPtr := flag.String("d", "./", "Location you wish files to be saved.")
	portPtr := flag.Int("p", 1234, "Port the server listens on.")
	flag.Parse()

	filename := ""
	var file []byte
	var command byte
	recPackets := 0
	totalPackets := 0
	fileHash := ""

	//Set receiving buffer size to 1500 as final packet may be bigger than 1472 due to headers.
	packetBuffer := make([]byte, 1500)

	addr := net.UDPAddr{
		Port: *portPtr,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}
	for {
		_, _, err := ser.ReadFromUDP(packetBuffer)

		//Test for control packet to say a new file is being received.
		if testEq(packetBuffer[0:10], []byte("!XxSENDxX!")) {
			//Get total number of packets expected to be received and file hash
			totalPackets, _ = strconv.Atoi(strings.TrimSpace(string(packetBuffer[10:20])))
			fileHash = string(packetBuffer[21:85])

			//Trim packet down by stripping null bytes and get receving file's filename
			trimPacket := bytes.Trim(packetBuffer[85:len(packetBuffer)], "\x00")
			filename = filepath.Base(string(trimPacket))

			//Set command byte and reset counters / receiving file
			command = 1
			file = nil
			recPackets = 0

			fmt.Printf("\n\nReceiving file %s\n", filename)
		}

		//Test for control packet to say file transfer is complete.
		if testEq(packetBuffer[0:10], []byte("!XxDONExX!")) {
			//Generate hash on received file
			hash := sha256.New()
			hash.Write(file)

			//Check file hash matches sending file
			if hex.EncodeToString(hash.Sum(nil)) == fileHash {
				fmt.Print("\nFile received. Hash check passed.\n")
				//Save file
				err := ioutil.WriteFile(*dirPtr+filename, file, 0644)
				if err != nil {
					fmt.Printf("Error  %v", err)
					continue
				}
				fmt.Printf("File saved to %s%s", *dirPtr, filename)
			} else {
				fmt.Print("\nFile transfer failed, hashes do not match.\n")
			}

			//Set command byte and reset counters / receiving file
			command = 1
			file = nil
			recPackets = 0
		}

		//Control packet to say packet is the final one
		if testEq(packetBuffer[0:10], []byte("!XxLASTxX!")) {

			//Convert final length string to an integer
			finalLengthString := strings.TrimSpace(string(packetBuffer[10:18]))
			finalLength, err := strconv.Atoi(finalLengthString)
			if err != nil {
				fmt.Printf("Error  %v", err)
				continue
			}
			//Append the correct part of the packet to the end of the file by offsetting the command header.
			file = append(file, packetBuffer[19:finalLength+19]...)

			//Set command byte to stop packet being added to file again
			command = 1
		}

		//If command byte is not set, add packet to receving file.
		if command != 1 {
			file = append(file, packetBuffer[0:1472]...)
			recPackets++
			total := int(float32(recPackets) / float32(totalPackets) * 100)
			fmt.Printf("\rFile %s: %d%% received", filename, total)
		}

		//Reset Command byte to zero and reset receiving packet buffer
		command = 0
		packetBuffer = make([]byte, 1500)

		if err != nil {
			fmt.Printf("Error  %v", err)
			continue
		}
		continue
	}
}

//testEq - Check if a byte array matches a second byte array
func testEq(a, b []byte) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
