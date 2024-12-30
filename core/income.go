package core

import (
	"time"
)

type Income struct {
	ID        int       `json:"id"                required:"true"`
	UserID    int       `json:"user_id"           required:"true"`
	Amount    float64   `json:"amount"            required:"true"`
	Comment   string    `json:"comment,omitempty" required:"false"`
	CreatedAt time.Time `json:"created_at"        required:"true"`
	UpdatedAt time.Time `json:"updated_at"        required:"true"`
}

type Incomes struct {
	Incomes []*Income `json:"incomes" required:"true"`
	Count   int       `json:"count"   required:"true"`
	Sum     float64   `json:"sum"     required:"true"`
}

type IncomeCreateInput struct {
	UserID    int        `json:"user_id"           required:"true"`
	Amount    float64    `json:"amount"            required:"true"`
	Comment   *string    `json:"comment,omitempty" required:"false"`
	CreatedAt *time.Time `json:"created_at"      required:"false"`
	UpdatedAt *time.Time `json:"updated_at"      required:"false"`
}

type IncomeGetInput struct {
	ID     int `json:"id"      required:"true"`
	UserID int `json:"user_id" required:"true"`
}

type IncomeDeleteInput struct {
	ID     int `json:"id"      required:"true"`
	UserID int `json:"user_id" required:"true"`
}

type IncomeUpdateInput struct {
	ID      int      `json:"id"                required:"true"`
	UserID  int      `json:"user_id"           required:"true"`
	Amount  *float64 `json:"amount"            required:"false"`
	Comment *string  `json:"comment,omitempty" required:"false"`
}

type GetIncomesInput struct {
	Limit     int        `json:"limit,omitempty"  default:"100" minimum:"1" maximum:"100" required:"true"`
	Offset    int        `json:"offset,omitempty" default:"0"   minimum:"0" maximum:"1"   required:"true"`
	MinDate   *time.Time `json:"min_date,omitempty"`
	MaxDate   *time.Time `json:"maxDate,omitempty"`
	MinAmount *float64   `json:"min_amount,omitempty"`
	MaxAmount *float64   `json:"maxAmount,omitempty"`
	Users     []int      `json:"users,omitempty"`
}

func (i *IncomeCreateInput) WithUserId(id int) *IncomeCreateInput {
	i.UserID = id
	return i
}

func (i *IncomeGetInput) WithUserId(id int) *IncomeGetInput {
	i.UserID = id
	return i
}

func (i *IncomeDeleteInput) WithUserId(id int) *IncomeDeleteInput {
	i.UserID = id
	return i
}

func (i *IncomeUpdateInput) WithUserId(id int) *IncomeUpdateInput {
	i.UserID = id
	return i
}

func (i *IncomeUpdateInput) WithId(id int) *IncomeUpdateInput {
	i.ID = id
	return i
}

func (i *GetIncomesInput) WithUserId(id ...int) *GetIncomesInput {
	i.Users = append(i.Users, id...)
	return i
}
