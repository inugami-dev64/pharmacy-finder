package service_test

import (
	"io"
	"net/http"
	"os"
	"pharmafinder/mock"
	"pharmafinder/service"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRecaptchaVerification_SuccessTrue(t *testing.T) {
	os.Setenv("ALLOWED_DOMAINS", "hrt.girlkisser.gay")

	ctrl := gomock.NewController(t)
	httpMock := mock.NewMockHttpClient(ctrl)
	httpMock.EXPECT().
		Do(gomock.Any()).
		DoAndReturn(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"success":true,"hostname":"hrt.girlkisser.gay"}`)),
			}, nil
		})

	verifier := service.ProvideRecaptchaVerifier(httpMock)
	assert.True(t, verifier.Verify("xD"))

	os.Clearenv()
}

func TestRecaptchaVerification_SuccessFalse(t *testing.T) {
	os.Setenv("ALLOWED_DOMAINS", "hrt.girlkisser.gay")

	ctrl := gomock.NewController(t)
	httpMock := mock.NewMockHttpClient(ctrl)
	httpMock.EXPECT().
		Do(gomock.Any()).
		DoAndReturn(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"success":false,"hostname":"hrt.girlkisser.gay"}`)),
			}, nil
		})

	verifier := service.ProvideRecaptchaVerifier(httpMock)
	assert.False(t, verifier.Verify("xD"))

	os.Clearenv()
}

func TestRecaptchaVerification_Non200StatusCode(t *testing.T) {
	os.Setenv("ALLOWED_DOMAINS", "hrt.girlkisser.gay")

	ctrl := gomock.NewController(t)
	httpMock := mock.NewMockHttpClient(ctrl)
	httpMock.EXPECT().
		Do(gomock.Any()).
		DoAndReturn(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader(`{"success":false,"hostname":"hrt.girlkisser.gay"}`)),
			}, nil
		})

	verifier := service.ProvideRecaptchaVerifier(httpMock)
	assert.False(t, verifier.Verify("xD"))

	os.Clearenv()
}
