package gopngchunks

// Helper library to retreive and implant tEXt chunks in to PNG files

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/murkland/pngchunks"
)

// 89 50 4E 47 0D 0A 1A 0A
var PNGHeader = "\x89\x50\x4E\x47\x0D\x0A\x1A\x0A"
var tEXtChunkType = "tEXt"
var IENDChunkType = "IEND"
var NULLBYTE = "\x00"
var tEXtChunkDataSpecification = "%s" + NULLBYTE + "%s"

type Chunk struct {
	Length int    // chunk data length
	CType  string // chunk type
	Data   []byte // chunk data
	Crc32  []byte // CRC32 of chunk data
}

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
		// t.Errorf("NewWriter(): %s", err)
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

			newComment := []byte("comment\x00hi everyone!")
			if err := pngw.WriteChunk(int32(len(newComment)), chunk.Type(), bytes.NewBuffer(newComment)); err != nil {
				return outputBytes, fmt.Errorf("WriteChunk(): %s", err)
			}
		}

		if err := chunk.Close(); err != nil {
			return outputBytes, fmt.Errorf("chunk.Close(): %s", err)
		}
	}
	return outputBytes, nil
}
