package repository




type UserRepository interface {
	ExistByEmail(string) (bool, error)
}


