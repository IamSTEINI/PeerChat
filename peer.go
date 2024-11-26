package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var connections []net.Conn
var connMutex sync.Mutex
var myConn string

func main() {
	log.SetOutput(io.Discard)
	fmt.Print("PeerChat by DXBY")
	fmt.Print("Select an option:\n1: Host a node\n2: Join \n>")
	reader := bufio.NewReader(os.Stdin)
	option, _ := reader.ReadString('\n')
	if strings.TrimSpace(option) == "1" {
		go listenAndAccept()
	} else if strings.TrimSpace(option) == "2" {
		fmt.Print("Enter the node's IP:PORT : ")
		ipin, _ := reader.ReadString('\n')
		if establishConnection(ipin) {
			fmt.Print("Creating own node for others...\n")
			go listenAndAccept()
		} else {
			fmt.Println("\n[!] Failed to connect. Please try again.")
		}
	} else {
		fmt.Println("Invalid option. Try again.")
	}
	go chat()
	go sendPeersPeriodically()
	select {}
}

func chat() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		message, _ := reader.ReadString('\n')
		if strings.TrimSpace(message) == "exit" {
			break
		}
		sendAll(message)
	}
}

var lastMessages = make(map[string]string)

func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		connMutex.Lock()
		connections = removeConnection(connections, conn)
		connMutex.Unlock()
		changeCmdTitle(fmt.Sprintf("Active peers: %d", len(connections)))
	}()
	go askForPeers(conn)
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection closed: %s", conn.RemoteAddr().String())
			return
		}
		message = strings.TrimSpace(message)
		log.Printf("[RAW RECEIVE] '%s' from %s", message, conn.RemoteAddr().String())
		if strings.HasPrefix(message, "GET_") {
			switch message {
			case "GET_PEERS":
				peers := append(getMyPeers(), myConn)
				log.Printf("[PEER REQUEST] Sharing peers (%s)", strings.Join(peers, ";"))
				send(conn, "NEW_PEERS_"+strings.Join(peers, ";"))
			default:
				// send(conn, "ERROR_unknown request '"+message+"'")
			}
		} else if strings.HasPrefix(message, "ERROR_") {
			log.Println(strings.Split(message, "_")[1])

		} else if strings.HasPrefix(message, "NEW_PEERS_") {
			peers := strings.Split(message, "_")[2]
			peersArray := strings.Split(peers, ";")
			newPeers := make([]string, 0)
			for _, peer := range peersArray {
				if peer != conn.RemoteAddr().String() && peer != myConn && !containsConnectionStr(connections, peer) {
					newPeers = append(newPeers, peer)
				}
			}
			log.Printf("[NEW PEERS] %v", newPeers)
			for _, peer := range newPeers {
				go func(p string) {
					establishConnection(p)
				}(peer)
			}
		} else {
			if strings.TrimSpace(message) != "" {
				if lastMessages[conn.RemoteAddr().String()] == message {
					continue
				}
				lastMessages[conn.RemoteAddr().String()] = message
				fmt.Printf("[%s] %s\n", conn.RemoteAddr().String(), message)
			}
		}
	}
}

func removeConnection(connections []net.Conn, conn net.Conn) []net.Conn {
	for i, c := range connections {
		if c == conn {
			newList := append(connections[:i], connections[i+1:]...)
			log.Printf("[UPDATED PEERS] Active peers: %v", getMyPeers())
			return newList
		}
	}
	return connections
}

func getMyPeers() []string {
	var peers []string
	for _, conn := range connections {
		peers = append(peers, conn.RemoteAddr().String())
	}
	return peers
}

func send(conn net.Conn, data string) error {
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(data + "\n")
	if err != nil {
		return err
	}
	return writer.Flush()
}

func sendAll(data string) error {
	for _, conn := range connections {
		writer := bufio.NewWriter(conn)
		_, err := writer.WriteString(data + "\n")
		if err != nil {
			return err
		}
		if err := writer.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func askForPeers(conn net.Conn) {
	err := send(conn, "GET_PEERS")
	if err != nil {
		log.Printf("Error sending request to peers: %s. Closing connection: %s", err, conn.RemoteAddr().String())
		conn.Close()
		return
	}
	log.Printf("[~] Sent 'GET_PEERS' -> %s", conn.RemoteAddr().String())
}

func establishConnection(connString string) bool {
	connString = strings.TrimSpace(connString)
	log.Printf("[~] Connecting to %s", connString)
	conn, err := net.Dial("tcp", connString)
	if err != nil {
		log.Printf("[ERROR] Cannot connect to %s", connString)
		return false
	}
	log.Printf("\033[2K\r[+] Connected to %s\n", connString)
	if !containsConnection(connections, conn) {
		connMutex.Lock()
		connections = append(connections, conn)
		connMutex.Unlock()
		changeCmdTitle(fmt.Sprintf("Active peers: %d", len(connections)))
	}
	go handleConnection(conn)
	return true
}

func changeCmdTitle(title string) {
	go func() {
		fmt.Printf("\x1b]2;%s\x07", title)
	}()
}

func listenAndAccept() {
	ln, err := net.Listen("tcp", getLocalIP()+":0")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("[+] Node started")
	tcpAddr := ln.Addr().(*net.TCPAddr)
	log.Printf("[~] Listening on port: %d", tcpAddr.Port)
	myConn = fmt.Sprintf("%s:%d", getLocalIP(), tcpAddr.Port)
	fmt.Printf("\nSHARE: %s\n___________________________________________\n", myConn)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		if !containsConnection(connections, conn) {
			connMutex.Lock()
			connections = append(connections, conn)
			connMutex.Unlock()
			changeCmdTitle(fmt.Sprintf("Active peers: %d", len(connections)))
		}
		log.Printf("[+] NEW connection #%d from %s\n", len(connections), conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}

func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip.To4() != nil {
				return ip.String()
			}
		}
	}
	return ""
}

func containsConnection(connections []net.Conn, conn net.Conn) bool {
	for _, c := range connections {
		if c == conn {
			return true
		}
	}
	return false
}

func containsConnectionStr(connections []net.Conn, conn string) bool {
	connMutex.Lock()
	defer connMutex.Unlock()
	for _, c := range connections {
		if c.RemoteAddr().String() == conn {
			return true
		}
	}
	return false
}

func sendPeersPeriodically() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			peers := append(getMyPeers(), myConn)
			log.Printf("[PERIODIC PEERS] Sharing peers (%s)", strings.Join(peers, ";"))
			for _, conn := range connections {
				send(conn, "NEW_PEERS_"+strings.Join(peers, ";"))
			}
		}
	}
}
