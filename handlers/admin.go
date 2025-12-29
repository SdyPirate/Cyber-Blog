package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"time"

	"bli.tf/middleware"
	"bli.tf/models"
	"golang.org/x/crypto/bcrypt"
)

const defaultPassword = "admin123"

// InitAdmin 初始化管理员账号
func InitAdmin() error {
	data, err := os.ReadFile("data/admin.json")
	if err != nil {
		return createDefaultAdmin()
	}

	var admin models.Admin
	if err := json.Unmarshal(data, &admin); err != nil {
		return createDefaultAdmin()
	}

	// 如果密码hash包含 "YourHash"，则重新生成
	if len(admin.PasswordHash) < 20 {
		return createDefaultAdmin()
	}

	return nil
}

func createDefaultAdmin() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.Admin{
		Username:     "admin",
		PasswordHash: string(hash),
	}

	data, err := json.MarshalIndent(admin, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("data/admin.json", data, 0644)
}

// LoadAdmin 加载管理员信息
func LoadAdmin() (*models.Admin, error) {
	data, err := os.ReadFile("data/admin.json")
	if err != nil {
		return nil, err
	}
	var admin models.Admin
	err = json.Unmarshal(data, &admin)
	return &admin, err
}

// LoginPageHandler 登录页面
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/admin/login.html")
	if err != nil {
		http.Error(w, "模板加载失败", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// LoginHandler 登录处理
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	admin, err := LoadAdmin()
	if err != nil {
		http.Error(w, `{"error":"服务器错误"}`, http.StatusInternalServerError)
		return
	}

	if username != admin.Username {
		http.Redirect(w, r, "/admin/login?error=1", http.StatusFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		http.Redirect(w, r, "/admin/login?error=1", http.StatusFound)
		return
	}

	token, err := middleware.GenerateToken(username)
	if err != nil {
		http.Error(w, "生成令牌失败", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
	})

	http.Redirect(w, r, "/admin", http.StatusFound)
}

// LogoutHandler 登出处理
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
	})
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

// DashboardHandler 后台管理面板
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	content, err := LoadContent()
	if err != nil {
		http.Error(w, "加载内容失败", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin/dashboard.html")
	if err != nil {
		http.Error(w, "模板加载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, content)
}

// GetContentHandler 获取内容 API
func GetContentHandler(w http.ResponseWriter, r *http.Request) {
	content, err := LoadContent()
	if err != nil {
		http.Error(w, `{"error":"加载失败"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

// SaveContentHandler 保存内容 API
func SaveContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"方法不允许"}`, http.StatusMethodNotAllowed)
		return
	}

	var content models.Content
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		http.Error(w, `{"error":"数据格式错误"}`, http.StatusBadRequest)
		return
	}

	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		http.Error(w, `{"error":"保存失败"}`, http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile("data/content.json", data, 0644); err != nil {
		http.Error(w, `{"error":"写入失败"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true,"message":"保存成功"}`))
}
