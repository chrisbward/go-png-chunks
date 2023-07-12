# go-png-chunks
Library/Utility to manage tEXt chunks inside PNG files


```go
package main

import (
    gopngchunks "github.com/chrisbward/go-png-chunks"
)

var helloWorld = "aGVsbG8gd29ybGQ="


func WritetEXtChunkToFile(inputFilePath string, outputFilePath string) error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("os.Open(): %s", err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll(): %s", err)
	}
	tEXtChunkToWrite := gopngchunks.TEXtChunk{
		Key:   "helloworld",
		Value: helloWorld,
	}
	w, err := gopngchunks.WritetEXtToPngBytes(data, tEXtChunkToWrite)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	out.Write(w.Bytes())

	return nil
}


func main(){
    WritetEXtChunkToFile("./in.png", "./out.png")
}

```