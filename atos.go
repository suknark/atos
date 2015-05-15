package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Config struct {
	Memcached     string
	Aerospike     string
	Elasticsearch string
}

func ReadConfig() (string, string, string) {

	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	c := Config{}
	err := decoder.Decode(&c)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	return c.Memcached, c.Aerospike, c.Elasticsearch
}

func StatsItems(command string, netAdr string) (string, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", netAdr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	if strings.Index(netAdr, ":9200") != -1 {
		command = strings.Replace(command, "\n", "", -1) + " HTTP/1.1\n"
	}
	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
	reply := make([]byte, 4096)
	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
	stats := string(reply)
	conn.Close()
	return stats, err

}

func ConnectResource(resource string) {

	memAddr, aeAddr, elAddr := ReadConfig()

memcached:

	if resource == "memcached" || resource == "me"{
		for {
			fmt.Print("goched> ")
			reader := bufio.NewReader(os.Stdin)
			cmd, _ := reader.ReadString('\n')
			if cmd == "exit\n" || cmd == "q\n" {
				break
			}
			if cmd == "aerospike\n" || cmd == "ae\n" {
				resource = "aerospike"
				goto aerospike
			}
			if cmd == "elasticsearch\n" || cmd == "el\n" {
				resource = "elasticsearch"
				goto elasticsearch
			}

			stat_items, _ := StatsItems(cmd, memAddr)
			fmt.Printf("%s\n", string(stat_items))
		}
	}

aerospike:

	if resource == "aerospike" || resource == "ae" {
		for {
			fmt.Print("gospike> ")
			reader := bufio.NewReader(os.Stdin)
			cmd, _ := reader.ReadString('\n')
			if cmd == "exit\n" || cmd == "q\n" {
				break
			}
			if cmd == "memcached\n" || cmd == "me\n" {
				resource = "memcached"
				goto memcached
			}
			if cmd == "elasticsearch\n" || cmd == "el\n" {
				resource = "elasticsearch"
				goto elasticsearch
			}

			it, _ := StatsItems(cmd, aeAddr)
			fmt.Printf("%s\n", it)
		}
	}

elasticsearch:
	if resource == "elasticsearch" || resource == "el" {
		for {
			fmt.Print("gostic> ")
			reader := bufio.NewReader(os.Stdin)
			cmd, _ := reader.ReadString('\n')
			if cmd == "exit\n" || cmd == "q\n" {
				break
			}
			if cmd == "memcached\n" || cmd == "me\n"{
				resource = "memcached"
				goto memcached
			}

			if cmd == "aerospike\n" || cmd == "ae\n" {
				resource = "erospike"
				goto aerospike
			}
			it, _ := StatsItems(cmd, elAddr)
			fmt.Printf("%s\n", it)

		}
	}

}

func PrintHelp() {
	fmt.Printf("Simple usage:\n type \"aerospike\" to use aerospike storage\n type \"memcached\" to use memcached-storage\n")
}

func main() {
	memAddr, aeAddr, elAddr := ReadConfig()
	var resource, adr, c string
	if len(os.Args[:]) > 1 {
		if len(os.Args[:]) > 2 {
			if os.Args[1] == "memcached" || os.Args[1] == "me"{
				adr = memAddr
			}
			if os.Args[1] == "aerospike" || os.Args[1] == "ae"{
				adr = aeAddr
			}
			if os.Args[1] == "elasticsearch" || os.Args[1] == "el"{
				adr = elAddr
			}
			for _, cc := range os.Args[2:] {
				c = c + " " + cc
			}
			it, _ := StatsItems(c, adr)
			fmt.Printf("%s\n", it)
			os.Exit(0)
		} else {
			ConnectResource(os.Args[1])
		}
	}
	for {
		fmt.Print("> ")
		fmt.Scanf("%s\n", &resource)
		if resource == "h" || resource == "help" {
			PrintHelp()
		}
		if resource == "q" || resource == "exit" {
			break
		}
		ConnectResource(resource)

	}

}
