// textwork.go
package textwork

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type TWORK struct {
	TEXT       string
	si, ei, ci int
}

func (tw *TWORK) SetFromUTF8(text string) {
	tw.TEXT = text
	tw.ci = 0
	tw.si = 0
	tw.ei = 0
}
func (tw *TWORK) OpenFile(fname string) (r bool) {
	var err error
	var data []byte

	data, err = ioutil.ReadFile(fname)
	if err == nil {
		tw.TEXT = (string(data))
		tw.ci = 0
		tw.si = 0
		tw.ei = 0

		return true
	}
	return false
}
func (tw *TWORK) SaveToFile(fname string) {
	da := []byte(string(tw.TEXT))
	ioutil.WriteFile(fname, da, 0777)
	os.Chmod(fname, 0777)

}

//add string
func (tw *TWORK) AddString(text string) {
	tw.TEXT += text
}
func (tw *TWORK) AddRunes(runet []rune) {
	str := string(runet)
	tw.TEXT += str
}
func (tw *TWORK) Add(ti ...interface{}) {
	for i := range ti {
		tw.TEXT += fmt.Sprint(ti[i])
	}
}

func (tw *TWORK) GetBlock(ss, es string) string {
	if !tw.Seek(ss) {
		return ""
	}
	si := tw.ci
	if !tw.Seek(es) {
		return ""
	}
	return string(tw.TEXT[si:tw.si])
}

func (tw *TWORK) SetBlock(startstring, endstring, newstring string) {
	if !tw.Seek(startstring) {
		return
	}
	si := tw.ci
	if !tw.Seek(endstring) {
		return
	}
	//fmt.Println(si, tw.si, tw.ci)
	//fmt.Println(tw.TEXT[si], tw.TEXT[tw.si], tw.TEXT[tw.ci])

	//fmt.Println(tw.TEXT[si], string(tw.TEXT[si]), tw.TEXT[tw.si], string(tw.TEXT[tw.si]), tw.TEXT[tw.ci], string(tw.TEXT[tw.ci]))

	sr := tw.TEXT[:si] + newstring
	tw.TEXT = sr + tw.TEXT[tw.si:]
	tw.ci = si + len(newstring) + len(endstring)

	//return string(tw.TEXT[si:tw.si])
}

func (tw *TWORK) Split(splstring string) []string {

	return strings.Split(tw.TEXT, splstring)
}

func (tw *TWORK) Seek(ss string) bool {
	//fmt.Println("Seek ->", ss)
	l := len(tw.TEXT)
	if tw.ci >= l {
		return false
	}
	//sr := []rune(ss)
	srl := len(ss)
	for tw.ci < l {
		//если первый символ совпадает, то проверить весь слайс
		if tw.TEXT[tw.ci] == ss[0] {
			tw.si = tw.ci

			//если в результате выйдем за границ текста, то и проверять нет смысла
			if tw.ci+srl > l {
				return false
			}
			if srl == 1 {
				tw.ei = tw.ci + 1
				return true
			}
			i := 1
			for i < srl {
				if ss[i] == tw.TEXT[tw.ci+i] {
					i++
				} else {
					break
				}
			}
			if i == srl {
				tw.ci += i
				tw.ei = tw.ci
				return true
			}

		}
		tw.ci++
	}
	return false
}

func (tw *TWORK) Replace(ss, ns string) {

	tw.TEXT = strings.Replace(string(tw.TEXT), ss, ns, 1)

}
func (tw *TWORK) ReplaceN(ss, ns string, count int) {

	tw.TEXT = strings.Replace(string(tw.TEXT), ss, ns, count)
}

//words work. ss and es = 32-> " " or 10->\n

func (tw *TWORK) NextWord() bool {
	l := len(tw.TEXT)
	if tw.ci >= l {
		return false
	}
	tw.si = tw.ci
	//делает шаг вперед если курсор не на старте
	/*
		if tw.si != 0 && (tw.TEXT[tw.ci] == 32 || tw.TEXT[tw.ci] == 10) {
			tw.si++
		}
	*/
	tw.ci++
	//перебирать символы пока курсор не дошел до конца массива
	for tw.ci < l {
		//если символ равен пробелу или переносу строки, то надо проверить слово ли это
		if (tw.TEXT[tw.ci] == 32 || tw.TEXT[tw.ci] == 10) && tw.ci > 0 {
			//если предыдущие смвол отличается от разделителя слов, то у нас слово
			if tw.TEXT[tw.ci-1] != 32 && tw.TEXT[tw.ci-1] != 10 {
				tw.ei = tw.ci
				if tw.ei == tw.si {
					return false
				}
				tw.ci++
				return true
			}

		}
		tw.ci++
	}
	tw.ei = tw.ci
	if tw.ei == tw.si {
		return false
	}
	return true
}
func (tw *TWORK) GoToWord(w string) bool {
	for tw.NextWord() {
		if tw.CurWord() == w {
			return true
		}
	}
	return false
}
func (tw *TWORK) GetWordsTo(ew string) string {

	si := tw.ci
	if !tw.GoToWord(ew) {
		return ""
	}
	return string(tw.TEXT[si : tw.si-1])
}
func (tw *TWORK) CurWord() string {
	return string(tw.TEXT[tw.si:tw.ei])
}
func (tw *TWORK) CurIndex() int {
	return tw.ci
}
func (tw *TWORK) SetIndex(ni int) {
	if ni >= 0 && ni < len(tw.TEXT) {
		tw.ci = ni
	}
}

func slicetostring(elems []string) string {
	rt := ""
	for i := range elems {
		rt += elems[i]
	}
	return rt
}
