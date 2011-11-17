package main

import "os"
import "fmt"
import "time"

var GAMEID = fmt.Sprintf("%v-%v", time.UTC().Format("2006-01-02T15:04:05Z"), os.Getpid())

func NewLog(name string, ext string) *os.File {
    dirname := fmt.Sprintf("../games/%v", GAMEID)
    basename := fmt.Sprintf("%v.%v", name, ext)
    filename := dirname + "/" + basename
    os.MkdirAll(dirname, 0700)
    file, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0600)
    if err != nil {
        panic("Could not create log file " + filename)
    }
    return file
}

func NewTurnLog(name string, ext string) *os.File {
    dirname := fmt.Sprintf("../games/%v/%v", GAMEID, name)
    basename := fmt.Sprintf("%04v.%v", turn, ext)
    filename := dirname + "/" + basename
    os.MkdirAll(dirname, 0700)
    file, err := os.Create(filename)
    if err != nil {
        panic("Could not create log file " + filename)
    }
    return file
}
