package custom_errors

// DB errors
const (
	DbError              = "Error while creating"
	RowDoesNotExistError = "Row does not exist"
)

// Service errors
const (
	UrlRequiredError   = "Url is required parameter"
	MakingRequestError = "Can't make request, invalid url"
)

// API client errors
const (
	CreateAPIRequestError    = "error creating request %s: %s"
	MakingAPIRequestError    = "error making request %s: %s"
	WrongResponseStatusError = "error status code %d for request %s"
	ReadBodyError            = "error reading body for request %s: %s"
)

// Rate limit errors
const (
	RateLimitError     = "Rate limit error"
	InvalidUserIdError = "Invalid user id"
)
