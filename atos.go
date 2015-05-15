package main

import (
	"encoding/json"
	"fmt"
	"github.com/fiorix/go-readline"
	"net"
	"os"
	"strings"
)

var coutsExamps = map[string][]string{
	"m": {"memcached"},
	"e": {"elasticsearch"},
	"a": {"aerospike"},
	"g": {"GET"},
	"d": {"DELETE"},
}

type Colors struct {
            Set string
            Reset string
}

type Config struct {
	Memcached     string
	Aerospike     string
	Elasticsearch string
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
	readline.SetCompletionFunction(completer)
	readline.ParseAndBind("TAB: menu-complete")
    	color := Colors{"\033[36m", "\033[0m"}
	resource = strings.Replace(resource, " ", "", -1)

memcached:

	if resource == "memcached"  {
		for {
			p := color.Set + "goched> " + color.Reset
			cmd := readline.Readline(&p)

			if *cmd == "exit" || *cmd == "q" {
				break
			}
			if strings.HasPrefix(*cmd, "aerospike")  {
				resource = "aerospike"
				goto aerospike
			}
			if strings.HasPrefix(*cmd, "elasticsearch") {

				resource = "elasticsearch"
				goto elasticsearch
			}

			stat_items, _ := StatsItems(*cmd, memAddr)
			fmt.Printf("%s\n", string(stat_items))
			readline.AddHistory(*cmd)
		}
	}

aerospike:

	if resource == "aerospike"  {
		for {
			p := color.Set + "gospike> " + color.Reset
			cmd := readline.Readline(&p)

			if *cmd == "exit" || *cmd == "q" {
				break
			}
			if strings.HasPrefix(*cmd, "memcached")  {
				resource = "memcached"
				goto memcached
			}
			if strings.HasPrefix(*cmd, "elasticsearch")  {
				resource = "elasticsearch"
				goto elasticsearch
			}

			it, _ := StatsItems(*cmd, aeAddr)
			fmt.Printf("%s\n", it)
			readline.AddHistory(*cmd)
		}
	}

elasticsearch:
	if resource == "elasticsearch"  {
		for {
			p := color.Set + "gostic> " + color.Reset
			cmd := readline.Readline(&p)

			if *cmd == "exit" || *cmd == "q" {
				break
			}
			if strings.HasPrefix(*cmd, "memcached") {
				resource = "memcached"
				goto memcached
			}

			if strings.HasPrefix(*cmd, "aerospike") {
				resource = "aerospike"
				goto aerospike
			}
			it, _ := StatsItems(*cmd, elAddr)
			fmt.Printf("%s\n", it)
			readline.AddHistory(*cmd)
		}
	}

}

func PrintHelp() {
	fmt.Printf("Simple usage:\n type \"aerospike\" to use aerospike storage\n type \"memcached\" to use memcached-storage\n")
}

func main() {
	memAddr, aeAddr, elAddr := ReadConfig()
	var adr, c string
	color := Colors{"\033[36m", "\033[0m"}
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
		promt := color.Set + "> " + color.Reset
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
