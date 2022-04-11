package forum

import "time"

type Page struct {
	Logged  bool
	Error   string
	Content interface{}
}

/*
	FORUM
*/
// Categories
type Forum struct {
	Categories []Category
	Error      string
}

type Category struct {
	Uuid    string
	Name    string
	Link    string
	Created time.Time
	Pinned  string
}

type Post struct {
	Uuid     string
	Title    string
	Content  string
	Created  string
	User     User
	Likes    int
	Dislikes int
	Category string
}

type Comment struct {
	Uuid     string
	Content  string
	Created  string
	User     User
	Post     string
	Likes    int
	Dislikes int
}

/*
	USER
*/

type User struct {
	Uuid        string
	Username    string
	Email       string
	Password    string
	Role        string
	Joined      string
	Description string
	Posts       []Post
	Comments    []Comment
}
