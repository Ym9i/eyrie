package spider

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/z-quan-tong/eyrie/pkg/spider/aliyundrive"
)

type AliYunConfig struct {
	RefreshToken string
}

var AliyunConfig = &AliYunConfig{}
var drive *aliyundrive.AliYunDrive

func ResolveAliyunDrive(client *Client, shareUrl string, sharePwd string, formats ...string) ([]string, error) {
	return resolveAliyunDrive(client, shareUrl, sharePwd, formats)
}

func resolveAliyunDrive(client *Client, shareUrl string, sharePwd string, formats []string) ([]string, error) {
	shareId := strings.TrimPrefix(shareUrl, "https://www.aliyundrive.com/s/")
	sharePwd = strings.TrimSpace(sharePwd)

	if drive == nil {
		drive = NewAliyunDrive(client, AliyunConfig)
	}

	return resolveShare(drive, shareId, sharePwd, formats)
}

func resolveShare(drive *aliyundrive.AliYunDrive, shareId string, sharePwd string, formats []string) ([]string, error) {
	token, err := drive.GetShareToken(shareId, sharePwd)
	log.Println(token, err)
	if err != nil {
		return nil, err
	}

	shareFiles, err := drive.GetShare(shareId, token.ShareToken)
	if err != nil {
		return nil, err
	}

	var links []string

	for item := range shareFiles {

		for _, format := range formats {
			if strings.EqualFold(item.FileExtension, format) {
				url, err := drive.GetFileDownloadUrl(token.ShareToken, shareId, item.FileId)
				if err != nil {
					return nil, err
				}
				log.Println("url", url)
				links = append(links, url)
			}
		}
	}

	return links, nil
}

func NewAliyunDrive(c *Client, aliConfig *AliYunConfig) *aliyundrive.AliYunDrive {
	client := resty.NewWithClient(c.client)
	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal

	client.SetTimeout(c.config.Timeout)
	client.SetRetryCount(c.config.Retry)
	client.SetDisableWarn(true)
	client.SetDebug(c.config.Debug)
	client.SetPreRequestHook(aliyundrive.HcHook)
	client.SetHeader(aliyundrive.UserAgent, c.config.UserAgent)
	client.SetHeader(aliyundrive.ContentType, aliyundrive.ContentTypeJSON)
	return &aliyundrive.AliYunDrive{
		Client:       client,
		RefreshToken: aliConfig.RefreshToken,
		Cache:        map[string]string{},
	}
}
