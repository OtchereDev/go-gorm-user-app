package serializers

type UserSignUpSerializer struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"email,min=1,max=255"`
	Password string `json:"password" validate:"required,min=1,max=255"`
}

type UserSignUpResponseSerializer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserLoginSerializer struct {
	Email    string `json:"email" validate:"email,min=1,max=255"`
	Password string `json:"password" validate:"required,min=1,max=255"`
}
