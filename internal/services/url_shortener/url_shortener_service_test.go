package url_shortener

import (
	"fmt"
	"testing"
	"url_shortener/configs"
	"url_shortener/internal/custom_errors"
	"url_shortener/internal/db"
	"url_shortener/internal/repositories"
	"url_shortener/internal/schemas"
	"url_shortener/internal/services/api_client"
)

type serviceTestData struct {
	name            string
	urlInput        *schemas.URLInput
	domain          string
	getNewLink      func(int) string
	shortenerResult *repositories.FakeShortenerRepository
	apiClientResult []byte
	apiClientErr    error
	resultUrl       string
	resultErr       error
}

func fakeNewService(testData serviceTestData) *URLShortenerService {
	return &URLShortenerService{
		config:              &configs.Config{DOMAIN: testData.domain},
		shortenerRepository: testData.shortenerResult,
		getNewLink:          testData.getNewLink,
		apiClient:           api_client.FakeAPIClient{Result: testData.apiClientResult, Err: testData.apiClientErr},
	}
}

func TestShortURL(t *testing.T) {
	tests := []serviceTestData{
		{
			name:            "Without url",
			urlInput:        &schemas.URLInput{Url: ""},
			domain:          "http://localhost",
			getNewLink:      func(int) string { return "" },
			shortenerResult: &repositories.FakeShortenerRepository{},
			resultUrl:       "",
			resultErr:       fmt.Errorf(custom_errors.UrlRequiredError),
		},
		{
			name:            "With url",
			urlInput:        &schemas.URLInput{Url: "http://localhost:8000/api/check"},
			domain:          "http://localhost",
			getNewLink:      func(int) string { return "1234567" },
			shortenerResult: &repositories.FakeShortenerRepository{InsertResult: &db.Shortener{NewLink: "1234567"}},
			resultUrl:       "http://localhost/api/1234567",
			resultErr:       nil,
		},
		{
			name:            "With url db error",
			urlInput:        &schemas.URLInput{Url: "http://localhost:8000/api/check"},
			domain:          "http://localhost",
			getNewLink:      func(int) string { return "1234567" },
			shortenerResult: &repositories.FakeShortenerRepository{InsertErr: fmt.Errorf("Some error")},
			resultUrl:       "",
			resultErr:       fmt.Errorf("Some error"),
		},
		{
			name:         "With url request error",
			urlInput:     &schemas.URLInput{Url: "http://localhost:8000/api/check"},
			domain:       "http://localhost",
			getNewLink:   func(int) string { return "1234567" },
			apiClientErr: fmt.Errorf("Error"),
			resultUrl:    "",
			resultErr:    fmt.Errorf(custom_errors.MakingRequestError),
		},
	}
	for _, test := range tests {
		service := fakeNewService(test)
		actual, err := service.ShortURL(test.urlInput)
		// log.Println("Started test - ", test.name)
		if actual != test.resultUrl {
			t.Error(fmt.Errorf("case: %s, actual result: %s, expected result: %s", test.name, actual, test.resultUrl))
		}
		if err != nil && test.resultErr != nil && err.Error() != test.resultErr.Error() || err != nil && test.resultErr == nil || err == nil && test.resultErr != nil {
			t.Error(fmt.Errorf("case: %s, actual err: %#v, expected err: %#v", test.name, err, test.resultErr))
		}
	}
}
