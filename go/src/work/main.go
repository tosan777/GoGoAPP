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
	return &Page{Title: title, Body: body}, nil
}

/* ------ viewHandler ------ */
/* http.ResponseWriterと*http.Requestを記述する必要がある（HandleFuncの定義元にジャンプすると引数として渡してあるため）*/
func viewHandler(w http.ResponseWriter, r *http.Request) {
	/* /view/test */
	/* http.Requestにリクエストを出した中身のデータが入っているため、r.URL.PathでURLが取得できる */
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	/*----- もしviewというURLが来た場合viewHandlerにいく */
	/* 自分のハンドラーを立ち上げたい場合は、ListenAndServeの前にhttp.HandleFuncのメソッドを使用し、ハンドラーを登録する必要がある */
	http.HandleFunc("/view/", viewHandler)
	/*----- webサーバーの立ち上げ ----- */
	/* ハンドルがnilの為、/view/以外にアクセスしようとするとnot foundが返る */
	log.Fatal(http.ListenAndServe(":8080", nil))
}
