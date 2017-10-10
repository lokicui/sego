package sego

import (
    //"fmt"
    "strings"
)

// 分词器结构体
type Segmenter struct {
}


func (seg *Segmenter) Segment(bytes []byte) []Segment {
    query := string(bytes)
    words, err := SegmentQuery(query, false)
    if err != nil {
    }
    segments := []Segment{}
    offset := 0
    for _, item := range words {
        word := SBC2DBC(item.Word) //半角
        start := strings.Index(query[offset:], word)
        end := offset
        if start != -1 {
            end = start + len(word)
        } else {
            start = strings.Index(query[offset:], item.Word)
            if start != -1 {
                end = start + len(item.Word)
            }
        }
        if start != -1 {
            //fmt.Printf("%d-%d-%d-%d\n", offset, start, end, len(word))
            text := []Text{}
            for _, i := range []rune(word) {
                text = append(text, Text(string(i)))
            }
            token := &Token{text:text}
            segment := Segment{start:offset+start, end:offset+end, token:token}
            offset += end
            //fmt.Printf("%d-%d-%s\n", offset, len(query), word)
            segments = append(segments, segment)
        }
    }
    return segments
}

func (seg *Segmenter) LoadDictionary(files string) {
}
