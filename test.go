package main

import "time"
import "fmt"

func main() {
    fmt.Println(time.UTC().Format("2006-01-02T15:04:05Z"))
}
