package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"github.com/fiorix/go-readline"
)

var coutsExamps = map[string][]string{
	"m": {"memcached", "me"},
	"e": {"elasticsearch", "el"},
	"a": {"aerospike", "ae"},
}

func completer(input, line string, start, end int) []string {
	if len(input) == 1 {
		letters, exists := coutsExamps[strings.ToLower(input)]
		if exists {
			return letters
		}
	}
	return []string{}
}

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
		command = "GET " + strings.Replace(command, "\n", "", -1) + " HTTP/1.1\n"
	}
	fmt.Println(command)
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

	if resource == "memcached" {
		for {
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
			if cmd == "elasticsearch\n" {
				resource = "elasticsearch"
				goto elasticsearch
			}

			stat_items, _ := StatsItems(cmd, memAddr)
			fmt.Printf("%s\n", string(stat_items))
		}
	}

aerospike:

	if resource == "aerospike" {
		for {
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
			if cmd == "elasticsearch\n" {
				resource = "elasticsearch"
				goto elasticsearch
			}

			it, _ := StatsItems(cmd, aeAddr)
			fmt.Printf("%s\n", it)
		}
	}

elasticsearch:
	if resource == "elasticsearch" {
		for {
			fmt.Print("gostic> ")
			reader := bufio.NewReader(os.Stdin)
			cmd, _ := reader.ReadString('\n')
			if cmd == "exit\n" || cmd == "q\n" {
				break
			}
			if cmd == "memcached\n" {
				resource = "memcached"
				goto memcached
			}

			if cmd == "memcached\n" {
				resource = "memcached"
				goto memcached
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
	var adr, c string
	readline.SetCompletionFunction(completer)
	readline.ParseAndBind("TAB: menu-complete")
	if len(os.Args[:]) > 1 {
		if len(os.Args[:]) > 2 {
			if os.Args[1] == "memcached" {
				adr = memAddr
			}
			if os.Args[1] == "aerospike" {
				adr = aeAddr
			}
			if os.Args[1] == "elasticsearch" {
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
		promt := "> "
		resource := readline.Readline(&promt)
		if *resource == "h" || *resource == "help" {
			PrintHelp()
		}
		if *resource == "q" || *resource == "exit" {
			break
		}
		ConnectResource(*resource)
		readline.AddHistory(*resource)

	}

}
