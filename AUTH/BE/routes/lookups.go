package routes

import (
	app_repo "auth_api/domain/app/repo"
	auth_repo "auth_api/domain/auth/repo"
)

// appKeyLookup menghubungkan repo aplikasi ke middleware.AppLookup.
type appKeyLookup struct {
	repo app_repo.AppRepoInterface
}

func (l appKeyLookup) GetActiveByKeyHash(keyHash string) (int, string, bool, error) {
	app, err := l.repo.GetActiveByKeyHash(keyHash)
	if err != nil || app == nil {
		return 0, "", false, err
	}
	return app.ID, app.Code, true, nil
}

// actingUserLookup menghubungkan repo auth ke middleware.UserLookup.
type actingUserLookup struct {
	repo auth_repo.AuthRepoInterface
}

func (l actingUserLookup) GetActingUser(id int) (string, bool, bool, bool, error) {
	u, err := l.repo.GetUserByID(id)
	if err != nil || u == nil {
		return "", false, false, false, err
	}
	return u.RoleName, u.IsActive, u.MustChangePassword, true, nil
}
