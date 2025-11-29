# Admin Backend Notes

## New APIs

- `POST /api/admin/posts/:id/hide` – hide a community post.
- `POST /api/admin/posts/:id/limit` – set a post visibility to `limited`.
- `POST /api/admin/posts/:id/unhide` / `POST /api/admin/posts/:id/unlimit` – restore public visibility.
- `GET /api/admin/posts` – list posts with query/visibility filters.
- `GET /api/admin/users` – search and paginate users.
- `POST /api/admin/users` – create a new user (optionally admin).
- `POST /api/admin/users/:id/ban` / `POST /api/admin/users/:id/unban` – toggle ban flag.
- `DELETE /api/admin/users/:id` – soft delete a user.
- `GET /api/admin/comments` – list comments with filters.
- `POST /api/admin/comments/:id/hide` – hide a comment.
- `DELETE /api/admin/comments/:id` – mark a comment as deleted.

All routes above require a valid admin token issued via `POST /api/admin/login`.

## Database Schema Changes

- `backend/migrations/0002_admin_moderation.sql` adds `users.is_banned`, `users.deleted_at`, and `community_comments.visibility` plus supporting indexes so moderation state can be persisted and queried efficiently.

## Authentication Impact

- `internal/service/auth` rejects logins for any user marked `is_banned = true` or `deleted_at IS NOT NULL`, ensuring suspended accounts (admin-created bans or deletions) cannot access the system.

## Environment Variables

- `AUTH_JWT_SECRET` – secret used to sign user JWTs (falls back to `JWT_SECRET` for backwards compatibility).
- `DEFAULT_MODEL_ID` – optional explicit model id that chat sessions will prefer when no model is requested.
