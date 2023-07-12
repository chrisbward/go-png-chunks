package gopngchunks

// Helper library to retreive and implant tEXt chunks in to PNG files

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/murkland/pngchunks"
)

// 89 50 4E 47 0D 0A 1A 0A
var PNGHeader = "\x89\x50\x4E\x47\x0D\x0A\x1A\x0A"
var tEXtChunkType = "tEXt"
var IENDChunkType = "IEND"
var NULLBYTE = "\x00"
var tEXtChunkDataSpecification = "%s" + NULLBYTE + "%s"

type TEXtChunk struct {
	Key   string
	Value string
}

func ContainsPNGMagicBytesHeader(data []byte) bool {
	return string(data) == PNGHeader
}

func WritetEXtToPngBytes(inputBytes []byte, TEXtChunkToWrite TEXtChunk) (outputBytes bytes.Buffer, err error) {

	isPng := ContainsPNGMagicBytesHeader(inputBytes[:8])
	if !isPng {
		return outputBytes, fmt.Errorf("ContainsPNGMagicBytesHeader(): %s", "Not a PNG file")
	}

	reader := bytes.NewReader(inputBytes)
	pngr, err := pngchunks.NewReader(reader)
	if err != nil {
		return outputBytes, fmt.Errorf("NewReader(): %s", err)
	}

	pngw, err := pngchunks.NewWriter(&outputBytes)
	if err != nil {
		return outputBytes, fmt.Errorf("NewWriter(): %s", err)
	}

	for {
		chunk, err := pngr.NextChunk()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return outputBytes, fmt.Errorf("NextChunk(): %s", err)
		}

		if chunk.Type() != tEXtChunkType {

			// IENDChunkType will only appear on the final iteration of a valid PNG
			if chunk.Type() == IENDChunkType {
				// This is where we inject tEXtChunkType as the penultimate chunk with the new value
				newtEXtChunk := []byte(fmt.Sprintf(tEXtChunkDataSpecification, TEXtChunkToWrite.Key, TEXtChunkToWrite.Value))
				if err := pngw.WriteChunk(int32(len(newtEXtChunk)), tEXtChunkType, bytes.NewBuffer(newtEXtChunk)); err != nil {
					return outputBytes, fmt.Errorf("WriteChunk(): %s", err)
				}

				newComment := []byte("encodedwith\x00github.com/chrisbward/go-png-chunks")
				if err := pngw.WriteChunk(int32(len(newComment)), chunk.Type(), bytes.NewBuffer(newComment)); err != nil {
					return outputBytes, fmt.Errorf("WriteChunk(): %s", err)
				}
				// Now we end the buffer with IENDChunkType chunk
				if err := pngw.WriteChunk(chunk.Length(), chunk.Type(), chunk); err != nil {
					return outputBytes, fmt.Errorf("WriteChunk(): %s", err)
				}
			} else {
				// writes back original chunk to buffer
				if err := pngw.WriteChunk(chunk.Length(), chunk.Type(), chunk); err != nil {
					return outputBytes, fmt.Errorf("WriteChunk(): %s", err)
				}
			}
		} else {
			if _, err := io.Copy(ioutil.Discard, chunk); err != nil {
				return outputBytes, fmt.Errorf("io.Copy(ioutil.Discard, chunk): %s", err)
			}
		}

		if err := chunk.Close(); err != nil {
			return outputBytes, fmt.Errorf("chunk.Close(): %s", err)
		}
	}
	return outputBytes, nil
}

func GetAlltEXtChunks(inputBytes []byte) (textChunks []TEXtChunk, err error) {

	reader := bytes.NewReader(inputBytes)
	pngr, err := pngchunks.NewReader(reader)
	if err != nil {
		// t.Errorf("NewReader(): %s", err)
		return textChunks, fmt.Errorf("pngchunks.NewReader(): %s", err)
	}

	for {
		chunk, err := pngr.NextChunk()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}

		if chunk.Type() != tEXtChunkType {
			if _, err := io.Copy(ioutil.Discard, chunk); err != nil {
				// t.Errorf("io.Copy(): %s", err)
				return textChunks, fmt.Errorf("io.Copy(): %s", err)
			}
		} else {
			buf, err := ioutil.ReadAll(chunk)
			if err != nil {
				return textChunks, fmt.Errorf("ioutil.ReadAll(): %s", err)
			}
			dataInChunk := string(buf)
			values := strings.Split(dataInChunk, "\x00")
			if len(values) == 2 {
				textChunks = append(textChunks, TEXtChunk{Key: values[0], Value: values[1]})
			}
		}

		if err := chunk.Close(); err != nil {
			return textChunks, fmt.Errorf("chunk.Close(): %s", err)
		}
	}

	return textChunks, nil
}
