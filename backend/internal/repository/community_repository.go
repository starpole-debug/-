package repository

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CommunityRepository deals with posts, comments and reactions.
type CommunityRepository struct {
	pool *pgxpool.Pool
	once sync.Once
	err  error
}

func NewCommunityRepository(pool *pgxpool.Pool) *CommunityRepository {
	return &CommunityRepository{pool: pool}
}

// ensureLinkColumns is a defensive guard to add link fields if migrations were skipped.
func (r *CommunityRepository) ensureLinkColumns(ctx context.Context) error {
	r.once.Do(func() {
		_, r.err = r.pool.Exec(ctx, `
			ALTER TABLE community_posts ADD COLUMN IF NOT EXISTS link_url TEXT;
			ALTER TABLE community_posts ADD COLUMN IF NOT EXISTS link_type TEXT;
		`)
	})
	return r.err
}

func (r *CommunityRepository) CreatePost(ctx context.Context, post *model.CommunityPost) error {
	if err := r.ensureLinkColumns(ctx); err != nil {
		return err
	}
	if post.ID == "" {
		post.ID = uuid.NewString()
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO community_posts(id, author_id, title, content, style_id, topic_ids, visibility, link_url, link_type)
        VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
        RETURNING created_at, updated_at
    `, post.ID, post.AuthorID, post.Title, post.Content, post.StyleID, post.TopicIDs, post.Visibility, post.LinkURL, post.LinkType)
	return row.Scan(&post.CreatedAt, &post.UpdatedAt)
}

func (r *CommunityRepository) ListPosts(ctx context.Context, limit int) ([]model.CommunityPost, error) {
	return r.ListPostsFiltered(ctx, limit, "latest", "all", "", "", "")
}

func (r *CommunityRepository) ListPostsFiltered(ctx context.Context, limit int, sort, filter, userID, authorID, search string) ([]model.CommunityPost, error) {
	if err := r.ensureLinkColumns(ctx); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 20
	}

	baseQuery := `
		SELECT
			p.id,
			p.author_id,
			u.nickname,
			u.avatar_url,
			p.title,
			p.content,
			p.style_id,
			p.topic_ids,
			p.visibility,
			COALESCE(p.link_url, '') AS link_url,
			COALESCE(p.link_type, '') AS link_type,
			p.created_at,
			p.updated_at
		FROM community_posts p
		LEFT JOIN users u ON p.author_id = u.id
	`
	whereClause := " WHERE p.visibility = 'public'"
	orderBy := " ORDER BY p.created_at DESC"
	args := []interface{}{}
	argIdx := 1

	if authorID != "" {
		whereClause += fmt.Sprintf(" AND p.author_id = $%d", argIdx)
		args = append(args, authorID)
		argIdx++
	}

	if filter == "following" && strings.TrimSpace(userID) == "" {
		filter = "all"
	}

	if filter == "following" && userID != "" {
		baseQuery += fmt.Sprintf(" JOIN follows f ON p.author_id = f.following_id AND f.follower_id = $%d", argIdx)
		args = append(args, userID)
		argIdx++
	}

	if sort == "hot" {
		// Hot: (likes * 2 + comments)
		orderBy = `
			ORDER BY (
				(SELECT COUNT(*) FROM community_reactions r WHERE r.post_id = p.id AND r.type = 'like') * 2 +
				(SELECT COUNT(*) FROM community_comments c WHERE c.post_id = p.id)
			) DESC, p.created_at DESC
		`
	}
	if q := strings.TrimSpace(search); q != "" {
		whereClause += fmt.Sprintf(" AND (p.title ILIKE $%d OR p.content ILIKE $%d)", argIdx, argIdx+1)
		pattern := "%" + q + "%"
		args = append(args, pattern, pattern)
		argIdx += 2
	}

	query := baseQuery + whereClause + orderBy + fmt.Sprintf(" LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.CommunityPost
	for rows.Next() {
		var post model.CommunityPost
		var topicIDs []string
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.AuthorName, &post.AuthorAvatar, &post.Title, &post.Content, &post.StyleID, &topicIDs, &post.Visibility, &post.LinkURL, &post.LinkType, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		post.TopicIDs = topicIDs
		posts = append(posts, post)
	}
	if len(posts) > 0 {
		attachments, _ := r.listAttachmentsByPosts(ctx, posts)
		for i := range posts {
			posts[i].Attachments = attachments[posts[i].ID]
		}
	}
	return posts, nil
}

func (r *CommunityRepository) FindPost(ctx context.Context, id string) (*model.CommunityPost, error) {
	if err := r.ensureLinkColumns(ctx); err != nil {
		return nil, err
	}
	row := r.pool.QueryRow(ctx, `
        SELECT
			p.id,
			p.author_id,
			u.nickname,
			u.avatar_url,
			p.title,
			p.content,
			p.style_id,
			p.topic_ids,
			p.visibility,
			COALESCE(p.link_url, '') AS link_url,
			COALESCE(p.link_type, '') AS link_type,
			p.created_at,
			p.updated_at
        FROM community_posts p
        LEFT JOIN users u ON p.author_id = u.id
        WHERE p.id = $1
    `, id)
	var post model.CommunityPost
	var topicIDs []string
	if err := row.Scan(&post.ID, &post.AuthorID, &post.AuthorName, &post.AuthorAvatar, &post.Title, &post.Content, &post.StyleID, &topicIDs, &post.Visibility, &post.LinkURL, &post.LinkType, &post.CreatedAt, &post.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	post.TopicIDs = topicIDs
	attachments, _ := r.listAttachmentsByPosts(ctx, []model.CommunityPost{post})
	post.Attachments = attachments[post.ID]
	return &post, nil
}

func (r *CommunityRepository) CreateComment(ctx context.Context, comment *model.CommunityComment) error {
	if comment.ID == "" {
		comment.ID = uuid.NewString()
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO community_comments(id, post_id, author_id, content)
        VALUES($1,$2,$3,$4)
        RETURNING created_at
    `, comment.ID, comment.PostID, comment.AuthorID, comment.Content)
	return row.Scan(&comment.CreatedAt)
}

func (r *CommunityRepository) ListComments(ctx context.Context, postID string) ([]model.CommunityComment, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT c.id, c.post_id, c.author_id, u.nickname, u.avatar_url, c.content, c.visibility, c.created_at
        FROM community_comments c
        LEFT JOIN users u ON c.author_id = u.id
        WHERE c.post_id = $1 AND c.visibility = 'public' 
        ORDER BY c.created_at ASC
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []model.CommunityComment
	for rows.Next() {
		var comment model.CommunityComment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.AuthorName, &comment.AuthorAvatar, &comment.Content, &comment.Visibility, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommunityRepository) ToggleReaction(ctx context.Context, postID, userID, reactionType string) (bool, error) {
	if strings.TrimSpace(postID) == "" || strings.TrimSpace(userID) == "" || strings.TrimSpace(reactionType) == "" {
		return false, errors.New("missing reaction attributes")
	}
	var exists int
	err := r.pool.QueryRow(ctx, `SELECT 1 FROM community_reactions WHERE post_id = $1 AND user_id = $2 AND type = $3`, postID, userID, reactionType).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		return false, err
	}

	if err == nil {
		_, err = r.pool.Exec(ctx, `DELETE FROM community_reactions WHERE post_id = $1 AND user_id = $2 AND type = $3`, postID, userID, reactionType)
		return false, err
	}

	_, err = r.pool.Exec(ctx, `
        INSERT INTO community_reactions(id, post_id, user_id, type, updated_at)
        VALUES($1,$2,$3,$4, now())
        ON CONFLICT (post_id, user_id, type) DO UPDATE
        SET updated_at = now()
    `, uuid.NewString(), postID, userID, reactionType)
	return err == nil, err
}

// CreateAttachments stores uploaded file metadata for a given post.
func (r *CommunityRepository) CreateAttachments(ctx context.Context, postID string, urls []string) error {
	if strings.TrimSpace(postID) == "" || len(urls) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, url := range urls {
		u := strings.TrimSpace(url)
		if u == "" {
			continue
		}
		filename := path.Base(u)
		batch.Queue(`
            INSERT INTO community_attachments(id, post_id, file_url, file_type, file_name, file_size)
            VALUES($1,$2,$3,$4,$5,$6)
        `, uuid.NewString(), postID, u, inferFileType(filename), filename, int64(0))
	}
	if batch.Len() == 0 {
		return nil
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	for i := 0; i < batch.Len(); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}

func (r *CommunityRepository) listAttachmentsByPosts(ctx context.Context, posts []model.CommunityPost) (map[string][]model.CommunityAttachment, error) {
	ids := make([]string, 0, len(posts))
	for _, p := range posts {
		if strings.TrimSpace(p.ID) != "" {
			ids = append(ids, p.ID)
		}
	}
	if len(ids) == 0 {
		return map[string][]model.CommunityAttachment{}, nil
	}
	rows, err := r.pool.Query(ctx, `
        SELECT id, post_id, file_url, file_type, file_name, file_size, created_at
        FROM community_attachments
        WHERE post_id = ANY($1)
        ORDER BY created_at ASC
    `, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string][]model.CommunityAttachment)
	for rows.Next() {
		var att model.CommunityAttachment
		if err := rows.Scan(&att.ID, &att.PostID, &att.FileURL, &att.FileType, &att.FileName, &att.FileSize, &att.CreatedAt); err != nil {
			return nil, err
		}
		result[att.PostID] = append(result[att.PostID], att)
	}
	return result, nil
}

func inferFileType(filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return "image"
	default:
		return "file"
	}
}

func (r *CommunityRepository) CountReactions(ctx context.Context, postID string) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT type, count(*) FROM community_reactions WHERE post_id = $1 GROUP BY type
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := make(map[string]int)
	for rows.Next() {
		var typ string
		var total int
		if err := rows.Scan(&typ, &total); err != nil {
			return nil, err
		}
		counts[typ] = total
	}
	return counts, nil
}

func (r *CommunityRepository) UpdateVisibility(ctx context.Context, id, visibility string) error {
	if err := r.ensureLinkColumns(ctx); err != nil {
		return err
	}
	_, err := r.pool.Exec(ctx, `UPDATE community_posts SET visibility = $2, updated_at = now() WHERE id = $1`, id, visibility)
	return err
}

func (r *CommunityRepository) ListPostsAdmin(ctx context.Context, query, visibility string, limit, offset int) ([]model.CommunityPost, error) {
	if err := r.ensureLinkColumns(ctx); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	base := `
        SELECT
			id,
			author_id,
			title,
			content,
			style_id,
			topic_ids,
			visibility,
			COALESCE(link_url, '') AS link_url,
			COALESCE(link_type, '') AS link_type,
			created_at,
			updated_at
        FROM community_posts
        WHERE 1=1
    `
	args := []interface{}{}
	sb := strings.Builder{}
	sb.WriteString(base)
	idx := 1
	if v := strings.TrimSpace(visibility); v != "" && v != "all" {
		sb.WriteString(fmt.Sprintf(" AND visibility = $%d", idx))
		args = append(args, v)
		idx++
	}
	if q := strings.TrimSpace(query); q != "" {
		sb.WriteString(fmt.Sprintf(" AND (title ILIKE $%d OR content ILIKE $%d OR id::TEXT ILIKE $%d)", idx, idx+1, idx+2))
		pattern := "%" + q + "%"
		args = append(args, pattern, pattern, pattern)
		idx += 3
	}
	sb.WriteString(fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1))
	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []model.CommunityPost
	for rows.Next() {
		var post model.CommunityPost
		var topicIDs []string
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.StyleID, &topicIDs, &post.Visibility, &post.LinkURL, &post.LinkType, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		post.TopicIDs = topicIDs
		posts = append(posts, post)
	}
	if len(posts) > 0 {
		attachments, _ := r.listAttachmentsByPosts(ctx, posts)
		for i := range posts {
			posts[i].Attachments = attachments[posts[i].ID]
		}
	}
	return posts, nil
}

func (r *CommunityRepository) ListCommentsAdmin(ctx context.Context, postID, visibility string, limit, offset int) ([]model.CommunityComment, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	base := `
        SELECT id, post_id, author_id, content, visibility, created_at
        FROM community_comments
        WHERE 1=1
    `
	args := []interface{}{}
	sb := strings.Builder{}
	sb.WriteString(base)
	idx := 1
	if post := strings.TrimSpace(postID); post != "" {
		sb.WriteString(fmt.Sprintf(" AND post_id = $%d", idx))
		args = append(args, post)
		idx++
	}
	if vis := strings.TrimSpace(visibility); vis != "" && vis != "all" {
		sb.WriteString(fmt.Sprintf(" AND visibility = $%d", idx))
		args = append(args, vis)
		idx++
	}
	sb.WriteString(fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1))
	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []model.CommunityComment
	for rows.Next() {
		var comment model.CommunityComment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.Visibility, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommunityRepository) UpdateCommentVisibility(ctx context.Context, id, visibility string) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE community_comments SET visibility = $2 WHERE id = $1
    `, id, visibility)
	return err
}

func (r *CommunityRepository) ListFavorites(ctx context.Context, userID string, limit int) ([]model.CommunityPost, error) {
	if err := r.ensureLinkColumns(ctx); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx, `
        SELECT
			p.id,
			p.author_id,
			u.nickname,
			u.avatar_url,
			p.title,
			p.content,
			p.style_id,
			p.topic_ids,
			p.visibility,
			COALESCE(p.link_url, '') AS link_url,
			COALESCE(p.link_type, '') AS link_type,
			p.created_at,
			p.updated_at
        FROM community_reactions fav
        JOIN community_posts p ON p.id = fav.post_id
        LEFT JOIN users u ON p.author_id = u.id
        WHERE fav.user_id = $1
          AND fav.type = 'favorite'
          AND p.visibility = 'public'
        ORDER BY fav.updated_at DESC
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]model.CommunityPost, 0)
	for rows.Next() {
		var post model.CommunityPost
		var topicIDs []string
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.AuthorName, &post.AuthorAvatar, &post.Title, &post.Content, &post.StyleID, &topicIDs, &post.Visibility, &post.LinkURL, &post.LinkType, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		post.TopicIDs = topicIDs
		posts = append(posts, post)
	}
	if len(posts) > 0 {
		attachments, _ := r.listAttachmentsByPosts(ctx, posts)
		for i := range posts {
			posts[i].Attachments = attachments[posts[i].ID]
		}
	}
	return posts, nil
}

func (r *CommunityRepository) RecordView(ctx context.Context, userID, postID string) error {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(postID) == "" {
		return nil
	}
	if err := r.ensureLinkColumns(ctx); err != nil {
		return err
	}
	_, err := r.pool.Exec(ctx, `
        INSERT INTO community_post_views(user_id, post_id)
        VALUES($1,$2)
        ON CONFLICT (user_id, post_id) DO UPDATE
        SET view_count = community_post_views.view_count + 1,
            last_viewed = now()
    `, userID, postID)
	return err
}

func (r *CommunityRepository) ListRecentViews(ctx context.Context, userID string, limit int) ([]model.CommunityPost, error) {
	if err := r.ensureLinkColumns(ctx); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.pool.Query(ctx, `
        SELECT
			p.id,
			p.author_id,
			u.nickname,
			u.avatar_url,
			p.title,
			p.content,
			p.style_id,
			p.topic_ids,
			p.visibility,
			COALESCE(p.link_url, '') AS link_url,
			COALESCE(p.link_type, '') AS link_type,
			p.created_at,
			p.updated_at
        FROM community_post_views v
        JOIN community_posts p ON p.id = v.post_id
        LEFT JOIN users u ON p.author_id = u.id
        WHERE v.user_id = $1 AND p.visibility = 'public'
        ORDER BY v.last_viewed DESC
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]model.CommunityPost, 0)
	for rows.Next() {
		var post model.CommunityPost
		var topicIDs []string
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.AuthorName, &post.AuthorAvatar, &post.Title, &post.Content, &post.StyleID, &topicIDs, &post.Visibility, &post.LinkURL, &post.LinkType, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		post.TopicIDs = topicIDs
		posts = append(posts, post)
	}
	if len(posts) > 0 {
		attachments, _ := r.listAttachmentsByPosts(ctx, posts)
		for i := range posts {
			posts[i].Attachments = attachments[posts[i].ID]
		}
	}
	return posts, nil
}

func (r *CommunityRepository) CountPostsByAuthor(ctx context.Context, authorID string) (int, error) {
	if err := r.ensureLinkColumns(ctx); err != nil {
		return 0, err
	}
	row := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM community_posts WHERE author_id = $1`, authorID)
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *CommunityRepository) ListPostsByAuthor(ctx context.Context, authorID string, limit int) ([]model.CommunityPost, error) {
	if err := r.ensureLinkColumns(ctx); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.pool.Query(ctx, `
        SELECT
			p.id,
			p.author_id,
			u.nickname,
			u.avatar_url,
			p.title,
			p.content,
			p.style_id,
			p.topic_ids,
			p.visibility,
			COALESCE(p.link_url, '') AS link_url,
			COALESCE(p.link_type, '') AS link_type,
			p.created_at,
			p.updated_at
        FROM community_posts p
        LEFT JOIN users u ON p.author_id = u.id
        WHERE p.author_id = $1
        ORDER BY p.created_at DESC
        LIMIT $2
    `, authorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []model.CommunityPost
	for rows.Next() {
		var post model.CommunityPost
		var topicIDs []string
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.AuthorName, &post.AuthorAvatar, &post.Title, &post.Content, &post.StyleID, &topicIDs, &post.Visibility, &post.LinkURL, &post.LinkType, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		post.TopicIDs = topicIDs
		posts = append(posts, post)
	}
	if len(posts) > 0 {
		attachments, _ := r.listAttachmentsByPosts(ctx, posts)
		for i := range posts {
			posts[i].Attachments = attachments[posts[i].ID]
		}
	}
	return posts, nil
}

func (r *CommunityRepository) FollowUser(ctx context.Context, followerID, followingID string) error {
	_, err := r.pool.Exec(ctx, `
        INSERT INTO follows(follower_id, following_id)
        VALUES($1,$2)
        ON CONFLICT DO NOTHING
    `, followerID, followingID)
	return err
}

func (r *CommunityRepository) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`, followerID, followingID)
	return err
}

func (r *CommunityRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	row := r.pool.QueryRow(ctx, `SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2`, followerID, followingID)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *CommunityRepository) CountFollowers(ctx context.Context, userID string) (int, error) {
	row := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM follows WHERE following_id = $1`, userID)
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *CommunityRepository) CountFollowing(ctx context.Context, userID string) (int, error) {
	row := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM follows WHERE follower_id = $1`, userID)
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *CommunityRepository) CountFavorites(ctx context.Context, userID string) (int, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM community_reactions fav
        JOIN community_posts p ON p.id = fav.post_id
        WHERE fav.user_id = $1
          AND fav.type = 'favorite'
          AND p.visibility = 'public'
    `, userID)
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
