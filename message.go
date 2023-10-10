package main

import (
	"strconv"
	"strings"
)

type Message struct {
	Subject string
	Text    string
	Count   int
	Filter  *Filter
}

func (msg *Message) BuildSubject() {
	msg.Subject = strings.Replace(msg.Filter.SubjectFormat, "%name", msg.Filter.Name, -1)
}

func (msg *Message) BuildText() {
	text := strings.Replace(msg.Filter.TextFormat, "%name", msg.Filter.Name, -1)
	text = strings.Replace(text, "%count", strconv.Itoa(msg.Count), -1)
	msg.Text = strings.Replace(text, "%text", msg.Text, -1)
}
