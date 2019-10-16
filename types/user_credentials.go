package types

type UserCredentials struct {
	Username string `validate:"nonzero" db:"username"`
	Password string `validate:"nonzero" db:"password"`
}
