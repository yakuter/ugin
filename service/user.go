package service

// FindByCredentials method
func FindByCredentials(email string, password string) bool {
	return email == "sample@domain.com" && password == "sample" //check your db
}
