package constants

var (
	InvalidSignUpRequestFormat    = "Invalid sign up request format: check your request and try again later!"
	InvalidUserQueryRequestFormat = "Invalid user query request format: check your request and try again later!"
	InvalidRoleRequestFormat      = "Invalid role request format: check your request and try again later!"
	InvalidAccessToken            = "Invalid accessToken: your accessToken can be damage or interrupt, try to sign out and sign in again to refresh accessToken!"

	MissingActivationCodeInQueryString = "Missing activation code in query string: activation code must be numeric and cannot be empty!"
	MissingBearerToken                 = "Missing bearer token: try to sign out and sign in again to refresh token!"

	ErrCreateNewRole                = "Error while create new role: make sure you has permission to do this and try again later!"
	ErrEncodeFileMultiPartHeader    = "Error while encode multipart file header from user request: multipart file so large or something went wrong!"
	ErrUploadSingleFileToCloudinary = "Error while upload single image file to cloudinary: invalid file type or url format!"
	ErrSignUpForNewUser             = "Error while sign up for new user: try again later or contact admin to get more information!"
	ErrSendActivationRequestEmail   = "Error while send activation email to user: make sure you send correct email address!"
	ErrActivateUser                 = "Error while activate user: make sure you click to right activation url"
	ErrSignIn                       = "Error while sign in: check your email and password and sign in again!"
	ErrFindDetailUser               = "Error while find detail user: make sure you have right permission and try again later!"

	CreateNewRoleSuccess    = "Create new role success: congrats!!!"
	SignUpForNewUserSuccess = "Sign up for new user success: congrats!!!"
	ActivateUserSuccess     = "Activate user success: congrats!!!"
	SignInSuccess           = "Sign in success: congrats!!!"
	FindDetailUserSuccess   = "Find detail user success: congrats!!!"
)
