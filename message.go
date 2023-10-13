package main

import (
	"strconv"
	"strings"
)

type Message struct {
	FileName string
	Subject  string
	Text     string
	Count    int
	Filter   *Filter
}

func (msg *Message) BuildSubject() {
	msg.Subject = strings.Replace(msg.Filter.SubjectFormat, "%filename", msg.FileName, -1)
	msg.Subject = strings.Replace(msg.Subject, "%filtername", msg.Filter.Name, -1)
	msg.Subject = strings.Replace(msg.Subject, "%count", strconv.Itoa(msg.Count), -1)
}

func (msg *Message) BuildText() {
	msg.Text = strings.Replace(msg.Filter.TextFormat, "%filename", msg.FileName, -1)
	msg.Text = strings.Replace(msg.Text, "%filtername", msg.Filter.Name, -1)
	msg.Text = strings.Replace(msg.Text, "%count", strconv.Itoa(msg.Count), -1)
	msg.Text = strings.Replace(msg.Text, "%text", msg.Text, -1)
}
