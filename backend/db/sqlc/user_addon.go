package db

import (
	"context"
	"database/sql"
	"time"
)

func (q *Queries) GetPhoneOrCreateUser(ctx context.Context, phone string) (User, error) {
	user, err := q.GetUserByPhone(ctx, phone)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exist, create a new one
			// You need to provide the necessary fields for a new user
			newUser := User{
				Phone:     phone,
				CreatedAt: time.Now(),
			}
			return q.CreateUser(ctx, CreateUserParams{
				Phone: newUser.Phone,
			})
		}
		// Some other error occurred
		return User{}, err
	}
	// User exists
	return user, nil
}
