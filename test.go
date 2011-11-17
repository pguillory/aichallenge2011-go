package main

import "time"
import "fmt"

func main() {
    fmt.Println(time.UTC().Format("2006-01-02T15:04:05Z"))
    
    slice := make([]int, 1, 3)
    slice[0] = 2
    //slice = slice[:len(slice) + 1]
    //slice = slice[:len(slice) + 1]
    //slice[1] = 4
    slice = append(slice, 4)
    slice = append(slice, 6)
    slice = append(slice, 8)
    fmt.Println(slice)
}
