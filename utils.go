package sego

import (
	"bytes"
	"fmt"
    "strings"
)

//全角->半角
func SBC2DBC(s string) string {
    r := [] string {}
    for _, i := range s {
        inside_code := i
        if inside_code == 0x3000 {
            inside_code = 0x0020
        } else if inside_code >= 0xff01 && inside_code <= 0xff5e {
            inside_code -= 0xfee0
        }
        r = append(r, string(inside_code))
    }
    return strings.Join(r, "")
}

//半角->全角
func DBC2SBC(s string) string {
    r := [] string {}
    for _, i := range s {
        inside_code := i
        if inside_code == 0x20 {
            inside_code = 0x3000
        } else if inside_code >= 0x20 && inside_code <= 0x7e {
            inside_code += 0xfee0
        }
        r = append(r, string(inside_code))
    }
    return strings.Join(r, "")
}

func tokenToString(token *Token) (output string) {
	for _, s := range token.segments {
		output += tokenToString(s.token)
	}
	output += fmt.Sprintf("%s/%s ", textSliceToString(token.text), token.pos)
	return
}

// 输出分词结果到一个字符串slice
//
// 有两种输出模式，以"中华人民共和国"为例
//
//  普通模式（searchMode=false）输出一个分词"[中华人民共和国]"
//  搜索模式（searchMode=true） 输出普通模式的再细致切分：
//      "[中华 人民 共和 共和国 人民共和国 中华人民共和国]"
//
// 搜索模式主要用于给搜索引擎提供尽可能多的关键字，详情请见Token结构体的注释。

func SegmentsToSlice(segs []Segment, searchMode bool) (output []string) {
	if searchMode {
		for _, seg := range segs {
			output = append(output, tokenToSlice(seg.token)...)
		}
	} else {
		for _, seg := range segs {
			output = append(output, seg.token.Text())
		}
	}
	return
}

func tokenToSlice(token *Token) (output []string) {
	for _, s := range token.segments {
		output = append(output, tokenToSlice(s.token)...)
	}
	output = append(output, textSliceToString(token.text))
	return output
}

// 将多个字元拼接一个字符串输出
func textSliceToString(text []Text) string {
	var output string
	for _, word := range text {
		output += string(word)
	}
	return output
}

// 返回多个字元的字节总长度
func textSliceByteLength(text []Text) (length int) {
	for _, word := range text {
		length += len(word)
	}
	return
}

func textSliceToBytes(text []Text) []byte {
	var buf bytes.Buffer
	for _, word := range text {
		buf.Write(word)
	}
	return buf.Bytes()
}
