package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"strings"
	"time"

	crypto "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	host "gx/ipfs/Qmf5yHzmWAyHSJRPAmZzfk3Yd7icydBLi7eec5741aov7v/go-libp2p-host"

	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"

	net "gx/ipfs/QmSTaEYUgDe1r581hxyd2u9582Hgp3KX4wGwYbRqz2u9Qh/go-libp2p-net"

	libp2p "github.com/libp2p/go-libp2p"
)

func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {
	fmt.Println("listen port", listenPort)
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	{
		fmt.Println(priv.Type)
		if err != nil {
			return nil, err
		}
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}
	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	fmt.Println("ID", basicHost.ID())

	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))
	fmt.Println("hostaddress", hostAddr)

	addr := basicHost.Addrs()[0]
	fulladdress := addr.Encapsulate(hostAddr)
	fmt.Println("basichostAll", basicHost.Addrs())
	fmt.Println("basichostAll", (basicHost.Addrs()[0]))
	fmt.Println("full", fulladdress)

	log.Printf("I am %s\n", fulladdress)
	if secio {
		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fulladdress)
	} else {
		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fulladdress)
	}
	return basicHost, nil
}
func readdata(rw *bufio.ReadWriter) string {
	str, err := rw.ReadString('\n')
	if err != nil {
		return "error"
	} else if str == "" {
		return "Empty"
	} else {
		ReadBlock := make([]Blockmember, 0)
		err := json.Unmarshal([]byte(str), &ReadBlock)
		if err != nil {
			return "error"
		}
		mutex.Lock()
		if len(ReadBlock) >= len(blocks) {
			/*if VerifiedincomingBlock(ReadBlock) == false {
				return "Block not Mathched"
			}*/
			blocks = ReadBlock
			mutex.Unlock()
		} else {
			return "Different Length"
		}
	}
	return "Success"
}
func writedata(rw *bufio.ReadWriter) {
	log.Println("added")
	//Goroutine for broadcasting for every 5 sec
	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			btyes, err := json.Marshal(blocks)
			if err != nil {
				log.Println("Error when Writing Blockchain")
			}
			rw.WriteString(fmt.Sprintf("%s\n", string(btyes)))
			log.Println(string(btyes))
			rw.Flush()
			mutex.Unlock()
		}
	}()
	//Buffer reader for Command console input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Can't read data throuhg Command line")
	}
	//Since read command line will add a new line character
	ordernum := strings.Replace(input, "\n", "", -1)
	newblock := GenerateBlock(len(blocks), blocks[len(blocks)-1].Previoushash, ordernum)
	newblock.Hash = GenerateHash(&newblock)
	mutex.Lock()
	if VerifyBlock() == true {
		blocks = append(blocks, newblock)
	}
	mutex.Unlock()
	bytes, err := json.Marshal(blocks)
	if err != nil {
		log.Println("error")
	}
	mutex.Lock()
	rw.WriteString(string(bytes))
	rw.Flush()
	mutex.Unlock()
}
func VerifiedincomingBlock(ReceviedBlocks []Blockmember) bool {
	for i := 0; i < len(blocks); i++ {
		if blocks[i].Previoushash != ReceviedBlocks[i].Previoushash || blocks[i].Hash != ReceviedBlocks[i].Hash {
			return false
		}
	}
	return true
}
func streamhandler(net net.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(net), bufio.NewWriter(net))
	go readdata(rw)
	go writedata(rw)
}
