package sessiondto

import (
	"incomster/backend/api/oas"
	"incomster/core"
)

func CoreToOas(in *core.Session) *oas.Session {
	return &oas.Session{
		ID:  in.UserID,
		Jwt: in.Token,
	}
}
