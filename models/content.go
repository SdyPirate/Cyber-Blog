package models

// Content 网站所有内容
type Content struct {
	Profile     Profile     `json:"profile"`
	Quote       Quote       `json:"quote"`
	Education   Education   `json:"education"`
	Experiences []Experience `json:"experiences"`
	Skills      Skills      `json:"skills"`
	Contact     Contact     `json:"contact"`
	Footer      Footer      `json:"footer"`
}

// Profile 个人信息
type Profile struct {
	Name      string `json:"name"`
	Subtitle  string `json:"subtitle"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
	Origin    string `json:"origin"`
	Status    string `json:"status"`
	AvatarUrl string `json:"avatarUrl"`
}

// Quote 语录
type Quote struct {
	Main string `json:"main"`
	Sub  string `json:"sub"`
}

// Education 教育背景
type Education struct {
	Degree      string `json:"degree"`
	School      string `json:"school"`
	Description string `json:"description"`
}

// Experience 工作经历
type Experience struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Period      string `json:"period"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// Skills 技能
type Skills struct {
	Languages []Language  `json:"languages"`
	TechStack []TechStack `json:"techStack"`
}

// Language 语言能力
type Language struct {
	Name    string `json:"name"`
	Level   string `json:"level"`
	Percent int    `json:"percent"`
}

// TechStack 技术栈
type TechStack struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Contact 联系方式
type Contact struct {
	Email    string `json:"email"`
	Github   string `json:"github"`
	Linkedin string `json:"linkedin"`
}

// Footer 页脚
type Footer struct {
	Status string `json:"status"`
}

// Admin 管理员
type Admin struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}
