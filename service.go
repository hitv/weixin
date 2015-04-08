package weixin

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

const (
	weixinAPI = "https://api.weixin.qq.com/cgi-bin"
)

type WeixinService interface {
	IsRequestValid(timestamp, nonce, signature string) bool
	ParseRequest(r io.Reader) (*Request, error)
	CreateMenu(menu *Menu) error
	GetMenu() (*Menu, error)
	DeleteMenu() error
}

type wxService struct {
	appId       string
	secret      string
	token       string
	accessToken *Token
	aesKey      []byte
	client      *http.Client
}

func NewWxService(appId, secret, token string, aesKey []byte) WeixinService {
	return &wxService{
		appId:       appId,
		secret:      secret,
		token:       token,
		accessToken: NewToken(appId, secret),
		aesKey:      aesKey,
		client:      &http.Client{},
	}
}

func (s *wxService) makeMsgSignature(timestamp, nonce, encrypt string) string {
	var (
		hash   = sha1.New()
		params = []string{s.token, timestamp, nonce, encrypt}
	)
	sort.Strings(params)
	io.WriteString(hash, strings.Join(params, ""))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (s *wxService) IsRequestValid(timestamp, nonce, signature string) bool {
	str := s.makeMsgSignature(timestamp, nonce, "")
	return str == signature
}

func (s *wxService) validAppId(plainData []byte) (data []byte, valid bool) {
	r := bytes.NewReader(plainData[16:20])

	var length int
	binary.Read(r, binary.BigEndian, &length)

	start := length + 20
	id := plainData[start : start+len(s.appId)]
	if string(id) != s.appId {
		valid = false
		return
	}
	data = plainData[20 : 20+length]
	return
}

func (s *wxService) ParseRequest(r io.Reader) (req *Request, err error) {
	if s.aesKey != nil {
		var (
			encReq    = &EncryptRequest{}
			plainData []byte
		)

		err = xml.NewDecoder(r).Decode(encReq)
		if err != nil {
			return
		}

		plainData, err = encReq.DecryptAES(s.aesKey)
		if err != nil {
			return
		}

		plainData, idValid := s.validAppId(plainData)
		if !idValid {
			err = errors.New("Appid is invalid")
			return
		}

		r = bytes.NewReader(plainData)
	}
	req = &Request{}
	decoder := xml.NewDecoder(r)
	err = decoder.Decode(req)
	return
}

func (s *wxService) CreateMenu(menu *Menu) (err error) {
	err = s.accessToken.Ensure()
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(menu)
	if err != nil {
		return
	}

	apiURL := weixinAPI + "/menu/create?access_token=" + s.accessToken.AccessToken
	resp, err := s.client.Post(apiURL, "application/x-www-form-urlencoded", buf)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	menuResp := &MenuResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(menuResp)
	if err != nil {
		return
	}

	err = menuResp.Error()
	return
}

func (s *wxService) GetMenu() (menu *Menu, err error) {
	err = s.accessToken.Ensure()
	if err != nil {
		return
	}

	apiURL := weixinAPI + "/menu/get?access_token=" + s.accessToken.AccessToken
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	menuResp := &MenuResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(menuResp)
	if err != nil {
		return
	}

	err = menuResp.Error()
	if err == nil {
		menu = menuResp.Menu
	}

	return
}

func (s *wxService) DeleteMenu() (err error) {
	err = s.accessToken.Ensure()
	if err != nil {
		return
	}

	apiURL := weixinAPI + "/menu/delete?access_token=" + s.accessToken.AccessToken
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	menuResp := &MenuResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(menuResp)
	if err != nil {
		return
	}

	err = menuResp.Error()
	return
}
