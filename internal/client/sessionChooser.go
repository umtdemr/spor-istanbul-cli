package client

import (
	"fmt"
	"golang.org/x/net/context"
)

func ChooseSession() {
	sessionId, isLoggedIn := Login("id", "password") // not real don't worry :)
	if !isLoggedIn {
		return
	}
	ctx := context.WithValue(context.Background(), ctx_session_id, sessionId)
	fmt.Println(ctx.Value(ctx_session_id))
}
