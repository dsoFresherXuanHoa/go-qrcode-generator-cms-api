package constants

var (
	InvalidSignUpRequestFormat = "Invalid sign up request format: check your request and try again later!"
	InvalidRoleRequestFormat   = "Invalid role request format: check your request and try again later!"

	ErrEncodeFileMultiPartHeader    = "Error while encode multipart file header from user request: multipart file so large or something went wrong!"
	ErrUploadSingleFileToCloudinary = "Error while upload single image file to cloudinary: invalid file type or url format!"
	ErrSignUpForNewUser             = "Error while sign up for new user: try again later or contact admin to get more information!"
	ErrSendActivationRequestEmail   = "Error while send activation email to user: make sure you send correct email address!"
	ErrCreateNewRole                = "Error while create new role: make sure you has permission to do this and try again later!"

	SignUpForNewUserSuccess = "Sign up for new user success: congrats!!!"
	CreateNewRoleSuccess    = "Create new role success: congrats!!!"
)
