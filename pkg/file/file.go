package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

// func main() {
// 	p := NewContent()
// 	file := "a.txt"
// 	w, err := NewFileWriter(file)
// 	if err != nil {
// 		panic(err)
// 	}
// 	w.write(p)
// }

func main() {
	file := "a.txt"
	readLine(file)
}

type ConcurWriter struct {
	mu     sync.Mutex
	target io.WriteCloser
}

func NewFileWriter(file string) (*ConcurWriter, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return &ConcurWriter{
		target: f,
	}, nil
}

func (cw *ConcurWriter) write(p []byte) error {
	defer cw.target.Close()
	_, err := cw.target.Write(p)
	return err
}

func NewContent() []byte {
	p := make([]byte, 0)
	for i := 0; i <= 100000; i++ {
		p = append(p, []byte(strconv.Itoa(i))...)
	}
	return p
}

func readLine(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		fmt.Println(len(line))
		if err == io.EOF {
			break
		}
	}
}
