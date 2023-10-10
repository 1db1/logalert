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
	Notifications []string
}

func NewFilter(cfg FilterConfig, hostname string) (*Filter, error) {
	lineReg, err := regexp.Compile(cfg.Pattern)
	if err != nil {
		return nil, fmt.Errorf("LogFile filter %s pattern compile error: %v", cfg.Name, err)
	}

	exceptRegs := make([]*regexp.Regexp, len(cfg.Exceptions))

	for _, exStr := range cfg.Exceptions {
		exRe, err := regexp.Compile(exStr)
		if err != nil {
			return nil, fmt.Errorf("LogFile filter %s exception pattern %s compile error: %v", cfg.Name, exStr, err)
		}
		exceptRegs = append(exceptRegs, exRe)
	}

	return &Filter{
		Name:          cfg.Name,
		LineReg:       lineReg,
		ExceptRegs:    exceptRegs,
		TextFormat:    strings.Replace(cfg.Message, "%hostname", hostname, -1),
		SubjectFormat: strings.Replace(cfg.Subject, "%hostname", hostname, -1),
		Notifications: removeDuplicates(cfg.Notifications),
	}, nil
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
