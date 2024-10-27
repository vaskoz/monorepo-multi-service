package app

// GetAllUsers returns a list of all users.
func GetAllUsers() []User {
	return []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
}
