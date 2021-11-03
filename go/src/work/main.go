package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

/* ------ viewHandler ------ */
/* http.ResponseWriterと*http.Requestを記述する必要がある（HandleFuncの定義元にジャンプすると引数として渡してあるため）*/
func viewHandler(w http.ResponseWriter, r *http.Request) {
	/* /view/test */
	/* http.Requestにリクエストを出した中身のデータが入っているため、r.URL.PathでURLが取得できる */
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

/* ------ editHandler ------ */
func editHandler(w http.ResponseWriter, r *http.Request) {
	/* /view/test */
	/* http.Requestにリクエストを出した中身のデータが入っているため、r.URL.PathでURLが取得できる */
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

/* --------- save -----------*/
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	/*----- もしviewというURLが来た場合viewHandlerにいく */
	/* 自分のハンドラーを立ち上げたい場合は、ListenAndServeの前にhttp.HandleFuncのメソッドを使用し、ハンドラーを登録する必要がある */
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	/*----- webサーバーの立ち上げ ----- */
	/* ハンドルがnilの為、/view/以外にアクセスしようとするとnot foundが返る */
	log.Fatal(http.ListenAndServe(":8080", nil))
}
