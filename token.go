package weixin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Token struct {
	appId       string    `json:"-"`
	secret      string    `json:"-"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	ExpireTime  time.Time `json:"-"`
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpireTime)
}

func (t *Token) Ensure() error {
	if t.IsExpired() {
		return t.Refresh()
	}
	return nil
}

func (t *Token) Refresh() error {
	tokenURL := fmt.Sprintf("%s/token?grant_type=client_credential&appid=%s&secret=%s", weixinAPI, t.appId, t.secret)
	resp, err := http.Get(tokenURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(t)
	if err != nil {
		return err
	}

	t.ExpireTime = time.Now().Add(time.Second * time.Duration(t.ExpiresIn-600))
	return nil
}

func NewToken(appId, secret string) *Token {
	return &Token{
		appId:  appId,
		secret: secret,
	}
}
