package users

import (
	"github.com/google/uuid"
)

type CreateUserReq struct {
	ID       uuid.UUID
	Username string
	Password string
}

type FindCriteria struct {
	ID       *uuid.UUID
	Username *string
}
