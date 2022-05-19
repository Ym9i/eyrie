package aliyundrive

import (
	"log"
)

func (ali AliYunDrive) GetAnonymousShare(shareId string) (*GetShareInfoResponse, error) {
	resp, err := ali.Client.R().
		SetBody(GetShareInfoRequest{ShareId: shareId}).
		SetResult(GetShareInfoResponse{}).
		SetError(ErrorResponse{}).
		Post(V3ShareLinkGetShareByAnonymous)

	if err != nil {
		return nil, err
	}

	r := resp.Result().(*GetShareInfoResponse)
	return r, nil
}

func (ali AliYunDrive) GetFileDownloadUrl(shareToken string, shareId string, fileId string) (string, error) {
	resp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetHeader(xShareToken, shareToken).
		SetBody(GetShareLinkDownloadUrlRequest{
			ShareId:   shareId,
			FileId:    fileId,
			ExpireSec: 600,
		}).
		SetResult(GetShareLinkDownloadUrlResponse{}).
		SetError(ErrorResponse{}).
		Post(V2FileGetShareLinkDownloadUrl)

	if err != nil {
		return "", err
	}

	log.Println(resp, shareId)

	i := resp.Result().(*GetShareLinkDownloadUrlResponse)
	return i.DownloadUrl, nil

}

func (ali AliYunDrive) GetShareToken(shareId string, sharePwd string) (*GetShareTokenResponse, error) {
	resp, err := ali.Client.R().
		SetBody(GetShareTokenRequest{ShareId: shareId, SharePwd: sharePwd}).
		SetResult(GetShareTokenResponse{}).
		SetError(ErrorResponse{}).
		Post(V2ShareLinkGetShareToken)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*GetShareTokenResponse), err
}

func (ali AliYunDrive) GetShare(shareId string, shareToken string) (data chan *BaseShareFile, err error) {
	result := make(chan *BaseShareFile, 100)

	go func() {
		err = ali.fileList(shareToken, shareId, result)
		if err != nil {
			log.Fatal(err)
		}
		close(result)
	}()
	return result, nil
}

func (ali AliYunDrive) fileList(shareToken string, shareId string, result chan *BaseShareFile) error {
	return ali.fileListByMarker(FileListParam{
		shareToken:   shareToken,
		shareId:      shareId,
		parentFileId: "root",
		marker:       "",
	}, result)
}

func (ali AliYunDrive) fileListByMarker(param FileListParam, result chan *BaseShareFile) error {
	resp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetHeader(xShareToken, param.shareToken).
		SetBody(GetShareFileListRequest{
			ShareId:        param.shareId,
			ParentFileId:   param.parentFileId,
			UrlExpireSec:   14400,
			OrderBy:        "name",
			OrderDirection: "DESC",
			Limit:          20,
			Marker:         param.marker,
		}).
		SetResult(GetShareFileListResponse{}).
		SetError(ErrorResponse{}).
		Post(V3FileList)
	if err != nil {
		return err
	}
	data := resp.Result().(*GetShareFileListResponse)
	for _, item := range data.Items {
		if item.FileType == "folder" {
			err := ali.fileListByMarker(FileListParam{
				shareToken:   param.shareToken,
				shareId:      param.shareId,
				parentFileId: item.FileId,
				marker:       "",
			}, result)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		result <- item
	}
	if data.NextMarker != "" {
		err := ali.fileListByMarker(FileListParam{
			shareToken:   param.shareToken,
			shareId:      param.shareId,
			parentFileId: param.parentFileId,
			marker:       data.NextMarker,
		}, result)
		if err != nil {
			return err
		}
	}
	return nil
}
