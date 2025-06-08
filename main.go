package main

import "fmt"

func main() {
    _, err := createClient()
    if err != nil {
        fmt.Println("Error: ", err)
        return
    }



}
