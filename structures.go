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

/*
	ADMIN
*/

type Admin struct {
	Stats Stats
}

type Stats struct {
	Categories    []AD_Categories
	CatActivities []AD_CatActivities
	Inscriptions  []AD_Inscription
}

type AD_Categories struct {
	Name  string
	Count int
}

type AD_CatActivities struct {
	Name string
}

type AD_Inscription struct {
	Month string
	Count int
}
