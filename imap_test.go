package main

import (
	"strings"
	"testing"
)

func TestIMAPWithTXT(t *testing.T) {
	imap, err := NewIMAP("pop3.lolipop.jp:143")
	if err != nil {
		t.Error(err.Error())
	}

	var table = []struct {
		cmd string
		req string
		res string
	}{
		{"NOOP", "NO1 NOOP", "NO1 OK NOOP completed"},
		{"CAPABILITY", "CA2 CAPABILITY", "* CAPABILITY IMAP4rev1"},
	}

	for _, tt := range table {
		cmd, err := imap.Cmd(tt.cmd)
		if err != nil {
			t.Error(err.Error())
		}
		req := strings.Join(cmd.req, " ")
		if req != tt.req {
			t.Errorf("Request(%s): %s", tt.cmd, req)
		}
		res := strings.Join(cmd.res, " ")
		if !strings.HasPrefix(res, tt.res) {
			t.Errorf("Response(%s): %s", tt.cmd, res)
		}
	}
}

func TestIMAPWithTLS(t *testing.T) {
	imap, err := NewIMAP("pop3.lolipop.jp:993")
	if err != nil {
		t.Error(err.Error())
	}

	var table = []struct {
		cmd string
		req string
		res string
	}{
		{"NOOP", "NO3 NOOP", "NO3 OK NOOP completed"},
		{"CAPABILITY", "CA4 CAPABILITY", "* CAPABILITY IMAP4rev1"},
	}

	for _, tt := range table {
		cmd, err := imap.Cmd(tt.cmd)
		if err != nil {
			t.Error(err.Error())
		}
		req := strings.Join(cmd.req, " ")
		if req != tt.req {
			t.Errorf("Request(%s): %s", tt.cmd, req)
		}
		res := strings.Join(cmd.res, " ")
		if !strings.HasPrefix(res, tt.res) {
			t.Errorf("Response(%s): %s", tt.cmd, res)
		}
	}
}
