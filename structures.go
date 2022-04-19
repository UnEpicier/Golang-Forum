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
	User         User
	CategoryId   int
	Pinned       int
	LastUpdate   interface{} // string or time.Time
	Likes        int
	Dislikes     int
	CommentNB    int
	Category     Category
}

type Comment struct {
	ID           string
	Content      string
	CreationDate interface{}
	User         User
	PostID       Post
	Pinned       int
	Likes        int
	Dislikes     int
}

type Write struct {
	Categories string
	Action     string
	Post       Post
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
	Stats   Stats
	Users   []User
	Reports Reports
}

type Stats struct {
	Categories    []AD_Categories
	CatActivities []AD_CatActivities
	Inscriptions  []AD_Inscription
	Users         AD_Users
	Forum         AD_Forum
}

type AD_Categories struct {
	Name  string
	Count int
}

type AD_CatActivities struct {
	Name     string
	Activity [12]int
}

type AD_Inscription struct {
	Month string
	Count int
}

type AD_Users struct {
	Total   int
	Admins  int
	Mods    int
	Members int
}

type AD_Forum struct {
	Categories int
	Posts      int
	Comments   int
}

/*
	REPORTS
*/
type Reports struct {
	Users    []Report
	Posts    []Report
	Comments []Report
}

type Report struct {
	ID           int
	Type         string
	Reason       string
	CreationDate interface{}
	User         User
	Post         Post
	Comment      Comment
}
