package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

/* ----- ファイルの読み込み ------*/
// test.txtが作成される
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: "test", Body: body}, nil
}

/* ------ viewHandler ------ */
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}
func main() {
	/*----- もしviewというURLが来た場合viewHandlerにいく */
	http.HandleFunc("/view/", viewHandler)
	/*----- webサーバーの立ち上げ ----- */
	log.Fatal(http.ListenAndServe(":8080", nil))
}
