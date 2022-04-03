package forum

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
	Uuid string
	Name string
	Link string
}

type Post struct {
	Uuid     string
	Title    string
	Content  string
	Created  string
	User     string
	Likes    int
	Dislikes int
	Category string
}

type Comment struct {
	Uuid     string
	Content  string
	Created  string
	User     string
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
