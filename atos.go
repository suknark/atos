package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Config struct {
	Memcached string
	Aerospike string
}

func ReadConfig() (string, string) {

	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	c := Config{}
	err := decoder.Decode(&c)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	return c.Memcached, c.Aerospike
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

	memAddr, aeAddr := ReadConfig()

memcached:

	if resource == "memcached" {
		for  {
			fmt.Print("goched> ")
			reader := bufio.NewReader(os.Stdin)
			cmd, _ := reader.ReadString('\n')
			if cmd == "exit\n" || cmd == "q\n" {
				break
			}
			if cmd == "aerospike\n" {
				resource = "aerospike"
				goto aerospike
			}
			stat_items, _ := StatsItems(cmd, memAddr)
			fmt.Printf("%s\n", string(stat_items))
		}
	}

aerospike:

	if resource == "aerospike" {
		for  {
			fmt.Print("gospike> ")
			reader := bufio.NewReader(os.Stdin)
			cmd, _ := reader.ReadString('\n')
			if cmd == "exit\n" || cmd == "q\n" {
				break
			}
			if cmd == "memcached\n" {
				resource = "memcached"
				goto memcached
			}
			it, _ := StatsItems(cmd, aeAddr)
			fmt.Printf("%s\n", it)
		}
	}

}

func PrintHelp() {
	fmt.Printf("Simple usage:\n type \"aerospike\" to use aerospike storage\n type \"memcached\" to use memcached-storage\n")
}

func main() {
	var resource string
	if len(os.Args[:]) > 1 {
		ConnectResource(os.Args[1])
	}
	for  {
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
