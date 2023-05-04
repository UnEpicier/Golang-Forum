/*
 ! SQLITE DATABASE
 * DB Name: forum.db
 */
CREATE TABLE user(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL,
  profile_pic TEXT DEFAULT '/assets/profile/default.png',
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  role TEXT NOT NULL,
  creation_date DATETIME NOT NULL,
  biography TEXT,
  last_seen DATETIME NOT NULL
);

CREATE TABLE category(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  creation_date DATETIME NOT NULL,
  pinned INTEGER DEFAULT 0,
  last_update DATETIME NOT NULL
);

CREATE TABLE post(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  creation_date DATETIME NOT NULL,
  user_id TEXT NOT NULL,
  category_id INTEGER NOT NULL,
  pinned INTEGER DEFAULT '0',
  last_update DATETIME NOT NULL,
  FOREIGN KEY(last_update) REFERENCES comment(id) FOREIGN KEY(user_id) REFERENCES user(id) FOREIGN KEY(category_id) REFERENCES category(id)
);

CREATE TABLE comment(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  content TEXT NOT NULL,
  creation_date DATETIME NOT NULL,
  user_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  pinned INTEGER DEFAULT '0',
  FOREIGN KEY(user_id) REFERENCES user(id) FOREIGN KEY(post_id) REFERENCES post(id)
);

CREATE TABLE vote (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  type TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  post_id INTEGER DEFAULT NULL,
  comment_id INTEGER DEFAULT NULL,
  FOREIGN KEY(user_id) REFERENCES user(id) FOREIGN KEY(post_id) REFERENCES post(id) FOREIGN KEY(comment_id) REFERENCES comment(id)
);

CREATE TABLE report (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  type TEXT NOT NULL,
  reason TEXT NOT NULL,
  creation_date DATETIME NOT NULL,
  user_id INTEGER NOT NULL,
  post_id INTEGER DEFAULT NULL,
  comment_id INTEGER DEFAULT NULL,
  FOREIGN KEY(user_id) REFERENCES user(id) FOREIGN KEY(post_id) REFERENCES post(id) FOREIGN KEY(comment_id) REFERENCES comment(id)
);
