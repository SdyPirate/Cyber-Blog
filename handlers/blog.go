package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strings"

	"bli.tf/models"
)

// LoadPosts 加载所有文章
func LoadPosts() (*models.PostsData, error) {
	data, err := os.ReadFile("data/posts.json")
	if err != nil {
		return &models.PostsData{Posts: []models.Post{}}, nil
	}
	var postsData models.PostsData
	err = json.Unmarshal(data, &postsData)
	if err != nil {
		return &models.PostsData{Posts: []models.Post{}}, nil
	}
	return &postsData, nil
}

// SavePosts 保存文章列表
func SavePosts(postsData *models.PostsData) error {
	data, err := json.MarshalIndent(postsData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("data/posts.json", data, 0644)
}

// GetPublishedPosts 获取已发布的文章，按时间降序
func GetPublishedPosts() []models.Post {
	postsData, _ := LoadPosts()
	var published []models.Post
	for _, p := range postsData.Posts {
		if p.Published {
			published = append(published, p)
		}
	}
	// 按发布时间降序排序
	sort.Slice(published, func(i, j int) bool {
		return published[i].CreatedAt > published[j].CreatedAt
	})
	return published
}

// GetAllTags 获取所有标签
func GetAllTags() []string {
	postsData, _ := LoadPosts()
	tagSet := make(map[string]bool)
	for _, p := range postsData.Posts {
		if p.Published {
			for _, tag := range p.Tags {
				tagSet[tag] = true
			}
		}
	}
	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// BlogListData 博客列表页数据
type BlogListData struct {
	Posts      []models.Post
	Tags       []string
	CurrentTag string
	Profile    models.Profile
	TotalPosts int
}

// BlogListHandler 博客列表页
func BlogListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/blog" && r.URL.Path != "/blog/" {
		// 处理单篇文章
		BlogPostHandler(w, r)
		return
	}

	posts := GetPublishedPosts()
	tags := GetAllTags()

	// 标签过滤
	currentTag := r.URL.Query().Get("tag")
	if currentTag != "" {
		var filtered []models.Post
		for _, p := range posts {
			for _, t := range p.Tags {
				if t == currentTag {
					filtered = append(filtered, p)
					break
				}
			}
		}
		posts = filtered
	}

	// 加载 profile 信息
	content, _ := LoadContent()

	data := BlogListData{
		Posts:      posts,
		Tags:       tags,
		CurrentTag: currentTag,
		Profile:    content.Profile,
		TotalPosts: len(posts),
	}

	tmpl, err := template.ParseFiles("templates/blog_list.html")
	if err != nil {
		http.Error(w, "模板加载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// BlogPostData 单篇文章页数据
type BlogPostData struct {
	Post     models.Post
	Content  template.HTML
	Profile  models.Profile
	PrevPost *models.Post
	NextPost *models.Post
}

// BlogPostHandler 单篇文章页
func BlogPostHandler(w http.ResponseWriter, r *http.Request) {
	// 从 URL 获取 slug
	path := strings.TrimPrefix(r.URL.Path, "/blog/")
	slug := strings.TrimSuffix(path, "/")

	if slug == "" || slug == "archive" {
		BlogArchiveHandler(w, r)
		return
	}

	posts := GetPublishedPosts()
	var post *models.Post
	var postIndex int

	for i, p := range posts {
		if p.Slug == slug {
			post = &posts[i]
			postIndex = i
			break
		}
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	content, _ := LoadContent()

	data := BlogPostData{
		Post:    *post,
		Content: template.HTML(markdownToHTML(post.Content)),
		Profile: content.Profile,
	}

	// 上一篇 / 下一篇
	if postIndex > 0 {
		data.NextPost = &posts[postIndex-1]
	}
	if postIndex < len(posts)-1 {
		data.PrevPost = &posts[postIndex+1]
	}

	tmpl, err := template.ParseFiles("templates/blog_post.html")
	if err != nil {
		http.Error(w, "模板加载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// BlogArchiveData 归档页数据
type BlogArchiveData struct {
	PostsByYear map[string][]models.Post
	Years       []string
	Tags        []string
	Profile     models.Profile
}

// BlogArchiveHandler 归档页
func BlogArchiveHandler(w http.ResponseWriter, r *http.Request) {
	posts := GetPublishedPosts()
	tags := GetAllTags()
	content, _ := LoadContent()

	// 按年份分组
	postsByYear := make(map[string][]models.Post)
	for _, p := range posts {
		year := p.CreatedAt[:4] // 获取年份
		postsByYear[year] = append(postsByYear[year], p)
	}

	// 年份列表降序
	var years []string
	for year := range postsByYear {
		years = append(years, year)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(years)))

	data := BlogArchiveData{
		PostsByYear: postsByYear,
		Years:       years,
		Tags:        tags,
		Profile:     content.Profile,
	}

	tmpl, err := template.ParseFiles("templates/blog_archive.html")
	if err != nil {
		http.Error(w, "模板加载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// 简单的 Markdown 转 HTML (基础实现)
func markdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var result strings.Builder
	inCodeBlock := false
	inList := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 代码块
		if strings.HasPrefix(trimmed, "```") {
			if inCodeBlock {
				result.WriteString("</code></pre>\n")
				inCodeBlock = false
			} else {
				lang := strings.TrimPrefix(trimmed, "```")
				result.WriteString("<pre><code class=\"language-" + lang + "\">\n")
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			result.WriteString(escapeHTML(line) + "\n")
			continue
		}

		// 空行
		if trimmed == "" {
			if inList {
				result.WriteString("</ul>\n")
				inList = false
			}
			result.WriteString("<br>\n")
			continue
		}

		// 标题
		if strings.HasPrefix(trimmed, "### ") {
			result.WriteString("<h3 class=\"text-xl font-bold text-white mt-6 mb-3\">" + strings.TrimPrefix(trimmed, "### ") + "</h3>\n")
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			result.WriteString("<h2 class=\"text-2xl font-bold text-white mt-8 mb-4\">" + strings.TrimPrefix(trimmed, "## ") + "</h2>\n")
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			result.WriteString("<h1 class=\"text-3xl font-bold text-white mt-8 mb-4\">" + strings.TrimPrefix(trimmed, "# ") + "</h1>\n")
			continue
		}

		// 列表
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			if !inList {
				result.WriteString("<ul class=\"list-disc list-inside space-y-2 text-slate-300\">\n")
				inList = true
			}
			result.WriteString("<li>" + processInline(strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "* ")) + "</li>\n")
			continue
		}

		// 普通段落
		if inList {
			result.WriteString("</ul>\n")
			inList = false
		}
		result.WriteString("<p class=\"text-slate-300 leading-relaxed mb-4\">" + processInline(trimmed) + "</p>\n")
	}

	if inList {
		result.WriteString("</ul>\n")
	}
	if inCodeBlock {
		result.WriteString("</code></pre>\n")
	}

	return result.String()
}

// 处理行内格式
func processInline(text string) string {
	// 粗体
	for strings.Contains(text, "**") {
		start := strings.Index(text, "**")
		end := strings.Index(text[start+2:], "**")
		if end == -1 {
			break
		}
		end += start + 2
		text = text[:start] + "<strong class=\"text-white\">" + text[start+2:end] + "</strong>" + text[end+2:]
	}

	// 斜体
	for strings.Contains(text, "_") {
		start := strings.Index(text, "_")
		end := strings.Index(text[start+1:], "_")
		if end == -1 {
			break
		}
		end += start + 1
		text = text[:start] + "<em>" + text[start+1:end] + "</em>" + text[end+1:]
	}

	// 行内代码
	for strings.Contains(text, "`") {
		start := strings.Index(text, "`")
		end := strings.Index(text[start+1:], "`")
		if end == -1 {
			break
		}
		end += start + 1
		text = text[:start] + "<code class=\"bg-slate-800 px-1 rounded text-neon-blue\">" + text[start+1:end] + "</code>" + text[end+1:]
	}

	// 链接 [text](url)
	for strings.Contains(text, "](") {
		startBracket := strings.Index(text, "[")
		if startBracket == -1 {
			break
		}
		endBracket := strings.Index(text[startBracket:], "]")
		if endBracket == -1 {
			break
		}
		endBracket += startBracket

		if endBracket+1 >= len(text) || text[endBracket+1] != '(' {
			break
		}

		endParen := strings.Index(text[endBracket:], ")")
		if endParen == -1 {
			break
		}
		endParen += endBracket

		linkText := text[startBracket+1 : endBracket]
		linkURL := text[endBracket+2 : endParen]
		text = text[:startBracket] + "<a href=\"" + linkURL + "\" class=\"text-neon-blue hover:underline\">" + linkText + "</a>" + text[endParen+1:]
	}

	return text
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
