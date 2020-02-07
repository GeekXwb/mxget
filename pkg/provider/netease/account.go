package netease

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/winterssy/ghttp"
)

// 邮箱登录
func (a *API) EmailLoginRaw(ctx context.Context, email string, password string) (*LoginResponse, error) {
	passwordHash := md5.Sum([]byte(password))
	password = hex.EncodeToString(passwordHash[:])
	data := map[string]interface{}{
		"username":      email,
		"password":      password,
		"rememberLogin": true,
	}

	resp := new(LoginResponse)
	err := a.SendRequest(ghttp.MethodPost, APIEmailLogin,
		ghttp.WithForm(weapi(data)),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("email login: %s", resp.errorMessage())
	}

	return resp, nil
}

// 手机登录
func (a *API) CellphoneLoginRaw(ctx context.Context, countryCode int, phone int, password string) (*LoginResponse, error) {
	passwordHash := md5.Sum([]byte(password))
	password = hex.EncodeToString(passwordHash[:])
	data := map[string]interface{}{
		"phone":         phone,
		"countrycode":   countryCode,
		"password":      password,
		"rememberLogin": true,
	}

	resp := new(LoginResponse)
	err := a.SendRequest(ghttp.MethodPost, APICellphoneLogin,
		ghttp.WithForm(weapi(data)),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("cellphone login: %s", resp.errorMessage())
	}

	return resp, nil
}

// 刷新登录状态
func (a *API) RefreshLoginRaw(ctx context.Context) (*CommonResponse, error) {
	resp := new(CommonResponse)
	err := a.SendRequest(ghttp.MethodPost, APIRefreshLogin,
		ghttp.WithForm(weapi(struct{}{})),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("refresh login: %s", resp.errorMessage())
	}

	return resp, nil
}

// 退出登录
func (a *API) LogoutRaw(ctx context.Context) (*CommonResponse, error) {
	resp := new(CommonResponse)
	err := a.SendRequest(ghttp.MethodPost, APILogout,
		ghttp.WithForm(weapi(struct{}{})),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("logout: %s", resp.errorMessage())
	}

	return resp, nil
}
