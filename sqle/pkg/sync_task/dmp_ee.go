//go:build enterprise
// +build enterprise

package sync_task

// import (
// 	"bytes"
// 	"context"
// 	"crypto/rand"
// 	"crypto/rsa"
// 	"crypto/sha256"
// 	"crypto/x509"
// 	"encoding/base64"
// 	"encoding/json"
// 	"encoding/pem"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"sync"
// 	"time"

// 	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"

// 	"github.com/actiontech/sqle/sqle/common"

// 	"github.com/actiontech/sqle/sqle/model"
// 	"github.com/sirupsen/logrus"
// )

// const (
// 	DataObjectSourceDMPSupportedVersion = "5.23.01.0"
// 	Token                               = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJoYXNoIjoiWzUyIDMgMTI5IDE4NCAyMzggNjcgMiAzMCA5NSAyMzkgNTMgMjExIDE1MyAxNzQgODkgMjQ3XSIsIm5iZiI6MTY3MDU2MTg1NSwic2VlZCI6ImdaQ0J4RFJxUm1DQzJpaG4iLCJ1c2VyIjoic3FsZSJ9.SxGEi6QP8Dtl3ChsetDeZQxbcYpqsXmQibmytRuDbsg"
// 	SqleTag                             = "SQLE项目"
// )

// type DmpSync struct {
// 	SyncTaskID     uint
// 	DmpVersion     string
// 	Url            string
// 	DbType         string
// 	RuleTempleName string
// 	token          string
// 	client         *http.Client
// 	L              *logrus.Entry
// }

// func NewDmpSync(log *logrus.Entry, id uint, url, dmpVersion, dbType, ruleTemplateName string) *DmpSync {
// 	return &DmpSync{
// 		SyncTaskID:     id,
// 		DmpVersion:     dmpVersion,
// 		Url:            url,
// 		DbType:         dbType,
// 		RuleTempleName: ruleTemplateName,
// 		client:         &http.Client{},
// 		token:          Token,
// 		L:              log,
// 	}
// }

// type GetDmpInstanceResp struct {
// 	// data
// 	Data []*ListService `json:"data"`

// 	// total number of data sources. 数据源列表总计
// 	TotalNums uint32 `json:"total_nums"`
// }

// type Tag struct {
// 	// tag attribute. 标签名
// 	TagAttribute string `json:"tag_attribute"`
// 	// tag value. 标签值
// 	TagValue string `json:"tag_value"`
// }

// type ListService struct {
// 	// data source id. 数据源ID, 例如实例组 ID
// 	DataSrcID string `json:"data_src_id,omitempty"`

// 	// data source encrypted password. 数据源密码（已加密）
// 	DataSrcPassword string `json:"data_src_password,omitempty"`

// 	// data source port. 数据源端口
// 	DataSrcPort string `json:"data_src_port,omitempty"`

// 	// data source ip. 数据源的SIP，即实例组的SIP，或者固定不变的实例IP
// 	DataSrcSip string `json:"data_src_sip,omitempty"`

// 	// user name. 数据源用户名
// 	DataSrcUser string `json:"data_src_user,omitempty"`

// 	// tags  业务标签.
// 	Tags []Tag `json:"tags,omitempty"`
// }

// func (d *DmpSync) GetSyncInstanceTaskFunc(ctx context.Context) func() {
// 	return func() {
// 		d.startSyncDmpData(ctx)
// 	}
// }

// func (d *DmpSync) startSyncDmpData(ctx context.Context) {
// 	s := model.GetStorage()
// 	isSyncSuccess := false
// 	defer func() {
// 		m := make(map[string]interface{})
// 		if isSyncSuccess {
// 			m["last_sync_status"] = model.SyncInstanceStatusSucceeded
// 			m["last_sync_success_time"] = time.Now()
// 		} else {
// 			m["last_sync_status"] = model.SyncInstanceStatusFailed
// 		}

// 		if err := s.UpdateSyncInstanceTaskById(d.SyncTaskID, m); err != nil {
// 			d.L.Errorf("update sync instance task failed, err: %v", err)
// 		}
// 	}()

// 	if d.DmpVersion < DataObjectSourceDMPSupportedVersion {
// 		d.L.Errorf("dmp version %s not supported", d.DmpVersion)
// 		return
// 	}

// 	dmpFilterType := getDmpFilterType(d.DbType)

// 	url := fmt.Sprintf("%s/v3/support/data_sources?filter_by_type=%s", d.Url, dmpFilterType)

// 	body, err := d.do(ctx, d.client, http.MethodGet, url, nil)
// 	if err != nil {
// 		d.L.Errorf("get dmp data source fail: %s", err)
// 		return
// 	}

// 	var getDmpInstanceResp GetDmpInstanceResp
// 	if err := json.Unmarshal([]byte(body), &getDmpInstanceResp); err != nil {
// 		d.L.Errorf("unmarshal dmp data source fail: %s", err)
// 		return
// 	}

// 	if getDmpInstanceResp.TotalNums < 1 {
// 		isSyncSuccess = true
// 		d.L.Info("dmp data source total nums less than 1")
// 		return
// 	}

// 	ruleTemplate, exist, err := s.GetGlobalAndProjectRuleTemplateByNameAndProjectId(d.RuleTempleName, model.ProjectIdForGlobalRuleTemplate)
// 	if err != nil {
// 		d.L.Errorf("get rule template fail: %s", err)
// 		return
// 	}
// 	if !exist {
// 		d.L.Errorf("rule template %s not exist", d.RuleTempleName)
// 		return
// 	}

// 	dmpInst := make(map[string] /*project name*/ map[string] /*数据源名*/ struct{})
// 	instProjectName := make(map[string] /*数据源名*/ string /*project name*/)

// 	var syncTaskInstance *model.SyncTaskInstance
// 	var instances []*model.Instance
// 	for _, dmpInstance := range getDmpInstanceResp.Data {
// 		if dmpInstance.DataSrcSip == "" {
// 			d.L.Errorf("dmp data source %s sip is empty", dmpInstance.DataSrcID)
// 			continue
// 		}

// 		var projectName string
// 		for _, tag := range dmpInstance.Tags {
// 			if tag.TagAttribute == SqleTag {
// 				projectName = tag.TagValue
// 			}
// 		}

// 		if projectName == "" {
// 			d.L.Infof("dmp data source %s not have SqleTag,skip record", dmpInstance.DataSrcID)
// 			continue
// 		}

// 		if _, ok := dmpInst[projectName]; !ok {
// 			dmpInst[projectName] = make(map[string]struct{})
// 		}

// 		instProjectName[dmpInstance.DataSrcID] = projectName

// 		project, exist, err := s.GetProjectByName(projectName)
// 		if err != nil {
// 			d.L.Errorf("get Instances by project name fail: %s", err)
// 			return
// 		}
// 		if !exist {
// 			d.L.Errorf("project %s not exist", projectName)
// 			return
// 		}
// 		if project.Status == model.ProjectStatusArchived {
// 			continue
// 		}

// 		password, err := DecryptPassword(dmpInstance.DataSrcPassword)
// 		if err != nil {
// 			d.L.Errorf("decrypt password fail: %s", err)
// 			return
// 		}

// 		m := make(map[string]*model.Instance)
// 		for _, instance := range project.Instances {
// 			m[instance.Name] = instance
// 		}

// 		var inst *model.Instance
// 		var ok bool
// 		// 如果数据源已经存在，检测是否需要更新；如果数据源不存在，新增数据源到sqle
// 		if inst, ok = m[dmpInstance.DataSrcID]; ok {
// 			if inst.Source != model.SyncTaskSourceActiontechDmp {
// 				d.L.Errorf("instance has already exist and  %s source is not dmp", inst.Name)
// 				continue
// 			}

// 			isHostOrPortDiff := dmpInstance.DataSrcSip != inst.Host || dmpInstance.DataSrcPort != inst.Port
// 			isUserOrPasswdDiff := dmpInstance.DataSrcUser != inst.User || password != inst.Password
// 			if isHostOrPortDiff || isUserOrPasswdDiff {
// 				inst.Host = dmpInstance.DataSrcSip
// 				inst.Port = dmpInstance.DataSrcPort
// 				inst.User = dmpInstance.DataSrcUser
// 				inst.Password = password
// 			} else {
// 				continue
// 			}
// 		} else {
// 			inst = &model.Instance{
// 				Name:               dmpInstance.DataSrcID,
// 				Host:               dmpInstance.DataSrcSip,
// 				Port:               dmpInstance.DataSrcPort,
// 				User:               dmpInstance.DataSrcUser,
// 				WorkflowTemplateId: ruleTemplate.ID,
// 				ProjectId:          project.ID,
// 				Password:           password,
// 				DbType:             d.DbType,
// 				Source:             model.SyncTaskSourceActiontechDmp,
// 				SyncInstanceTaskID: d.SyncTaskID,
// 			}
// 		}

// 		instances = append(instances, inst)
// 	}

// 	for instName, projectName := range instProjectName {
// 		dmpInst[projectName][instName] = struct{}{}
// 	}

// 	canDeletedInstances, err := d.getNeedDeletedInstList(s, dmpInst)
// 	if err != nil {
// 		d.L.Errorf("get need deleted instance list fail: %s", err)
// 		return
// 	}

// 	syncTaskInstance = &model.SyncTaskInstance{
// 		Instances:            instances,
// 		RuleTemplate:         ruleTemplate,
// 		NeedDeletedInstances: canDeletedInstances,
// 	}

// 	if err := s.BatchUpdateSyncTask(syncTaskInstance); err != nil {
// 		d.L.Errorf("batch insert instance template fail: %s", err)
// 		return
// 	} else {
// 		isSyncSuccess = true
// 	}
// }

// func (d *DmpSync) getNeedDeletedInstList(s *model.Storage, dmpInst map[string]map[string]struct{}) ([]*model.Instance, error) {
// 	projectList, err := s.GetProjectListBySyncTaskId(d.SyncTaskID)
// 	if err != nil {
// 		return nil, fmt.Errorf("get project name list failed: %s", err)
// 	}

// 	var canDeletedInstList []*model.Instance
// 	var needDeletedInstList []*model.Instance
// 	for _, project := range projectList {
// 		// 项目不存在，删除该项目下所有该同步任务的数据源
// 		if _, ok := dmpInst[project.Name]; !ok {
// 			instanceList, err := s.GetInstancesBySyncTaskId(project.ID, d.SyncTaskID)
// 			if err != nil {
// 				return nil, fmt.Errorf("get instance list failed: %s", err)
// 			}
// 			needDeletedInstList = append(needDeletedInstList, instanceList...)
// 		} else {
// 			// 删除在sqle中存在，在dmp中不存在的数据源
// 			for _, instance := range project.Instances {
// 				if _, ok = dmpInst[project.Name][instance.Name]; !ok {
// 					needDeletedInstList = append(needDeletedInstList, instance)
// 				}
// 			}
// 		}
// 	}

// 	for _, instance := range needDeletedInstList {
// 		if err := common.CheckDeleteInstance(instance.ID); err == nil {
// 			d.L.Errorf("instance %s not exist in dmp, delete it", instance.Name)
// 			canDeletedInstList = append(canDeletedInstList, instance)
// 		}
// 	}

// 	return canDeletedInstList, nil
// }

// func getDmpFilterType(dbType string) string {
// 	switch dbType {
// 	case driverV2.DriverTypeMySQL:
// 		return "mysql"
// 	default:
// 		return ""
// 	}
// }

// func (d *DmpSync) do(ctx context.Context, client *http.Client, method, url string, body []byte) (respBody string, err error) {
// 	return retry(ctx, d.L, requester(ctx, client, method, url, d.token, body))
// }

// const maxRetries = 3

// func retry(ctx context.Context, l *logrus.Entry, f func() (response *http.Response, err error)) (respBody string, err error) {
// 	for i := 0; i < maxRetries; i++ {
// 		select {
// 		case <-ctx.Done():
// 			return "", ctx.Err()
// 		default:
// 		}

// 		resp, err := f()
// 		if err != nil {
// 			l.Errorf("request fail: %s", err)
// 			continue
// 		}

// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return "", fmt.Errorf("read response body fail: %s", err)
// 		}

// 		return string(body), nil
// 	}

// 	return "", fmt.Errorf("retry %d times fail", maxRetries)
// }

// func requester(ctx context.Context, client *http.Client, method, url string, token string, body []byte) func() (*http.Response, error) {
// 	return func() (*http.Response, error) {
// 		reader := bytes.NewReader(body)
// 		req, err := http.NewRequestWithContext(ctx, method, url, reader)
// 		if err != nil {
// 			return nil, fmt.Errorf("construct %s %s fail", method, url)
// 		}

// 		req.Header.Set("Content-Type", "application/json")
// 		req.Header.Add("Authorization", token)

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			return nil, fmt.Errorf("%s %s fail", method, url)
// 		}

// 		return resp, nil
// 	}
// }

// var buildinCaPrimaryKey = `-----BEGIN RSA PRIVATE KEY-----
// MIICXgIBAAKBgQDPr3VKlbvgP6NCRM4kuW+6wDEythhfxzVGgHG9Gthu3KUNxnd2
// RiR8RwntEgsWum+dn1oYnsV5TdaA3Wg/mJrP7U3XFmX/OVl2UyqHMd+k/bFAJh8v
// XKC2w9BtamzaSHdYro30skZ8jnxpjCD87+JtYgANySD0wRvlWwcD3slZ6wIDAQAB
// AoGBAIPz/7q2rdrJtAm7u5n7s7BcwiVtKslXwVKc8ybqMo8lYz0AVxBvemj3nafh
// aegz5gyonU69OcxblyjjA4Q8ikbhS6GOYxy27Oe6fYFjWzOWIMFkHe9QY7cqBkOL
// 8jzL0lzGyCWw+57l4h2tHPCctw/rfNi6NMFZUL2678H7u9rRAkEA/LlQ1WbZjHD2
// dy6+xpPiut0dHrOG034uqHTv6P8NchT/eHjvy+5SDKiJc/nk1POZmwXWJXKcvFC6
// zxwjvOAX2QJBANJgrbZRAIpwMJPNuOI+0ti3gIR7mpPZzwn0p2dJ4ORh48zp2b1M
// QghibXGJMEHKcfU39H7v50/H/lZx0f0liWMCQQDR+ldDN/VBTwo49EnmTDFx+Q2c
// 2KUJTCoQJTjAakoNo4yv2CvFUPozMkUia1rJ5KyXtT28V4IKpTjRpBu9bqPhAkEA
// nnIAAzMotBthCsDDQWrNhDlYiu9I8Zf2zem8dxd2UKvFVRy/SEn55bSz9vG7LaHa
// iDS3aS8oSLc4wESDQiSWPwJADu2OQHUDDPE7hGU788Dsess2gY0xmJR6z36mWftD
// Zz/GX75HZYICZBr6JjOVHHkLpByAWr5xonTLRyBhDqB7dg==
// -----END RSA PRIVATE KEY-----`

// var label = []byte("TVMvGp6rbgaBbWTU")

// func getCaPrimaryKey() (string, error) {
// 	// if secure.IsSecurityEnabled() {
// 	// 	caMu.Lock()
// 	// 	ret := caPrimaryKeyCache
// 	// 	caMu.Unlock()
// 	// 	if "" == ret {
// 	// 		security := NewSecurityConfig()
// 	// 		if err := orm.Get(security); nil != err {
// 	// 			return "", err
// 	// 		}
// 	// 		ret = security.Private
// 	// 		caMu.Lock()
// 	// 		caPrimaryKeyCache = security.Private
// 	// 		caMu.Unlock()
// 	// 	}
// 	// 	return ret, nil
// 	// } else {
// 	return buildinCaPrimaryKey, nil
// 	// }
// }

// var (
// 	decryptCache   = map[string]string{}
// 	decryptCacheMu = sync.RWMutex{}
// )

// func DecryptPassword(encrypted string) (string, error) {
// 	decryptCacheMu.RLock()
// 	if val, ok := decryptCache[encrypted]; ok {
// 		decryptCacheMu.RUnlock()
// 		return val, nil
// 	}
// 	decryptCacheMu.RUnlock()

// 	caKey, err := getCaPrimaryKey()
// 	if nil != err {
// 		return "", err
// 	}

// 	block, _ := pem.Decode([]byte(caKey))
// 	if nil == block {
// 		return "", errors.New("public GetKey error")
// 	}

// 	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

// 	secretMessage, err := base64.StdEncoding.DecodeString(encrypted)
// 	if nil != err {
// 		return "", err
// 	}

// 	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, secretMessage, label)
// 	if nil != err {
// 		return "", err
// 	}

// 	password := string(decrypted)

// 	decryptCacheMu.Lock()
// 	if 10240 < len(decryptCache) {
// 		decryptCache = map[string]string{encrypted: password}

// 	} else {
// 		decryptCache[encrypted] = password
// 	}
// 	decryptCacheMu.Unlock()

// 	return password, nil
// }
