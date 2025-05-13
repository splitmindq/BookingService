package pgx

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserCreate   = errors.New("failed to create user")
	ErrInvalidInput = errors.New("invalid input data")
)
