package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"bli.tf/models"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

// LoadContent 加载内容
func LoadContent() (*models.Content, error) {
	data, err := os.ReadFile("data/content.json")
	if err != nil {
		return nil, err
	}
	var content models.Content
	err = json.Unmarshal(data, &content)
	return &content, err
}

// IndexHandler 首页处理
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	content, err := LoadContent()
	if err != nil {
		http.Error(w, "加载内容失败", http.StatusInternalServerError)
		return
	}
	
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "模板加载失败", http.StatusInternalServerError)
		return
	}
	
	tmpl.Execute(w, content)
}
