package main

import (
	"strconv"
	"strings"
)

type Message struct {
	Text    string
	Subject string
	Count   int
	Match   *LogMatch
}

func (msg *Message) FormatSubject() {
	msg.Subject = strings.Replace(msg.Match.SubjectFormat, "%name", msg.Match.Name, -1)
}

func (msg *Message) FormatText() {
	msg.Text = strings.Replace(msg.Match.TextFormat, "%name", msg.Match.Name, -1)
	msg.Text = strings.Replace(msg.Text, "%count", strconv.Itoa(msg.Count), -1)
	msg.Text = strings.Replace(msg.Text, "%text", msg.Match.Text, -1)
}
