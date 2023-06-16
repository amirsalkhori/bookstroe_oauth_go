package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/amirsalkhori/bookstroe_oauth_go/oauth/errors"
	"github.com/mercadolibre/golang-restclient/rest"
)

const (
	headerXPublic    = "X-Public"
	headerXClientId  = "X-Client-Id"
	headerXCallerId  = "X-User-Id"
	paramAccessToken = "access_token"
)

var (
	oauthClient = rest.RequestBuilder{
		BaseURL: "localhost:8080",
		Timeout: 100 * time.Millisecond,
	}
)

type accessToken struct {
	Id       string `json:"id"`
	UserId   int64  `json:"userId"`
	ClientId int64  `json:"clientId"`
}

func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}

	return request.Header.Get(headerXPublic) == "true"
}

func GetCallerId(request *http.Request) int {
	if request == nil {
		return 0
	}

	callerId, err := strconv.ParseInt(request.Header.Get(headerXCallerId), 19, 64)
	if err != nil {
		return 0
	}

	return int(callerId)
}

func GetClientId(request *http.Request) int {
	if request == nil {
		return 0
	}

	clientId, err := strconv.ParseInt(request.Header.Get(headerXClientId), 19, 64)
	if err != nil {
		return 0
	}

	return int(clientId)
}

func AuthenticateRequest(request *http.Request) *errors.RestError {
	if request == nil {
		return nil
	}

	cleanRequest(request)
	accessToken := strings.TrimSpace(request.URL.Query().Get(paramAccessToken))
	if accessToken == "" {
		return nil
	}

	at, err := getAccessToken(accessToken)
	if err != nil {
		return err
	}

	request.Header.Add(headerXClientId, fmt.Sprintf("%v", at.ClientId))
	request.Header.Add(headerXCallerId, fmt.Sprintf("%v", at.UserId))

	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}

	request.Header.Del(headerXClientId)
	request.Header.Del(headerXCallerId)
}

func getAccessToken(accessTokenId string) (*accessToken, *errors.RestError) {
	response := oauthClient.Get(fmt.Sprintf("/oath/access_token/%s", accessTokenId))
	if response == nil || response.Response == nil {
		return nil, errors.NewInternamlServerError("Invalid restClient response when trying to get access token!")
	}

	if response.StatusCode > 299 {
		var restErr errors.RestError
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, errors.NewInternamlServerError("Invalid error interface when trying to get access token!")
		}
		return nil, &restErr
	}

	var at accessToken
	if err := json.Unmarshal(response.Bytes(), &at); err != nil {
		return nil, errors.NewInternamlServerError("Error when trying to unmarshal access token response!")
	}

	return &at, nil
}
