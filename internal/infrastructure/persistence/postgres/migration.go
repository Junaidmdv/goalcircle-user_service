package postgres

import "github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"

func (db *postgressDB) Migration() error {
	if err := db.DB.AutoMigrate(&entity.User{}, &entity.TempUser{}, &entity.Otp{}); err != nil {
		return err
	}
	return nil
}
