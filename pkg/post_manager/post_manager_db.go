package post_manager

import (
	"errors"
	"time"

	"github.com/ani5msr/microservices-project/pkg/db_utils"
	om "github.com/ani5msr/microservices-project/pkg/object_model"

	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

const pageSize = 10

type DbPostStore struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

const dbName = "post_manager"

func NewDbPostStore(host string, port int, username string, password string) (store *DbPostStore, err error) {
	db, err := db_utils.EnsureDB(host, port, username, password, dbName)
	if err != nil {
		return
	}

	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	err = createSchema(db)
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}
	store = &DbPostStore{db, sb}
	return
}

func createSchema(db *sql.DB) (err error) {
	schema := `
		DO $$ BEGIN
    		CREATE TYPE post_status AS ENUM ('pending', 'valid', 'invalid');
		EXCEPTION
    		WHEN duplicate_object THEN null;
		END $$;		
        CREATE TABLE IF NOT EXISTS posts (
          id SERIAL   PRIMARY KEY,
		  username    TEXT,
          url TEXT    NOT NULL,
          title TEXT  NOT NULL,
		  description TEXT,
	      status      link_status NOT NULL DEFAULT 'pending',		
		  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP          
        );
		CREATE INDEX IF NOT EXISTS posts_username_idx ON posts(username);
        CREATE TABLE IF NOT EXISTS tags (
          id SERIAL PRIMARY KEY,
          post_id   INTEGER REFERENCES posts(id) ON DELETE CASCADE,			
          name      TEXT		  
        );
        CREATE INDEX IF NOT EXISTS tags_name_idx ON tags(name);
    `

	_, err = db.Exec(schema)
	return
}

func (s *DbPostStore) GetPost(request om.GetPostsRequest) (result om.GetPostsResult, err error) {
	q := s.sb.Select("*").From("posts")
	if request.Tag != "" {
		q = q.Join("tags ON post.id = tags.post_id")
	}
	q = q.Where(sq.Eq{"username": request.Username}).OrderBy("created_at").Limit(pageSize)
	if request.StartToken != "" {
		var createdAt time.Time
		createdAt, err = time.Parse(time.RFC3339, request.StartToken)
		if err != nil {
			return
		}

		q = q.Where(sq.Gt{"created_at": createdAt})
	}
	if request.Tag != "" {
		q = q.Where(sq.Eq{"tag": request.Tag})
	}

	rows, err := q.RunWith(s.db).Query()
	if err != nil {
		return result, err
	}

	posts := map[string]om.Post{}

	var post om.Post
	var id int
	var tag_id int
	var tag_name string
	var username string
	for rows.Next() {
		if request.Tag != "" {
			err = rows.Scan(&id, &username, &post.Url, &post.Title, &post.Description, &post.Status, &post.CreatedAt, &post.UpdatedAt, &tag_id, &id, &tag_name)
		} else {
			err = rows.Scan(&id, &username, &post.Url, &post.Title, &post.Description, &post.Status, &post.CreatedAt, &post.UpdatedAt)
		}
		if err != nil {
			return
		}

		_, ok := posts[post.Url]
		if !ok {
			posts[post.Url] = post
			result.Posts = append(result.Posts, post)
		}
	}

	if len(result.Posts) == pageSize {
		result.NextPageToken = post.CreatedAt.UTC().Format(time.RFC3339)
	}
	return
}

func (s *DbPostStore) AddPost(request om.AddPostRequest) (post *om.Post, err error) {
	post = &om.Post{
		Tags: map[string]bool{},
	}
	cmd := s.sb.Insert("posts").Columns("username", "url", "title", "description").
		Values(request.Username, request.Url, request.Title, request.Description)
	_, err = cmd.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	q := s.sb.Select("*").From("posts").Where(sq.Eq{"username": request.Username, "url": request.Url})
	var post_id int
	var username string
	row := q.RunWith(s.db).QueryRow()
	err = row.Scan(&post_id, &username, &post.Url, &post.Title, &post.Description, &post.Status, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return
	}

	for t, _ := range request.Tags {
		cmd := s.sb.Insert("tags").Columns("post_id", "name").Values(post_id, t)
		_, err = cmd.RunWith(s.db).Exec()
		if err != nil {
			return
		}

		post.Tags[t] = true
	}

	return
}

func (s *DbPostStore) UpdatePost(request om.UpdatePostRequest) (post *om.Post, err error) {
	q := s.sb.Update("posts").Where(sq.Eq{"username": request.Username, "url": request.Url})
	if request.Title != "" {
		q = q.Set("title", request.Title)
	}

	if request.Description != "" {
		q = q.Set("description", request.Description)
	}

	q = q.Suffix("RETURNING \"id\"")
	res, err := q.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected != 1 {
		err = errors.New("update failed")
	}

	var post_id int
	q.QueryRow().Scan(&post_id)

	for t, _ := range request.RemoveTags {
		_, err = s.sb.Delete("tags").Where(sq.Eq{"post_id": post_id, "name": t}).RunWith(s.db).Exec()
		if err != nil {
			return
		}
	}

	for t, _ := range request.AddTags {
		_, err = s.sb.Insert("tags").Columns("post_id", "name").Values(post_id, t).RunWith(s.db).Exec()
		if err != nil {
			return
		}
	}

	return
}

func (s *DbPostStore) DeletePost(username string, url string) (err error) {
	_, err = s.sb.Delete("posts").Where(sq.Eq{"username": username, "url": url}).RunWith(s.db).Exec()
	return
}

func (s *DbPostStore) SetPostStatus(username string, url string, status om.PostStatus) (err error) {
	m := map[om.PostStatus]string{
		om.PostStatusPending: "pending",
		om.PostStatusValid:   "valid",
		om.PostStatusInvalid: "invalid",
	}

	q := s.sb.Update("posts").Where(sq.Eq{"username": username, "url": url})
	q = q.Set("status", m[status])
	res, err := q.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected != 1 {
		err = errors.New("update failed")
	}

	return
}
