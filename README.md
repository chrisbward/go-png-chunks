# go-png-chunks
Library/Utility to manage tEXt chunks inside PNG files


```go

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
	tEXtChunkToWrite := TEXtChunk{
		Key:   "helloworld",
		Value: helloWorld,
	}
	w, err := TestWritePNG(data, tEXtChunkToWrite)
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