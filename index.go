package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type dat struct {
	Records int
	Sub     string
	Dat     int
}
type data struct {
	Url string
	X   []dat
}
type record struct {
	Url string
	Sub string
	Num string
}

func main() {
	buf, _ := ioutil.ReadFile("x.json")
	var d []data
	json.Unmarshal(buf, &d)
	newR := []record{}

	re := regexp.MustCompile("(\\d+)\\.dat<>(.*) \\((\\d+)\\)")
	reurl := regexp.MustCompile("^(https?://.+/)(.*)/$")
	//next:=new([]data)
	wg := new(sync.WaitGroup)
	for ii, value := range d {
		wg.Add(1)
		go func(i int) {
            x := []dat{}
            old := d[i].X
			for _, l := range hRead(value.Url + "subject.txt") {
				s := re.FindStringSubmatch(l)
				if s == nil {
					break
				}
				n := dat{parse(s[3]), s[2], parse(s[1])}
				x = append(x, n)
				f := true
				for _, ll := range old {
					if (ll.Dat == n.Dat) && (ll.Sub == n.Sub) {
						if ll.Records < n.Records {
							r := new(record)
							s0 := reurl.FindStringSubmatch(value.Url)
							r.Url = s0[1] + "test/read.cgi/" + s0[2] + "/" + strconv.Itoa(n.Dat) + "/"
							r.Sub = ll.Sub
							r.Num = strconv.Itoa(ll.Records) + "-" + strconv.Itoa(n.Records)
							newR = append(newR, *r)
						}
						f = false
						break
					}
				}
				if f {
					r := new(record)
					s0 := reurl.FindStringSubmatch(value.Url)
					r.Url = s0[1] + "test/read.cgi/" + s0[2] + "/" + strconv.Itoa(n.Dat) + "/"
					r.Sub = n.Sub
					r.Num = "0-" + strconv.Itoa(n.Records)
					newR = append(newR, *r)
				}
			}
            d[i].X = x
            fmt.Println(i)
			wg.Done()
		}(ii)
	}
	wg.Wait()
	tmpl, err := template.New("master").Parse("<html><head></head><body>{{range .}}<a href= \"{{ .Url }}\">{{ .Sub }}({{ .Num}})</a>{{end}}</body></html>")
	if err != nil {
		panic(err)
	}
    f, err2 := os.OpenFile("index.html", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
    if err2 != nil {
		panic(err2)
	}
    err = tmpl.Execute(f, newR)
    err = f.Close()
	//fmt.Printf("%+v\n", newR)
	result, _ := json.MarshalIndent(d, "", "  ")
	ioutil.WriteFile("x.json", result, os.ModePerm)
}

func parse(str string) int {
	r, _ := strconv.Atoi(str)
	return r
}

func hRead(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
        fmt.Print(err)
		return nil
	}
	defer resp.Body.Close()
	res := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
	b, err2 := ioutil.ReadAll(res)
	if err2 == nil {
		return strings.Split(string(b), "\n")
	}
	return nil
}

/*func ayncmap(arr A[],f func){

}*/
