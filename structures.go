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
}
