package postgres

import "github.com/junaidmdv/goalcirlcle/user_service/internal/domain/entity"

func (db *postgressDB) Migration() error {
	err := db.DB.AutoMigrate(&entity.User{}, &entity.TempUser{})
	if err != nil {
		return err
	}
	return nil
}
