package main

import (
	"bufio"
	crand "crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/seehuhn/mt19937"
)

type Tag struct {
	id  []byte
	seq int
}

type Client struct {
	addr         string
	count        int
	conn         net.Conn
	tag          *Tag
	LastSent     string
	LastReceived string
}

func NewClient(addr string) (*Client, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
	if err != nil {
		return nil, err
	}

	c := &Client{
		addr:  addr,
		count: 0,
		conn:  conn,
		tag:   NewTag(6),
	}

	if port == "993" {
		t := tls.Client(conn, &tls.Config{ServerName: host})
		c.conn = t
	}

	reader := bufio.NewReader(c.conn)
	_, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Send(name string, args ...string) error {
	if name == "" {
		return fmt.Errorf("IMAP command is empty")
	}

	tag := c.tag.Gen()
	var received []string
	command := strings.Join(append([]string{tag, name}, args...), " ") + "\n"
	c.LastSent = command
	_, err := c.conn.Write([]byte(command))
	if err != nil {
		return err
	}

	upper := strings.ToUpper(name)
	reader := bufio.NewReader(c.conn)

	for {
		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			c.LastReceived = strings.Join(received, " ")
			return err
		}
		received = append(received, msg)
		if strings.HasPrefix(msg, tag) {
			break
		}
		if upper == "CAPABILITY" || upper == "LOGOUT" {
			break
		}
	}

	c.LastReceived = strings.Join(received, " ")
	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func NewTag(length int) *Tag {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rng := rand.New(mt19937.New())
	rng.Seed(seed.Int64())
	id := make([]byte, length, length)
	for i, v := range rng.Perm(26)[:length] {
		id[i] = 'A' + byte(v)
	}
	return &Tag{
		id:  id,
		seq: 0,
	}
}

func (t *Tag) Gen() string {
	t.seq++
	return fmt.Sprintf("%s%d", t.id, t.seq)
}
