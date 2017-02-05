package main

import (
	"bufio"
	"container/list"
	"os"
	"strconv"
)

// Debugger is the debugger struct
type Debugger struct {
	messages   *list.List
	maxSize    int
	fileWriter *bufio.Writer
	file       *os.File
}

// Anything that wants to be printed
type Debuggable interface {
	ToString() string
}

// Debug with a list of messages
func (d *Debugger) Debug(messages ...interface{}) {
	finalMessage := ""
	for _, message := range messages {
		switch m := message.(type) {
		case string:
			finalMessage += m
		case int:
			finalMessage += strconv.Itoa(m)
		case Debuggable:
			finalMessage += m.ToString()
		}
	}
	d.messages.PushFront(finalMessage)
	if d.messages.Len() > d.maxSize {
		d.messages.Remove(d.messages.Back())
	}
	if d.fileWriter != nil {
		d.fileWriter.WriteString(finalMessage + "\n")
		d.fileWriter.Flush()
	}
}

// Close the debugger
func (d *Debugger) Close() {
	d.file.Close()
}

var debugger *Debugger

func debug(messages ...interface{}) {
	if debugger != nil {
		debugger.Debug(messages...)
	}
}

// InitDebug init the debugger and return the Debugger
func InitDebug() (*Debugger, error) {
	debugger = &Debugger{}
	debugger.messages = list.New()
	file, err := os.OpenFile("debug.log", os.O_APPEND|os.O_WRONLY, 644)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create("debug.log")
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	debugger.fileWriter = bufio.NewWriter(file)
	debugger.file = file

	return debugger, nil
}
