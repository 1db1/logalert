package main

import (
	"strconv"
	"strings"
)

type Message struct {
	Subject string
	Text    string
	Count   int
	Match   *LogMatch
}

func (msg *Message) BuildSubject() {
	msg.Subject = strings.Replace(msg.Match.SubjectFormat, "%name", msg.Match.Name, -1)
}

func (msg *Message) BuildText() {
	text := strings.Replace(msg.Match.TextFormat, "%name", msg.Match.Name, -1)
	text = strings.Replace(text, "%count", strconv.Itoa(msg.Count), -1)
	msg.Text = strings.Replace(text, "%text", msg.Text, -1)
}
