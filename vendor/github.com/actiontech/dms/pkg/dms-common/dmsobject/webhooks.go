package dmsobject

import (
	"context"
	"fmt"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

func WebHookSendMessage(ctx context.Context, dmsAddr string, req *dmsV1.WebHookSendMessageReq) (err error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reply := &dmsV1.NotificationReply{}
	url := fmt.Sprintf("%v%v", dmsAddr, dmsV1.GetWebHooksRouter())

	if err := pkgHttp.POST(ctx, url, header, req, reply); err != nil {
		return fmt.Errorf("failed to notify by %v: %v", url, err)
	}
	if reply.Code != 0 {
		return fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return nil
}
