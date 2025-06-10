package main

import (
    "fmt"
    "sync"
)

func main() {
    pgnFilesList, err := getPgns()
    if err != nil {
        fmt.Println("error in getting pgnfile: %w", err)
        return 
    }
    chunksize := 10
    var wg *sync.WaitGroup
    for _, pgnFile := range pgnFilesList {
        gamesList, err := parsePgnFile(pgnFile)
        if err != nil {
            fmt.Println("error in parsing pgnfile (inside loop): %w", err)
            return
        }
        chunks := chunker(gamesList, chunksize)
        wg.Add(1)
        go func(chunks [][]Pgn) {
            defer wg.Done()
            doChunks(chunks)
        } (chunks)
    }
    wg.Wait()
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

func doChunk(chunk []Pgn, wg *sync.WaitGroup) {
    defer wg.Done()
    operableChunk := pgnsToGames(chunk)
    insertGames(operableChunk)
}

func doChunks(chunks [][]Pgn) {
    var wg sync.WaitGroup

    for _, chunk := range chunks {
        wg.Add(1)
        go doChunk(chunk, &wg)
    }
    wg.Wait()
}