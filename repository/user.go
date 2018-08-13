package repository

import (
	"github.com/tetsuyanh/c2c-demo/model"
)

const (
	UserTableName           = "users"
	SessionTableName        = "sessions"
	AuthenticationTableName = "authentications"
)

var (
	userRepo UserRepo
)

type (
	UserRepo interface {
		CreateUser(u *model.User) error
		CreateSession(s *model.Session) error
		CreateAuth(a *model.Authentication) error
		FindUserBySettionToken(token string) (*model.User, error)
		FindAuthByEmail(email string) (*model.Authentication, error)
		UpdateAuthEnable(token string) error
	}

	userRepoImpl struct{}
)

func GetUserRepo() UserRepo {
	if userRepo == nil {
		userRepo = &userRepoImpl{}
	}
	return userRepo
}

func (r *userRepoImpl) CreateUser(u *model.User) error {
	if err := repo.dbMap.Insert(u); err != nil {
		return err
	}
	return nil
}

func (r *userRepoImpl) CreateSession(s *model.Session) error {
	if err := repo.dbMap.Insert(s); err != nil {
		return err
	}
	return nil
}

func (r *userRepoImpl) CreateAuth(a *model.Authentication) error {
	if err := repo.dbMap.Insert(a); err != nil {
		return err
	}
	return nil
}

func (ur *userRepoImpl) FindUserBySettionToken(token string) (*model.User, error) {
	u := &model.User{}
	if err := repo.dbMap.SelectOne(u, `
		select
			u.id,
			u.created_at
		from
			users u
			left join sessions s on u.id = s.user_id
		where
			s.token = $1
		`, token); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *userRepoImpl) FindAuthByEmail(email string) (*model.Authentication, error) {
	u := &model.Authentication{}
	if err := repo.dbMap.SelectOne(u, `
		select
			*
		from
			authentications
		where
			email = $1
		`, email); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *userRepoImpl) UpdateAuthEnable(token string) error {
	return resultExec(repo.dbMap.Exec(`
		update
			authentications
		set
			enabled = true,
			updated_at = now()
		where
			token = $1 and
			enabled = false
			`, token))
}
