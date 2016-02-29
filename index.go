package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "regexp"
    "strconv"
    "os"
    "golang.org/x/text/encoding/japanese"
    "golang.org/x/text/transform"
)

type dat struct {
    Records  int64
    Sub   string
    Dat int64
}
type data struct{
    Url string
    X []dat
}
type record struct{
    Url string
    Sub string
    Num string
}

func main() {
    buf, _ := ioutil.ReadFile("x.json")
    var d []data
    json.Unmarshal(buf,&d)
    newR:=[]record{}
    
    re := regexp.MustCompile("(\\d+)\\.dat<>(.*) \\((\\d+)\\)")
    //next:=new([]data)
    for i, value := range d {
        x:= []dat{}
        old:=d[i].X
        for _, l := range hRead(value.Url+"subject.txt") {
            s:=re.FindStringSubmatch(l)
            if(s==nil){break}
            n:=new(dat)
            n.Records=parse(s[3])
            n.Sub=s[2]
            n.Dat=parse(s[1])
            x=append(x,*n)
            for _,ll := range old{
                if((ll.Sub==n.Sub)&&(ll.Dat==n.Dat)){
                   if(ll.Records<n.Records){
                    r:=new(record)
                    r.Url=value.Url+strconv.FormatInt(n.Dat, 10)+"/"
                    r.Sub=ll.Sub
                    r.Num=strconv.FormatInt(ll.Records, 10)+"-"+strconv.FormatInt(n.Records, 10)
                    newR=append(newR,*r)
                   }
                   break
                }
            }
        }
        d[i].X=x
    }
    
    fmt.Printf("%+v\n", newR)
    result , _ :=json.Marshal(d)
    ioutil.WriteFile("x.json", result, os.ModePerm)
}

func parse(str string) int64 {
    r,_:=strconv.ParseInt(str, 10, 32)
    return r
}

func hRead(url string) []string{
    resp, err := http.Get(url)
    if err != nil {panic(err)}
    defer resp.Body.Close()
    res:=transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
    b, err2 := ioutil.ReadAll(res)
    if err2 == nil {
        return strings.Split(string(b),"\n")
    }
    return nil;
}