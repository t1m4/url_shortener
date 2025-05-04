package custom_errors

// DB errors
const (
	DbError              = "error while creating"
	RowDoesNotExistError = "row does not exist"
)

// Service errors
const (
	UrlRequiredError   = "url is required parameter"
	MakingRequestError = "can't make request, invalid url"
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
	RateLimitError     = "rate limit error"
	InvalidUserIdError = "invalid user id"
)
