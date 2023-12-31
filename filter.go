package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Filter struct {
	Name          string
	LineReg       *regexp.Regexp
	ExceptRegs    []*regexp.Regexp
	TextFormat    string
	SubjectFormat string
	Notifiers     []Notifier
}

func NewFilter(cfg FilterConfig, hostname string, notifiers []Notifier) (*Filter, error) {
	lineReg, err := regexp.Compile(cfg.Pattern)
	if err != nil {
		return nil, fmt.Errorf("LogFile filter %s pattern compile error: %v", cfg.Name, err)
	}

	exceptRegs := make([]*regexp.Regexp, 0, len(cfg.Exceptions))

	for _, exStr := range cfg.Exceptions {
		exReg, err := regexp.Compile(exStr)
		if err != nil {
			return nil, fmt.Errorf("LogFile filter %s exception pattern %s compile error: %v", cfg.Name, exStr, err)
		}
		exceptRegs = append(exceptRegs, exReg)
	}

	cfg.Notifications = removeDuplicates(cfg.Notifications)

	f := &Filter{
		Name:          cfg.Name,
		LineReg:       lineReg,
		ExceptRegs:    exceptRegs,
		TextFormat:    strings.Replace(cfg.Message, "%hostname", hostname, -1),
		SubjectFormat: strings.Replace(cfg.Subject, "%hostname", hostname, -1),
		Notifiers:     make([]Notifier, 0, len(cfg.Notifications)),
	}

	for _, notifName := range cfg.Notifications {
		for _, notif := range notifiers {
			if notif.Name() == notifName {
				f.Notifiers = append(f.Notifiers, notif)
			}
		}
	}

	return f, nil
}

func (f *Filter) Match(str string) bool {
	if !f.LineReg.MatchString(str) {
		return false
	}

	for _, exReg := range f.ExceptRegs {
		if exReg.MatchString(str) {
			return false
		}
	}

	return true
}
