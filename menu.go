package weixin

import (
	"errors"
)

type ButtonType string

const (
	ButtonTypeClick           ButtonType = "click"
	ButtonTypeView                       = "view"
	ButtonTypeScanCodeWaitMsg            = "scancode_waitmsg"
	ButtonTypeScanCodePush               = "scancode_push"
	ButtonTypePicSysPhoto                = "pic_sysphoto"
	ButtonTypePicPhotoOrAlbum            = "pic_photo_or_album"
	ButtonTypePicWeixin                  = "pic_weixin"
	ButtonTypeLocationSelect             = "location_select"
)

type MenuButton struct {
	Type      ButtonType    `json:"type"`
	Name      string        `json:"name"`
	URL       string        `json:"url"`
	Key       string        `json:"key"`
	SubButton []*MenuButton `json:"sub_button"`
}

type MenuResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Menu    *Menu  `json:"menu"`
}

func (r *MenuResponse) Error() error {
	if r.ErrCode != 0 {
		return errors.New(r.ErrMsg)
	}
	return nil
}

type Menu struct {
	Button []*MenuButton `json:"button"`
}
