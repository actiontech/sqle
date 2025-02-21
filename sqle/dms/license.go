package dms

import (
	"context"
	"fmt"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

type GetLicenseReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Content string `json:"content"` // 只需要Content因此只解析Content
}

// GET /v1/dms/configurations/license Configuration GetLicense
func GetLicense(ctx context.Context, dmsAddr string) (*GetLicenseReply, error) {
	// Prepare the header with authorization token (assuming it's similar to the other request)
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken, // Replace with your token if needed
	}

	// Define the response structure to hold the data
	reply := &GetLicenseReply{}

	// Construct the URL for the GET request
	url := fmt.Sprintf("%v/v1/dms/configurations/license", dmsAddr)

	// Perform the GET request
	if err := pkgHttp.Get(ctx, url, header, nil, reply); err != nil {
		return nil, fmt.Errorf("failed to get license from %v: %v", url, err)
	}

	// Handle the response code (if the reply structure contains such validation)
	if reply.Code != 0 { // Assuming a Code field in the reply, modify accordingly
		return nil, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	// Return the retrieved license data
	return reply, nil
}
