package handlers

import (
	"file-analyzer/internals/embeddings"
	"fmt"
	"io"
	"log"
	"net/http"
)

func FileHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error in File ")
	}
	const MAX_CHUNKS = 96
	fmt.Println(header.Filename)
	buff := make([]byte, 4096)
	var chunk string
	var chunkBuffer []string
	go func() {
		for {
			n, err := file.Read(buff)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			chunk += string(buff[:n])
			if len(chunk) >= 400 {
				chunkBuffer = append(chunkBuffer, chunk)
			}
		}
		len := len(chunkBuffer)
		for i:= 0; i < len; i += MAX_CHUNKS {
			end := min(i + MAX_CHUNKS, len)
			embeddings.ProcessChunksToEmbeddings(chunkBuffer[i:end])
		}
	}()
	// after it store in db (doc id , user id ,file meta info ) maybe

	// fmt.Println("Data ", string(data))
	if err != nil {
		fmt.Println("ERROR")
	}
}
