package main

import (

)

func main() {
    chunksize := 10

}

func chunker(slice []Pgn, chunksize int) [][]Pgn {
    var chunks [][]Pgn
    for i := 0; i < len(slice); i += chunksize {
        end := i + chunksize
        if end > len(slice) {
            end = len(slice)
        }
        chunks = append(chunks, slice[i:end])
    }
    return chunks
}