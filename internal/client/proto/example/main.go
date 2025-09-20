package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/JIIL07/jcloud/api/protobuf"

	"google.golang.org/protobuf/proto"
)

func main() {
	f, err := os.Open("go.mod")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer f.Close()

	var buffer bytes.Buffer
	n, err := io.Copy(&buffer, f)
	if err != nil {
		log.Fatal("Error copying file to buffer:", err)
	}
	fmt.Printf("Copied %d bytes into buffer\n", n)

	if buffer.Len() == 0 {
		log.Fatal("Buffer is empty after reading file")
	}

	var compressedBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBuffer)
	_, err = gzipWriter.Write(buffer.Bytes())
	if err != nil {
		log.Fatal("Error compressing data:", err)
	}
	gzipWriter.Close()

	file := &protobuf.File{
		Id: 1,
		Metadata: &protobuf.FileMetadata{
			Filename:  strings.Split(f.Name(), ".")[0],
			Extension: strings.Split(f.Name(), ".")[1],
			Filesize:  int64(compressedBuffer.Len()),
		},
		Status: "upload",
		Data:   compressedBuffer.Bytes(),
	}

	data, err := proto.Marshal(file)
	if err != nil {
		log.Fatal("Marshaling error:", err)
	}

	err = os.WriteFile("file.bin", data, 0600)
	if err != nil {
		log.Fatal("Error writing file:", err)
	}

	in, err := os.ReadFile("file.bin")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	newFile := &protobuf.File{}
	err = proto.Unmarshal(in, newFile)
	if err != nil {
		log.Fatal("Unmarshaling error:", err)
	}

	var decompressedBuffer bytes.Buffer
	gzipReader, err := gzip.NewReader(bytes.NewReader(newFile.Data))
	if err != nil {
		log.Fatal("Error creating gzip reader:", err)
	}

	maxDecompressedSize := 10 * 1024 * 1024 * 1024 * 1024
	limitedReader := io.LimitReader(gzipReader, int64(maxDecompressedSize))

	n, err = io.Copy(&decompressedBuffer, limitedReader)
	if err != nil {
		return
	}
	if n == int64(maxDecompressedSize) {
		return
	}
	gzipReader.Close()

	fmt.Printf("Deserialized File: %+v\n", newFile)
	fmt.Printf("File Content: %s\n", decompressedBuffer.String())
}
