package uaac

import "context"

type UaacRepository interface {
	FindAccount(ctx context.Context, username string) (*Account, error)
	FindProfile(ctx context.Context, username string) (*Profile, error)

	UpdatePassword(ctx context.Context, username string, password string) error
	UpdateProfile(ctx context.Context, profile *Profile) error
}
