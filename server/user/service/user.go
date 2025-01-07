package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"zmall/common/bizerr"
	"zmall/common/constant"

	"zmall/server/user/domain"
	"zmall/server/user/repository"
)

var (
	redirectUri   = url.PathEscape("http://zmall.top/api/user/callback")
	defaultAvatar = ""
)

type UserService struct {
	repo   *repository.UserRepo
	client *http.Client
	appId  string
	appKey string
}

type Result struct {
	ErrCode int64  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenId  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionId string `json:"unionid"`
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) AuthUrl(ctx context.Context, state string) (string, error) {
	format := "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s"
	return fmt.Sprintf(format, svc.appId, redirectUri, state), nil
}

func (svc *UserService) VerifyCode(ctx context.Context, code string) (domain.WeChatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, svc.appId, svc.appKey, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, nil)
	if err != nil {
		return domain.WeChatInfo{}, err
	}

	var resp *http.Response
	resp, err = svc.client.Do(req)
	if err != nil {
		return domain.WeChatInfo{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WeChatInfo{}, err
	}

	if res.ErrCode != 0 {
		return domain.WeChatInfo{}, fmt.Errorf("微信返回错误码 %d，错误信息 %s", res.ErrCode, res.ErrMsg)
	}

	return domain.WeChatInfo{
		OpenID:  res.OpenId,
		UnionID: res.UnionId,
	}, nil
}

func (svc *UserService) FindOrCreateUser(ctx context.Context, info domain.WeChatInfo) (domain.User, error) {
	user, err := svc.repo.FindByWechat(ctx, info.OpenID)
	if err == nil {
		return user, nil
	}

	var code string
	code, err = svc.GenerateCode()
	user = domain.User{
		NickName: code[:15] + "-" + code[15:],
		Avatar:   defaultAvatar,
		OpenId:   info.OpenID,
	}

	err = svc.repo.CreateUser(ctx, user)
	if err != nil {
		return user, bizerr.NewBizError(constant.Success)
	}

	return svc.repo.FindByWechat(ctx, info.OpenID)
}

func (svc *UserService) GenerateCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.Grow(20)

	for _, b := range bytes {
		sb.WriteByte(charset[int(b)%len(charset)])
	}

	return sb.String(), nil
}
