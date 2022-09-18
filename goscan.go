package main

import (
	"fmt"
	"moul.io/banner"
	"net"
	"strconv"
	"time"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	print_banner()
	if len(os.Args) <= 1 {
		fmt.Println("Usage:")
		fmt.Println("goscan IP")
		return
	}
	ip := os.Args[1]
	fmt.Println(ip)
	number_of_ports := 10000
	ports := make(chan int, 65535)
	open_ports_channel := make(chan int, 1000)
	var open_ports []int
	var wg sync.WaitGroup
	
	
	//1000 workers
	
	for i:= 1; i <= 1000; i++ {
		wg.Add(1)
		
		go func() {
			defer wg.Done()
			scanner(i, ip, ports, open_ports_channel)
		}()
	}
	
	//Sending jobs
	
	for i:=1; i<= number_of_ports; i++ {
		ports <- i
	}
	close(ports)
	wg.Wait()
	// Will hang if not close
	close(open_ports_channel)
	//Getting open ports
	for open_port := range open_ports_channel {
		open_ports = append(open_ports, open_port)
	}	
	//fmt.Println(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(open_ports)), ","), "[]"))
	nmap(ip, open_ports)
}


func scanner(id int, ip string, ports <- chan int, open_ports chan int) {
	for port := range ports {
		address := ip + ":" + strconv.Itoa(port)
		for i:=1; i<=2; i++ {
			_, err := net.DialTimeout("tcp", address, time.Millisecond * 200)
			if err == nil {
				fmt.Println("Open ", port)
				open_ports <- port
				break
			}
		}
	}
}

func nmap(ip string, open_ports []int) {

	ports := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(open_ports)), ","), "[]")
	command := "nmap " + ip + " -p" + ports + " -sC -sV -T4 -oN " + ip + ".nmap"
	fmt.Println("Running Command:", command)
	out, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Command Execute Success!")
	output := string(out[:])
	fmt.Println(output)
}

func print_banner(){
	fmt.Println(banner.Inline("goGun"))
}
