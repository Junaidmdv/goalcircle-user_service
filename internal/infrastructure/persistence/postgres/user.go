package postgres

type userRepository struct {
	Db *postgressDB
}

func NewUserRepository(db *postgressDB) *userRepository {
	return &userRepository{
		Db: db,
	}
}

func (d *userRepository) ExistByEmail(email string) (bool, error) {
	return true, nil
}
 

