package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"bli.tf/models"
)

// PostsAPIHandler 文章列表 API
func PostsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// 获取所有文章（包括未发布的）
		postsData, err := LoadPosts()
		if err != nil {
			http.Error(w, `{"error":"加载失败"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(postsData)

	case http.MethodPost:
		// 创建新文章
		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, `{"error":"数据格式错误"}`, http.StatusBadRequest)
			return
		}

		// 生成 ID 和时间戳
		post.ID = generateID()
		if post.Slug == "" {
			post.Slug = generateSlug(post.Title)
		}
		now := time.Now().Format("2006-01-02T15:04:05")
		post.CreatedAt = now
		post.UpdatedAt = now

		postsData, _ := LoadPosts()
		postsData.Posts = append(postsData.Posts, post)

		if err := SavePosts(postsData); err != nil {
			http.Error(w, `{"error":"保存失败"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"post":    post,
		})

	default:
		http.Error(w, `{"error":"方法不允许"}`, http.StatusMethodNotAllowed)
	}
}

// PostAPIHandler 单篇文章 API
func PostAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 从 URL 获取文章 ID
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postID := strings.TrimSuffix(path, "/")

	postsData, err := LoadPosts()
	if err != nil {
		http.Error(w, `{"error":"加载失败"}`, http.StatusInternalServerError)
		return
	}

	// 查找文章
	postIndex := -1
	for i, p := range postsData.Posts {
		if p.ID == postID {
			postIndex = i
			break
		}
	}

	switch r.Method {
	case http.MethodGet:
		if postIndex == -1 {
			http.Error(w, `{"error":"文章不存在"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(postsData.Posts[postIndex])

	case http.MethodPut:
		var updatedPost models.Post
		if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
			http.Error(w, `{"error":"数据格式错误"}`, http.StatusBadRequest)
			return
		}

		if postIndex == -1 {
			http.Error(w, `{"error":"文章不存在"}`, http.StatusNotFound)
			return
		}

		// 保留原始 ID 和创建时间
		updatedPost.ID = postsData.Posts[postIndex].ID
		updatedPost.CreatedAt = postsData.Posts[postIndex].CreatedAt
		updatedPost.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")

		if updatedPost.Slug == "" {
			updatedPost.Slug = generateSlug(updatedPost.Title)
		}

		postsData.Posts[postIndex] = updatedPost

		if err := SavePosts(postsData); err != nil {
			http.Error(w, `{"error":"保存失败"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"post":    updatedPost,
		})

	case http.MethodDelete:
		if postIndex == -1 {
			http.Error(w, `{"error":"文章不存在"}`, http.StatusNotFound)
			return
		}

		// 删除文章
		postsData.Posts = append(postsData.Posts[:postIndex], postsData.Posts[postIndex+1:]...)

		if err := SavePosts(postsData); err != nil {
			http.Error(w, `{"error":"删除失败"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "删除成功",
		})

	default:
		http.Error(w, `{"error":"方法不允许"}`, http.StatusMethodNotAllowed)
	}
}

// 生成唯一 ID
func generateID() string {
	return time.Now().Format("20060102150405") + randomString(4)
}

// 生成 slug
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	// 移除特殊字符
	allowed := "abcdefghijklmnopqrstuvwxyz0123456789-"
	var result strings.Builder
	for _, c := range slug {
		if strings.ContainsRune(allowed, c) {
			result.WriteRune(c)
		}
	}
	return result.String()
}

// 随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
