package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	//Set flags
	filePtr := flag.String("f", "", "File you wish to send.")
	ipPtr := flag.String("l", "127.0.0.1", "IP to transfer a file to.")
	portPtr := flag.String("p", "1234", "Port to transfer a file to")

	flag.Parse()

	if *filePtr == "" {
		fmt.Print("Please enter a file to send\n")
		os.Exit(1)
	}

	//Initialise UDP Connection
	conn, err := net.Dial("udp", *ipPtr+":"+*portPtr)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}

	//Load file to memory
	file, _ := ioutil.ReadFile(*filePtr)

	//Generate hash
	hash := sha256.New()
	hash.Write(file)

	//Calculate total number of packets.
	totalPackets := int(len(file) / 1472)

	//Send control packet indicating a file transfer is in progress
	var startPacket []byte
	startPacket = append(startPacket, "!XxSENDxX!"...)
	startPacket = append(startPacket, fmt.Sprintf("%10d\n", totalPackets)...)
	startPacket = append(startPacket, hex.EncodeToString(hash.Sum(nil))...)
	startPacket = append(startPacket, filepath.Base(string([]byte(*filePtr)))...)
	conn.Write(startPacket)
	fmt.Printf("MD5 Sum: %s\n", hex.EncodeToString(hash.Sum(nil)))
	fmt.Printf("Total packets: %d\n", totalPackets)

	//Send main body of packets
	for i := 0; i < totalPackets; i++ {
		conn.Write(file[i*1472 : (i*1472)+1472])

		//Sleep a bit per packet to allow for transfer time
		time.Sleep(500 * time.Microsecond)

		//Calculate total percentage sent so far and output to screen
		total := int(float32(i) / float32(totalPackets) * 100)
		fmt.Printf("\rFile %s: %d%% sent", *filePtr, total)
	}

	//Send extra packet if it's needed at the end of a file transfer (Because file does not divide exactly by 1472 bytes)
	if len(file)%1472 != 0 {
		var finalPacket []byte
		finalPacket = append(finalPacket, "!XxLASTxX!"...)
		finalPacket = append(finalPacket, fmt.Sprintf("%8d\n", len(file[totalPackets*1472:len(file)]))...)
		finalPacket = append(finalPacket, file[totalPackets*1472:len(file)+1]...)
		fmt.Printf("\rFile %s: 100%% sent", *filePtr)

		fmt.Print("\nFile sent.")
		conn.Write(finalPacket)
	}

	//Send control packet indicating file transfer is complete
	conn.Write([]byte("!XxDONExX!"))
	conn.Close()
}
