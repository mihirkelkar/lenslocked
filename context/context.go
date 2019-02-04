package context

import (
	"context"

	"github.com/mihirkelkar/lenslocked.com/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

//This function accepts a context package variable, then sets user as
//value and then returns the context back. This whole package is essentailly
//a wrapper on the standard context package since the context package does not
//offer type safety.
func WithUser(ctx context.Context, user *models.User) context.Context {
	//why did we define a userKey and a privateKey data type?
	//the reason is context uses a key's data type and value as the "key"
	//so to be absolutely sure no one with the "user" key name can override us,
	//we created a datatype that is backed by string but isn't actually string
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
