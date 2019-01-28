package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/lokicui/sego"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	flag.Parse()
	fname := os.Args[1]
	log.SetOutput(os.Stdout)
	fd, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	scanner.Buffer([]byte{}, bufio.MaxScanTokenSize*100)
	tokens := make(chan struct{}, 30)
	waitgroup := new(sync.WaitGroup)
	for scanner.Scan() {
		waitgroup.Add(1)
		tokens <- struct{}{}
		line := scanner.Text()
		go func() {
			line_items := strings.Split(line, "\t")
			query := strings.Trim(line_items[0], " ")
			items, err := sego.SegmentQuery(query, false)
			if err != nil {
				items, err = sego.SegmentQuery(query, false)
			}
			if err != nil {
				items, err = sego.SegmentQuery(query, false)
			}
			if err != nil {
				items, err = sego.SegmentQuery(query, false)
			}
			if err != nil {
				items, err = sego.SegmentQuery(query, false)
			}
			if err != nil {
				items, err = sego.SegmentQuery(query, false)
			}
			pieces := []string{query}
			if err == nil {
				words := []string{}
				for _, item := range items {
					str := fmt.Sprintf("%s:%d:%d", sego.SBC2DBC(item.Word), item.Term_NImps, item.Weight)
					words = append(words, str)
				}
				pieces = append(pieces, strings.Join(words, " "))
				pieces = append(pieces, line_items[1:]...)
				log.Println(strings.Join(pieces, "\t"))
			}
			waitgroup.Done()
			<-tokens
		}()
	}
	waitgroup.Wait()
}
