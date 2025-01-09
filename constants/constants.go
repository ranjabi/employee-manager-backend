package constants

const (
	SALT_ROUND int = 10
	JWT_SECRET string = "secret"
	HASH_ALG string = "HS256"
	
	UNIQUE_VIOLATION_ERROR_CODE string = "23505"
	FOREIGN_KEY_CONSTRAINT_VIOLATION_ERROR_CODE string = "23503"
)