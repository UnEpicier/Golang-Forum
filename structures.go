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
	User         User
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
