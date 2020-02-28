package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	rest_errors "github.com/ravayak/utils-go/rest_errors"

	"github.com/mercadolibre/golang-restclient/rest"
)

const (
	// Header determining private or public informations
	// A public request is defined as a request made from the
	// "outside" world compared to a private request made
	// from our internal network
	headerXPublic        = "X-Public"
	headerXClientId      = "X-Client-Id"
	headerXCallerId      = "X-Caller-Id"
	parameterAccessToken = "access_token"
)

var (
	oauthRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8082",
		Timeout: 200 * time.Millisecond,
	}
)

type oauthClient struct{}
type oauthInterface interface{}

type accessToken struct {
	UserID int64 `json:"user_id"`
	// what kind of clients requests the access token ?
	// A web frontend, an Android APP... The id shouldn't
	// be the same. A client might have a longer expiration
	// time for a token than another one for ex. We want
	// to limit access to some api for a certain type of client...
	ClientID int64 `json:"client_id"`
	Expires  int64 `json:"expires"`
}

func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}
	return request.Header.Get(headerXPublic) == "true"
}

func GetCallerID(req *http.Request) int64 {
	if req == nil {
		return 0
	}

	callerID, err := strconv.ParseInt(req.Header.Get(headerXCallerId), 10, 64)
	if err != nil {
		return 0
	}

	return callerID
}

func GetClientID(req *http.Request) int64 {
	if req == nil {
		return 0
	}

	callerID, err := strconv.ParseInt(req.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}

	return callerID
}

func AuthenticateRequest(req *http.Request) rest_errors.RestError {
	if req == nil {
		return nil
	}
	cleanRequest(req)

	// http://api/resource?access_token=abc123
	accessTokenID := strings.TrimSpace(req.URL.Query().Get(parameterAccessToken))
	if accessTokenID == "" {
		return nil
	}

	at, err := getAccessToken(accessTokenID)
	if err != nil {
		if err.Status() == http.StatusNotFound {
			return nil
		}
		return err
	}

	req.Header.Add(headerXCallerId, fmt.Sprintf("%v", at.UserID))
	req.Header.Add(headerXClientId, fmt.Sprintf("%v", at.ClientID))

	return nil
}

func cleanRequest(req *http.Request) {
	if req == nil {
		return
	}

	req.Header.Del(headerXClientId)
	req.Header.Del(headerXCallerId)
}

func getAccessToken(accessTokenID string) (*accessToken, rest_errors.RestError) {
	res := oauthRestClient.Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenID))

	if res == nil || res.Response == nil {
		// Timeout
		return nil, rest_errors.NewInternalServerError("invalid rest client response while trying to get access token",
			errors.New("Invalid restclient response"))
	}

	if res.StatusCode > 299 {
		// Error
		var restErr rest_errors.RestError
		err := json.Unmarshal(res.Bytes(), &restErr)
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface while trying to unmarshal access token", err)
		}

		return nil, restErr
	}

	var at accessToken
	if err := json.Unmarshal(res.Bytes(), &at); err != nil {
		return nil, rest_errors.NewInternalServerError("erro while trying to unmarshal access token", err)
	}

	return &at, nil

}
