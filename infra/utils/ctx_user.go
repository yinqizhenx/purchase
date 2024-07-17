package utils

import "context"

const CurrentUserKey = "current_user"

func GetCurrentUser(ctx context.Context) string {
	user := ctx.Value(CurrentUserKey)
	if user == nil {
		return "fake_user"
	}
	return user.(string)
}

func SetCurrentUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, CurrentUserKey, user)
}
