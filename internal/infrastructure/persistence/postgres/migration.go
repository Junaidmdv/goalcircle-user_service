package postgres

import "github.com/junaidmdv/goalcirlcle/user_service/internal/domain/entity"

func (db *postgressDB) Migration() error {
	if err := db.DB.AutoMigrate(&entity.User{}, &entity.TempUser{}); err != nil {
		return err
	}
	return nil
}
