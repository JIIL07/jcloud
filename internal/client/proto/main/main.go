package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	protobuf "github.com/JIIL07/jcloud/internal/client/proto"
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
)

func main() {
	// Открываем файл
	f, err := os.Open("server.exe")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer f.Close()

	// Создаем буфер
	var buffer bytes.Buffer

	// Читаем файл в буфер
	n, err := io.Copy(&buffer, f)
	if err != nil {
		log.Fatal("Error copying file to buffer:", err)
	}

	// Проверяем количество прочитанных байт
	fmt.Printf("Copied %d bytes into buffer\n", n)

	// Если буфер пуст, выдаем предупреждение
	if buffer.Len() == 0 {
		log.Fatal("Buffer is empty after reading file")
	}

	var b bytes.Buffer
	gzipWriter := gzip.NewWriter(&b)
	gzipWriter.Write(buffer.Bytes())
	gzipWriter.Close()

	// Создаем структуру protobuf.File
	file := &protobuf.File{
		Id: 1,
		Metadata: &protobuf.FileMetadata{
			Filename:  f.Name(),
			Extension: "exe",
			Filesize:  int32(buffer.Len()),
		},
		Status: "active",
		Data:   b.Bytes(),
	}

	// Сериализуем структуру в бинарный формат
	data, err := proto.Marshal(file)
	if err != nil {
		log.Fatal("Marshaling error:", err)
	}

	// Записываем бинарные данные в файл
	err = os.WriteFile("file.bin", data, 0644)
	if err != nil {
		log.Fatal("Error writing file:", err)
	}

	// Читаем бинарные данные из файла
	in, err := os.ReadFile("file.bin")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// Десериализуем данные обратно в структуру
	newFile := &protobuf.File{}
	err = proto.Unmarshal(in, newFile)
	if err != nil {
		log.Fatal("Unmarshaling error:", err)
	}

	// Выводим десериализованные данные
	fmt.Printf("Deserialized File: %+v\n", newFile)
}
