package main

import (
	"sync"
	"net"
	"bufio"
	"log"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	hostname = "localhost"
	port = "7070"
	defaultDir = "."
)

var wg sync.WaitGroup

var CRLF = []byte("\n")



func printFile(workingDirectory string, writer *bufio.Writer, path string) {
	bytes, err := ioutil.ReadFile(workingDirectory+path)
	if err != nil{
		log.Fatal(err)
		return
	}
	writer.Write(bytes)
}


func listDirectory(writer *bufio.Writer, dir *directory) {
	dir.iterate(func(name string, locType string) {
		resource := dir.path+"/"+name
		if dir.path == "/" {
			resource = dir.path[1:] + "/" + name
		}
		directoryString := fmt.Sprintf("%s%s\t%s\t%s\t%s\r\n", locType, name, resource, hostname, port)
		writer.Write([]byte(directoryString))
	})
}

func handleConnection(conn net.Conn) {

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	defer conn.Close()
	defer wg.Done()

	t, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	providedPath := t[:len(t)-2]
	providedDirectory := GetDirectoryAtPath(defaultDir, providedPath)


	if providedDirectory == nil {
		printFile(defaultDir, writer, providedPath)
	} else {
		index, err := providedDirectory.getIndex(defaultDir)
		if err == nil {
			if index[len(index)-1] != '\n' {
				index = append(index, CRLF...)
			}
			writer.Write(index)
		}

		listDirectory(writer, providedDirectory)
	}

	writer.Flush()
}

func main() {

	argsError := "Your arguments should look something like: <working dir> <hostname> <port>"

	if len(os.Args) != 4 {
		log.Fatal(argsError)
	}

	if os.Args[1] != "" {
		hostname = os.Args[1]
	} else {
		log.Fatal(argsError)
	}

	if os.Args[2] != "" {
		hostname = os.Args[2]
	} else {
		log.Fatal(argsError)
	}

	if os.Args[3] != "" {
		port = os.Args[3]
	} else {
		log.Fatal(argsError)
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			break
		}
		wg.Add(1)
		go handleConnection(conn)
	}
	wg.Wait()

}

