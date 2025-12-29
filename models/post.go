package models

// Post 博客文章
type Post struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Slug       string   `json:"slug"`
	Content    string   `json:"content"` // Markdown 内容
	Summary    string   `json:"summary"` // 摘要
	Tags       []string `json:"tags"`
	CoverImage string   `json:"coverImage"` // 封面图
	Published  bool     `json:"published"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
}

// PostsData 文章数据包装
type PostsData struct {
	Posts []Post `json:"posts"`
}
