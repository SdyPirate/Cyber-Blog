package main

import (
	"fmt"
	"log"
	"net/http"

	"bli.tf/handlers"
	"bli.tf/middleware"
)

func main() {
	// 初始化管理员账号
	if err := handlers.InitAdmin(); err != nil {
		log.Printf("警告: 初始化管理员账号失败: %v", err)
	}

	// 静态文件服务
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	// 公开路由
	http.HandleFunc("/", handlers.IndexHandler)

	// 博客公开路由
	http.HandleFunc("/blog", handlers.BlogListHandler)
	http.HandleFunc("/blog/", handlers.BlogPostHandler)

	// 后台路由
	http.HandleFunc("/admin/login", handlers.LoginPageHandler)
	http.HandleFunc("/admin/auth", handlers.LoginHandler)
	http.HandleFunc("/admin/logout", handlers.LogoutHandler)
	http.HandleFunc("/admin", middleware.AuthMiddleware(handlers.DashboardHandler))

	// API 路由
	http.HandleFunc("/api/content", middleware.AuthAPIMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetContentHandler(w, r)
		} else if r.Method == http.MethodPost {
			handlers.SaveContentHandler(w, r)
		} else {
			http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		}
	}))

	// 博客 API 路由
	http.HandleFunc("/api/posts", middleware.AuthAPIMiddleware(handlers.PostsAPIHandler))
	http.HandleFunc("/api/posts/", middleware.AuthAPIMiddleware(handlers.PostAPIHandler))

	port := ":12302"
	fmt.Printf("🚀 服务器启动于 http://localhost%s\n", port)
	fmt.Println("📝 后台管理: http://localhost:12302/admin")
	fmt.Println("👤 默认账号: admin")
	fmt.Println("🔑 默认密码: admin123")

	log.Fatal(http.ListenAndServe(port, nil))
}
