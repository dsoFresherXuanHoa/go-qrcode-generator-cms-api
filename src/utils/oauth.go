package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
	"io"
	"net/http"
)

var (
	ErrReadByteFromHttpResponse = errors.New("read all bytes from http response failure")
)

type oauthUtil struct{}

func NewOAuthUtil() *oauthUtil {
	return &oauthUtil{}
}

// TODO: Do not use default gender, set gender value from Google Authentication Service
// FIXME: Do not use default password (agent)
func (oauthUtil) Response2User(res *http.Response) (*entity.UserCreatable, error) {
	if content, err := io.ReadAll(res.Body); err != nil {
		fmt.Println("Error while read all bytes from http response: " + err.Error())
		return nil, ErrReadByteFromHttpResponse
	} else {
		resUser := struct {
			Email      string `json:"email"`
			FamilyName string `json:"family_name"`
			GivenName  string `json:"given_name"`
			Picture    string `json:"picture"`
		}{}
		json.Unmarshal(content, &resUser)
		defaultRoleId := uint(2)
		defaultPassword := "nil"
		usr := entity.UserCreatable{RoleID: &defaultRoleId, FirstName: &resUser.FamilyName, LastName: &resUser.GivenName, Email: &resUser.Email, Password: &defaultPassword, AvatarURL: resUser.Picture}
		return &usr, nil
	}
}
