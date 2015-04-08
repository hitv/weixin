package weixin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
)

type MsgType string

const (
	MsgTypeText     MsgType = "text"
	MsgTypeImage            = "image"
	MsgTypeVoice            = "voice"
	MsgTypeVideo            = "video"
	MsgTypeLocation         = "location"
	MsgTypeLink             = "link"
	MsgTypeUnknow           = "unknow"
	MsgTypeEvent            = "event"
	MsgTypeMusic            = "music"
	MsgTypeNews             = "news"
)

type EventType string

const (
	EventTypeSubscribe   EventType = "subscribe"
	EventTypeUnsubscribe           = "unsubscribe"
	EventTypeScan                  = "SCAN"
	EventTypeLocation              = "LOCATION"
	EventTypeClick                 = "CLICK"
	EventTypeView                  = "VIEW"
)

type MsgBase struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      MsgType
}

type RespTextMsg struct {
	*MsgBase
	Content string
}

type RespImageMsg struct {
	*MsgBase
	MediaId string `xml:"Image>MediaId"`
}

type RespVoiceMsg struct {
	*MsgBase
	MediaId string `xml:"Voice>MediaId"`
}

type RespVideoMsg struct {
	*MsgBase
	MediaId     string `xml:"Video>MediaId"`
	Title       string `xml:"Video>Title"`
	Description string `xml:"Video>Description"`
}

type RespArticleItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}

type RespNewsMsg struct {
	*MsgBase
	ArticleCount int                `xml:"ArticleCount"`
	ArticleItems []*RespArticleItem `xml:"Articles>item"`
}

type RespMusicMsg struct {
	*MsgBase
	Title        string `xml:"Music>Title"`
	Description  string `xml:"Music>Description"`
	MusicUrl     string `xml:"Music>MusicUrl"`
	HQMusicUrl   string `xml:"Music>HQMusicUrl"`
	ThumbMediaId string `xml:"Music>ThumbMediaId"`
}

type ReqEventMsg struct {
	Event    EventType
	EventKey string
	Ticket   string
}

type ReqTextMsg struct {
	Content string
}

type ReqImageMsg struct {
	MediaId string
	PicUrl  string
}

type ReqVoiceMsg struct {
	MediaId string
	Format  string
}

type ReqVideoMsg struct {
	MediaId      string
	ThumbMediaId string
}

type ReqLocationMsg struct {
	LocationX float32 `xml:"Location_X"`
	LocationY float32 `xml:"Location_Y"`
	Scale     int
	Label     string
}

type ReqLinkMsg struct {
	Title       string
	Description string
	Url         string
}

type ReqLocationEventMsg struct {
	Latitude  float64
	Longitude float64
	Precision float64
}

type Request struct {
	MsgBase
	MsgId int64
	ReqEventMsg
	ReqTextMsg
	ReqImageMsg
}

type EncryptRequest struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	Encrypt    []byte
}

func (r *EncryptRequest) DecryptAES(aesKey []byte) (plainData []byte, err error) {
	var binData []byte
	_, err = base64.StdEncoding.Decode(binData, r.Encrypt)
	if err != nil {
		return
	}

	k := len(aesKey)
	if len(binData)%k != 0 {
		err = fmt.Errorf("crypto/cipher: cipher data size is not multiple of aes key length")
		return
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return
	}

	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainData = make([]byte, len(binData))
	blockMode.CryptBlocks(plainData, binData)
	return
}

func (req *Request) Reply() *Reply {
	return &Reply{
		toUserName:   req.FromUserName,
		fromUserName: req.ToUserName,
	}
}
