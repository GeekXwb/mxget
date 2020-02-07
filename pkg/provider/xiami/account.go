package xiami

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/winterssy/ghttp"
)

// 登录接口，account 可为邮箱/手机号码
func (a *API) LoginRaw(ctx context.Context, account string, password string) (*LoginResponse, error) {
	token, err := a.getToken(APILogin)
	if err != nil {
		return nil, err
	}

	passwordHash := md5.Sum([]byte(password))
	password = hex.EncodeToString(passwordHash[:])
	model := map[string]string{
		"account":  account,
		"password": password,
	}
	params := signPayload(token, model)

	resp := new(LoginResponse)
	err = a.SendRequest(ghttp.MethodGet, APILogin,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	return resp, nil
}
