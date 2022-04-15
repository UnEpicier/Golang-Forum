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
	ID           string
	Name         string
	CreationDate interface{} // string or time.Time
	Pinned       int
	LastUpdate   interface{} // string or time.Time
}

type Post struct {
	ID           int
	Title        string
	Content      string
	CreationDate interface{} // string or time.Time
	UserID       User
	UpVotes      string
	DownVotes    string
	CategoryId   int
	Pinned       int
	LastUpdate   interface{} // string or time.Time
	Category     Category
}

type Comment struct {
	ID           string
	Content      string
	CreationDate interface{}
	UserID       User
	PostID       Post
	UpVotes      string
	DownVotes    string
	Pinned       int
}

/*
	USER
*/

type User struct {
	ID           int
	Uuid         string
	Username     string
	Email        string
	Password     string
	Role         string
	CreationDate interface{} // string or time.Time
	Biography    string
	LastSeen     interface{} // string or time.Time
	Posts        []Post
	Comments     []Comment
}
