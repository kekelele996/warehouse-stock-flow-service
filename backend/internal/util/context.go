package util

import "context"

type contextKey string

const currentUserKey contextKey = "current_user"

// WithUser 将 JWT 声明写入 context。
func WithUser(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, currentUserKey, claims)
}

// CurrentUser 从 context 读取当前用户声明。
func CurrentUser(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(currentUserKey).(*Claims)
	return claims, ok
}
