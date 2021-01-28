package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {

	if len(os.Args) == 3 {

		serverAddress := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])

		l, err := net.Listen("tcp", serverAddress)

		if err != nil {
			log.Fatal(err)
		}

		defer l.Close()

		fmt.Println("Listening on: " + serverAddress)
		fmt.Println()

		conn, err := l.Accept()

		if err != nil {
			sendMessageError(conn)
			log.Fatal(err)
		}

		defer conn.Close()

		fmt.Println("alguien ha ingresado")

		err = sendMessageEntry(conn)

		if err != nil {
			sendMessageError(conn)
			log.Fatal(err)
		}

		for {
			err = readCommand(conn)

			if err != nil {
				sendMessageError(conn)
				log.Fatal(err)
			}
		}
	}

}

func sendMessageEntry(conn net.Conn) (err error) {

	_, err = conn.Write([]byte("---Has ingresado---\n"))

	if err != nil {
		return
	}

	err = sendLineBreak(conn)

	if err != nil {
		return
	}
	return

}

func readCommand(conn net.Conn) (err error) {

	_, err = conn.Write([]byte("> "))

	if err != nil {
		return
	}

	reader := bufio.NewReader(conn)

	s, err := reader.ReadString('\n')

	if err != nil {
		return
	}

	s = strings.Replace(s, "\n", "", -1)

	comando := strings.Split(s, " ")

	if len(comando) > 1 {

		message, err := extractResultWithOneString(comando)

		if err != nil {
			return err
		}

		err = sendMessage(conn, message)

		if err != nil {
			return err
		}

	} else {

		message, err := extractResultWithStrings(comando)

		if err != nil {
			return err
		}

		err = sendMessage(conn, message)

		if err != nil {
			return err
		}

	}
	return
}

func extractResultWithOneString(comando []string) (message []byte, err error) {

	primer := comando[0]

	comando = append(comando[:0], comando[1:]...)

	cmd := exec.Command(primer, comando...)

	message, err = cmd.CombinedOutput()

	if err != nil {
		return
	}

	return
}

func extractResultWithStrings(comando []string) (message []byte, err error) {

	cmd := exec.Command(comando[0])

	message, err = cmd.CombinedOutput()

	if err != nil {
		return
	}

	return
}

func sendMessage(conn net.Conn, message []byte) (err error) {

	_, err = conn.Write(message)

	if err != nil {
		return
	}

	sendLineBreak(conn)

	return
}

func sendMessageError(conn net.Conn) {

	_, err := conn.Write([]byte("*Ocurrio un error*"))

	if err != nil {
		log.Fatal(err)
	}

}

func sendLineBreak(conn net.Conn) (err error) {

	_, err = conn.Write([]byte("\n"))

	if err != nil {
		return
	}
	return
}
