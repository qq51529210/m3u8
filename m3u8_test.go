package m3u8

// import (
// 	"os"
// 	"testing"
// 	"time"
// )

// func TestReadListFrom(t *testing.T) {
// 	f, e := os.Open("list1.m3u8")
// 	if nil != e {
// 		t.Fatal(e)
// 	}
// 	defer f.Close()

// 	l, e := ReadListFrom(f)
// 	if nil != e {
// 		t.Fatal(e)
// 	}
// 	t.Log(l.Version)
// 	t.Log(l.Sequence)
// 	t.Log(l.TargetDuration)
// 	for ele := l.Segment.Front(); nil != ele; ele = ele.Next() {
// 		s := ele.Value.(*Segment)
// 		t.Log(s.Duration, s.URI)
// 	}
// 	t.Log(l.End)
// }

// func TestWriteListTo(t *testing.T) {
// 	l := new(List)
// 	l.Version = "3"
// 	l.Sequence = 10
// 	l.End = true
// 	l.TargetDuration = time.Second * 5
// 	WriteListTo(l, os.Stderr)
// }
