package utils

import "context"

const CurrentUserKey = "current_user"

func GetCurrentUser(ctx context.Context) string {
	user := ctx.Value(CurrentUserKey).(string)
	if user == "" {
		user = "fake_user"
	}
	return user
}

func SetCurrentUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, CurrentUserKey, user)
}
