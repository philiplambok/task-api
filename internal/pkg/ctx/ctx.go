package ctx

import "context"

// contextKey is a private type for context keys to prevent collisions
type contextKey string

// Context keys for storing values in request context
const (
	userEmailKey contextKey = "user_email"
	userIDKey    contextKey = "user_id"
)

// InjectUserEmail adds user email to the context
func InjectUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, userEmailKey, email)
}

// ExtractUserEmail retrieves the user email from the context
// Returns the email and a boolean indicating whether it was found
func ExtractUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}

// InjectUserID adds user ID to the context
func InjectUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// ExtractUserID retrieves the user ID from the context
// Returns the ID and a boolean indicating whether it was found
func ExtractUserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}
