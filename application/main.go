package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "html/template"
    "os"
)

type Page struct {
    Title string
    Body []byte
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func main() {
    http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf(r.Method + "\n\n")
        title := r.URL.Path[len("/view/"):]
        p, err1 := loadPage(title)

        if err1 != nil {
            p = &Page{ Title: title }
        }

        t, err := template.ParseFiles("view.html")
        if err != nil {}
        t.Execute(w, p)
    })
    http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            r.ParseForm()
            name := r.FormValue("name")
            body := []byte(r.FormValue("body"))
            fmt.Printf(name)
            p := &Page{Title: name, Body: body}
            p.save()
            redirect_page := "/view/" + name
            fmt.Printf("Redirecting you to your view.")
            http.Redirect(w, r, redirect_page, http.StatusFound)
        }
    });

    http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf(r.Method)
        if r.Method != "DELETE" {
            p := &Page{}
            t, _ := template.ParseFiles("error.html")
            t.Execute(w, p)
        }
        title := r.URL.Path[len("/delete/"):]
        filename := title + ".txt"
        os.Remove(filename)
        http.Redirect(w, r, "", http.StatusOK)
    });
    http.ListenAndServe(":8080", nil)
}
