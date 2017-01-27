# Go-Wordseg

Thai word segmentation with Golang.

## Algorithm in use

Dictionary-based longest matching

## Dictionary

Dictionary ```dict/full.dict``` contains all words listed in official Thai dictionary.

## Example

```
package main

import (
	"fmt"
	"strings"

	"./wordseg"
)

func main() {
	s := wordseg.NewSeg(nil)

	s.UseDictFile("./dict/full.dict")

	t := s.SegmentText("ประสบการณ์ครั้งใหม่ตอบโจทย์ทุกคนเออดีนะ")
	fmt.Println(strings.Join(t, "-"))

	t = s.SegmentText("test")
	fmt.Println(strings.Join(t, "-"))

	t = s.SegmentText("เออ")
	fmt.Println(strings.Join(t, "-"))

	t = s.SegmentText("this is a working software. และมันเป็นเรื่อง normal นะ นะ AB.123ข้อความ")
	fmt.Println(strings.Join(t, "-"))
}
```