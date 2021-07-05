// textwork.go
package textwork

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type TWORK struct {
	Runes      []rune
	si, ei, ci int
}

func (tw *TWORK) SetFromUTF8(text string) {
	tw.Runes = []rune(text)
	tw.ci = 0
	tw.si = 0
	tw.ei = 0
}
func (tw *TWORK) OpenFile(fname string) (r bool) {
	var err error
	var data []byte

	data, err = ioutil.ReadFile(fname)
	if err == nil {
		tw.Runes = []rune(string(data))
		tw.ci = 0
		tw.si = 0
		tw.ei = 0

		return true
	}
	return false
}
func (tw *TWORK) SaveToFile(fname string) {
	da := []byte(string(tw.Runes))
	ioutil.WriteFile(fname, da, 0777)
	os.Chmod(fname, 0777)

}

//add string
func (tw *TWORK) AddString(text string) {
	tw.Runes = append(tw.Runes, []rune(text)...)
}
func (tw *TWORK) AddRunes(runet []rune) {
	tw.Runes = append(tw.Runes, runet...)
}
func (tw *TWORK) Add(ti ...interface{}) {
	for i := range ti {
		tw.Runes = append(tw.Runes, []rune(fmt.Sprint(ti[i]))...)
	}
}
func (tw *TWORK) GetAsString() string {
	return string(tw.Runes)
}

func (tw *TWORK) GetBlock(ss, es string) string {
	if !tw.Seek(ss) {
		return ""
	}
	si := tw.ci
	if !tw.Seek(es) {
		return ""
	}
	return string(tw.Runes[si:tw.si])
}
func (tw *TWORK) GetWordsTo(ew string) string {

	si := tw.ci
	if !tw.GoToWord(ew) {
		return ""
	}
	return string(tw.Runes[si : tw.si-1])
}
func (tw *TWORK) SetBlock(startstring, endstring, newstring string) {
	if !tw.Seek(startstring) {
		return
	}
	si := tw.ci
	if !tw.Seek(endstring) {
		return
	}
	rr := []rune(newstring)
	sr := append(tw.Runes[:si], rr...)
	tw.Runes = append(sr, tw.Runes[tw.si:]...)
	//return string(tw.Runes[si:tw.si])
}

func (tw *TWORK) Split(splstring string) []string {
	var r []string
	tw.ci = 0
	tw.si = 0
	si := 0
	for tw.Seek(splstring) {
		ts := string(tw.Runes[si:tw.si])
		//fmt.Println("'" + ts + "'")
		r = append(r, ts)
		si = tw.ei
		tw.ci = si
	}

	return r
}

func (tw *TWORK) Seek(ss string) bool {
	l := len(tw.Runes)
	if tw.ci >= l {
		return false
	}
	sr := []rune(ss)
	srl := len(sr)
	for tw.ci < l {
		//если первый символ совпадает, то проверить весь слайс
		if tw.Runes[tw.ci] == sr[0] {
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
				if sr[i] == tw.Runes[tw.ci+i] {
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
func (tw *TWORK) GoToWord(w string) bool {
	for tw.NextWord() {
		if tw.CurWord() == w {
			return true
		}
	}
	return false
}
func (tw *TWORK) Replace(ss, ns string) bool {
	/*
		if tw.Seek(ss) {
			rr := []rune(ns)
			sr := append(tw.Runes[:tw.si], rr...)
			tw.Runes = append(sr, tw.Runes[tw.ei:]...)
			return true
		}
		return false
	*/
	tw.Runes = []rune(strings.Replace(string(tw.Runes), ss, ns, 1))
	return true
}
func (tw *TWORK) ReplaceN(ss, ns string, count int) bool {

	for count > 0 {
		if tw.Seek(ss) {
			rr := []rune(ns)
			sr := append(tw.Runes[:tw.si], rr...)
			tw.Runes = append(sr, tw.Runes[tw.ei:]...)
			count--
		} else {
			break
		}
	}
	return false
}

//words work. ss and es = 32-> " "

func (tw *TWORK) NextWord() bool {
	l := len(tw.Runes)
	if tw.ci >= l {
		return false
	}
	tw.si = tw.ci
	//делает шаг вперед если курсор не на старте
	/*
		if tw.si != 0 && (tw.Runes[tw.ci] == 32 || tw.Runes[tw.ci] == 10) {
			tw.si++
		}
	*/
	tw.ci++
	//перебирать символы пока курсор не дошел до конца массива
	for tw.ci < l {
		//если символ равен пробелу или переносу строки, то надо проверить слово ли это
		if (tw.Runes[tw.ci] == 32 || tw.Runes[tw.ci] == 10) && tw.ci > 0 {
			//если предыдущие смвол отличается от разделителя слов, то у нас слово
			if tw.Runes[tw.ci-1] != 32 && tw.Runes[tw.ci-1] != 10 {
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
func (tw *TWORK) CurWord() string {
	return string(tw.Runes[tw.si:tw.ei])
}
func (tw *TWORK) CurIndex() int {
	return tw.ci
}
func (tw *TWORK) SetIndex(ni int) {
	if ni >= 0 && ni < len(tw.Runes) {
		tw.ci = ni
	}
}
