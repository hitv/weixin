package weixin

import (
	"encoding/xml"
	"io"
	"time"
)

type ReplyMsgType string

const (
	ReplyMsgTypeText              = "text"
	ReplyMsgTypeNews ReplyMsgType = "news"
)

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

type Reply struct {
	toUserName   string
	fromUserName string
	msg          interface{}
}

func (r *Reply) replyMsgBase(msgType MsgType) *MsgBase {
	return &MsgBase{
		ToUserName:   r.fromUserName,
		FromUserName: r.toUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      msgType,
	}
}

func (r *Reply) Write(w io.Writer) error {
	return xml.NewEncoder(w).Encode(r.msg)
}

func (r *Reply) TextMsg(content string) *Reply {
	r.msg = &RespTextMsg{
		MsgBase: r.replyMsgBase(MsgTypeText),
		Content: content,
	}
	return r
}

func (r *Reply) ImageMsg(mediaId string) *Reply {
	r.msg = &RespImageMsg{
		MsgBase: r.replyMsgBase(MsgTypeImage),
		MediaId: mediaId,
	}
	return r
}

func (r *Reply) VoiceMsg(mediaId string) *Reply {
	r.msg = &RespVoiceMsg{
		MsgBase: r.replyMsgBase(MsgTypeVoice),
		MediaId: mediaId,
	}
	return r
}

func (r *Reply) VideoMsg(mediaId, title, description string) *Reply {
	r.msg = &RespVideoMsg{
		MsgBase:     r.replyMsgBase(MsgTypeVideo),
		MediaId:     mediaId,
		Title:       title,
		Description: description,
	}
	return r
}

func (r *Reply) MusicMsg(description, musicURL, hqMusicURL, thumbMediaId string) *Reply {
	r.msg = &RespMusicMsg{
		MsgBase:      r.replyMsgBase(MsgTypeMusic),
		Description:  description,
		MusicUrl:     musicURL,
		HQMusicUrl:   hqMusicURL,
		ThumbMediaId: thumbMediaId,
	}
	return r
}

func (r *Reply) NewsMsg(items []*RespArticleItem) *Reply {
	r.msg = &RespNewsMsg{
		MsgBase:      r.replyMsgBase(MsgTypeNews),
		ArticleCount: len(items),
		ArticleItems: items,
	}
	return r
}
