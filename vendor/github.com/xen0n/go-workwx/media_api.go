package workwx

import (
	"strconv"
	"time"
)

// XXX: 由于 sdkcodegen 目前不支持生成 `import` 语句，这个模型不能用 sdkcodegen 生成
// import "time"

// MediaUploadResult 临时素材上传结果
type MediaUploadResult struct {
	// Type 媒体文件类型，分别有图片（image）、语音（voice）、视频（video），普通文件(file)
	Type string
	// MediaID 媒体文件上传后获取的唯一标识，3天内有效
	MediaID string
	// CreatedAt 媒体文件上传时间戳
	CreatedAt time.Time
}

func (x respMediaUpload) intoMediaUploadResult() (MediaUploadResult, error) {
	createdAtInt, err := strconv.ParseInt(x.CreatedAt, 10, 64)
	if err != nil {
		return MediaUploadResult{}, err
	}
	createdAt := time.Unix(createdAtInt, 0)

	return MediaUploadResult{
		Type:      x.Type,
		MediaID:   x.MediaID,
		CreatedAt: createdAt,
	}, nil
}

//
// API 接口
//

// mediaUpload 上传临时素材
//
// NOTE: 因为名字很难听，所以不直接暴露给用户使用
func (c *WorkwxApp) mediaUpload(typ string, media *Media) (*MediaUploadResult, error) {
	resp, err := c.execMediaUpload(reqMediaUpload{
		Type:  typ,
		Media: media,
	})
	if err != nil {
		return nil, err
	}

	obj, err := resp.intoMediaUploadResult()
	if err != nil {
		return nil, err
	}

	// TODO: return bare T instead of &T?
	return &obj, nil
}

// mediaUploadImg 上传永久图片
//
// NOTE: 因为名字很难听，所以不直接暴露给用户使用
func (c *WorkwxApp) mediaUploadImg(media *Media) (url string, err error) {
	resp, err := c.execMediaUploadImg(reqMediaUploadImg{
		Media: media,
	})
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}

//
// convenient wrappers
//

const (
	tempMediaTypeImage = "image"
	tempMediaTypeVoice = "voice"
	tempMediaTypeVideo = "video"
	tempMediaTypeFile  = "file"
)

// UploadTempImageMedia 上传临时图片素材
func (c *WorkwxApp) UploadTempImageMedia(media *Media) (*MediaUploadResult, error) {
	result, err := c.mediaUpload(tempMediaTypeImage, media)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UploadTempVoiceMedia 上传临时语音素材
func (c *WorkwxApp) UploadTempVoiceMedia(media *Media) (*MediaUploadResult, error) {
	result, err := c.mediaUpload(tempMediaTypeVoice, media)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UploadTempVideoMedia 上传临时视频素材
func (c *WorkwxApp) UploadTempVideoMedia(media *Media) (*MediaUploadResult, error) {
	result, err := c.mediaUpload(tempMediaTypeVideo, media)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UploadTempFileMedia 上传临时文件素材
func (c *WorkwxApp) UploadTempFileMedia(media *Media) (*MediaUploadResult, error) {
	result, err := c.mediaUpload(tempMediaTypeFile, media)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UploadPermanentImageMedia 上传永久图片素材
func (c *WorkwxApp) UploadPermanentImageMedia(media *Media) (url string, err error) {
	url, err = c.mediaUploadImg(media)
	if err != nil {
		return "", err
	}

	return url, nil
}
