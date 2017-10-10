package sego

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	iconv "github.com/djimenez/iconv-go"
	"github.com/lokicui/mlt/interface/gen-go/wenwen_seg"
	"math/rand"
	"time"
)

type WordInfo struct {
	Word     string
	IsEntity bool
	wenwen_seg.QueryTermInfo
}

func SegmentQuery(query string, useEntity bool) (words []WordInfo, err error) {
	addrs := []string{
		"10.134.104.53:30001",
		"10.134.104.54:30001",
		"10.134.104.53:30002",
		"10.134.104.54:30002",
		"10.134.45.60:30001",
		"10.134.45.63:30001",
		"10.134.45.64:30001",
		"10.134.92.21:30001",
		"10.134.100.115:30001",
		//"10.134.45.59:30001"
	}
	addr := addrs[rand.Intn(len(addrs))]
	timeout := time.Duration(50) //50ms
	socket, err := thrift.NewTSocketTimeout(addr, timeout*time.Millisecond)
	if err != nil {
		fmt.Println("Error opening socket:", err)
		return
	}
	var transportFactory thrift.TTransportFactory
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport := transportFactory.GetTransport(socket)
	if err = transport.Open(); err != nil {
		fmt.Println("Error opening socket:", err)
		return
	}
	defer transport.Close()
	client := wenwen_seg.NewSegServiceClientFactory(transport, protocolFactory)

	for _, q := range splitQuery(query) {
		wds, err := segment(q, useEntity, client)
		//try 3 times
		if err != nil {
			wds, err = segment(q, useEntity, client)
		}
		if err != nil {
			wds, err = segment(q, useEntity, client)
		}
		if err != nil {
			wds, err = segment(q, useEntity, client)
		}
		if err != nil {
			return words, err
		}
		for _, w := range wds {
			words = append(words, w)
		}
	}
	return
}

func hasSecond(termid int32) bool {
	t := termid >> 28
	flag := (termid >> 24) & 0xf
	if flag == 0xf && (3 == t || 4 == t || 6 == t || 10 == t) {
		return true
	}
	return false
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

//分句
func splitQuery(query string) (pieces []string) {
	delims := []string{" ", "，", "。", "　"}
	query_rune := []rune(query)

	//strtoken
	pos_list := []int{}
	for i, a := range query_rune {
		for _, d := range delims {
			if string(a) == d {
				pos_list = append(pos_list, i)
				break
			}
		}
	}
	pos_list = append(pos_list, len(query_rune))

	threshold := 64
	start := 0
	for i, pos := range pos_list {
		j := i + 1
		if j < len(pos_list) {
			jpos := pos_list[j]
			if jpos-start > threshold {
				piece := string(query_rune[start : pos+1])
				pieces = append(pieces, piece)
				start = pos + 1
			}
		} else if pos-start > 0 {
			piece := string(query_rune[start:pos])
			pieces = append(pieces, piece)
		}
	}
	return
}

func segment(query string, useEntity bool, client *wenwen_seg.SegServiceClient) (words []WordInfo, err error) {
	query_gbk, err := iconv.ConvertString(query, "utf-8", "gbk")
	resp, err := client.QuerySegment(query_gbk, 0)
	if err != nil {
		return
	}
	terms := resp.GetTerms()
	used := make([]bool, len(terms), len(terms))
	l := make([]WordInfo, 0, len(terms))
	entity_map := make(map[int16]int16)
	if useEntity {
		for _, info := range resp.GetEntityWords() {
			b := info.TermBeg
			e := If(info.TermEnd < int16(len(resp.GetTerms())), info.TermEnd, int16(len(resp.GetTerms()))).(int16)
			v, ok := entity_map[b]
			if ok {
				if v < e {
					entity_map[b] = e
					v = e
					for j := int(b); j <= int(v); j++ {
						used[j] = true
					}
				}
			} else if used[b] == false {
				entity_map[b] = e
				v = e
				for j := int(b); j <= int(v); j++ {
					used[j] = true
				}
			}
		}
	}
	for i := 0; i < len(terms); i++ {
		term := terms[i]
		termid := term.GetTermID()
		word := WordInfo{}
		e, ok := entity_map[int16(i)]
		if ok {
            start := int(resp.Terms[i].Pos*2)
            if start > len(resp.QueryGbkSbc) {
                start = len(resp.QueryGbkSbc)
            }
            end := int(resp.Terms[e].Pos*2+resp.Terms[e].Len*2)
            if end > len(resp.QueryGbkSbc) {
                end = len(resp.QueryGbkSbc)
            }
			termstr := resp.QueryGbkSbc[start:end]
			word.IsEntity = true
			word.Word, err = iconv.ConvertString(termstr, "gbk", "utf-8")
			i = int(e)
		} else {
			if hasSecond(termid) && i+1 < len(resp.GetTerms()) {
                start := int(resp.Terms[i].Pos*2)
                if start > len(resp.QueryGbkSbc) {
                    start = len(resp.QueryGbkSbc)
                }
                end := int(resp.Terms[i+1].Pos*2+resp.Terms[i+1].Len*2)
                if end > len(resp.QueryGbkSbc) {
                    end = len(resp.QueryGbkSbc)
                }
				termstr := resp.QueryGbkSbc[start:end]
				word.Word, err = iconv.ConvertString(termstr, "gbk", "utf-8")
				if err != nil {
					continue
				}
				word.QueryTermInfo = *term
				word.IsEntity = false
			} else {
                start := int(resp.Terms[i].Pos*2)
                if start > len(resp.QueryGbkSbc) {
                    start = len(resp.QueryGbkSbc)
                }
                end := int(resp.Terms[i].Pos*2+resp.Terms[i].Len*2)
                if end > len(resp.QueryGbkSbc) {
                    end = len(resp.QueryGbkSbc)
                }
				termstr := resp.QueryGbkSbc[start:end]
				word.Word, err = iconv.ConvertString(termstr, "gbk", "utf-8")
				if err != nil {
					continue
				}
				word.QueryTermInfo = *term
				word.IsEntity = false
			}
		}
		l = append(l, word)
	}
	return l, err
}
