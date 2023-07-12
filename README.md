# go-png-chunks
Library/Utility to manage tEXt chunks inside PNG files

## Example for reading all tEXt chunks from file

```go
import (
	"fmt"
	"io/ioutil"
	"os"
	gopngchunks "github.com/chrisbward/go-png-chunks"
)

func ReadtEXtChunksFromFile(inputFilePath string) (tEXtChunks []gopngchunks.TEXtChunk, err error) {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return tEXtChunks, fmt.Errorf("os.Open(): %s", err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return tEXtChunks, fmt.Errorf("ioutil.ReadAll(): %s", err)
	}
	tEXtChunks, err = gopngchunks.GetAlltEXtChunks(data)
	if err != nil {
		panic(err)
	}
	return tEXtChunks, nil
}


```


## Example for writing to PNG
```go
import (
	"fmt"
	"io/ioutil"
	"os"

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

```