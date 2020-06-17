package main

import (
	"strings"
	"testing"
)

func TestClientByPlainText(t *testing.T) {
	c, err := NewClient("pop3.lolipop.jp:143")
	if err != nil {
		t.Error(err.Error())
	}
	noStateCmdTest(c, t)
}

func TestClientByTLS(t *testing.T) {
	c, err := NewClient("pop3.lolipop.jp:993")
	if err != nil {
		t.Error(err.Error())
	}
	noStateCmdTest(c, t)
}

func noStateCmdTest(c *Client, t *testing.T) {
	var table = []struct {
		cmd string
		req string
		res string
	}{
		{"NOOP", "NOOP", "OK NOOP completed"},
		{"CAPABILITY", "CAPABILITY", "* CAPABILITY IMAP4rev1"},
	}

	for _, tt := range table {
		if err := c.Send(tt.cmd); err != nil {
			t.Error(err.Error())
		}
		sent := c.LastSent
		if !strings.Contains(sent, tt.req) {
			t.Errorf("Sent(%s): %s", tt.cmd, sent)
		}
		received := c.LastReceived
		if !strings.Contains(received, tt.res) {
			t.Errorf("Received(%s): %s", tt.cmd, received)
		}
	}
}

func TestTag(t *testing.T) {
	tag := NewTag(6)
	result := tag.Gen()
	if len(result) != 7 {
		t.Errorf("Tag length expects 7 but %d", len(result))
	}
	result2 := tag.Gen()
	if result == result2 {
		t.Errorf("Use the same tag twice: %s", result)
	}
}
