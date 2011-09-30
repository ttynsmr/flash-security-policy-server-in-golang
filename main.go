package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"net"
//	"io/ioutil"
	"flag"
	"strconv"
)

var (
	socketpolicy []byte
	port int = 843
	socketpolicyfile string = "socketpolicy.xml"
)

func signal_arm () {
	fmt.Println("signal armed")

	go func(){
		for {
			sig := (<-signal.Incoming).(signal.UnixSignal)
			
			fmt.Println("%s", sig.String())
			
			switch sig {
			case syscall.SIGINT:
				os.Exit(0)
			default:
				//os.Exit(-1)
			}
		}
	}()
}

func loadPolicyFile() bool {
	file, err := os.Open( socketpolicyfile, os.O_RDONLY, 0 );
	if err == nil {
		stat, _ := file.Stat()
		socketpolicy = make([]byte, stat.Size + 1)
		file.Read(socketpolicy)
		socketpolicy[stat.Size] = 0
		file.Close()
		file = nil
	} else {
		fmt.Println(">>>> socket policy file open error:", err.String())
		return false
	}
	return true
}

func main() {
	signal_arm()

	// Parse args
	flag.IntVar(&port, "port", port, "listen port number")
	flag.StringVar(&socketpolicyfile, "file", socketpolicyfile, "socket policy file name")
	flag.Parse()

	fmt.Println("config listen port   =", port)
	fmt.Println("config socket policy =", socketpolicyfile)

	if loadPolicyFile() == false {
		return
	}

	accepts()
}

func accepts() {
	strPort := ":" + strconv.Itoa(port)
	fmt.Println("open port", strPort)
	addr, err := net.ResolveTCPAddr(strPort)
	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		fmt.Println(">>>> listen failed: ", err.String())
		return
	}
	addr = nil
	for {
		fmt.Println("Accept ready")
		session, err := l.AcceptTCP()
		if err != nil {
			//return
			fmt.Println("Accept error:", err.String())
			continue
		}

		go session_process(session)
		session = nil
	}

	l = nil
}

func session_process(s *net.TCPConn) {
	fmt.Println("Accepted session start")
	if( recieve_request(s) ) {
		send_response(s)
	}
	s.Close()
	s = nil
	fmt.Println("session done")
}

func recieve_request(s *net.TCPConn) bool {
	// 超適当なのでリクエストの受信完了は待たない
	/*
	_, err := ioutil.ReadAll(s)
	if( err == os.EOF ) {
		// write request validation here!
		return true
	}
	return  false
	 */

	return true
}

func send_response(s *net.TCPConn) {
	s.Write(socketpolicy)
}
