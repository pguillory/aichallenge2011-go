package main

import "os"
import "fmt"
import "time"

var GAMEID = fmt.Sprintf("%v-%v", time.UTC().Format("2006-01-02T15:04:05Z"), os.Getpid())

type Log struct {
    name string
    ext string
}

func NewLog(name string, ext string) *Log {
    this := new(Log)
    this.name = name
    this.ext = ext
    return this
}

func (this *Log) File() *os.File {
    dirname := fmt.Sprintf("games/%v", GAMEID)
    basename := fmt.Sprintf("%04v.%v", this.name, this.ext)
    filename := dirname + "/" + basename
    os.MkdirAll(dirname, 0700)
    file, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0600)
    if err != nil {
        panic("Could not create log file " + filename)
    }
    return file
}

func (this *Log) TurnFile() *os.File {
    dirname := fmt.Sprintf("games/%v/%v", GAMEID, this.name)
    basename := fmt.Sprintf("%04v.%v", turn, this.ext)
    filename := dirname + "/" + basename
    os.MkdirAll(dirname, 0700)
    file, err := os.Create(filename)
    if err != nil {
        panic("Could not create log file " + filename)
    }
    return file
}
