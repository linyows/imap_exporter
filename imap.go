package main

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type Cmd struct {
	req []string
	res []string
}

type IMAP struct {
	addr  string
	count int
	conn  net.Conn
}

func NewIMAP(addr string) (*IMAP, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
	if err != nil {
		return nil, err
	}

	if port == "993" {
		t := tls.Client(conn, &tls.Config{ServerName: host})
		return &IMAP{addr: addr, count: 0, conn: t}, nil
	}

	return &IMAP{addr: addr, count: 0, conn: conn}, nil
}

func (i *IMAP) Cmd(name string, args ...string) (*Cmd, error) {
	tag := GenTag(name)
	cmd := &Cmd{req: append([]string{tag, name}, args...)}
	_, err := i.conn.Write([]byte(strings.Join(cmd.req, " ") + "\n"))
	if err != nil {
		return cmd, err
	}
	reader := bufio.NewReader(i.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return cmd, err
		}
		cmd.res = append(cmd.res, msg)
		if strings.HasPrefix(msg, tag) {
			break
		}
	}
	return cmd, nil
}

func (i *IMAP) Close() {
	i.conn.Close()
}

var count = 0

func GenTag(k string) string {
	count += 1
	return string([]rune(k)[:2]) + strconv.Itoa(count)
}
