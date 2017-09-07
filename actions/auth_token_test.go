package actions_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client/clientfakes"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token", func() {
	var (
		subject    actions.ServerInfo
		httpClient clientfakes.FakeHttpClient
		testConfig config.Config
	)

	BeforeEach(func() {
		testConfig = config.Config{AuthURL: "example.com"}
		subject = actions.NewAuthToken(&httpClient, testConfig)
	})

	Describe("GetAuthTokenByClientCredential", func() {
		It("returns the token from the authorization server using client credential", func() {
			responseObj := http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"access_token":"3YotnFZFEjr1zCsicMWpAA",
					"token_type":"bearer",
					"expires_in":3600}`)),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				return &responseObj, nil
			}

			token, _ := subject.GetAuthTokenByClientCredential("test_client", "test_secret")
			Expect(token.AccessToken).To(Equal("3YotnFZFEjr1zCsicMWpAA"))
		})
	})

	Describe("VerifyAuthServerConnection", func() {
		It("succeeds if the request succeeds", func() {
			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				return nil, nil
			}

			error := actions.VerifyAuthServerConnection(&httpClient, testConfig)
			Expect(error).To(BeNil())
		})

		It("returns an error if no valid caCert exists for the auth server", func() {
			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				return nil, errors.NewCatchAllError()
			}

			error := actions.VerifyAuthServerConnection(&httpClient, testConfig)
			Expect(error).ToNot(BeNil())
		})
	})
})
