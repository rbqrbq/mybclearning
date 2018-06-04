package Block

import (
	"fmt"
	"net"
	"log"
	"encoding/gob"
	"bytes"
	"io/ioutil"
)

type Version struct {
	Version    int
	BestHeight int
	AddFrom    string
}

type GetBlocks struct {
	AddFrom string
}

type Inv struct {
	AddFrom string
	Type    string
	Items   [][]byte
}

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:3000"}

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := NewBlockchain(nodeID)
	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func sendVersion(addr string, bc *BlockChain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(Version{nodeVersion, bestHeight, nodeAddress})
	request := append(commandToBytes("version"), payload...)
	sendData(addr, request)
}

func handleConnection(conn net.Conn, bc *BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)
	switch command {
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	defer conn.Close()

}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	encoded := gob.NewEncoder(&buff)
	err := encoded.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte
	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func bytesToCommand(data []byte) string {
	var command []byte
	for _, b := range data {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

func sendData(addr string, payload []byte) {

}

func handleVersion(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload Version
	buff.Write(request[:commandLength])
	decode := gob.NewDecoder(&buff)
	err := decode.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight
	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddFrom, bc)
	}

	if (!nodeIsKnown(payload.AddFrom)) {
		knownNodes = append(knownNodes, payload.AddFrom)
	}
}

func sendGetBlocks(address string) {

}

func nodeIsKnown(address string) bool {
	return true
}

func handleGetBlocks(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload GetBlocks
	buff.Write(request[commandLength:])
	decode := gob.NewDecoder(&buff)
	err := decode.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blocks := bc.GetBlockHashes()
	sendInv(payload.AddFrom, "block", blocks)
}

func sendInv(address, kind string, blocks [][]byte) {

}

func handleInv(request []byte, bc *BlockChain) {

}
