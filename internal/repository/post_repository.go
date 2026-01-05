package repository

import (
	"database/sql"
	"fmt"
	"time"

	"traingolang/internal/util"
)

type Post struct {
	ID          int64
	ImageID     sql.NullInt64
	Name        string
	Description string
	Topic       string
	Prompt      sql.NullString
	HotLevel    int8
	HotAt       sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type PostSearchResponse struct {
	ID          int64   `json:"id"`
	ImageURL    *string `json:"image_url"`
	BlurURL     *string `json:"blur_url"`
	TinyBlurURL *string `json:"tiny_blur_url"`

	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Topic       string     `json:"topic"`
	HotLevel    int8       `json:"hot_level"`
	HotAt       *time.Time `json:"hot_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
type PostOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type PostRepo interface {
	CreatePost(post *Post) (int64, error)

	SearchPosts(
		name string,
		Topic string,
		hotLevel *int8,
		page, pageSize int,
	) (*util.PaginatedResponse[PostSearchResponse], error)

	GetByID(id int64) (*Post, error)
	UpdatePost(post *Post) error
	DeletePost(id int64) error
	GetPostOptions() ([]PostOption, error)
	// ExistsByTopic(topic string, excludeID *int64) (bool, error)
}

type postRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) PostRepo {
	return &postRepo{db: db}
}

func (r *postRepo) CreatePost(post *Post) (int64, error) {
	query := `
		INSERT INTO posts (
			image_id,
			name,
			description,
			topic,
			prompt,
			hot_level,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(
		query,
		post.ImageID,
		post.Name,
		post.Description,
		post.Topic,
		post.Prompt,
		post.HotLevel,
		post.CreatedAt,
		post.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	post.ID = id
	return id, nil
}

func (r *postRepo) SearchPosts(
	name string,
	topic string,
	hotLevel *int8,
	page, pageSize int,
) (*util.PaginatedResponse[PostSearchResponse], error) {

	offset, limit := util.NewPagination(page, pageSize)

	where := "1=1"
	args := []interface{}{}

	if name != "" {
		where += fmt.Sprintf(" AND p.name ILIKE $%d", len(args)+1)
		args = append(args, "%"+name+"%")
	}
	if topic != "" {
		where += fmt.Sprintf(" AND p.topic = $%d", len(args)+1)
		args = append(args, topic)
	}

	if hotLevel != nil {
		where += fmt.Sprintf(" AND p.hot_level = $%d", len(args)+1)
		args = append(args, *hotLevel)
	}

	limitIdx := len(args) + 1
	offsetIdx := len(args) + 2

	query := fmt.Sprintf(`
		SELECT
			p.id,
			i.image_url,
			i.blur_url,
			i.tiny_blur_url,
			p.name,
			p.description,
			p.topic,
			p.hot_level,
			p.hot_at,
			p.created_at
		FROM posts p
		LEFT JOIN images i ON p.image_id = i.id
		WHERE %s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, limitIdx, offsetIdx)

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM posts p
		WHERE %s
	`, where)

	return util.Paginate(
		r.db,
		query,
		countQuery,
		args,
		offset,
		limit,
		func(rows *sql.Rows) (*PostSearchResponse, error) {
			var p PostSearchResponse
			err := rows.Scan(
				&p.ID,
				&p.ImageURL,
				&p.BlurURL,
				&p.TinyBlurURL,
				&p.Name,
				&p.Description,
				&p.Topic,
				&p.HotLevel,
				&p.HotAt,
				&p.CreatedAt,
			)
			if err != nil {
				return nil, err
			}
			return &p, nil
		},
	)
}

func (r *postRepo) GetByID(id int64) (*Post, error) {
	query := `
		SELECT
			id,
			image_id,
			name,
			description,
			topic,
			prompt,
			hot_level,
			hot_at,
			created_at,
			updated_at
		FROM posts
		WHERE id = $1
	`

	var p Post
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.ImageID,
		&p.Name,
		&p.Description,
		&p.Topic,
		&p.Prompt,
		&p.HotLevel,
		&p.HotAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *postRepo) UpdatePost(post *Post) error {
	query := `
		UPDATE posts
		SET
			image_id = $1,
			name = $2,
			description = $3,
			topic = $4,
			prompt = $5,
			hot_level = $6,
			updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.Exec(
		query,
		post.ImageID,
		post.Name,
		post.Description,
		post.Topic,
		post.Prompt,
		post.HotLevel,
		post.UpdatedAt,
		post.ID,
	)

	return err
}

func (r *postRepo) DeletePost(id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *postRepo) GetPostOptions() ([]PostOption, error) {
	query := `
		SELECT DISTINCT topic
		FROM posts
		WHERE topic IS NOT NULL
		ORDER BY topic ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []PostOption

	for rows.Next() {
		var topic string
		if err := rows.Scan(&topic); err != nil {
			return nil, err
		}

		result = append(result, PostOption{
			Value: topic,
			Label: topic,
		})
	}

	return result, nil
}

func (r *postRepo) ExistsByTopic(topic string, excludeID *int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM posts
			WHERE topic = $1
	`

	args := []interface{}{topic}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	query += ")"

	var exists bool
	err := r.db.QueryRow(query, args...).Scan(&exists)
	return exists, err
}
