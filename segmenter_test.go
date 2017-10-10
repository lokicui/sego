package sego

import (
    "fmt"
    "testing"
)

var (
    prodSeg = Segmenter{}
)

func TestLargeDictionary(t *testing.T) {
    // 分词
    // text := []byte("中华人民共和国中央人民政府")
    text := []byte("abc haha中华人民共和国 中央人民政府,搜狗科技")
    segments := prodSeg.Segment(text)

    // 处理分词结果
    // 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
    //fmt.Println(sego.SegmentsToString(segments, false))
    for _, item := range segments {
        fmt.Printf("%#v\n", item)
        fmt.Printf("%#v\n", item.Token().Text())
    }
}
