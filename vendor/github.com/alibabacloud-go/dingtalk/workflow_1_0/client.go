// This file is auto-generated, don't edit it. Thanks.
/**
 *
 */
package workflow_1_0

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type AvaliableTemplate struct {
	// 表单名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 表单模板processCode
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s AvaliableTemplate) String() string {
	return tea.Prettify(s)
}

func (s AvaliableTemplate) GoString() string {
	return s.String()
}

func (s *AvaliableTemplate) SetName(v string) *AvaliableTemplate {
	s.Name = &v
	return s
}

func (s *AvaliableTemplate) SetProcessCode(v string) *AvaliableTemplate {
	s.ProcessCode = &v
	return s
}

type FormComponent struct {
	// 子控件集合
	Children []*FormComponent `json:"children,omitempty" xml:"children,omitempty" type:"Repeated"`
	// 控件类型
	ComponentType *string `json:"componentType,omitempty" xml:"componentType,omitempty"`
	// 控件属性
	Props *FormComponentProps `json:"props,omitempty" xml:"props,omitempty"`
}

func (s FormComponent) String() string {
	return tea.Prettify(s)
}

func (s FormComponent) GoString() string {
	return s.String()
}

func (s *FormComponent) SetChildren(v []*FormComponent) *FormComponent {
	s.Children = v
	return s
}

func (s *FormComponent) SetComponentType(v string) *FormComponent {
	s.ComponentType = &v
	return s
}

func (s *FormComponent) SetProps(v *FormComponentProps) *FormComponent {
	s.Props = v
	return s
}

type FormComponentProps struct {
	// 明细控件按钮显示文案
	ActionName *string `json:"actionName,omitempty" xml:"actionName,omitempty"`
	// 地址控件模式city省市,district省市区,street省市区街道
	AddressModel *string `json:"addressModel,omitempty" xml:"addressModel,omitempty"`
	// 文字提示控件显示方式:top|middle|bottom
	Align *string `json:"align,omitempty" xml:"align,omitempty"`
	// 套件中控件是否可设置为分条件字段
	AsyncCondition *bool `json:"asyncCondition,omitempty" xml:"asyncCondition,omitempty"`
	// 关联审批单控件限定模板列表
	AvailableTemplates []*AvaliableTemplate `json:"availableTemplates,omitempty" xml:"availableTemplates,omitempty" type:"Repeated"`
	// 业务别名
	BizAlias *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	// 套件的业务标识
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 联系人控件是否支持多选，1多选，0单选
	Choice *string `json:"choice,omitempty" xml:"choice,omitempty"`
	// 自定义控件渲染标识
	CommonBizType *string `json:"commonBizType,omitempty" xml:"commonBizType,omitempty"`
	// 控件表单内唯一id
	ComponentId *string `json:"componentId,omitempty" xml:"componentId,omitempty"`
	// 说明文字控件内容
	Content *string `json:"content,omitempty" xml:"content,omitempty"`
	// 关联数据源配置
	DataSource *FormDataSource `json:"dataSource,omitempty" xml:"dataSource,omitempty"`
	// 是否不可编辑
	Disabled *bool `json:"disabled,omitempty" xml:"disabled,omitempty"`
	// 是否自动计算时长
	Duration *bool `json:"duration,omitempty" xml:"duration,omitempty"`
	// 时间格式
	Format *string `json:"format,omitempty" xml:"format,omitempty"`
	// 公式
	Formula *string `json:"formula,omitempty" xml:"formula,omitempty"`
	// 是否隐藏字段
	Invisible *bool `json:"invisible,omitempty" xml:"invisible,omitempty"`
	// 控件标题
	Label *string `json:"label,omitempty" xml:"label,omitempty"`
	// 评分控件上限
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// 说明文字控件链接地址
	Link *string `json:"link,omitempty" xml:"link,omitempty"`
	// 文本控件支持的最大长度
	MaxLength *int32 `json:"maxLength,omitempty" xml:"maxLength,omitempty"`
	// 电话控件支持的类型
	Mode *string `json:"mode,omitempty" xml:"mode,omitempty"`
	// 部门控件是否可多选
	Multiple *bool `json:"multiple,omitempty" xml:"multiple,omitempty"`
	// 单选多选控件选项列表
	Options []*SelectOption `json:"options,omitempty" xml:"options,omitempty" type:"Repeated"`
	// 输入提示
	Placeholder *string `json:"placeholder,omitempty" xml:"placeholder,omitempty"`
	// 字段是否可打印，1打印，0不打印，默认打印
	Print *string `json:"print,omitempty" xml:"print,omitempty"`
	// 是否必填
	Required *bool `json:"required,omitempty" xml:"required,omitempty"`
	// 明细控件数据汇总统计
	StatField []*FormComponentPropsStatField `json:"statField,omitempty" xml:"statField,omitempty" type:"Repeated"`
	// 明细填写方式，table（表格）、list（列表）
	TableViewMode *string `json:"tableViewMode,omitempty" xml:"tableViewMode,omitempty"`
	// 时间单位（天、小时）
	Unit *string `json:"unit,omitempty" xml:"unit,omitempty"`
	// 金额控件是否需要大写，1不需要大写，其他需要大写
	Upper *string `json:"upper,omitempty" xml:"upper,omitempty"`
	// 明细打印方式false横向，true纵向
	VerticalPrint *bool `json:"verticalPrint,omitempty" xml:"verticalPrint,omitempty"`
}

func (s FormComponentProps) String() string {
	return tea.Prettify(s)
}

func (s FormComponentProps) GoString() string {
	return s.String()
}

func (s *FormComponentProps) SetActionName(v string) *FormComponentProps {
	s.ActionName = &v
	return s
}

func (s *FormComponentProps) SetAddressModel(v string) *FormComponentProps {
	s.AddressModel = &v
	return s
}

func (s *FormComponentProps) SetAlign(v string) *FormComponentProps {
	s.Align = &v
	return s
}

func (s *FormComponentProps) SetAsyncCondition(v bool) *FormComponentProps {
	s.AsyncCondition = &v
	return s
}

func (s *FormComponentProps) SetAvailableTemplates(v []*AvaliableTemplate) *FormComponentProps {
	s.AvailableTemplates = v
	return s
}

func (s *FormComponentProps) SetBizAlias(v string) *FormComponentProps {
	s.BizAlias = &v
	return s
}

func (s *FormComponentProps) SetBizType(v string) *FormComponentProps {
	s.BizType = &v
	return s
}

func (s *FormComponentProps) SetChoice(v string) *FormComponentProps {
	s.Choice = &v
	return s
}

func (s *FormComponentProps) SetCommonBizType(v string) *FormComponentProps {
	s.CommonBizType = &v
	return s
}

func (s *FormComponentProps) SetComponentId(v string) *FormComponentProps {
	s.ComponentId = &v
	return s
}

func (s *FormComponentProps) SetContent(v string) *FormComponentProps {
	s.Content = &v
	return s
}

func (s *FormComponentProps) SetDataSource(v *FormDataSource) *FormComponentProps {
	s.DataSource = v
	return s
}

func (s *FormComponentProps) SetDisabled(v bool) *FormComponentProps {
	s.Disabled = &v
	return s
}

func (s *FormComponentProps) SetDuration(v bool) *FormComponentProps {
	s.Duration = &v
	return s
}

func (s *FormComponentProps) SetFormat(v string) *FormComponentProps {
	s.Format = &v
	return s
}

func (s *FormComponentProps) SetFormula(v string) *FormComponentProps {
	s.Formula = &v
	return s
}

func (s *FormComponentProps) SetInvisible(v bool) *FormComponentProps {
	s.Invisible = &v
	return s
}

func (s *FormComponentProps) SetLabel(v string) *FormComponentProps {
	s.Label = &v
	return s
}

func (s *FormComponentProps) SetLimit(v int32) *FormComponentProps {
	s.Limit = &v
	return s
}

func (s *FormComponentProps) SetLink(v string) *FormComponentProps {
	s.Link = &v
	return s
}

func (s *FormComponentProps) SetMaxLength(v int32) *FormComponentProps {
	s.MaxLength = &v
	return s
}

func (s *FormComponentProps) SetMode(v string) *FormComponentProps {
	s.Mode = &v
	return s
}

func (s *FormComponentProps) SetMultiple(v bool) *FormComponentProps {
	s.Multiple = &v
	return s
}

func (s *FormComponentProps) SetOptions(v []*SelectOption) *FormComponentProps {
	s.Options = v
	return s
}

func (s *FormComponentProps) SetPlaceholder(v string) *FormComponentProps {
	s.Placeholder = &v
	return s
}

func (s *FormComponentProps) SetPrint(v string) *FormComponentProps {
	s.Print = &v
	return s
}

func (s *FormComponentProps) SetRequired(v bool) *FormComponentProps {
	s.Required = &v
	return s
}

func (s *FormComponentProps) SetStatField(v []*FormComponentPropsStatField) *FormComponentProps {
	s.StatField = v
	return s
}

func (s *FormComponentProps) SetTableViewMode(v string) *FormComponentProps {
	s.TableViewMode = &v
	return s
}

func (s *FormComponentProps) SetUnit(v string) *FormComponentProps {
	s.Unit = &v
	return s
}

func (s *FormComponentProps) SetUpper(v string) *FormComponentProps {
	s.Upper = &v
	return s
}

func (s *FormComponentProps) SetVerticalPrint(v bool) *FormComponentProps {
	s.VerticalPrint = &v
	return s
}

type FormComponentPropsStatField struct {
	// 需要统计的明细控件内子控件id
	ComponentId *string `json:"componentId,omitempty" xml:"componentId,omitempty"`
	// 子控件标题
	Label *string `json:"label,omitempty" xml:"label,omitempty"`
	// 金额控件是否需要大写，1不需要大写，其他需要大写
	Upper *string `json:"upper,omitempty" xml:"upper,omitempty"`
}

func (s FormComponentPropsStatField) String() string {
	return tea.Prettify(s)
}

func (s FormComponentPropsStatField) GoString() string {
	return s.String()
}

func (s *FormComponentPropsStatField) SetComponentId(v string) *FormComponentPropsStatField {
	s.ComponentId = &v
	return s
}

func (s *FormComponentPropsStatField) SetLabel(v string) *FormComponentPropsStatField {
	s.Label = &v
	return s
}

func (s *FormComponentPropsStatField) SetUpper(v string) *FormComponentPropsStatField {
	s.Upper = &v
	return s
}

type FormDataSource struct {
	// 关联表单信息
	Target *FormDataSourceTarget `json:"target,omitempty" xml:"target,omitempty" type:"Struct"`
	// 关联类型，form关联表单
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s FormDataSource) String() string {
	return tea.Prettify(s)
}

func (s FormDataSource) GoString() string {
	return s.String()
}

func (s *FormDataSource) SetTarget(v *FormDataSourceTarget) *FormDataSource {
	s.Target = v
	return s
}

func (s *FormDataSource) SetType(v string) *FormDataSource {
	s.Type = &v
	return s
}

type FormDataSourceTarget struct {
	// 表单类型，0流程表单
	AppType *int32 `json:"appType,omitempty" xml:"appType,omitempty"`
	// 应用appUuid
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 关联表单业务标识
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 关联表单的formCode
	FormCode *string `json:"formCode,omitempty" xml:"formCode,omitempty"`
}

func (s FormDataSourceTarget) String() string {
	return tea.Prettify(s)
}

func (s FormDataSourceTarget) GoString() string {
	return s.String()
}

func (s *FormDataSourceTarget) SetAppType(v int32) *FormDataSourceTarget {
	s.AppType = &v
	return s
}

func (s *FormDataSourceTarget) SetAppUuid(v string) *FormDataSourceTarget {
	s.AppUuid = &v
	return s
}

func (s *FormDataSourceTarget) SetBizType(v string) *FormDataSourceTarget {
	s.BizType = &v
	return s
}

func (s *FormDataSourceTarget) SetFormCode(v string) *FormDataSourceTarget {
	s.FormCode = &v
	return s
}

type SelectOption struct {
	// 每一个选项的唯一键
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// 每一个选项的值
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s SelectOption) String() string {
	return tea.Prettify(s)
}

func (s SelectOption) GoString() string {
	return s.String()
}

func (s *SelectOption) SetKey(v string) *SelectOption {
	s.Key = &v
	return s
}

func (s *SelectOption) SetValue(v string) *SelectOption {
	s.Value = &v
	return s
}

type AddApproveDentryAuthHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s AddApproveDentryAuthHeaders) String() string {
	return tea.Prettify(s)
}

func (s AddApproveDentryAuthHeaders) GoString() string {
	return s.String()
}

func (s *AddApproveDentryAuthHeaders) SetCommonHeaders(v map[string]*string) *AddApproveDentryAuthHeaders {
	s.CommonHeaders = v
	return s
}

func (s *AddApproveDentryAuthHeaders) SetXAcsDingtalkAccessToken(v string) *AddApproveDentryAuthHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type AddApproveDentryAuthRequest struct {
	// 授权的钉盘文件信息列表。支持批量授权，最大列表长度：10。
	FileInfos []*AddApproveDentryAuthRequestFileInfos `json:"fileInfos,omitempty" xml:"fileInfos,omitempty" type:"Repeated"`
	// 授权的用户userid。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s AddApproveDentryAuthRequest) String() string {
	return tea.Prettify(s)
}

func (s AddApproveDentryAuthRequest) GoString() string {
	return s.String()
}

func (s *AddApproveDentryAuthRequest) SetFileInfos(v []*AddApproveDentryAuthRequestFileInfos) *AddApproveDentryAuthRequest {
	s.FileInfos = v
	return s
}

func (s *AddApproveDentryAuthRequest) SetUserId(v string) *AddApproveDentryAuthRequest {
	s.UserId = &v
	return s
}

type AddApproveDentryAuthRequestFileInfos struct {
	// 文件ID。
	FileId *string `json:"fileId,omitempty" xml:"fileId,omitempty"`
	// 钉盘空间spaceId。
	SpaceId *int64 `json:"spaceId,omitempty" xml:"spaceId,omitempty"`
}

func (s AddApproveDentryAuthRequestFileInfos) String() string {
	return tea.Prettify(s)
}

func (s AddApproveDentryAuthRequestFileInfos) GoString() string {
	return s.String()
}

func (s *AddApproveDentryAuthRequestFileInfos) SetFileId(v string) *AddApproveDentryAuthRequestFileInfos {
	s.FileId = &v
	return s
}

func (s *AddApproveDentryAuthRequestFileInfos) SetSpaceId(v int64) *AddApproveDentryAuthRequestFileInfos {
	s.SpaceId = &v
	return s
}

type AddApproveDentryAuthResponseBody struct {
	// 返回结果。
	Result *bool `json:"result,omitempty" xml:"result,omitempty"`
	// 接口调用是否成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s AddApproveDentryAuthResponseBody) String() string {
	return tea.Prettify(s)
}

func (s AddApproveDentryAuthResponseBody) GoString() string {
	return s.String()
}

func (s *AddApproveDentryAuthResponseBody) SetResult(v bool) *AddApproveDentryAuthResponseBody {
	s.Result = &v
	return s
}

func (s *AddApproveDentryAuthResponseBody) SetSuccess(v bool) *AddApproveDentryAuthResponseBody {
	s.Success = &v
	return s
}

type AddApproveDentryAuthResponse struct {
	Headers map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *AddApproveDentryAuthResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s AddApproveDentryAuthResponse) String() string {
	return tea.Prettify(s)
}

func (s AddApproveDentryAuthResponse) GoString() string {
	return s.String()
}

func (s *AddApproveDentryAuthResponse) SetHeaders(v map[string]*string) *AddApproveDentryAuthResponse {
	s.Headers = v
	return s
}

func (s *AddApproveDentryAuthResponse) SetBody(v *AddApproveDentryAuthResponseBody) *AddApproveDentryAuthResponse {
	s.Body = v
	return s
}

type AddProcessInstanceCommentHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s AddProcessInstanceCommentHeaders) String() string {
	return tea.Prettify(s)
}

func (s AddProcessInstanceCommentHeaders) GoString() string {
	return s.String()
}

func (s *AddProcessInstanceCommentHeaders) SetCommonHeaders(v map[string]*string) *AddProcessInstanceCommentHeaders {
	s.CommonHeaders = v
	return s
}

func (s *AddProcessInstanceCommentHeaders) SetXAcsDingtalkAccessToken(v string) *AddProcessInstanceCommentHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type AddProcessInstanceCommentRequest struct {
	// 评论人的userid。
	CommentUserId *string `json:"commentUserId,omitempty" xml:"commentUserId,omitempty"`
	// 文件。
	File *AddProcessInstanceCommentRequestFile `json:"file,omitempty" xml:"file,omitempty" type:"Struct"`
	// 审批实例ID，可通过调用获取审批实例ID列表接口获取。
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 评论的内容。
	Text *string `json:"text,omitempty" xml:"text,omitempty"`
}

func (s AddProcessInstanceCommentRequest) String() string {
	return tea.Prettify(s)
}

func (s AddProcessInstanceCommentRequest) GoString() string {
	return s.String()
}

func (s *AddProcessInstanceCommentRequest) SetCommentUserId(v string) *AddProcessInstanceCommentRequest {
	s.CommentUserId = &v
	return s
}

func (s *AddProcessInstanceCommentRequest) SetFile(v *AddProcessInstanceCommentRequestFile) *AddProcessInstanceCommentRequest {
	s.File = v
	return s
}

func (s *AddProcessInstanceCommentRequest) SetProcessInstanceId(v string) *AddProcessInstanceCommentRequest {
	s.ProcessInstanceId = &v
	return s
}

func (s *AddProcessInstanceCommentRequest) SetText(v string) *AddProcessInstanceCommentRequest {
	s.Text = &v
	return s
}

type AddProcessInstanceCommentRequestFile struct {
	// 附件列表。
	Attachments []*AddProcessInstanceCommentRequestFileAttachments `json:"attachments,omitempty" xml:"attachments,omitempty" type:"Repeated"`
	// 图片URL地址。
	Photos []*string `json:"photos,omitempty" xml:"photos,omitempty" type:"Repeated"`
}

func (s AddProcessInstanceCommentRequestFile) String() string {
	return tea.Prettify(s)
}

func (s AddProcessInstanceCommentRequestFile) GoString() string {
	return s.String()
}

func (s *AddProcessInstanceCommentRequestFile) SetAttachments(v []*AddProcessInstanceCommentRequestFileAttachments) *AddProcessInstanceCommentRequestFile {
	s.Attachments = v
	return s
}

func (s *AddProcessInstanceCommentRequestFile) SetPhotos(v []*string) *AddProcessInstanceCommentRequestFile {
	s.Photos = v
	return s
}

type AddProcessInstanceCommentRequestFileAttachments struct {
	// 文件ID。
	FileId *string `json:"fileId,omitempty" xml:"fileId,omitempty"`
	// 文件名称。
	FileName *string `json:"fileName,omitempty" xml:"fileName,omitempty"`
	// 文件大小。
	FileSize *string `json:"fileSize,omitempty" xml:"fileSize,omitempty"`
	// 文件类型。
	FileType *string `json:"fileType,omitempty" xml:"fileType,omitempty"`
	// 钉盘空间ID。
	SpaceId *string `json:"spaceId,omitempty" xml:"spaceId,omitempty"`
}

func (s AddProcessInstanceCommentRequestFileAttachments) String() string {
	return tea.Prettify(s)
}

func (s AddProcessInstanceCommentRequestFileAttachments) GoString() string {
	return s.String()
}

func (s *AddProcessInstanceCommentRequestFileAttachments) SetFileId(v string) *AddProcessInstanceCommentRequestFileAttachments {
	s.FileId = &v
	return s
}

func (s *AddProcessInstanceCommentRequestFileAttachments) SetFileName(v string) *AddProcessInstanceCommentRequestFileAttachments {
	s.FileName = &v
	return s
}

func (s *AddProcessInstanceCommentRequestFileAttachments) SetFileSize(v string) *AddProcessInstanceCommentRequestFileAttachments {
	s.FileSize = &v
	return s
}

func (s *AddProcessInstanceCommentRequestFileAttachments) SetFileType(v string) *AddProcessInstanceCommentRequestFileAttachments {
	s.FileType = &v
	return s
}

func (s *AddProcessInstanceCommentRequestFileAttachments) SetSpaceId(v string) *AddProcessInstanceCommentRequestFileAttachments {
	s.SpaceId = &v
	return s
}

type AddProcessInstanceCommentResponseBody struct {
	// 评论是否成功。
	Result *bool `json:"result,omitempty" xml:"result,omitempty"`
	// 接口调用是否成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s AddProcessInstanceCommentResponseBody) String() string {
	return tea.Prettify(s)
}

func (s AddProcessInstanceCommentResponseBody) GoString() string {
	return s.String()
}

func (s *AddProcessInstanceCommentResponseBody) SetResult(v bool) *AddProcessInstanceCommentResponseBody {
	s.Result = &v
	return s
}

func (s *AddProcessInstanceCommentResponseBody) SetSuccess(v bool) *AddProcessInstanceCommentResponseBody {
	s.Success = &v
	return s
}

type AddProcessInstanceCommentResponse struct {
	Headers map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *AddProcessInstanceCommentResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s AddProcessInstanceCommentResponse) String() string {
	return tea.Prettify(s)
}

func (s AddProcessInstanceCommentResponse) GoString() string {
	return s.String()
}

func (s *AddProcessInstanceCommentResponse) SetHeaders(v map[string]*string) *AddProcessInstanceCommentResponse {
	s.Headers = v
	return s
}

func (s *AddProcessInstanceCommentResponse) SetBody(v *AddProcessInstanceCommentResponseBody) *AddProcessInstanceCommentResponse {
	s.Body = v
	return s
}

type BatchUpdateProcessInstanceHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s BatchUpdateProcessInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s BatchUpdateProcessInstanceHeaders) GoString() string {
	return s.String()
}

func (s *BatchUpdateProcessInstanceHeaders) SetCommonHeaders(v map[string]*string) *BatchUpdateProcessInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *BatchUpdateProcessInstanceHeaders) SetXAcsDingtalkAccessToken(v string) *BatchUpdateProcessInstanceHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type BatchUpdateProcessInstanceRequest struct {
	// 实列列表。
	UpdateProcessInstanceRequests []*BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests `json:"updateProcessInstanceRequests,omitempty" xml:"updateProcessInstanceRequests,omitempty" type:"Repeated"`
}

func (s BatchUpdateProcessInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s BatchUpdateProcessInstanceRequest) GoString() string {
	return s.String()
}

func (s *BatchUpdateProcessInstanceRequest) SetUpdateProcessInstanceRequests(v []*BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests) *BatchUpdateProcessInstanceRequest {
	s.UpdateProcessInstanceRequests = v
	return s
}

type BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests struct {
	// 抄送列表，注意最大抄送人数量不能超过30。
	Notifiers []*BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequestsNotifiers `json:"notifiers,omitempty" xml:"notifiers,omitempty" type:"Repeated"`
	// 实例id
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 实例结果：
	// 实例状态是COMPLETED，必须设置代表以下含义。
	// agree：同意
	// refuse：拒绝
	// 实例状态为TERMINATED，必须设置代表含义，result取值agree和refuse均代表撤销审批流。
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 实例状态：
	// COMPLETED：结束审批流
	// TERMINATED：终止审批流
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests) String() string {
	return tea.Prettify(s)
}

func (s BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests) GoString() string {
	return s.String()
}

func (s *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests) SetNotifiers(v []*BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequestsNotifiers) *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests {
	s.Notifiers = v
	return s
}

func (s *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests) SetProcessInstanceId(v string) *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests {
	s.ProcessInstanceId = &v
	return s
}

func (s *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests) SetResult(v string) *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests {
	s.Result = &v
	return s
}

func (s *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests) SetStatus(v string) *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequests {
	s.Status = &v
	return s
}

type BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequestsNotifiers struct {
	// 抄送接收人用户userId。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequestsNotifiers) String() string {
	return tea.Prettify(s)
}

func (s BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequestsNotifiers) GoString() string {
	return s.String()
}

func (s *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequestsNotifiers) SetUserId(v string) *BatchUpdateProcessInstanceRequestUpdateProcessInstanceRequestsNotifiers {
	s.UserId = &v
	return s
}

type BatchUpdateProcessInstanceResponseBody struct {
	// 成功标识
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s BatchUpdateProcessInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s BatchUpdateProcessInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *BatchUpdateProcessInstanceResponseBody) SetSuccess(v bool) *BatchUpdateProcessInstanceResponseBody {
	s.Success = &v
	return s
}

type BatchUpdateProcessInstanceResponse struct {
	Headers map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *BatchUpdateProcessInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s BatchUpdateProcessInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s BatchUpdateProcessInstanceResponse) GoString() string {
	return s.String()
}

func (s *BatchUpdateProcessInstanceResponse) SetHeaders(v map[string]*string) *BatchUpdateProcessInstanceResponse {
	s.Headers = v
	return s
}

func (s *BatchUpdateProcessInstanceResponse) SetBody(v *BatchUpdateProcessInstanceResponseBody) *BatchUpdateProcessInstanceResponse {
	s.Body = v
	return s
}

type CancelIntegratedTaskHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s CancelIntegratedTaskHeaders) String() string {
	return tea.Prettify(s)
}

func (s CancelIntegratedTaskHeaders) GoString() string {
	return s.String()
}

func (s *CancelIntegratedTaskHeaders) SetCommonHeaders(v map[string]*string) *CancelIntegratedTaskHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CancelIntegratedTaskHeaders) SetXAcsDingtalkAccessToken(v string) *CancelIntegratedTaskHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type CancelIntegratedTaskRequest struct {
	// 待办组ID，需要在调用创建待办接口时，主动设置该值。
	ActivityId *string `json:"activityId,omitempty" xml:"activityId,omitempty"`
	// 待办组ID列表，用于批量取消。
	ActivityIds []*string `json:"activityIds,omitempty" xml:"activityIds,omitempty" type:"Repeated"`
	// OA审批流程实例ID
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
}

func (s CancelIntegratedTaskRequest) String() string {
	return tea.Prettify(s)
}

func (s CancelIntegratedTaskRequest) GoString() string {
	return s.String()
}

func (s *CancelIntegratedTaskRequest) SetActivityId(v string) *CancelIntegratedTaskRequest {
	s.ActivityId = &v
	return s
}

func (s *CancelIntegratedTaskRequest) SetActivityIds(v []*string) *CancelIntegratedTaskRequest {
	s.ActivityIds = v
	return s
}

func (s *CancelIntegratedTaskRequest) SetProcessInstanceId(v string) *CancelIntegratedTaskRequest {
	s.ProcessInstanceId = &v
	return s
}

type CancelIntegratedTaskResponseBody struct {
	// 是否更新成功
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s CancelIntegratedTaskResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CancelIntegratedTaskResponseBody) GoString() string {
	return s.String()
}

func (s *CancelIntegratedTaskResponseBody) SetSuccess(v bool) *CancelIntegratedTaskResponseBody {
	s.Success = &v
	return s
}

type CancelIntegratedTaskResponse struct {
	Headers map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *CancelIntegratedTaskResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CancelIntegratedTaskResponse) String() string {
	return tea.Prettify(s)
}

func (s CancelIntegratedTaskResponse) GoString() string {
	return s.String()
}

func (s *CancelIntegratedTaskResponse) SetHeaders(v map[string]*string) *CancelIntegratedTaskResponse {
	s.Headers = v
	return s
}

func (s *CancelIntegratedTaskResponse) SetBody(v *CancelIntegratedTaskResponseBody) *CancelIntegratedTaskResponse {
	s.Body = v
	return s
}

type CleanProcessDataHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s CleanProcessDataHeaders) String() string {
	return tea.Prettify(s)
}

func (s CleanProcessDataHeaders) GoString() string {
	return s.String()
}

func (s *CleanProcessDataHeaders) SetCommonHeaders(v map[string]*string) *CleanProcessDataHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CleanProcessDataHeaders) SetXAcsDingtalkAccessToken(v string) *CleanProcessDataHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type CleanProcessDataRequest struct {
	// 企业的corpId。
	CorpId *string `json:"corpId,omitempty" xml:"corpId,omitempty"`
	// 模板唯一码。
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s CleanProcessDataRequest) String() string {
	return tea.Prettify(s)
}

func (s CleanProcessDataRequest) GoString() string {
	return s.String()
}

func (s *CleanProcessDataRequest) SetCorpId(v string) *CleanProcessDataRequest {
	s.CorpId = &v
	return s
}

func (s *CleanProcessDataRequest) SetProcessCode(v string) *CleanProcessDataRequest {
	s.ProcessCode = &v
	return s
}

type CleanProcessDataResponseBody struct {
	// 是否调用成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s CleanProcessDataResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CleanProcessDataResponseBody) GoString() string {
	return s.String()
}

func (s *CleanProcessDataResponseBody) SetSuccess(v bool) *CleanProcessDataResponseBody {
	s.Success = &v
	return s
}

type CleanProcessDataResponse struct {
	Headers map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *CleanProcessDataResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CleanProcessDataResponse) String() string {
	return tea.Prettify(s)
}

func (s CleanProcessDataResponse) GoString() string {
	return s.String()
}

func (s *CleanProcessDataResponse) SetHeaders(v map[string]*string) *CleanProcessDataResponse {
	s.Headers = v
	return s
}

func (s *CleanProcessDataResponse) SetBody(v *CleanProcessDataResponseBody) *CleanProcessDataResponse {
	s.Body = v
	return s
}

type CopyProcessHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s CopyProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s CopyProcessHeaders) GoString() string {
	return s.String()
}

func (s *CopyProcessHeaders) SetCommonHeaders(v map[string]*string) *CopyProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CopyProcessHeaders) SetXAcsDingtalkAccessToken(v string) *CopyProcessHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type CopyProcessRequest struct {
	// 复制选项
	CopyOptions  *CopyProcessRequestCopyOptions `json:"copyOptions,omitempty" xml:"copyOptions,omitempty" type:"Struct"`
	SourceCorpId *string                        `json:"sourceCorpId,omitempty" xml:"sourceCorpId,omitempty"`
	// 源模版列表
	SourceProcessVOList []*CopyProcessRequestSourceProcessVOList `json:"sourceProcessVOList,omitempty" xml:"sourceProcessVOList,omitempty" type:"Repeated"`
}

func (s CopyProcessRequest) String() string {
	return tea.Prettify(s)
}

func (s CopyProcessRequest) GoString() string {
	return s.String()
}

func (s *CopyProcessRequest) SetCopyOptions(v *CopyProcessRequestCopyOptions) *CopyProcessRequest {
	s.CopyOptions = v
	return s
}

func (s *CopyProcessRequest) SetSourceCorpId(v string) *CopyProcessRequest {
	s.SourceCorpId = &v
	return s
}

func (s *CopyProcessRequest) SetSourceProcessVOList(v []*CopyProcessRequestSourceProcessVOList) *CopyProcessRequest {
	s.SourceProcessVOList = v
	return s
}

type CopyProcessRequestCopyOptions struct {
	// 默认为1
	// COPE_TYPE_DEFAULT = 1 默认会使用groupId进行隔离分组，审批首页不可见
	// COPE_TYPE_INCLUDE_FLOW = 2 使用dingtalk 2作为隔离分组，审批首页可见
	CopyType *int32 `json:"copyType,omitempty" xml:"copyType,omitempty"`
}

func (s CopyProcessRequestCopyOptions) String() string {
	return tea.Prettify(s)
}

func (s CopyProcessRequestCopyOptions) GoString() string {
	return s.String()
}

func (s *CopyProcessRequestCopyOptions) SetCopyType(v int32) *CopyProcessRequestCopyOptions {
	s.CopyType = &v
	return s
}

type CopyProcessRequestSourceProcessVOList struct {
	// 套件业务标识
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 模板名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 模板code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s CopyProcessRequestSourceProcessVOList) String() string {
	return tea.Prettify(s)
}

func (s CopyProcessRequestSourceProcessVOList) GoString() string {
	return s.String()
}

func (s *CopyProcessRequestSourceProcessVOList) SetBizType(v string) *CopyProcessRequestSourceProcessVOList {
	s.BizType = &v
	return s
}

func (s *CopyProcessRequestSourceProcessVOList) SetName(v string) *CopyProcessRequestSourceProcessVOList {
	s.Name = &v
	return s
}

func (s *CopyProcessRequestSourceProcessVOList) SetProcessCode(v string) *CopyProcessRequestSourceProcessVOList {
	s.ProcessCode = &v
	return s
}

type CopyProcessResponseBody struct {
	Result []*CopyProcessResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Repeated"`
}

func (s CopyProcessResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CopyProcessResponseBody) GoString() string {
	return s.String()
}

func (s *CopyProcessResponseBody) SetResult(v []*CopyProcessResponseBodyResult) *CopyProcessResponseBody {
	s.Result = v
	return s
}

type CopyProcessResponseBodyResult struct {
	// 业务标识
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 模板名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 模板code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s CopyProcessResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s CopyProcessResponseBodyResult) GoString() string {
	return s.String()
}

func (s *CopyProcessResponseBodyResult) SetBizType(v string) *CopyProcessResponseBodyResult {
	s.BizType = &v
	return s
}

func (s *CopyProcessResponseBodyResult) SetName(v string) *CopyProcessResponseBodyResult {
	s.Name = &v
	return s
}

func (s *CopyProcessResponseBodyResult) SetProcessCode(v string) *CopyProcessResponseBodyResult {
	s.ProcessCode = &v
	return s
}

type CopyProcessResponse struct {
	Headers map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *CopyProcessResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CopyProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s CopyProcessResponse) GoString() string {
	return s.String()
}

func (s *CopyProcessResponse) SetHeaders(v map[string]*string) *CopyProcessResponse {
	s.Headers = v
	return s
}

func (s *CopyProcessResponse) SetBody(v *CopyProcessResponseBody) *CopyProcessResponse {
	s.Body = v
	return s
}

type CreateIntegratedTaskHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s CreateIntegratedTaskHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateIntegratedTaskHeaders) GoString() string {
	return s.String()
}

func (s *CreateIntegratedTaskHeaders) SetCommonHeaders(v map[string]*string) *CreateIntegratedTaskHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CreateIntegratedTaskHeaders) SetXAcsDingtalkAccessToken(v string) *CreateIntegratedTaskHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type CreateIntegratedTaskRequest struct {
	// 待办组ID，调用方提供自定义唯一分组标识
	ActivityId *string `json:"activityId,omitempty" xml:"activityId,omitempty"`
	// OA审批实例ID，通过创建实例接口获取。
	//
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 任务列表
	Tasks []*CreateIntegratedTaskRequestTasks `json:"tasks,omitempty" xml:"tasks,omitempty" type:"Repeated"`
}

func (s CreateIntegratedTaskRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateIntegratedTaskRequest) GoString() string {
	return s.String()
}

func (s *CreateIntegratedTaskRequest) SetActivityId(v string) *CreateIntegratedTaskRequest {
	s.ActivityId = &v
	return s
}

func (s *CreateIntegratedTaskRequest) SetProcessInstanceId(v string) *CreateIntegratedTaskRequest {
	s.ProcessInstanceId = &v
	return s
}

func (s *CreateIntegratedTaskRequest) SetTasks(v []*CreateIntegratedTaskRequestTasks) *CreateIntegratedTaskRequest {
	s.Tasks = v
	return s
}

type CreateIntegratedTaskRequestTasks struct {
	// 待办事项跳转URL。
	// 创建审批实例接口里的url，实现的是钉钉审批应用里的审批单跳转。这个接口里的url，实现的是钉钉待办页面，对应的待办卡片的跳转。两种跳转场景不同。需要注意的是，钉钉的待办页面，也同时支持移动端和PC端，所以这个接口里传的url参数，它所对应的页面也要适配好移动端和PC端。
	//
	Url *string `json:"url,omitempty" xml:"url,omitempty"`
	// OA审批任务执行人的用户ID
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s CreateIntegratedTaskRequestTasks) String() string {
	return tea.Prettify(s)
}

func (s CreateIntegratedTaskRequestTasks) GoString() string {
	return s.String()
}

func (s *CreateIntegratedTaskRequestTasks) SetUrl(v string) *CreateIntegratedTaskRequestTasks {
	s.Url = &v
	return s
}

func (s *CreateIntegratedTaskRequestTasks) SetUserId(v string) *CreateIntegratedTaskRequestTasks {
	s.UserId = &v
	return s
}

type CreateIntegratedTaskResponseBody struct {
	Result []*CreateIntegratedTaskResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Repeated"`
	// 是否创建成功
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s CreateIntegratedTaskResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateIntegratedTaskResponseBody) GoString() string {
	return s.String()
}

func (s *CreateIntegratedTaskResponseBody) SetResult(v []*CreateIntegratedTaskResponseBodyResult) *CreateIntegratedTaskResponseBody {
	s.Result = v
	return s
}

func (s *CreateIntegratedTaskResponseBody) SetSuccess(v bool) *CreateIntegratedTaskResponseBody {
	s.Success = &v
	return s
}

type CreateIntegratedTaskResponseBodyResult struct {
	// OA审批任务ID
	TaskId *int64 `json:"taskId,omitempty" xml:"taskId,omitempty"`
	// OA审批任务执行人用户ID
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s CreateIntegratedTaskResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s CreateIntegratedTaskResponseBodyResult) GoString() string {
	return s.String()
}

func (s *CreateIntegratedTaskResponseBodyResult) SetTaskId(v int64) *CreateIntegratedTaskResponseBodyResult {
	s.TaskId = &v
	return s
}

func (s *CreateIntegratedTaskResponseBodyResult) SetUserId(v string) *CreateIntegratedTaskResponseBodyResult {
	s.UserId = &v
	return s
}

type CreateIntegratedTaskResponse struct {
	Headers map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *CreateIntegratedTaskResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateIntegratedTaskResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateIntegratedTaskResponse) GoString() string {
	return s.String()
}

func (s *CreateIntegratedTaskResponse) SetHeaders(v map[string]*string) *CreateIntegratedTaskResponse {
	s.Headers = v
	return s
}

func (s *CreateIntegratedTaskResponse) SetBody(v *CreateIntegratedTaskResponseBody) *CreateIntegratedTaskResponse {
	s.Body = v
	return s
}

type DeleteProcessHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s DeleteProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteProcessHeaders) GoString() string {
	return s.String()
}

func (s *DeleteProcessHeaders) SetCommonHeaders(v map[string]*string) *DeleteProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteProcessHeaders) SetXAcsDingtalkAccessToken(v string) *DeleteProcessHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type DeleteProcessRequest struct {
	// 是否清理正在运行的任务
	CleanRunningTask *bool `json:"cleanRunningTask,omitempty" xml:"cleanRunningTask,omitempty"`
	// 模板code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s DeleteProcessRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteProcessRequest) GoString() string {
	return s.String()
}

func (s *DeleteProcessRequest) SetCleanRunningTask(v bool) *DeleteProcessRequest {
	s.CleanRunningTask = &v
	return s
}

func (s *DeleteProcessRequest) SetProcessCode(v string) *DeleteProcessRequest {
	s.ProcessCode = &v
	return s
}

type DeleteProcessResponseBody struct {
	Result *DeleteProcessResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s DeleteProcessResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteProcessResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteProcessResponseBody) SetResult(v *DeleteProcessResponseBodyResult) *DeleteProcessResponseBody {
	s.Result = v
	return s
}

type DeleteProcessResponseBodyResult struct {
	// 模板code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s DeleteProcessResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s DeleteProcessResponseBodyResult) GoString() string {
	return s.String()
}

func (s *DeleteProcessResponseBodyResult) SetProcessCode(v string) *DeleteProcessResponseBodyResult {
	s.ProcessCode = &v
	return s
}

type DeleteProcessResponse struct {
	Headers map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DeleteProcessResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteProcessResponse) GoString() string {
	return s.String()
}

func (s *DeleteProcessResponse) SetHeaders(v map[string]*string) *DeleteProcessResponse {
	s.Headers = v
	return s
}

func (s *DeleteProcessResponse) SetBody(v *DeleteProcessResponseBody) *DeleteProcessResponse {
	s.Body = v
	return s
}

type ExecuteProcessInstanceHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s ExecuteProcessInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s ExecuteProcessInstanceHeaders) GoString() string {
	return s.String()
}

func (s *ExecuteProcessInstanceHeaders) SetCommonHeaders(v map[string]*string) *ExecuteProcessInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ExecuteProcessInstanceHeaders) SetXAcsDingtalkAccessToken(v string) *ExecuteProcessInstanceHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type ExecuteProcessInstanceRequest struct {
	// 操作人userid，可通过调用获取审批实例详情接口获取。
	ActionerUserId *string `json:"actionerUserId,omitempty" xml:"actionerUserId,omitempty"`
	// 文件。
	File *ExecuteProcessInstanceRequestFile `json:"file,omitempty" xml:"file,omitempty" type:"Struct"`
	// 审批实例ID，可通过调用获取审批实例ID列表接口获取。
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 审批意见，可为空。
	Remark *string `json:"remark,omitempty" xml:"remark,omitempty"`
	// 审批操作，取值。
	//
	// agree：同意
	//
	// refuse：拒绝
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 任务节点id，可通过调用获取审批实例详情接口获取。
	TaskId *int64 `json:"taskId,omitempty" xml:"taskId,omitempty"`
}

func (s ExecuteProcessInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s ExecuteProcessInstanceRequest) GoString() string {
	return s.String()
}

func (s *ExecuteProcessInstanceRequest) SetActionerUserId(v string) *ExecuteProcessInstanceRequest {
	s.ActionerUserId = &v
	return s
}

func (s *ExecuteProcessInstanceRequest) SetFile(v *ExecuteProcessInstanceRequestFile) *ExecuteProcessInstanceRequest {
	s.File = v
	return s
}

func (s *ExecuteProcessInstanceRequest) SetProcessInstanceId(v string) *ExecuteProcessInstanceRequest {
	s.ProcessInstanceId = &v
	return s
}

func (s *ExecuteProcessInstanceRequest) SetRemark(v string) *ExecuteProcessInstanceRequest {
	s.Remark = &v
	return s
}

func (s *ExecuteProcessInstanceRequest) SetResult(v string) *ExecuteProcessInstanceRequest {
	s.Result = &v
	return s
}

func (s *ExecuteProcessInstanceRequest) SetTaskId(v int64) *ExecuteProcessInstanceRequest {
	s.TaskId = &v
	return s
}

type ExecuteProcessInstanceRequestFile struct {
	// 附件列表。
	Attachments []*ExecuteProcessInstanceRequestFileAttachments `json:"attachments,omitempty" xml:"attachments,omitempty" type:"Repeated"`
	// 图片URL地址。
	Photos []*string `json:"photos,omitempty" xml:"photos,omitempty" type:"Repeated"`
}

func (s ExecuteProcessInstanceRequestFile) String() string {
	return tea.Prettify(s)
}

func (s ExecuteProcessInstanceRequestFile) GoString() string {
	return s.String()
}

func (s *ExecuteProcessInstanceRequestFile) SetAttachments(v []*ExecuteProcessInstanceRequestFileAttachments) *ExecuteProcessInstanceRequestFile {
	s.Attachments = v
	return s
}

func (s *ExecuteProcessInstanceRequestFile) SetPhotos(v []*string) *ExecuteProcessInstanceRequestFile {
	s.Photos = v
	return s
}

type ExecuteProcessInstanceRequestFileAttachments struct {
	// 文件ID。
	FileId *string `json:"fileId,omitempty" xml:"fileId,omitempty"`
	// 文件名称。
	FileName *string `json:"fileName,omitempty" xml:"fileName,omitempty"`
	// 文件大小。
	FileSize *string `json:"fileSize,omitempty" xml:"fileSize,omitempty"`
	// 文件类型。
	FileType *string `json:"fileType,omitempty" xml:"fileType,omitempty"`
	// 钉盘空间ID。
	SpaceId *string `json:"spaceId,omitempty" xml:"spaceId,omitempty"`
}

func (s ExecuteProcessInstanceRequestFileAttachments) String() string {
	return tea.Prettify(s)
}

func (s ExecuteProcessInstanceRequestFileAttachments) GoString() string {
	return s.String()
}

func (s *ExecuteProcessInstanceRequestFileAttachments) SetFileId(v string) *ExecuteProcessInstanceRequestFileAttachments {
	s.FileId = &v
	return s
}

func (s *ExecuteProcessInstanceRequestFileAttachments) SetFileName(v string) *ExecuteProcessInstanceRequestFileAttachments {
	s.FileName = &v
	return s
}

func (s *ExecuteProcessInstanceRequestFileAttachments) SetFileSize(v string) *ExecuteProcessInstanceRequestFileAttachments {
	s.FileSize = &v
	return s
}

func (s *ExecuteProcessInstanceRequestFileAttachments) SetFileType(v string) *ExecuteProcessInstanceRequestFileAttachments {
	s.FileType = &v
	return s
}

func (s *ExecuteProcessInstanceRequestFileAttachments) SetSpaceId(v string) *ExecuteProcessInstanceRequestFileAttachments {
	s.SpaceId = &v
	return s
}

type ExecuteProcessInstanceResponseBody struct {
	// 同意或拒绝结果。
	Result *bool `json:"result,omitempty" xml:"result,omitempty"`
	// 接口调用是否成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s ExecuteProcessInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ExecuteProcessInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *ExecuteProcessInstanceResponseBody) SetResult(v bool) *ExecuteProcessInstanceResponseBody {
	s.Result = &v
	return s
}

func (s *ExecuteProcessInstanceResponseBody) SetSuccess(v bool) *ExecuteProcessInstanceResponseBody {
	s.Success = &v
	return s
}

type ExecuteProcessInstanceResponse struct {
	Headers map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ExecuteProcessInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ExecuteProcessInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s ExecuteProcessInstanceResponse) GoString() string {
	return s.String()
}

func (s *ExecuteProcessInstanceResponse) SetHeaders(v map[string]*string) *ExecuteProcessInstanceResponse {
	s.Headers = v
	return s
}

func (s *ExecuteProcessInstanceResponse) SetBody(v *ExecuteProcessInstanceResponseBody) *ExecuteProcessInstanceResponse {
	s.Body = v
	return s
}

type FormCreateHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s FormCreateHeaders) String() string {
	return tea.Prettify(s)
}

func (s FormCreateHeaders) GoString() string {
	return s.String()
}

func (s *FormCreateHeaders) SetCommonHeaders(v map[string]*string) *FormCreateHeaders {
	s.CommonHeaders = v
	return s
}

func (s *FormCreateHeaders) SetXAcsDingtalkAccessToken(v string) *FormCreateHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type FormCreateRequest struct {
	// 表单模板描述
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// 表单控件列表
	FormComponents []*FormComponent `json:"formComponents,omitempty" xml:"formComponents,omitempty" type:"Repeated"`
	// 表单模板名称
	Name        *string `json:"name,omitempty" xml:"name,omitempty"`
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 模板配置信息
	TemplateConfig *FormCreateRequestTemplateConfig `json:"templateConfig,omitempty" xml:"templateConfig,omitempty" type:"Struct"`
}

func (s FormCreateRequest) String() string {
	return tea.Prettify(s)
}

func (s FormCreateRequest) GoString() string {
	return s.String()
}

func (s *FormCreateRequest) SetDescription(v string) *FormCreateRequest {
	s.Description = &v
	return s
}

func (s *FormCreateRequest) SetFormComponents(v []*FormComponent) *FormCreateRequest {
	s.FormComponents = v
	return s
}

func (s *FormCreateRequest) SetName(v string) *FormCreateRequest {
	s.Name = &v
	return s
}

func (s *FormCreateRequest) SetProcessCode(v string) *FormCreateRequest {
	s.ProcessCode = &v
	return s
}

func (s *FormCreateRequest) SetTemplateConfig(v *FormCreateRequestTemplateConfig) *FormCreateRequest {
	s.TemplateConfig = v
	return s
}

type FormCreateRequestTemplateConfig struct {
	// 更新后模板目录id
	DirId *string `json:"dirId,omitempty" xml:"dirId,omitempty"`
	// 禁用模板删除按钮
	DisableDeleteProcess *bool `json:"disableDeleteProcess,omitempty" xml:"disableDeleteProcess,omitempty"`
	// 禁用表单编辑
	DisableFormEdit *bool `json:"disableFormEdit,omitempty" xml:"disableFormEdit,omitempty"`
	// 首页工作台是否可见
	DisableHomepage *bool `json:"disableHomepage,omitempty" xml:"disableHomepage,omitempty"`
	// 禁用再次提交
	DisableResubmit *bool `json:"disableResubmit,omitempty" xml:"disableResubmit,omitempty"`
	// 禁用停止按钮
	DisableStopProcessButton *bool `json:"disableStopProcessButton,omitempty" xml:"disableStopProcessButton,omitempty"`
	// 审批场景内隐藏模板
	Hidden *bool `json:"hidden,omitempty" xml:"hidden,omitempty"`
	// 源模板目录id
	OriginDirId *string `json:"originDirId,omitempty" xml:"originDirId,omitempty"`
}

func (s FormCreateRequestTemplateConfig) String() string {
	return tea.Prettify(s)
}

func (s FormCreateRequestTemplateConfig) GoString() string {
	return s.String()
}

func (s *FormCreateRequestTemplateConfig) SetDirId(v string) *FormCreateRequestTemplateConfig {
	s.DirId = &v
	return s
}

func (s *FormCreateRequestTemplateConfig) SetDisableDeleteProcess(v bool) *FormCreateRequestTemplateConfig {
	s.DisableDeleteProcess = &v
	return s
}

func (s *FormCreateRequestTemplateConfig) SetDisableFormEdit(v bool) *FormCreateRequestTemplateConfig {
	s.DisableFormEdit = &v
	return s
}

func (s *FormCreateRequestTemplateConfig) SetDisableHomepage(v bool) *FormCreateRequestTemplateConfig {
	s.DisableHomepage = &v
	return s
}

func (s *FormCreateRequestTemplateConfig) SetDisableResubmit(v bool) *FormCreateRequestTemplateConfig {
	s.DisableResubmit = &v
	return s
}

func (s *FormCreateRequestTemplateConfig) SetDisableStopProcessButton(v bool) *FormCreateRequestTemplateConfig {
	s.DisableStopProcessButton = &v
	return s
}

func (s *FormCreateRequestTemplateConfig) SetHidden(v bool) *FormCreateRequestTemplateConfig {
	s.Hidden = &v
	return s
}

func (s *FormCreateRequestTemplateConfig) SetOriginDirId(v string) *FormCreateRequestTemplateConfig {
	s.OriginDirId = &v
	return s
}

type FormCreateResponseBody struct {
	// 表单模板信息
	Result *FormCreateResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s FormCreateResponseBody) String() string {
	return tea.Prettify(s)
}

func (s FormCreateResponseBody) GoString() string {
	return s.String()
}

func (s *FormCreateResponseBody) SetResult(v *FormCreateResponseBodyResult) *FormCreateResponseBody {
	s.Result = v
	return s
}

type FormCreateResponseBodyResult struct {
	// 保存或更新的表单code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s FormCreateResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s FormCreateResponseBodyResult) GoString() string {
	return s.String()
}

func (s *FormCreateResponseBodyResult) SetProcessCode(v string) *FormCreateResponseBodyResult {
	s.ProcessCode = &v
	return s
}

type FormCreateResponse struct {
	Headers map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *FormCreateResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s FormCreateResponse) String() string {
	return tea.Prettify(s)
}

func (s FormCreateResponse) GoString() string {
	return s.String()
}

func (s *FormCreateResponse) SetHeaders(v map[string]*string) *FormCreateResponse {
	s.Headers = v
	return s
}

func (s *FormCreateResponse) SetBody(v *FormCreateResponseBody) *FormCreateResponse {
	s.Body = v
	return s
}

type GetAttachmentSpaceHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GetAttachmentSpaceHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetAttachmentSpaceHeaders) GoString() string {
	return s.String()
}

func (s *GetAttachmentSpaceHeaders) SetCommonHeaders(v map[string]*string) *GetAttachmentSpaceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetAttachmentSpaceHeaders) SetXAcsDingtalkAccessToken(v string) *GetAttachmentSpaceHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GetAttachmentSpaceRequest struct {
	// 应用的agentid。
	AgentId *int64 `json:"agentId,omitempty" xml:"agentId,omitempty"`
	// 用户的userid
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s GetAttachmentSpaceRequest) String() string {
	return tea.Prettify(s)
}

func (s GetAttachmentSpaceRequest) GoString() string {
	return s.String()
}

func (s *GetAttachmentSpaceRequest) SetAgentId(v int64) *GetAttachmentSpaceRequest {
	s.AgentId = &v
	return s
}

func (s *GetAttachmentSpaceRequest) SetUserId(v string) *GetAttachmentSpaceRequest {
	s.UserId = &v
	return s
}

type GetAttachmentSpaceResponseBody struct {
	// 返回结果。
	Result *GetAttachmentSpaceResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
	// 接口调用是否成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s GetAttachmentSpaceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAttachmentSpaceResponseBody) GoString() string {
	return s.String()
}

func (s *GetAttachmentSpaceResponseBody) SetResult(v *GetAttachmentSpaceResponseBodyResult) *GetAttachmentSpaceResponseBody {
	s.Result = v
	return s
}

func (s *GetAttachmentSpaceResponseBody) SetSuccess(v bool) *GetAttachmentSpaceResponseBody {
	s.Success = &v
	return s
}

type GetAttachmentSpaceResponseBodyResult struct {
	// 钉盘空间ID。
	SpaceId *int64 `json:"spaceId,omitempty" xml:"spaceId,omitempty"`
}

func (s GetAttachmentSpaceResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s GetAttachmentSpaceResponseBodyResult) GoString() string {
	return s.String()
}

func (s *GetAttachmentSpaceResponseBodyResult) SetSpaceId(v int64) *GetAttachmentSpaceResponseBodyResult {
	s.SpaceId = &v
	return s
}

type GetAttachmentSpaceResponse struct {
	Headers map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetAttachmentSpaceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetAttachmentSpaceResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAttachmentSpaceResponse) GoString() string {
	return s.String()
}

func (s *GetAttachmentSpaceResponse) SetHeaders(v map[string]*string) *GetAttachmentSpaceResponse {
	s.Headers = v
	return s
}

func (s *GetAttachmentSpaceResponse) SetBody(v *GetAttachmentSpaceResponseBody) *GetAttachmentSpaceResponse {
	s.Body = v
	return s
}

type GetConditionFormComponentHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GetConditionFormComponentHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetConditionFormComponentHeaders) GoString() string {
	return s.String()
}

func (s *GetConditionFormComponentHeaders) SetCommonHeaders(v map[string]*string) *GetConditionFormComponentHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetConditionFormComponentHeaders) SetXAcsDingtalkAccessToken(v string) *GetConditionFormComponentHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GetConditionFormComponentRequest struct {
	// 应用ID (三方应用为AppID)，可在开发者管理后台后台的应用详情页面获取。
	AgentId *int64 `json:"agentId,omitempty" xml:"agentId,omitempty"`
	// 审批模板ID。
	//
	// processCode需要OA管理后台，在编辑审批表单的URL中获取。
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s GetConditionFormComponentRequest) String() string {
	return tea.Prettify(s)
}

func (s GetConditionFormComponentRequest) GoString() string {
	return s.String()
}

func (s *GetConditionFormComponentRequest) SetAgentId(v int64) *GetConditionFormComponentRequest {
	s.AgentId = &v
	return s
}

func (s *GetConditionFormComponentRequest) SetProcessCode(v string) *GetConditionFormComponentRequest {
	s.ProcessCode = &v
	return s
}

type GetConditionFormComponentResponseBody struct {
	// 返回结果。
	Result []*GetConditionFormComponentResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Repeated"`
}

func (s GetConditionFormComponentResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetConditionFormComponentResponseBody) GoString() string {
	return s.String()
}

func (s *GetConditionFormComponentResponseBody) SetResult(v []*GetConditionFormComponentResponseBodyResult) *GetConditionFormComponentResponseBody {
	s.Result = v
	return s
}

type GetConditionFormComponentResponseBodyResult struct {
	// 表单ID。
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 表单名称。
	Label *string `json:"label,omitempty" xml:"label,omitempty"`
}

func (s GetConditionFormComponentResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s GetConditionFormComponentResponseBodyResult) GoString() string {
	return s.String()
}

func (s *GetConditionFormComponentResponseBodyResult) SetId(v string) *GetConditionFormComponentResponseBodyResult {
	s.Id = &v
	return s
}

func (s *GetConditionFormComponentResponseBodyResult) SetLabel(v string) *GetConditionFormComponentResponseBodyResult {
	s.Label = &v
	return s
}

type GetConditionFormComponentResponse struct {
	Headers map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetConditionFormComponentResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetConditionFormComponentResponse) String() string {
	return tea.Prettify(s)
}

func (s GetConditionFormComponentResponse) GoString() string {
	return s.String()
}

func (s *GetConditionFormComponentResponse) SetHeaders(v map[string]*string) *GetConditionFormComponentResponse {
	s.Headers = v
	return s
}

func (s *GetConditionFormComponentResponse) SetBody(v *GetConditionFormComponentResponseBody) *GetConditionFormComponentResponse {
	s.Body = v
	return s
}

type GetCrmProcCodesHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GetCrmProcCodesHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetCrmProcCodesHeaders) GoString() string {
	return s.String()
}

func (s *GetCrmProcCodesHeaders) SetCommonHeaders(v map[string]*string) *GetCrmProcCodesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetCrmProcCodesHeaders) SetXAcsDingtalkAccessToken(v string) *GetCrmProcCodesHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GetCrmProcCodesResponseBody struct {
	// 模板code列表。
	Result []*string `json:"result,omitempty" xml:"result,omitempty" type:"Repeated"`
}

func (s GetCrmProcCodesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetCrmProcCodesResponseBody) GoString() string {
	return s.String()
}

func (s *GetCrmProcCodesResponseBody) SetResult(v []*string) *GetCrmProcCodesResponseBody {
	s.Result = v
	return s
}

type GetCrmProcCodesResponse struct {
	Headers map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetCrmProcCodesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetCrmProcCodesResponse) String() string {
	return tea.Prettify(s)
}

func (s GetCrmProcCodesResponse) GoString() string {
	return s.String()
}

func (s *GetCrmProcCodesResponse) SetHeaders(v map[string]*string) *GetCrmProcCodesResponse {
	s.Headers = v
	return s
}

func (s *GetCrmProcCodesResponse) SetBody(v *GetCrmProcCodesResponseBody) *GetCrmProcCodesResponse {
	s.Body = v
	return s
}

type GetManageProcessByStaffIdHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GetManageProcessByStaffIdHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetManageProcessByStaffIdHeaders) GoString() string {
	return s.String()
}

func (s *GetManageProcessByStaffIdHeaders) SetCommonHeaders(v map[string]*string) *GetManageProcessByStaffIdHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetManageProcessByStaffIdHeaders) SetXAcsDingtalkAccessToken(v string) *GetManageProcessByStaffIdHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GetManageProcessByStaffIdRequest struct {
	// 用户的userid。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s GetManageProcessByStaffIdRequest) String() string {
	return tea.Prettify(s)
}

func (s GetManageProcessByStaffIdRequest) GoString() string {
	return s.String()
}

func (s *GetManageProcessByStaffIdRequest) SetUserId(v string) *GetManageProcessByStaffIdRequest {
	s.UserId = &v
	return s
}

type GetManageProcessByStaffIdResponseBody struct {
	// 返回结果列表。
	Result []*GetManageProcessByStaffIdResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Repeated"`
	// 接口调用是否成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s GetManageProcessByStaffIdResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetManageProcessByStaffIdResponseBody) GoString() string {
	return s.String()
}

func (s *GetManageProcessByStaffIdResponseBody) SetResult(v []*GetManageProcessByStaffIdResponseBodyResult) *GetManageProcessByStaffIdResponseBody {
	s.Result = v
	return s
}

func (s *GetManageProcessByStaffIdResponseBody) SetSuccess(v bool) *GetManageProcessByStaffIdResponseBody {
	s.Success = &v
	return s
}

type GetManageProcessByStaffIdResponseBodyResult struct {
	// 关联考勤类型，取值。
	//
	// 0：无
	// 1：补卡申请
	// 2：请假
	AttendanceType *int32 `json:"attendanceType,omitempty" xml:"attendanceType,omitempty"`
	// 模版名称。
	FlowTitle *string `json:"flowTitle,omitempty" xml:"flowTitle,omitempty"`
	// 修改时间。
	GmtModified *string `json:"gmtModified,omitempty" xml:"gmtModified,omitempty"`
	// 模板图标名。
	IconName *string `json:"iconName,omitempty" xml:"iconName,omitempty"`
	// 图标URL地址。
	IconUrl *string `json:"iconUrl,omitempty" xml:"iconUrl,omitempty"`
	// 是否新模版。
	NewProcess *bool `json:"newProcess,omitempty" xml:"newProcess,omitempty"`
	// 模版code。
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s GetManageProcessByStaffIdResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s GetManageProcessByStaffIdResponseBodyResult) GoString() string {
	return s.String()
}

func (s *GetManageProcessByStaffIdResponseBodyResult) SetAttendanceType(v int32) *GetManageProcessByStaffIdResponseBodyResult {
	s.AttendanceType = &v
	return s
}

func (s *GetManageProcessByStaffIdResponseBodyResult) SetFlowTitle(v string) *GetManageProcessByStaffIdResponseBodyResult {
	s.FlowTitle = &v
	return s
}

func (s *GetManageProcessByStaffIdResponseBodyResult) SetGmtModified(v string) *GetManageProcessByStaffIdResponseBodyResult {
	s.GmtModified = &v
	return s
}

func (s *GetManageProcessByStaffIdResponseBodyResult) SetIconName(v string) *GetManageProcessByStaffIdResponseBodyResult {
	s.IconName = &v
	return s
}

func (s *GetManageProcessByStaffIdResponseBodyResult) SetIconUrl(v string) *GetManageProcessByStaffIdResponseBodyResult {
	s.IconUrl = &v
	return s
}

func (s *GetManageProcessByStaffIdResponseBodyResult) SetNewProcess(v bool) *GetManageProcessByStaffIdResponseBodyResult {
	s.NewProcess = &v
	return s
}

func (s *GetManageProcessByStaffIdResponseBodyResult) SetProcessCode(v string) *GetManageProcessByStaffIdResponseBodyResult {
	s.ProcessCode = &v
	return s
}

type GetManageProcessByStaffIdResponse struct {
	Headers map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetManageProcessByStaffIdResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetManageProcessByStaffIdResponse) String() string {
	return tea.Prettify(s)
}

func (s GetManageProcessByStaffIdResponse) GoString() string {
	return s.String()
}

func (s *GetManageProcessByStaffIdResponse) SetHeaders(v map[string]*string) *GetManageProcessByStaffIdResponse {
	s.Headers = v
	return s
}

func (s *GetManageProcessByStaffIdResponse) SetBody(v *GetManageProcessByStaffIdResponseBody) *GetManageProcessByStaffIdResponse {
	s.Body = v
	return s
}

type GetProcessCodeByNameHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GetProcessCodeByNameHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetProcessCodeByNameHeaders) GoString() string {
	return s.String()
}

func (s *GetProcessCodeByNameHeaders) SetCommonHeaders(v map[string]*string) *GetProcessCodeByNameHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetProcessCodeByNameHeaders) SetXAcsDingtalkAccessToken(v string) *GetProcessCodeByNameHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GetProcessCodeByNameRequest struct {
	// 模板名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
}

func (s GetProcessCodeByNameRequest) String() string {
	return tea.Prettify(s)
}

func (s GetProcessCodeByNameRequest) GoString() string {
	return s.String()
}

func (s *GetProcessCodeByNameRequest) SetName(v string) *GetProcessCodeByNameRequest {
	s.Name = &v
	return s
}

type GetProcessCodeByNameResponseBody struct {
	// 表单模板信息
	Result *GetProcessCodeByNameResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s GetProcessCodeByNameResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetProcessCodeByNameResponseBody) GoString() string {
	return s.String()
}

func (s *GetProcessCodeByNameResponseBody) SetResult(v *GetProcessCodeByNameResponseBodyResult) *GetProcessCodeByNameResponseBody {
	s.Result = v
	return s
}

type GetProcessCodeByNameResponseBodyResult struct {
	// 保存或更新的表单code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s GetProcessCodeByNameResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s GetProcessCodeByNameResponseBodyResult) GoString() string {
	return s.String()
}

func (s *GetProcessCodeByNameResponseBodyResult) SetProcessCode(v string) *GetProcessCodeByNameResponseBodyResult {
	s.ProcessCode = &v
	return s
}

type GetProcessCodeByNameResponse struct {
	Headers map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetProcessCodeByNameResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetProcessCodeByNameResponse) String() string {
	return tea.Prettify(s)
}

func (s GetProcessCodeByNameResponse) GoString() string {
	return s.String()
}

func (s *GetProcessCodeByNameResponse) SetHeaders(v map[string]*string) *GetProcessCodeByNameResponse {
	s.Headers = v
	return s
}

func (s *GetProcessCodeByNameResponse) SetBody(v *GetProcessCodeByNameResponseBody) *GetProcessCodeByNameResponse {
	s.Body = v
	return s
}

type GetProcessConfigHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GetProcessConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigHeaders) GoString() string {
	return s.String()
}

func (s *GetProcessConfigHeaders) SetCommonHeaders(v map[string]*string) *GetProcessConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetProcessConfigHeaders) SetXAcsDingtalkAccessToken(v string) *GetProcessConfigHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GetProcessConfigRequest struct {
	// 模板code
	ProcCode *string `json:"procCode,omitempty" xml:"procCode,omitempty"`
}

func (s GetProcessConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigRequest) GoString() string {
	return s.String()
}

func (s *GetProcessConfigRequest) SetProcCode(v string) *GetProcessConfigRequest {
	s.ProcCode = &v
	return s
}

type GetProcessConfigResponseBody struct {
	Result *GetProcessConfigResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s GetProcessConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigResponseBody) GoString() string {
	return s.String()
}

func (s *GetProcessConfigResponseBody) SetResult(v *GetProcessConfigResponseBodyResult) *GetProcessConfigResponseBody {
	s.Result = v
	return s
}

type GetProcessConfigResponseBodyResult struct {
	// 自定义摘要信息
	AbstractGenRule []*string `json:"abstractGenRule,omitempty" xml:"abstractGenRule,omitempty" type:"Repeated"`
	// 表单节点权限
	ActivityAuth *string `json:"activityAuth,omitempty" xml:"activityAuth,omitempty"`
	// 是否允许撤销
	AllowRevoke *bool `json:"allowRevoke,omitempty" xml:"allowRevoke,omitempty"`
	// 是否允许加签
	AppendEnable *bool `json:"appendEnable,omitempty" xml:"appendEnable,omitempty"`
	// 如果审批人和发起人是同一个人，则去重
	AutoExecuteOriginatorTasks *bool `json:"autoExecuteOriginatorTasks,omitempty" xml:"autoExecuteOriginatorTasks,omitempty"`
	// 流程表单业务标识
	BizCategoryId *string `json:"bizCategoryId,omitempty" xml:"bizCategoryId,omitempty"`
	// 纯表单业务标识
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 评论配置
	CommentConf *GetProcessConfigResponseBodyResultCommentConf `json:"commentConf,omitempty" xml:"commentConf,omitempty" type:"Struct"`
	// 审批人自动去重
	DuplicateRemoval *string `json:"duplicateRemoval,omitempty" xml:"duplicateRemoval,omitempty"`
	// 表单配置
	FormSchema *string `json:"formSchema,omitempty" xml:"formSchema,omitempty"`
	// 手写签名配置
	HandSignConf *GetProcessConfigResponseBodyResultHandSignConf `json:"handSignConf,omitempty" xml:"handSignConf,omitempty" type:"Struct"`
	// 表单管理员
	Managers []*string `json:"managers,omitempty" xml:"managers,omitempty" type:"Repeated"`
	// 表单名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 是否流程表单
	ProcessAppType *bool `json:"processAppType,omitempty" xml:"processAppType,omitempty"`
	// 流程配置
	ProcessConfig *string `json:"processConfig,omitempty" xml:"processConfig,omitempty"`
	// 是否静态流程
	StaticProc *bool `json:"staticProc,omitempty" xml:"staticProc,omitempty"`
	// 代提交配置
	SubstituteSubmitConf *GetProcessConfigResponseBodyResultSubstituteSubmitConf `json:"substituteSubmitConf,omitempty" xml:"substituteSubmitConf,omitempty" type:"Struct"`
	// 自定义标题规则
	TitleGenRule *GetProcessConfigResponseBodyResultTitleGenRule `json:"titleGenRule,omitempty" xml:"titleGenRule,omitempty" type:"Struct"`
	// 模板可见性
	Visibility []*GetProcessConfigResponseBodyResultVisibility `json:"visibility,omitempty" xml:"visibility,omitempty" type:"Repeated"`
}

func (s GetProcessConfigResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigResponseBodyResult) GoString() string {
	return s.String()
}

func (s *GetProcessConfigResponseBodyResult) SetAbstractGenRule(v []*string) *GetProcessConfigResponseBodyResult {
	s.AbstractGenRule = v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetActivityAuth(v string) *GetProcessConfigResponseBodyResult {
	s.ActivityAuth = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetAllowRevoke(v bool) *GetProcessConfigResponseBodyResult {
	s.AllowRevoke = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetAppendEnable(v bool) *GetProcessConfigResponseBodyResult {
	s.AppendEnable = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetAutoExecuteOriginatorTasks(v bool) *GetProcessConfigResponseBodyResult {
	s.AutoExecuteOriginatorTasks = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetBizCategoryId(v string) *GetProcessConfigResponseBodyResult {
	s.BizCategoryId = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetBizType(v string) *GetProcessConfigResponseBodyResult {
	s.BizType = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetCommentConf(v *GetProcessConfigResponseBodyResultCommentConf) *GetProcessConfigResponseBodyResult {
	s.CommentConf = v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetDuplicateRemoval(v string) *GetProcessConfigResponseBodyResult {
	s.DuplicateRemoval = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetFormSchema(v string) *GetProcessConfigResponseBodyResult {
	s.FormSchema = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetHandSignConf(v *GetProcessConfigResponseBodyResultHandSignConf) *GetProcessConfigResponseBodyResult {
	s.HandSignConf = v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetManagers(v []*string) *GetProcessConfigResponseBodyResult {
	s.Managers = v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetName(v string) *GetProcessConfigResponseBodyResult {
	s.Name = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetProcessAppType(v bool) *GetProcessConfigResponseBodyResult {
	s.ProcessAppType = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetProcessConfig(v string) *GetProcessConfigResponseBodyResult {
	s.ProcessConfig = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetStaticProc(v bool) *GetProcessConfigResponseBodyResult {
	s.StaticProc = &v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetSubstituteSubmitConf(v *GetProcessConfigResponseBodyResultSubstituteSubmitConf) *GetProcessConfigResponseBodyResult {
	s.SubstituteSubmitConf = v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetTitleGenRule(v *GetProcessConfigResponseBodyResultTitleGenRule) *GetProcessConfigResponseBodyResult {
	s.TitleGenRule = v
	return s
}

func (s *GetProcessConfigResponseBodyResult) SetVisibility(v []*GetProcessConfigResponseBodyResultVisibility) *GetProcessConfigResponseBodyResult {
	s.Visibility = v
	return s
}

type GetProcessConfigResponseBodyResultCommentConf struct {
	// 提示内容
	CommentDescription *string `json:"commentDescription,omitempty" xml:"commentDescription,omitempty"`
	// 评论对发起人不可见
	CommentHiddenForProposer *bool `json:"commentHiddenForProposer,omitempty" xml:"commentHiddenForProposer,omitempty"`
	// 评论必填
	CommentRequired *bool `json:"commentRequired,omitempty" xml:"commentRequired,omitempty"`
}

func (s GetProcessConfigResponseBodyResultCommentConf) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigResponseBodyResultCommentConf) GoString() string {
	return s.String()
}

func (s *GetProcessConfigResponseBodyResultCommentConf) SetCommentDescription(v string) *GetProcessConfigResponseBodyResultCommentConf {
	s.CommentDescription = &v
	return s
}

func (s *GetProcessConfigResponseBodyResultCommentConf) SetCommentHiddenForProposer(v bool) *GetProcessConfigResponseBodyResultCommentConf {
	s.CommentHiddenForProposer = &v
	return s
}

func (s *GetProcessConfigResponseBodyResultCommentConf) SetCommentRequired(v bool) *GetProcessConfigResponseBodyResultCommentConf {
	s.CommentRequired = &v
	return s
}

type GetProcessConfigResponseBodyResultHandSignConf struct {
	// 开启手写签名
	HandSignEnable *bool `json:"handSignEnable,omitempty" xml:"handSignEnable,omitempty"`
	// 是否使用上次签名
	ResignEnable *bool `json:"resignEnable,omitempty" xml:"resignEnable,omitempty"`
}

func (s GetProcessConfigResponseBodyResultHandSignConf) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigResponseBodyResultHandSignConf) GoString() string {
	return s.String()
}

func (s *GetProcessConfigResponseBodyResultHandSignConf) SetHandSignEnable(v bool) *GetProcessConfigResponseBodyResultHandSignConf {
	s.HandSignEnable = &v
	return s
}

func (s *GetProcessConfigResponseBodyResultHandSignConf) SetResignEnable(v bool) *GetProcessConfigResponseBodyResultHandSignConf {
	s.ResignEnable = &v
	return s
}

type GetProcessConfigResponseBodyResultSubstituteSubmitConf struct {
	// 是否允许代提交
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 代提交人
	SubmitterList []*GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList `json:"submitterList,omitempty" xml:"submitterList,omitempty" type:"Repeated"`
}

func (s GetProcessConfigResponseBodyResultSubstituteSubmitConf) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigResponseBodyResultSubstituteSubmitConf) GoString() string {
	return s.String()
}

func (s *GetProcessConfigResponseBodyResultSubstituteSubmitConf) SetEnable(v bool) *GetProcessConfigResponseBodyResultSubstituteSubmitConf {
	s.Enable = &v
	return s
}

func (s *GetProcessConfigResponseBodyResultSubstituteSubmitConf) SetSubmitterList(v []*GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList) *GetProcessConfigResponseBodyResultSubstituteSubmitConf {
	s.SubmitterList = v
	return s
}

type GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList struct {
	// 名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 类型
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 员工staffId/部门id
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList) GoString() string {
	return s.String()
}

func (s *GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList) SetName(v string) *GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList {
	s.Name = &v
	return s
}

func (s *GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList) SetType(v string) *GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList {
	s.Type = &v
	return s
}

func (s *GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList) SetValue(v string) *GetProcessConfigResponseBodyResultSubstituteSubmitConfSubmitterList {
	s.Value = &v
	return s
}

type GetProcessConfigResponseBodyResultTitleGenRule struct {
	// 规则表达式
	Express *string `json:"express,omitempty" xml:"express,omitempty"`
	// 规则类型
	Type *int32 `json:"type,omitempty" xml:"type,omitempty"`
}

func (s GetProcessConfigResponseBodyResultTitleGenRule) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigResponseBodyResultTitleGenRule) GoString() string {
	return s.String()
}

func (s *GetProcessConfigResponseBodyResultTitleGenRule) SetExpress(v string) *GetProcessConfigResponseBodyResultTitleGenRule {
	s.Express = &v
	return s
}

func (s *GetProcessConfigResponseBodyResultTitleGenRule) SetType(v int32) *GetProcessConfigResponseBodyResultTitleGenRule {
	s.Type = &v
	return s
}

type GetProcessConfigResponseBodyResultVisibility struct {
	// 类型
	Type *int32 `json:"type,omitempty" xml:"type,omitempty"`
	// 员工staffId/部门id
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s GetProcessConfigResponseBodyResultVisibility) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigResponseBodyResultVisibility) GoString() string {
	return s.String()
}

func (s *GetProcessConfigResponseBodyResultVisibility) SetType(v int32) *GetProcessConfigResponseBodyResultVisibility {
	s.Type = &v
	return s
}

func (s *GetProcessConfigResponseBodyResultVisibility) SetValue(v string) *GetProcessConfigResponseBodyResultVisibility {
	s.Value = &v
	return s
}

type GetProcessConfigResponse struct {
	Headers map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetProcessConfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetProcessConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s GetProcessConfigResponse) GoString() string {
	return s.String()
}

func (s *GetProcessConfigResponse) SetHeaders(v map[string]*string) *GetProcessConfigResponse {
	s.Headers = v
	return s
}

func (s *GetProcessConfigResponse) SetBody(v *GetProcessConfigResponseBody) *GetProcessConfigResponse {
	s.Body = v
	return s
}

type GetProcessInstanceHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GetProcessInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetProcessInstanceHeaders) GoString() string {
	return s.String()
}

func (s *GetProcessInstanceHeaders) SetCommonHeaders(v map[string]*string) *GetProcessInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetProcessInstanceHeaders) SetXAcsDingtalkAccessToken(v string) *GetProcessInstanceHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GetProcessInstanceRequest struct {
	// 审批实例ID企业内部应用可通过获取审批实例ID列表接口获取。钉钉三方企业应用可以通过推送的审批事件中获取，参考biz_type=22。
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
}

func (s GetProcessInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s GetProcessInstanceRequest) GoString() string {
	return s.String()
}

func (s *GetProcessInstanceRequest) SetProcessInstanceId(v string) *GetProcessInstanceRequest {
	s.ProcessInstanceId = &v
	return s
}

type GetProcessInstanceResponseBody struct {
	// 返回结果。
	Result *GetProcessInstanceResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
	// 调用是否成功。
	Success *string `json:"success,omitempty" xml:"success,omitempty"`
}

func (s GetProcessInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetProcessInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *GetProcessInstanceResponseBody) SetResult(v *GetProcessInstanceResponseBodyResult) *GetProcessInstanceResponseBody {
	s.Result = v
	return s
}

func (s *GetProcessInstanceResponseBody) SetSuccess(v string) *GetProcessInstanceResponseBody {
	s.Success = &v
	return s
}

type GetProcessInstanceResponseBodyResult struct {
	// 审批人userid列表。
	ApproverUserIds []*string `json:"approverUserIds,omitempty" xml:"approverUserIds,omitempty" type:"Repeated"`
	// 审批附属实例列表，当已经通过的审批实例被修改或撤销，会生成一个新的实例，作为原有审批实例的附属。  如果想知道当前已经通过的审批实例的状态，可以依次遍历它的附属列表，查询里面每个实例的biz_action。
	AttachedProcessInstanceIds []*string `json:"attachedProcessInstanceIds,omitempty" xml:"attachedProcessInstanceIds,omitempty" type:"Repeated"`
	// 审批实例业务动作：  MODIFY：表示该审批实例是基于原来的实例修改而来  REVOKE：表示该审批实例是由原来的实例撤销后重新发起的  NONE表示正常发起
	BizAction *string `json:"bizAction,omitempty" xml:"bizAction,omitempty"`
	// 审批实例业务编号。
	BusinessId *string `json:"businessId,omitempty" xml:"businessId,omitempty"`
	// 抄送人userid列表。
	CcUserIds []*string `json:"ccUserIds,omitempty" xml:"ccUserIds,omitempty" type:"Repeated"`
	// 创建时间。
	CreateTime *string `json:"createTime,omitempty" xml:"createTime,omitempty"`
	// 结束时间。
	FinishTime *string `json:"finishTime,omitempty" xml:"finishTime,omitempty"`
	// 表单详情列表。
	FormComponentValues []*GetProcessInstanceResponseBodyResultFormComponentValues `json:"formComponentValues,omitempty" xml:"formComponentValues,omitempty" type:"Repeated"`
	// 主流程实例标识。
	MainProcessInstanceId *string `json:"mainProcessInstanceId,omitempty" xml:"mainProcessInstanceId,omitempty"`
	// 操作记录列表。
	OperationRecords []*GetProcessInstanceResponseBodyResultOperationRecords `json:"operationRecords,omitempty" xml:"operationRecords,omitempty" type:"Repeated"`
	// 发起人的部门。-1表示根部门。
	OriginatorDeptId *string `json:"originatorDeptId,omitempty" xml:"originatorDeptId,omitempty"`
	// 发起人的部门名。
	OriginatorDeptName *string `json:"originatorDeptName,omitempty" xml:"originatorDeptName,omitempty"`
	// 发起人的userid。
	OriginatorUserId *string `json:"originatorUserId,omitempty" xml:"originatorUserId,omitempty"`
	// 审批结果：  agree：同意  refuse：拒绝。 说明 status为COMPLETED且result为同意时，表示审批单完结并审批通过。
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 审批状态：  NEW：新创建  RUNNING：审批中  TERMINATED：被终止  COMPLETED：完成  CANCELED：取消
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// 任务列表。
	Tasks []*GetProcessInstanceResponseBodyResultTasks `json:"tasks,omitempty" xml:"tasks,omitempty" type:"Repeated"`
	// 审批实例标题。
	Title *string `json:"title,omitempty" xml:"title,omitempty"`
}

func (s GetProcessInstanceResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s GetProcessInstanceResponseBodyResult) GoString() string {
	return s.String()
}

func (s *GetProcessInstanceResponseBodyResult) SetApproverUserIds(v []*string) *GetProcessInstanceResponseBodyResult {
	s.ApproverUserIds = v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetAttachedProcessInstanceIds(v []*string) *GetProcessInstanceResponseBodyResult {
	s.AttachedProcessInstanceIds = v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetBizAction(v string) *GetProcessInstanceResponseBodyResult {
	s.BizAction = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetBusinessId(v string) *GetProcessInstanceResponseBodyResult {
	s.BusinessId = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetCcUserIds(v []*string) *GetProcessInstanceResponseBodyResult {
	s.CcUserIds = v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetCreateTime(v string) *GetProcessInstanceResponseBodyResult {
	s.CreateTime = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetFinishTime(v string) *GetProcessInstanceResponseBodyResult {
	s.FinishTime = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetFormComponentValues(v []*GetProcessInstanceResponseBodyResultFormComponentValues) *GetProcessInstanceResponseBodyResult {
	s.FormComponentValues = v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetMainProcessInstanceId(v string) *GetProcessInstanceResponseBodyResult {
	s.MainProcessInstanceId = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetOperationRecords(v []*GetProcessInstanceResponseBodyResultOperationRecords) *GetProcessInstanceResponseBodyResult {
	s.OperationRecords = v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetOriginatorDeptId(v string) *GetProcessInstanceResponseBodyResult {
	s.OriginatorDeptId = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetOriginatorDeptName(v string) *GetProcessInstanceResponseBodyResult {
	s.OriginatorDeptName = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetOriginatorUserId(v string) *GetProcessInstanceResponseBodyResult {
	s.OriginatorUserId = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetResult(v string) *GetProcessInstanceResponseBodyResult {
	s.Result = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetStatus(v string) *GetProcessInstanceResponseBodyResult {
	s.Status = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetTasks(v []*GetProcessInstanceResponseBodyResultTasks) *GetProcessInstanceResponseBodyResult {
	s.Tasks = v
	return s
}

func (s *GetProcessInstanceResponseBodyResult) SetTitle(v string) *GetProcessInstanceResponseBodyResult {
	s.Title = &v
	return s
}

type GetProcessInstanceResponseBodyResultFormComponentValues struct {
	// 组件别名。
	BizAlias *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	// 组件类型。
	ComponentType *string `json:"componentType,omitempty" xml:"componentType,omitempty"`
	// 标签扩展值。
	ExtValue *string `json:"extValue,omitempty" xml:"extValue,omitempty"`
	// 组件ID。
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 组件名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 标签值。
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s GetProcessInstanceResponseBodyResultFormComponentValues) String() string {
	return tea.Prettify(s)
}

func (s GetProcessInstanceResponseBodyResultFormComponentValues) GoString() string {
	return s.String()
}

func (s *GetProcessInstanceResponseBodyResultFormComponentValues) SetBizAlias(v string) *GetProcessInstanceResponseBodyResultFormComponentValues {
	s.BizAlias = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultFormComponentValues) SetComponentType(v string) *GetProcessInstanceResponseBodyResultFormComponentValues {
	s.ComponentType = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultFormComponentValues) SetExtValue(v string) *GetProcessInstanceResponseBodyResultFormComponentValues {
	s.ExtValue = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultFormComponentValues) SetId(v string) *GetProcessInstanceResponseBodyResultFormComponentValues {
	s.Id = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultFormComponentValues) SetName(v string) *GetProcessInstanceResponseBodyResultFormComponentValues {
	s.Name = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultFormComponentValues) SetValue(v string) *GetProcessInstanceResponseBodyResultFormComponentValues {
	s.Value = &v
	return s
}

type GetProcessInstanceResponseBodyResultOperationRecords struct {
	// 评论附件列表。
	Attachments []*GetProcessInstanceResponseBodyResultOperationRecordsAttachments `json:"attachments,omitempty" xml:"attachments,omitempty" type:"Repeated"`
	// 抄送人userIds列表
	CcUserIds []*string `json:"ccUserIds,omitempty" xml:"ccUserIds,omitempty" type:"Repeated"`
	// 操作时间。
	Date *string `json:"date,omitempty" xml:"date,omitempty"`
	// 评论内容。  审批操作附带评论时才返回该字段。
	Remark *string `json:"remark,omitempty" xml:"remark,omitempty"`
	// 操作结果：  AGREE：同意  REFUSE：拒绝  NONE
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 操作类型：  EXECUTE_TASK_NORMAL：正常执行任务  EXECUTE_TASK_AGENT：代理人执行任务  APPEND_TASK_BEFORE：前加签任务  APPEND_TASK_AFTER：后加签任务  REDIRECT_TASK：转交任务  START_PROCESS_INSTANCE：发起流程实例  TERMINATE_PROCESS_INSTANCE：终止(撤销)流程实例  FINISH_PROCESS_INSTANCE：结束流程实例  ADD_REMARK：添加评论  REDIRECT_PROCESS：审批退回  PROCESS_CC：抄送
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 操作人userid。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s GetProcessInstanceResponseBodyResultOperationRecords) String() string {
	return tea.Prettify(s)
}

func (s GetProcessInstanceResponseBodyResultOperationRecords) GoString() string {
	return s.String()
}

func (s *GetProcessInstanceResponseBodyResultOperationRecords) SetAttachments(v []*GetProcessInstanceResponseBodyResultOperationRecordsAttachments) *GetProcessInstanceResponseBodyResultOperationRecords {
	s.Attachments = v
	return s
}

func (s *GetProcessInstanceResponseBodyResultOperationRecords) SetCcUserIds(v []*string) *GetProcessInstanceResponseBodyResultOperationRecords {
	s.CcUserIds = v
	return s
}

func (s *GetProcessInstanceResponseBodyResultOperationRecords) SetDate(v string) *GetProcessInstanceResponseBodyResultOperationRecords {
	s.Date = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultOperationRecords) SetRemark(v string) *GetProcessInstanceResponseBodyResultOperationRecords {
	s.Remark = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultOperationRecords) SetResult(v string) *GetProcessInstanceResponseBodyResultOperationRecords {
	s.Result = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultOperationRecords) SetType(v string) *GetProcessInstanceResponseBodyResultOperationRecords {
	s.Type = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultOperationRecords) SetUserId(v string) *GetProcessInstanceResponseBodyResultOperationRecords {
	s.UserId = &v
	return s
}

type GetProcessInstanceResponseBodyResultOperationRecordsAttachments struct {
	// 附件ID。
	FileId *string `json:"fileId,omitempty" xml:"fileId,omitempty"`
	// 附件名称。
	FileName *string `json:"fileName,omitempty" xml:"fileName,omitempty"`
	// 附件大小。
	FileSize *string `json:"fileSize,omitempty" xml:"fileSize,omitempty"`
	// 附件类型。
	FileType *string `json:"fileType,omitempty" xml:"fileType,omitempty"`
}

func (s GetProcessInstanceResponseBodyResultOperationRecordsAttachments) String() string {
	return tea.Prettify(s)
}

func (s GetProcessInstanceResponseBodyResultOperationRecordsAttachments) GoString() string {
	return s.String()
}

func (s *GetProcessInstanceResponseBodyResultOperationRecordsAttachments) SetFileId(v string) *GetProcessInstanceResponseBodyResultOperationRecordsAttachments {
	s.FileId = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultOperationRecordsAttachments) SetFileName(v string) *GetProcessInstanceResponseBodyResultOperationRecordsAttachments {
	s.FileName = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultOperationRecordsAttachments) SetFileSize(v string) *GetProcessInstanceResponseBodyResultOperationRecordsAttachments {
	s.FileSize = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultOperationRecordsAttachments) SetFileType(v string) *GetProcessInstanceResponseBodyResultOperationRecordsAttachments {
	s.FileType = &v
	return s
}

type GetProcessInstanceResponseBodyResultTasks struct {
	// 任务节点ID。
	ActivityId *string `json:"activityId,omitempty" xml:"activityId,omitempty"`
	// 开始时间。
	CreateTime *string `json:"createTime,omitempty" xml:"createTime,omitempty"`
	// 结束时间。
	FinishTime *string `json:"finishTime,omitempty" xml:"finishTime,omitempty"`
	// 移动端任务URL。
	MobileUrl *string `json:"mobileUrl,omitempty" xml:"mobileUrl,omitempty"`
	// PC端任务URL。
	PcUrl *string `json:"pcUrl,omitempty" xml:"pcUrl,omitempty"`
	// 实例ID。
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 结果：  AGREE：同意  REFUSE：拒绝  REDIRECTED：转交
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 任务状态：  NEW：未启动  RUNNING：处理中  PAUSED：暂停  CANCELED：取消  COMPLETED：完成  TERMINATED：终止
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// 任务ID。
	TaskId *int64 `json:"taskId,omitempty" xml:"taskId,omitempty"`
	// 任务处理人。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s GetProcessInstanceResponseBodyResultTasks) String() string {
	return tea.Prettify(s)
}

func (s GetProcessInstanceResponseBodyResultTasks) GoString() string {
	return s.String()
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetActivityId(v string) *GetProcessInstanceResponseBodyResultTasks {
	s.ActivityId = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetCreateTime(v string) *GetProcessInstanceResponseBodyResultTasks {
	s.CreateTime = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetFinishTime(v string) *GetProcessInstanceResponseBodyResultTasks {
	s.FinishTime = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetMobileUrl(v string) *GetProcessInstanceResponseBodyResultTasks {
	s.MobileUrl = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetPcUrl(v string) *GetProcessInstanceResponseBodyResultTasks {
	s.PcUrl = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetProcessInstanceId(v string) *GetProcessInstanceResponseBodyResultTasks {
	s.ProcessInstanceId = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetResult(v string) *GetProcessInstanceResponseBodyResultTasks {
	s.Result = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetStatus(v string) *GetProcessInstanceResponseBodyResultTasks {
	s.Status = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetTaskId(v int64) *GetProcessInstanceResponseBodyResultTasks {
	s.TaskId = &v
	return s
}

func (s *GetProcessInstanceResponseBodyResultTasks) SetUserId(v string) *GetProcessInstanceResponseBodyResultTasks {
	s.UserId = &v
	return s
}

type GetProcessInstanceResponse struct {
	Headers map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetProcessInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetProcessInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s GetProcessInstanceResponse) GoString() string {
	return s.String()
}

func (s *GetProcessInstanceResponse) SetHeaders(v map[string]*string) *GetProcessInstanceResponse {
	s.Headers = v
	return s
}

func (s *GetProcessInstanceResponse) SetBody(v *GetProcessInstanceResponseBody) *GetProcessInstanceResponse {
	s.Body = v
	return s
}

type GetSpaceWithDownloadAuthHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GetSpaceWithDownloadAuthHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetSpaceWithDownloadAuthHeaders) GoString() string {
	return s.String()
}

func (s *GetSpaceWithDownloadAuthHeaders) SetCommonHeaders(v map[string]*string) *GetSpaceWithDownloadAuthHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetSpaceWithDownloadAuthHeaders) SetXAcsDingtalkAccessToken(v string) *GetSpaceWithDownloadAuthHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GetSpaceWithDownloadAuthRequest struct {
	// 应用的agentid。
	AgentId *int64 `json:"agentId,omitempty" xml:"agentId,omitempty"`
	// 审批附件ID。
	FileId *string `json:"fileId,omitempty" xml:"fileId,omitempty"`
	// 附件ID列表，支持批量授权，最大列表长度：20。
	FileIdList []*string `json:"fileIdList,omitempty" xml:"fileIdList,omitempty" type:"Repeated"`
	// 实例ID。
	//
	// 企业内部应用
	//
	// 可通过获取审批实例ID列表接口获取。
	//
	// 第三方企业应用
	//
	// 可以通过推送的审批事件中获取，参考biz_type=22。
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 授权允许预览附件的用户userid。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s GetSpaceWithDownloadAuthRequest) String() string {
	return tea.Prettify(s)
}

func (s GetSpaceWithDownloadAuthRequest) GoString() string {
	return s.String()
}

func (s *GetSpaceWithDownloadAuthRequest) SetAgentId(v int64) *GetSpaceWithDownloadAuthRequest {
	s.AgentId = &v
	return s
}

func (s *GetSpaceWithDownloadAuthRequest) SetFileId(v string) *GetSpaceWithDownloadAuthRequest {
	s.FileId = &v
	return s
}

func (s *GetSpaceWithDownloadAuthRequest) SetFileIdList(v []*string) *GetSpaceWithDownloadAuthRequest {
	s.FileIdList = v
	return s
}

func (s *GetSpaceWithDownloadAuthRequest) SetProcessInstanceId(v string) *GetSpaceWithDownloadAuthRequest {
	s.ProcessInstanceId = &v
	return s
}

func (s *GetSpaceWithDownloadAuthRequest) SetUserId(v string) *GetSpaceWithDownloadAuthRequest {
	s.UserId = &v
	return s
}

type GetSpaceWithDownloadAuthResponseBody struct {
	// 返回结果。
	Result *GetSpaceWithDownloadAuthResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
	// 接口调用是否成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s GetSpaceWithDownloadAuthResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetSpaceWithDownloadAuthResponseBody) GoString() string {
	return s.String()
}

func (s *GetSpaceWithDownloadAuthResponseBody) SetResult(v *GetSpaceWithDownloadAuthResponseBodyResult) *GetSpaceWithDownloadAuthResponseBody {
	s.Result = v
	return s
}

func (s *GetSpaceWithDownloadAuthResponseBody) SetSuccess(v bool) *GetSpaceWithDownloadAuthResponseBody {
	s.Success = &v
	return s
}

type GetSpaceWithDownloadAuthResponseBodyResult struct {
	// 钉盘空间ID。
	SpaceId *int64 `json:"spaceId,omitempty" xml:"spaceId,omitempty"`
}

func (s GetSpaceWithDownloadAuthResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s GetSpaceWithDownloadAuthResponseBodyResult) GoString() string {
	return s.String()
}

func (s *GetSpaceWithDownloadAuthResponseBodyResult) SetSpaceId(v int64) *GetSpaceWithDownloadAuthResponseBodyResult {
	s.SpaceId = &v
	return s
}

type GetSpaceWithDownloadAuthResponse struct {
	Headers map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetSpaceWithDownloadAuthResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetSpaceWithDownloadAuthResponse) String() string {
	return tea.Prettify(s)
}

func (s GetSpaceWithDownloadAuthResponse) GoString() string {
	return s.String()
}

func (s *GetSpaceWithDownloadAuthResponse) SetHeaders(v map[string]*string) *GetSpaceWithDownloadAuthResponse {
	s.Headers = v
	return s
}

func (s *GetSpaceWithDownloadAuthResponse) SetBody(v *GetSpaceWithDownloadAuthResponseBody) *GetSpaceWithDownloadAuthResponse {
	s.Body = v
	return s
}

type GetUserTodoTaskSumHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GetUserTodoTaskSumHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetUserTodoTaskSumHeaders) GoString() string {
	return s.String()
}

func (s *GetUserTodoTaskSumHeaders) SetCommonHeaders(v map[string]*string) *GetUserTodoTaskSumHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetUserTodoTaskSumHeaders) SetXAcsDingtalkAccessToken(v string) *GetUserTodoTaskSumHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GetUserTodoTaskSumRequest struct {
	// 要查询的用户userid。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s GetUserTodoTaskSumRequest) String() string {
	return tea.Prettify(s)
}

func (s GetUserTodoTaskSumRequest) GoString() string {
	return s.String()
}

func (s *GetUserTodoTaskSumRequest) SetUserId(v string) *GetUserTodoTaskSumRequest {
	s.UserId = &v
	return s
}

type GetUserTodoTaskSumResponseBody struct {
	// 待处理的审批数量。
	Result *int32 `json:"result,omitempty" xml:"result,omitempty"`
}

func (s GetUserTodoTaskSumResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetUserTodoTaskSumResponseBody) GoString() string {
	return s.String()
}

func (s *GetUserTodoTaskSumResponseBody) SetResult(v int32) *GetUserTodoTaskSumResponseBody {
	s.Result = &v
	return s
}

type GetUserTodoTaskSumResponse struct {
	Headers map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetUserTodoTaskSumResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetUserTodoTaskSumResponse) String() string {
	return tea.Prettify(s)
}

func (s GetUserTodoTaskSumResponse) GoString() string {
	return s.String()
}

func (s *GetUserTodoTaskSumResponse) SetHeaders(v map[string]*string) *GetUserTodoTaskSumResponse {
	s.Headers = v
	return s
}

func (s *GetUserTodoTaskSumResponse) SetBody(v *GetUserTodoTaskSumResponseBody) *GetUserTodoTaskSumResponse {
	s.Body = v
	return s
}

type GrantCspaceAuthorizationHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GrantCspaceAuthorizationHeaders) String() string {
	return tea.Prettify(s)
}

func (s GrantCspaceAuthorizationHeaders) GoString() string {
	return s.String()
}

func (s *GrantCspaceAuthorizationHeaders) SetCommonHeaders(v map[string]*string) *GrantCspaceAuthorizationHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GrantCspaceAuthorizationHeaders) SetXAcsDingtalkAccessToken(v string) *GrantCspaceAuthorizationHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GrantCspaceAuthorizationRequest struct {
	// 权限有效时间，单位为秒。
	DurationSeconds *int64 `json:"durationSeconds,omitempty" xml:"durationSeconds,omitempty"`
	// 审批控件 id。
	SpaceId *string `json:"spaceId,omitempty" xml:"spaceId,omitempty"`
	// 权限类型。
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 用户 id。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s GrantCspaceAuthorizationRequest) String() string {
	return tea.Prettify(s)
}

func (s GrantCspaceAuthorizationRequest) GoString() string {
	return s.String()
}

func (s *GrantCspaceAuthorizationRequest) SetDurationSeconds(v int64) *GrantCspaceAuthorizationRequest {
	s.DurationSeconds = &v
	return s
}

func (s *GrantCspaceAuthorizationRequest) SetSpaceId(v string) *GrantCspaceAuthorizationRequest {
	s.SpaceId = &v
	return s
}

func (s *GrantCspaceAuthorizationRequest) SetType(v string) *GrantCspaceAuthorizationRequest {
	s.Type = &v
	return s
}

func (s *GrantCspaceAuthorizationRequest) SetUserId(v string) *GrantCspaceAuthorizationRequest {
	s.UserId = &v
	return s
}

type GrantCspaceAuthorizationResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s GrantCspaceAuthorizationResponse) String() string {
	return tea.Prettify(s)
}

func (s GrantCspaceAuthorizationResponse) GoString() string {
	return s.String()
}

func (s *GrantCspaceAuthorizationResponse) SetHeaders(v map[string]*string) *GrantCspaceAuthorizationResponse {
	s.Headers = v
	return s
}

type GrantProcessInstanceForDownloadFileHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s GrantProcessInstanceForDownloadFileHeaders) String() string {
	return tea.Prettify(s)
}

func (s GrantProcessInstanceForDownloadFileHeaders) GoString() string {
	return s.String()
}

func (s *GrantProcessInstanceForDownloadFileHeaders) SetCommonHeaders(v map[string]*string) *GrantProcessInstanceForDownloadFileHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GrantProcessInstanceForDownloadFileHeaders) SetXAcsDingtalkAccessToken(v string) *GrantProcessInstanceForDownloadFileHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type GrantProcessInstanceForDownloadFileRequest struct {
	// 文件id，调用获取审批实例详情接口获取。
	FileId *string `json:"fileId,omitempty" xml:"fileId,omitempty"`
	// 实例ID。
	//
	// 调用获取审批实例详情接口获取。
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
}

func (s GrantProcessInstanceForDownloadFileRequest) String() string {
	return tea.Prettify(s)
}

func (s GrantProcessInstanceForDownloadFileRequest) GoString() string {
	return s.String()
}

func (s *GrantProcessInstanceForDownloadFileRequest) SetFileId(v string) *GrantProcessInstanceForDownloadFileRequest {
	s.FileId = &v
	return s
}

func (s *GrantProcessInstanceForDownloadFileRequest) SetProcessInstanceId(v string) *GrantProcessInstanceForDownloadFileRequest {
	s.ProcessInstanceId = &v
	return s
}

type GrantProcessInstanceForDownloadFileResponseBody struct {
	// 返回结果。
	Result *GrantProcessInstanceForDownloadFileResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
	// 接口调用是否成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s GrantProcessInstanceForDownloadFileResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GrantProcessInstanceForDownloadFileResponseBody) GoString() string {
	return s.String()
}

func (s *GrantProcessInstanceForDownloadFileResponseBody) SetResult(v *GrantProcessInstanceForDownloadFileResponseBodyResult) *GrantProcessInstanceForDownloadFileResponseBody {
	s.Result = v
	return s
}

func (s *GrantProcessInstanceForDownloadFileResponseBody) SetSuccess(v bool) *GrantProcessInstanceForDownloadFileResponseBody {
	s.Success = &v
	return s
}

type GrantProcessInstanceForDownloadFileResponseBodyResult struct {
	// 文件下载地址。
	DownloadUri *string `json:"downloadUri,omitempty" xml:"downloadUri,omitempty"`
	// 文件ID。
	FileId *string `json:"fileId,omitempty" xml:"fileId,omitempty"`
	// 钉盘空间ID。
	SpaceId *int64 `json:"spaceId,omitempty" xml:"spaceId,omitempty"`
}

func (s GrantProcessInstanceForDownloadFileResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s GrantProcessInstanceForDownloadFileResponseBodyResult) GoString() string {
	return s.String()
}

func (s *GrantProcessInstanceForDownloadFileResponseBodyResult) SetDownloadUri(v string) *GrantProcessInstanceForDownloadFileResponseBodyResult {
	s.DownloadUri = &v
	return s
}

func (s *GrantProcessInstanceForDownloadFileResponseBodyResult) SetFileId(v string) *GrantProcessInstanceForDownloadFileResponseBodyResult {
	s.FileId = &v
	return s
}

func (s *GrantProcessInstanceForDownloadFileResponseBodyResult) SetSpaceId(v int64) *GrantProcessInstanceForDownloadFileResponseBodyResult {
	s.SpaceId = &v
	return s
}

type GrantProcessInstanceForDownloadFileResponse struct {
	Headers map[string]*string                               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GrantProcessInstanceForDownloadFileResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GrantProcessInstanceForDownloadFileResponse) String() string {
	return tea.Prettify(s)
}

func (s GrantProcessInstanceForDownloadFileResponse) GoString() string {
	return s.String()
}

func (s *GrantProcessInstanceForDownloadFileResponse) SetHeaders(v map[string]*string) *GrantProcessInstanceForDownloadFileResponse {
	s.Headers = v
	return s
}

func (s *GrantProcessInstanceForDownloadFileResponse) SetBody(v *GrantProcessInstanceForDownloadFileResponseBody) *GrantProcessInstanceForDownloadFileResponse {
	s.Body = v
	return s
}

type InstallAppHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s InstallAppHeaders) String() string {
	return tea.Prettify(s)
}

func (s InstallAppHeaders) GoString() string {
	return s.String()
}

func (s *InstallAppHeaders) SetCommonHeaders(v map[string]*string) *InstallAppHeaders {
	s.CommonHeaders = v
	return s
}

func (s *InstallAppHeaders) SetXAcsDingtalkAccessToken(v string) *InstallAppHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type InstallAppRequest struct {
	// 业务分组
	BizGroup *string `json:"bizGroup,omitempty" xml:"bizGroup,omitempty"`
	// 安装选项
	//
	InstallOption *InstallAppRequestInstallOption `json:"installOption,omitempty" xml:"installOption,omitempty" type:"Struct"`
	// 安装的目录名称
	SourceDirName *string `json:"sourceDirName,omitempty" xml:"sourceDirName,omitempty"`
}

func (s InstallAppRequest) String() string {
	return tea.Prettify(s)
}

func (s InstallAppRequest) GoString() string {
	return s.String()
}

func (s *InstallAppRequest) SetBizGroup(v string) *InstallAppRequest {
	s.BizGroup = &v
	return s
}

func (s *InstallAppRequest) SetInstallOption(v *InstallAppRequestInstallOption) *InstallAppRequest {
	s.InstallOption = v
	return s
}

func (s *InstallAppRequest) SetSourceDirName(v string) *InstallAppRequest {
	s.SourceDirName = &v
	return s
}

type InstallAppRequestInstallOption struct {
	// 是否同步，目前只有同步
	IsSync *bool `json:"isSync,omitempty" xml:"isSync,omitempty"`
}

func (s InstallAppRequestInstallOption) String() string {
	return tea.Prettify(s)
}

func (s InstallAppRequestInstallOption) GoString() string {
	return s.String()
}

func (s *InstallAppRequestInstallOption) SetIsSync(v bool) *InstallAppRequestInstallOption {
	s.IsSync = &v
	return s
}

type InstallAppResponseBody struct {
	// 返回对象列表
	Result []*InstallAppResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Repeated"`
}

func (s InstallAppResponseBody) String() string {
	return tea.Prettify(s)
}

func (s InstallAppResponseBody) GoString() string {
	return s.String()
}

func (s *InstallAppResponseBody) SetResult(v []*InstallAppResponseBodyResult) *InstallAppResponseBody {
	s.Result = v
	return s
}

type InstallAppResponseBodyResult struct {
	// 套件业务类型
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 模版名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 模版processcode
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s InstallAppResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s InstallAppResponseBodyResult) GoString() string {
	return s.String()
}

func (s *InstallAppResponseBodyResult) SetBizType(v string) *InstallAppResponseBodyResult {
	s.BizType = &v
	return s
}

func (s *InstallAppResponseBodyResult) SetName(v string) *InstallAppResponseBodyResult {
	s.Name = &v
	return s
}

func (s *InstallAppResponseBodyResult) SetProcessCode(v string) *InstallAppResponseBodyResult {
	s.ProcessCode = &v
	return s
}

type InstallAppResponse struct {
	Headers map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *InstallAppResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s InstallAppResponse) String() string {
	return tea.Prettify(s)
}

func (s InstallAppResponse) GoString() string {
	return s.String()
}

func (s *InstallAppResponse) SetHeaders(v map[string]*string) *InstallAppResponse {
	s.Headers = v
	return s
}

func (s *InstallAppResponse) SetBody(v *InstallAppResponseBody) *InstallAppResponse {
	s.Body = v
	return s
}

type ListProcessInstanceIdsHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s ListProcessInstanceIdsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListProcessInstanceIdsHeaders) GoString() string {
	return s.String()
}

func (s *ListProcessInstanceIdsHeaders) SetCommonHeaders(v map[string]*string) *ListProcessInstanceIdsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListProcessInstanceIdsHeaders) SetXAcsDingtalkAccessToken(v string) *ListProcessInstanceIdsHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type ListProcessInstanceIdsRequest struct {
	// 审批实例结束时间，Unix时间戳，单位毫秒。  例如：获取审批单发起时间在2020.4.10-2020.4.14之间审批单，该值传2020.4.14 23:59:59对应的时间戳1586879999000。
	EndTime *int64 `json:"endTime,omitempty" xml:"endTime,omitempty"`
	// 分页参数，每页大小，最多传20。
	MaxResults *int64 `json:"maxResults,omitempty" xml:"maxResults,omitempty"`
	// 分页查询的游标，最开始传0，后续传返回参数中的nextToken值。
	NextToken *int64 `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// 审批流的唯一码。
	//
	// processCode在审批模板编辑页面的URL中获取。
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 审批实例开始时间。Unix时间戳，单位毫秒。
	//
	// 例如：获取审批单发起时间在2020.4.10-2020.4.14之间审批单，该值传2020.4.10 00:00:00对应的时间戳1586448000000。
	StartTime *int64 `json:"startTime,omitempty" xml:"startTime,omitempty"`
	// 发起userid列表，最大列表长度为10。
	UserIds []*string `json:"userIds,omitempty" xml:"userIds,omitempty" type:"Repeated"`
}

func (s ListProcessInstanceIdsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListProcessInstanceIdsRequest) GoString() string {
	return s.String()
}

func (s *ListProcessInstanceIdsRequest) SetEndTime(v int64) *ListProcessInstanceIdsRequest {
	s.EndTime = &v
	return s
}

func (s *ListProcessInstanceIdsRequest) SetMaxResults(v int64) *ListProcessInstanceIdsRequest {
	s.MaxResults = &v
	return s
}

func (s *ListProcessInstanceIdsRequest) SetNextToken(v int64) *ListProcessInstanceIdsRequest {
	s.NextToken = &v
	return s
}

func (s *ListProcessInstanceIdsRequest) SetProcessCode(v string) *ListProcessInstanceIdsRequest {
	s.ProcessCode = &v
	return s
}

func (s *ListProcessInstanceIdsRequest) SetStartTime(v int64) *ListProcessInstanceIdsRequest {
	s.StartTime = &v
	return s
}

func (s *ListProcessInstanceIdsRequest) SetUserIds(v []*string) *ListProcessInstanceIdsRequest {
	s.UserIds = v
	return s
}

type ListProcessInstanceIdsResponseBody struct {
	// 返回结果。
	Result *ListProcessInstanceIdsResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
	// 接口请求是否成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s ListProcessInstanceIdsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListProcessInstanceIdsResponseBody) GoString() string {
	return s.String()
}

func (s *ListProcessInstanceIdsResponseBody) SetResult(v *ListProcessInstanceIdsResponseBodyResult) *ListProcessInstanceIdsResponseBody {
	s.Result = v
	return s
}

func (s *ListProcessInstanceIdsResponseBody) SetSuccess(v bool) *ListProcessInstanceIdsResponseBody {
	s.Success = &v
	return s
}

type ListProcessInstanceIdsResponseBodyResult struct {
	// 审批实例ID列表。
	List []*string `json:"list,omitempty" xml:"list,omitempty" type:"Repeated"`
	// 表示下次查询的游标，当返回结果没有该字段时表示没有更多数据了。
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListProcessInstanceIdsResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s ListProcessInstanceIdsResponseBodyResult) GoString() string {
	return s.String()
}

func (s *ListProcessInstanceIdsResponseBodyResult) SetList(v []*string) *ListProcessInstanceIdsResponseBodyResult {
	s.List = v
	return s
}

func (s *ListProcessInstanceIdsResponseBodyResult) SetNextToken(v string) *ListProcessInstanceIdsResponseBodyResult {
	s.NextToken = &v
	return s
}

type ListProcessInstanceIdsResponse struct {
	Headers map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ListProcessInstanceIdsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListProcessInstanceIdsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListProcessInstanceIdsResponse) GoString() string {
	return s.String()
}

func (s *ListProcessInstanceIdsResponse) SetHeaders(v map[string]*string) *ListProcessInstanceIdsResponse {
	s.Headers = v
	return s
}

func (s *ListProcessInstanceIdsResponse) SetBody(v *ListProcessInstanceIdsResponseBody) *ListProcessInstanceIdsResponse {
	s.Body = v
	return s
}

type ListTodoWorkRecordsHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s ListTodoWorkRecordsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListTodoWorkRecordsHeaders) GoString() string {
	return s.String()
}

func (s *ListTodoWorkRecordsHeaders) SetCommonHeaders(v map[string]*string) *ListTodoWorkRecordsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListTodoWorkRecordsHeaders) SetXAcsDingtalkAccessToken(v string) *ListTodoWorkRecordsHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type ListTodoWorkRecordsRequest struct {
	// 分页大小，最大值50。
	MaxResults *int32 `json:"maxResults,omitempty" xml:"maxResults,omitempty"`
	// 分页游标。
	//
	// 如果是首次调用，该参数传0。
	// 如果是非首次调用，该参数传上次调用时返回的nextToken。
	//
	NextToken *int32 `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// 待办事项的状态：
	//
	// 0：待处理
	//
	// -1：已经移除
	//
	Status *int32 `json:"status,omitempty" xml:"status,omitempty"`
	// 要查询的执行人userid。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s ListTodoWorkRecordsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListTodoWorkRecordsRequest) GoString() string {
	return s.String()
}

func (s *ListTodoWorkRecordsRequest) SetMaxResults(v int32) *ListTodoWorkRecordsRequest {
	s.MaxResults = &v
	return s
}

func (s *ListTodoWorkRecordsRequest) SetNextToken(v int32) *ListTodoWorkRecordsRequest {
	s.NextToken = &v
	return s
}

func (s *ListTodoWorkRecordsRequest) SetStatus(v int32) *ListTodoWorkRecordsRequest {
	s.Status = &v
	return s
}

func (s *ListTodoWorkRecordsRequest) SetUserId(v string) *ListTodoWorkRecordsRequest {
	s.UserId = &v
	return s
}

type ListTodoWorkRecordsResponseBody struct {
	// 查询结果。
	Result *ListTodoWorkRecordsResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s ListTodoWorkRecordsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListTodoWorkRecordsResponseBody) GoString() string {
	return s.String()
}

func (s *ListTodoWorkRecordsResponseBody) SetResult(v *ListTodoWorkRecordsResponseBodyResult) *ListTodoWorkRecordsResponseBody {
	s.Result = v
	return s
}

type ListTodoWorkRecordsResponseBodyResult struct {
	// 待办事项列表。
	List []*ListTodoWorkRecordsResponseBodyResultList `json:"list,omitempty" xml:"list,omitempty" type:"Repeated"`
	// 分页游标。不为空表示有数据。
	NextToken *int64 `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListTodoWorkRecordsResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s ListTodoWorkRecordsResponseBodyResult) GoString() string {
	return s.String()
}

func (s *ListTodoWorkRecordsResponseBodyResult) SetList(v []*ListTodoWorkRecordsResponseBodyResultList) *ListTodoWorkRecordsResponseBodyResult {
	s.List = v
	return s
}

func (s *ListTodoWorkRecordsResponseBodyResult) SetNextToken(v int64) *ListTodoWorkRecordsResponseBodyResult {
	s.NextToken = &v
	return s
}

type ListTodoWorkRecordsResponseBodyResultList struct {
	// 表单列表。
	Forms []*ListTodoWorkRecordsResponseBodyResultListForms `json:"forms,omitempty" xml:"forms,omitempty" type:"Repeated"`
	// 实例ID。
	InstanceId *string `json:"instanceId,omitempty" xml:"instanceId,omitempty"`
	// 待办任务ID。
	TaskId *int64 `json:"taskId,omitempty" xml:"taskId,omitempty"`
	// 待办标题。
	Title *string `json:"title,omitempty" xml:"title,omitempty"`
	// 待办跳转链接。
	Url *string `json:"url,omitempty" xml:"url,omitempty"`
}

func (s ListTodoWorkRecordsResponseBodyResultList) String() string {
	return tea.Prettify(s)
}

func (s ListTodoWorkRecordsResponseBodyResultList) GoString() string {
	return s.String()
}

func (s *ListTodoWorkRecordsResponseBodyResultList) SetForms(v []*ListTodoWorkRecordsResponseBodyResultListForms) *ListTodoWorkRecordsResponseBodyResultList {
	s.Forms = v
	return s
}

func (s *ListTodoWorkRecordsResponseBodyResultList) SetInstanceId(v string) *ListTodoWorkRecordsResponseBodyResultList {
	s.InstanceId = &v
	return s
}

func (s *ListTodoWorkRecordsResponseBodyResultList) SetTaskId(v int64) *ListTodoWorkRecordsResponseBodyResultList {
	s.TaskId = &v
	return s
}

func (s *ListTodoWorkRecordsResponseBodyResultList) SetTitle(v string) *ListTodoWorkRecordsResponseBodyResultList {
	s.Title = &v
	return s
}

func (s *ListTodoWorkRecordsResponseBodyResultList) SetUrl(v string) *ListTodoWorkRecordsResponseBodyResultList {
	s.Url = &v
	return s
}

type ListTodoWorkRecordsResponseBodyResultListForms struct {
	// 表单内容。
	Content *string `json:"content,omitempty" xml:"content,omitempty"`
	// 表单标题。
	Title *string `json:"title,omitempty" xml:"title,omitempty"`
}

func (s ListTodoWorkRecordsResponseBodyResultListForms) String() string {
	return tea.Prettify(s)
}

func (s ListTodoWorkRecordsResponseBodyResultListForms) GoString() string {
	return s.String()
}

func (s *ListTodoWorkRecordsResponseBodyResultListForms) SetContent(v string) *ListTodoWorkRecordsResponseBodyResultListForms {
	s.Content = &v
	return s
}

func (s *ListTodoWorkRecordsResponseBodyResultListForms) SetTitle(v string) *ListTodoWorkRecordsResponseBodyResultListForms {
	s.Title = &v
	return s
}

type ListTodoWorkRecordsResponse struct {
	Headers map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ListTodoWorkRecordsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListTodoWorkRecordsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListTodoWorkRecordsResponse) GoString() string {
	return s.String()
}

func (s *ListTodoWorkRecordsResponse) SetHeaders(v map[string]*string) *ListTodoWorkRecordsResponse {
	s.Headers = v
	return s
}

func (s *ListTodoWorkRecordsResponse) SetBody(v *ListTodoWorkRecordsResponseBody) *ListTodoWorkRecordsResponse {
	s.Body = v
	return s
}

type ListUserVisibleBpmsProcessesHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s ListUserVisibleBpmsProcessesHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListUserVisibleBpmsProcessesHeaders) GoString() string {
	return s.String()
}

func (s *ListUserVisibleBpmsProcessesHeaders) SetCommonHeaders(v map[string]*string) *ListUserVisibleBpmsProcessesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListUserVisibleBpmsProcessesHeaders) SetXAcsDingtalkAccessToken(v string) *ListUserVisibleBpmsProcessesHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type ListUserVisibleBpmsProcessesRequest struct {
	// 分页大小，最大可设置成100。
	MaxResults *int64 `json:"maxResults,omitempty" xml:"maxResults,omitempty"`
	// 分页游标，从0开始。根据返回结果里的nextToken是否为空来判断是否还有下一页，且再次调用时设置成nextToken的最新值。
	NextToken *int64 `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// 要查询的员工的userid。不传表示查询企业下所有审批表单。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s ListUserVisibleBpmsProcessesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListUserVisibleBpmsProcessesRequest) GoString() string {
	return s.String()
}

func (s *ListUserVisibleBpmsProcessesRequest) SetMaxResults(v int64) *ListUserVisibleBpmsProcessesRequest {
	s.MaxResults = &v
	return s
}

func (s *ListUserVisibleBpmsProcessesRequest) SetNextToken(v int64) *ListUserVisibleBpmsProcessesRequest {
	s.NextToken = &v
	return s
}

func (s *ListUserVisibleBpmsProcessesRequest) SetUserId(v string) *ListUserVisibleBpmsProcessesRequest {
	s.UserId = &v
	return s
}

type ListUserVisibleBpmsProcessesResponseBody struct {
	// 返回结果。
	Result *ListUserVisibleBpmsProcessesResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s ListUserVisibleBpmsProcessesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListUserVisibleBpmsProcessesResponseBody) GoString() string {
	return s.String()
}

func (s *ListUserVisibleBpmsProcessesResponseBody) SetResult(v *ListUserVisibleBpmsProcessesResponseBodyResult) *ListUserVisibleBpmsProcessesResponseBody {
	s.Result = v
	return s
}

type ListUserVisibleBpmsProcessesResponseBodyResult struct {
	// 下一次分页调用的值，当返回结果里没有nextToken时，表示分页结束。
	NextToken *int64 `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// 可见表单列表。
	ProcessList []*ListUserVisibleBpmsProcessesResponseBodyResultProcessList `json:"processList,omitempty" xml:"processList,omitempty" type:"Repeated"`
}

func (s ListUserVisibleBpmsProcessesResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s ListUserVisibleBpmsProcessesResponseBodyResult) GoString() string {
	return s.String()
}

func (s *ListUserVisibleBpmsProcessesResponseBodyResult) SetNextToken(v int64) *ListUserVisibleBpmsProcessesResponseBodyResult {
	s.NextToken = &v
	return s
}

func (s *ListUserVisibleBpmsProcessesResponseBodyResult) SetProcessList(v []*ListUserVisibleBpmsProcessesResponseBodyResultProcessList) *ListUserVisibleBpmsProcessesResponseBodyResult {
	s.ProcessList = v
	return s
}

type ListUserVisibleBpmsProcessesResponseBodyResultProcessList struct {
	// 图标URL。
	IconUrl *string `json:"iconUrl,omitempty" xml:"iconUrl,omitempty"`
	// 表单名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 表单唯一标识。
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 表单URL。
	Url *string `json:"url,omitempty" xml:"url,omitempty"`
}

func (s ListUserVisibleBpmsProcessesResponseBodyResultProcessList) String() string {
	return tea.Prettify(s)
}

func (s ListUserVisibleBpmsProcessesResponseBodyResultProcessList) GoString() string {
	return s.String()
}

func (s *ListUserVisibleBpmsProcessesResponseBodyResultProcessList) SetIconUrl(v string) *ListUserVisibleBpmsProcessesResponseBodyResultProcessList {
	s.IconUrl = &v
	return s
}

func (s *ListUserVisibleBpmsProcessesResponseBodyResultProcessList) SetName(v string) *ListUserVisibleBpmsProcessesResponseBodyResultProcessList {
	s.Name = &v
	return s
}

func (s *ListUserVisibleBpmsProcessesResponseBodyResultProcessList) SetProcessCode(v string) *ListUserVisibleBpmsProcessesResponseBodyResultProcessList {
	s.ProcessCode = &v
	return s
}

func (s *ListUserVisibleBpmsProcessesResponseBodyResultProcessList) SetUrl(v string) *ListUserVisibleBpmsProcessesResponseBodyResultProcessList {
	s.Url = &v
	return s
}

type ListUserVisibleBpmsProcessesResponse struct {
	Headers map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ListUserVisibleBpmsProcessesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListUserVisibleBpmsProcessesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListUserVisibleBpmsProcessesResponse) GoString() string {
	return s.String()
}

func (s *ListUserVisibleBpmsProcessesResponse) SetHeaders(v map[string]*string) *ListUserVisibleBpmsProcessesResponse {
	s.Headers = v
	return s
}

func (s *ListUserVisibleBpmsProcessesResponse) SetBody(v *ListUserVisibleBpmsProcessesResponseBody) *ListUserVisibleBpmsProcessesResponse {
	s.Body = v
	return s
}

type ProcessForecastHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s ProcessForecastHeaders) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastHeaders) GoString() string {
	return s.String()
}

func (s *ProcessForecastHeaders) SetCommonHeaders(v map[string]*string) *ProcessForecastHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ProcessForecastHeaders) SetXAcsDingtalkAccessToken(v string) *ProcessForecastHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type ProcessForecastRequest struct {
	// 部门ID
	DeptId *int32 `json:"deptId,omitempty" xml:"deptId,omitempty"`
	// 表单数据内容，控件列表
	FormComponentValues []*ProcessForecastRequestFormComponentValues `json:"formComponentValues,omitempty" xml:"formComponentValues,omitempty" type:"Repeated"`
	// 审批流的唯一码
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 审批发起人的userId
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s ProcessForecastRequest) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastRequest) GoString() string {
	return s.String()
}

func (s *ProcessForecastRequest) SetDeptId(v int32) *ProcessForecastRequest {
	s.DeptId = &v
	return s
}

func (s *ProcessForecastRequest) SetFormComponentValues(v []*ProcessForecastRequestFormComponentValues) *ProcessForecastRequest {
	s.FormComponentValues = v
	return s
}

func (s *ProcessForecastRequest) SetProcessCode(v string) *ProcessForecastRequest {
	s.ProcessCode = &v
	return s
}

func (s *ProcessForecastRequest) SetUserId(v string) *ProcessForecastRequest {
	s.UserId = &v
	return s
}

type ProcessForecastRequestFormComponentValues struct {
	// 控件别名
	BizAlias *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	// 控件类型
	ComponentType *string                                             `json:"componentType,omitempty" xml:"componentType,omitempty"`
	Details       []*ProcessForecastRequestFormComponentValuesDetails `json:"details,omitempty" xml:"details,omitempty" type:"Repeated"`
	// 控件扩展值
	ExtValue *string `json:"extValue,omitempty" xml:"extValue,omitempty"`
	// 控件id
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 控件名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 控件值
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s ProcessForecastRequestFormComponentValues) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastRequestFormComponentValues) GoString() string {
	return s.String()
}

func (s *ProcessForecastRequestFormComponentValues) SetBizAlias(v string) *ProcessForecastRequestFormComponentValues {
	s.BizAlias = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValues) SetComponentType(v string) *ProcessForecastRequestFormComponentValues {
	s.ComponentType = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValues) SetDetails(v []*ProcessForecastRequestFormComponentValuesDetails) *ProcessForecastRequestFormComponentValues {
	s.Details = v
	return s
}

func (s *ProcessForecastRequestFormComponentValues) SetExtValue(v string) *ProcessForecastRequestFormComponentValues {
	s.ExtValue = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValues) SetId(v string) *ProcessForecastRequestFormComponentValues {
	s.Id = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValues) SetName(v string) *ProcessForecastRequestFormComponentValues {
	s.Name = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValues) SetValue(v string) *ProcessForecastRequestFormComponentValues {
	s.Value = &v
	return s
}

type ProcessForecastRequestFormComponentValuesDetails struct {
	// 控件别名
	BizAlias *string                                                    `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	Details  []*ProcessForecastRequestFormComponentValuesDetailsDetails `json:"details,omitempty" xml:"details,omitempty" type:"Repeated"`
	// 控件扩展值
	ExtValue *string `json:"extValue,omitempty" xml:"extValue,omitempty"`
	// 控件id
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 控件名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 控件值
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s ProcessForecastRequestFormComponentValuesDetails) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastRequestFormComponentValuesDetails) GoString() string {
	return s.String()
}

func (s *ProcessForecastRequestFormComponentValuesDetails) SetBizAlias(v string) *ProcessForecastRequestFormComponentValuesDetails {
	s.BizAlias = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetails) SetDetails(v []*ProcessForecastRequestFormComponentValuesDetailsDetails) *ProcessForecastRequestFormComponentValuesDetails {
	s.Details = v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetails) SetExtValue(v string) *ProcessForecastRequestFormComponentValuesDetails {
	s.ExtValue = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetails) SetId(v string) *ProcessForecastRequestFormComponentValuesDetails {
	s.Id = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetails) SetName(v string) *ProcessForecastRequestFormComponentValuesDetails {
	s.Name = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetails) SetValue(v string) *ProcessForecastRequestFormComponentValuesDetails {
	s.Value = &v
	return s
}

type ProcessForecastRequestFormComponentValuesDetailsDetails struct {
	// 控件别名
	BizAlias      *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	ComponentType *string `json:"componentType,omitempty" xml:"componentType,omitempty"`
	// 控件扩展值
	ExtValue *string `json:"extValue,omitempty" xml:"extValue,omitempty"`
	// 控件id
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 控件名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 控件值
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s ProcessForecastRequestFormComponentValuesDetailsDetails) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastRequestFormComponentValuesDetailsDetails) GoString() string {
	return s.String()
}

func (s *ProcessForecastRequestFormComponentValuesDetailsDetails) SetBizAlias(v string) *ProcessForecastRequestFormComponentValuesDetailsDetails {
	s.BizAlias = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetailsDetails) SetComponentType(v string) *ProcessForecastRequestFormComponentValuesDetailsDetails {
	s.ComponentType = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetailsDetails) SetExtValue(v string) *ProcessForecastRequestFormComponentValuesDetailsDetails {
	s.ExtValue = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetailsDetails) SetId(v string) *ProcessForecastRequestFormComponentValuesDetailsDetails {
	s.Id = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetailsDetails) SetName(v string) *ProcessForecastRequestFormComponentValuesDetailsDetails {
	s.Name = &v
	return s
}

func (s *ProcessForecastRequestFormComponentValuesDetailsDetails) SetValue(v string) *ProcessForecastRequestFormComponentValuesDetailsDetails {
	s.Value = &v
	return s
}

type ProcessForecastResponseBody struct {
	// 返回结果
	Result *ProcessForecastResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s ProcessForecastResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastResponseBody) GoString() string {
	return s.String()
}

func (s *ProcessForecastResponseBody) SetResult(v *ProcessForecastResponseBodyResult) *ProcessForecastResponseBody {
	s.Result = v
	return s
}

type ProcessForecastResponseBodyResult struct {
	// 是否预测成功
	IsForecastSuccess *bool `json:"isForecastSuccess,omitempty" xml:"isForecastSuccess,omitempty"`
	// 是否静态流程
	IsStaticWorkflow *bool `json:"isStaticWorkflow,omitempty" xml:"isStaticWorkflow,omitempty"`
	// 流程 code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 流程 id
	ProcessId *int64 `json:"processId,omitempty" xml:"processId,omitempty"`
	// 用户 id
	UserId                *string                                                   `json:"userId,omitempty" xml:"userId,omitempty"`
	WorkflowActivityRules []*ProcessForecastResponseBodyResultWorkflowActivityRules `json:"workflowActivityRules,omitempty" xml:"workflowActivityRules,omitempty" type:"Repeated"`
	WorkflowForecastNodes []*ProcessForecastResponseBodyResultWorkflowForecastNodes `json:"workflowForecastNodes,omitempty" xml:"workflowForecastNodes,omitempty" type:"Repeated"`
}

func (s ProcessForecastResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastResponseBodyResult) GoString() string {
	return s.String()
}

func (s *ProcessForecastResponseBodyResult) SetIsForecastSuccess(v bool) *ProcessForecastResponseBodyResult {
	s.IsForecastSuccess = &v
	return s
}

func (s *ProcessForecastResponseBodyResult) SetIsStaticWorkflow(v bool) *ProcessForecastResponseBodyResult {
	s.IsStaticWorkflow = &v
	return s
}

func (s *ProcessForecastResponseBodyResult) SetProcessCode(v string) *ProcessForecastResponseBodyResult {
	s.ProcessCode = &v
	return s
}

func (s *ProcessForecastResponseBodyResult) SetProcessId(v int64) *ProcessForecastResponseBodyResult {
	s.ProcessId = &v
	return s
}

func (s *ProcessForecastResponseBodyResult) SetUserId(v string) *ProcessForecastResponseBodyResult {
	s.UserId = &v
	return s
}

func (s *ProcessForecastResponseBodyResult) SetWorkflowActivityRules(v []*ProcessForecastResponseBodyResultWorkflowActivityRules) *ProcessForecastResponseBodyResult {
	s.WorkflowActivityRules = v
	return s
}

func (s *ProcessForecastResponseBodyResult) SetWorkflowForecastNodes(v []*ProcessForecastResponseBodyResultWorkflowForecastNodes) *ProcessForecastResponseBodyResult {
	s.WorkflowForecastNodes = v
	return s
}

type ProcessForecastResponseBodyResultWorkflowActivityRules struct {
	// 节点 id
	ActivityId *string `json:"activityId,omitempty" xml:"activityId,omitempty"`
	// 节点名称
	ActivityName *string `json:"activityName,omitempty" xml:"activityName,omitempty"`
	// 规则类型
	ActivityType *string `json:"activityType,omitempty" xml:"activityType,omitempty"`
	// 是否自选审批节点
	IsTargetSelect *bool `json:"isTargetSelect,omitempty" xml:"isTargetSelect,omitempty"`
	// 流程中前一个节点的 id
	PrevActivityId *string `json:"prevActivityId,omitempty" xml:"prevActivityId,omitempty"`
	// 节点操作人信息
	WorkflowActor *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor `json:"workflowActor,omitempty" xml:"workflowActor,omitempty" type:"Struct"`
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRules) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRules) GoString() string {
	return s.String()
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRules) SetActivityId(v string) *ProcessForecastResponseBodyResultWorkflowActivityRules {
	s.ActivityId = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRules) SetActivityName(v string) *ProcessForecastResponseBodyResultWorkflowActivityRules {
	s.ActivityName = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRules) SetActivityType(v string) *ProcessForecastResponseBodyResultWorkflowActivityRules {
	s.ActivityType = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRules) SetIsTargetSelect(v bool) *ProcessForecastResponseBodyResultWorkflowActivityRules {
	s.IsTargetSelect = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRules) SetPrevActivityId(v string) *ProcessForecastResponseBodyResultWorkflowActivityRules {
	s.PrevActivityId = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRules) SetWorkflowActor(v *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) *ProcessForecastResponseBodyResultWorkflowActivityRules {
	s.WorkflowActor = v
	return s
}

type ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor struct {
	// 节点激活类型
	ActorActivateType *string `json:"actorActivateType,omitempty" xml:"actorActivateType,omitempty"`
	// 节点操作人 key
	ActorKey *string `json:"actorKey,omitempty" xml:"actorKey,omitempty"`
	// 节点操作人选择范围
	ActorSelectionRange *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRange `json:"actorSelectionRange,omitempty" xml:"actorSelectionRange,omitempty" type:"Struct"`
	// 节点操作人选择范围类型
	ActorSelectionType *string `json:"actorSelectionType,omitempty" xml:"actorSelectionType,omitempty"`
	// 节点操作人类型
	ActorType *string `json:"actorType,omitempty" xml:"actorType,omitempty"`
	// 是否允许多选，还是仅允许选一人
	AllowedMulti *bool `json:"allowedMulti,omitempty" xml:"allowedMulti,omitempty"`
	// 节点审批方式
	ApprovalMethod *string `json:"approvalMethod,omitempty" xml:"approvalMethod,omitempty"`
	// 节点审批类型
	ApprovalType *string `json:"approvalType,omitempty" xml:"approvalType,omitempty"`
	// 该审批人节点在发起审批时是否必填
	Required *bool `json:"required,omitempty" xml:"required,omitempty"`
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) GoString() string {
	return s.String()
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) SetActorActivateType(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor {
	s.ActorActivateType = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) SetActorKey(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor {
	s.ActorKey = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) SetActorSelectionRange(v *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRange) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor {
	s.ActorSelectionRange = v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) SetActorSelectionType(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor {
	s.ActorSelectionType = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) SetActorType(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor {
	s.ActorType = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) SetAllowedMulti(v bool) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor {
	s.AllowedMulti = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) SetApprovalMethod(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor {
	s.ApprovalMethod = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) SetApprovalType(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor {
	s.ApprovalType = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor) SetRequired(v bool) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActor {
	s.Required = &v
	return s
}

type ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRange struct {
	// 审批指定成员
	Approvals []*ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeApprovals `json:"approvals,omitempty" xml:"approvals,omitempty" type:"Repeated"`
	// 审批指定角色
	Labels []*ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeLabels `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRange) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRange) GoString() string {
	return s.String()
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRange) SetApprovals(v []*ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeApprovals) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRange {
	s.Approvals = v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRange) SetLabels(v []*ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeLabels) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRange {
	s.Labels = v
	return s
}

type ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeApprovals struct {
	// 员工姓名
	UserName *string `json:"userName,omitempty" xml:"userName,omitempty"`
	// 员工 userId
	WorkNo *string `json:"workNo,omitempty" xml:"workNo,omitempty"`
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeApprovals) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeApprovals) GoString() string {
	return s.String()
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeApprovals) SetUserName(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeApprovals {
	s.UserName = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeApprovals) SetWorkNo(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeApprovals {
	s.WorkNo = &v
	return s
}

type ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeLabels struct {
	// 角色名字
	LabelNames *string `json:"labelNames,omitempty" xml:"labelNames,omitempty"`
	// 角色 id
	Labels *string `json:"labels,omitempty" xml:"labels,omitempty"`
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeLabels) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeLabels) GoString() string {
	return s.String()
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeLabels) SetLabelNames(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeLabels {
	s.LabelNames = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeLabels) SetLabels(v string) *ProcessForecastResponseBodyResultWorkflowActivityRulesWorkflowActorActorSelectionRangeLabels {
	s.Labels = &v
	return s
}

type ProcessForecastResponseBodyResultWorkflowForecastNodes struct {
	// 节点 id
	ActivityId *string `json:"activityId,omitempty" xml:"activityId,omitempty"`
	// 节点出线 id
	OutId *string `json:"outId,omitempty" xml:"outId,omitempty"`
}

func (s ProcessForecastResponseBodyResultWorkflowForecastNodes) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastResponseBodyResultWorkflowForecastNodes) GoString() string {
	return s.String()
}

func (s *ProcessForecastResponseBodyResultWorkflowForecastNodes) SetActivityId(v string) *ProcessForecastResponseBodyResultWorkflowForecastNodes {
	s.ActivityId = &v
	return s
}

func (s *ProcessForecastResponseBodyResultWorkflowForecastNodes) SetOutId(v string) *ProcessForecastResponseBodyResultWorkflowForecastNodes {
	s.OutId = &v
	return s
}

type ProcessForecastResponse struct {
	Headers map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ProcessForecastResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ProcessForecastResponse) String() string {
	return tea.Prettify(s)
}

func (s ProcessForecastResponse) GoString() string {
	return s.String()
}

func (s *ProcessForecastResponse) SetHeaders(v map[string]*string) *ProcessForecastResponse {
	s.Headers = v
	return s
}

func (s *ProcessForecastResponse) SetBody(v *ProcessForecastResponseBody) *ProcessForecastResponse {
	s.Body = v
	return s
}

type QueryAllFormInstancesHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s QueryAllFormInstancesHeaders) String() string {
	return tea.Prettify(s)
}

func (s QueryAllFormInstancesHeaders) GoString() string {
	return s.String()
}

func (s *QueryAllFormInstancesHeaders) SetCommonHeaders(v map[string]*string) *QueryAllFormInstancesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *QueryAllFormInstancesHeaders) SetXAcsDingtalkAccessToken(v string) *QueryAllFormInstancesHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type QueryAllFormInstancesRequest struct {
	// 应用搭建id
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 表单模板id
	FormCode *string `json:"formCode,omitempty" xml:"formCode,omitempty"`
	// 翻页size
	MaxResults *int32 `json:"maxResults,omitempty" xml:"maxResults,omitempty"`
	// 分页游标，第一次调用传空或者null
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s QueryAllFormInstancesRequest) String() string {
	return tea.Prettify(s)
}

func (s QueryAllFormInstancesRequest) GoString() string {
	return s.String()
}

func (s *QueryAllFormInstancesRequest) SetAppUuid(v string) *QueryAllFormInstancesRequest {
	s.AppUuid = &v
	return s
}

func (s *QueryAllFormInstancesRequest) SetFormCode(v string) *QueryAllFormInstancesRequest {
	s.FormCode = &v
	return s
}

func (s *QueryAllFormInstancesRequest) SetMaxResults(v int32) *QueryAllFormInstancesRequest {
	s.MaxResults = &v
	return s
}

func (s *QueryAllFormInstancesRequest) SetNextToken(v string) *QueryAllFormInstancesRequest {
	s.NextToken = &v
	return s
}

type QueryAllFormInstancesResponseBody struct {
	// 分页结果
	Result *QueryAllFormInstancesResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s QueryAllFormInstancesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s QueryAllFormInstancesResponseBody) GoString() string {
	return s.String()
}

func (s *QueryAllFormInstancesResponseBody) SetResult(v *QueryAllFormInstancesResponseBodyResult) *QueryAllFormInstancesResponseBody {
	s.Result = v
	return s
}

type QueryAllFormInstancesResponseBodyResult struct {
	// 是否有更多数据
	HasMore *bool `json:"hasMore,omitempty" xml:"hasMore,omitempty"`
	// 分页大小
	MaxResults *int64 `json:"maxResults,omitempty" xml:"maxResults,omitempty"`
	// 下一页的游标
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// 表单列表
	Values []*QueryAllFormInstancesResponseBodyResultValues `json:"values,omitempty" xml:"values,omitempty" type:"Repeated"`
}

func (s QueryAllFormInstancesResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s QueryAllFormInstancesResponseBodyResult) GoString() string {
	return s.String()
}

func (s *QueryAllFormInstancesResponseBodyResult) SetHasMore(v bool) *QueryAllFormInstancesResponseBodyResult {
	s.HasMore = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResult) SetMaxResults(v int64) *QueryAllFormInstancesResponseBodyResult {
	s.MaxResults = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResult) SetNextToken(v string) *QueryAllFormInstancesResponseBodyResult {
	s.NextToken = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResult) SetValues(v []*QueryAllFormInstancesResponseBodyResultValues) *QueryAllFormInstancesResponseBodyResult {
	s.Values = v
	return s
}

type QueryAllFormInstancesResponseBodyResultValues struct {
	// 应用搭建id
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 扩展信息
	Attributes map[string]interface{} `json:"attributes,omitempty" xml:"attributes,omitempty"`
	// 创建时间
	CreateTimestamp *int64 `json:"createTimestamp,omitempty" xml:"createTimestamp,omitempty"`
	// 创建人
	Creator *string `json:"creator,omitempty" xml:"creator,omitempty"`
	// 表单模板code
	FormCode *string `json:"formCode,omitempty" xml:"formCode,omitempty"`
	// 表单实例数据
	FormInstDataList []*QueryAllFormInstancesResponseBodyResultValuesFormInstDataList `json:"formInstDataList,omitempty" xml:"formInstDataList,omitempty" type:"Repeated"`
	// 表单实例id
	FormInstanceId *string `json:"formInstanceId,omitempty" xml:"formInstanceId,omitempty"`
	// 修改人
	Modifier *string `json:"modifier,omitempty" xml:"modifier,omitempty"`
	// 修改时间
	ModifyTimestamp *int64 `json:"modifyTimestamp,omitempty" xml:"modifyTimestamp,omitempty"`
	// 外部业务编码
	OutBizCode *string `json:"outBizCode,omitempty" xml:"outBizCode,omitempty"`
	// 外部实例编码
	OutInstanceId *string `json:"outInstanceId,omitempty" xml:"outInstanceId,omitempty"`
	// 标题
	Title *string `json:"title,omitempty" xml:"title,omitempty"`
}

func (s QueryAllFormInstancesResponseBodyResultValues) String() string {
	return tea.Prettify(s)
}

func (s QueryAllFormInstancesResponseBodyResultValues) GoString() string {
	return s.String()
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetAppUuid(v string) *QueryAllFormInstancesResponseBodyResultValues {
	s.AppUuid = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetAttributes(v map[string]interface{}) *QueryAllFormInstancesResponseBodyResultValues {
	s.Attributes = v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetCreateTimestamp(v int64) *QueryAllFormInstancesResponseBodyResultValues {
	s.CreateTimestamp = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetCreator(v string) *QueryAllFormInstancesResponseBodyResultValues {
	s.Creator = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetFormCode(v string) *QueryAllFormInstancesResponseBodyResultValues {
	s.FormCode = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetFormInstDataList(v []*QueryAllFormInstancesResponseBodyResultValuesFormInstDataList) *QueryAllFormInstancesResponseBodyResultValues {
	s.FormInstDataList = v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetFormInstanceId(v string) *QueryAllFormInstancesResponseBodyResultValues {
	s.FormInstanceId = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetModifier(v string) *QueryAllFormInstancesResponseBodyResultValues {
	s.Modifier = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetModifyTimestamp(v int64) *QueryAllFormInstancesResponseBodyResultValues {
	s.ModifyTimestamp = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetOutBizCode(v string) *QueryAllFormInstancesResponseBodyResultValues {
	s.OutBizCode = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetOutInstanceId(v string) *QueryAllFormInstancesResponseBodyResultValues {
	s.OutInstanceId = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValues) SetTitle(v string) *QueryAllFormInstancesResponseBodyResultValues {
	s.Title = &v
	return s
}

type QueryAllFormInstancesResponseBodyResultValuesFormInstDataList struct {
	// 控件别名
	BizAlias *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	// 控件类型
	ComponentType *string `json:"componentType,omitempty" xml:"componentType,omitempty"`
	// 表单控件扩展数据
	ExtendValue *string `json:"extendValue,omitempty" xml:"extendValue,omitempty"`
	// 控件唯一id
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// 控件名称
	Label *string `json:"label,omitempty" xml:"label,omitempty"`
	// 控件填写的数据
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s QueryAllFormInstancesResponseBodyResultValuesFormInstDataList) String() string {
	return tea.Prettify(s)
}

func (s QueryAllFormInstancesResponseBodyResultValuesFormInstDataList) GoString() string {
	return s.String()
}

func (s *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList) SetBizAlias(v string) *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList {
	s.BizAlias = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList) SetComponentType(v string) *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList {
	s.ComponentType = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList) SetExtendValue(v string) *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList {
	s.ExtendValue = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList) SetKey(v string) *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList {
	s.Key = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList) SetLabel(v string) *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList {
	s.Label = &v
	return s
}

func (s *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList) SetValue(v string) *QueryAllFormInstancesResponseBodyResultValuesFormInstDataList {
	s.Value = &v
	return s
}

type QueryAllFormInstancesResponse struct {
	Headers map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *QueryAllFormInstancesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s QueryAllFormInstancesResponse) String() string {
	return tea.Prettify(s)
}

func (s QueryAllFormInstancesResponse) GoString() string {
	return s.String()
}

func (s *QueryAllFormInstancesResponse) SetHeaders(v map[string]*string) *QueryAllFormInstancesResponse {
	s.Headers = v
	return s
}

func (s *QueryAllFormInstancesResponse) SetBody(v *QueryAllFormInstancesResponseBody) *QueryAllFormInstancesResponse {
	s.Body = v
	return s
}

type QueryAllProcessInstancesHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s QueryAllProcessInstancesHeaders) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesHeaders) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesHeaders) SetCommonHeaders(v map[string]*string) *QueryAllProcessInstancesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *QueryAllProcessInstancesHeaders) SetXAcsDingtalkAccessToken(v string) *QueryAllProcessInstancesHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type QueryAllProcessInstancesRequest struct {
	// 应用编码
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 结束时间
	EndTimeInMills *int64 `json:"endTimeInMills,omitempty" xml:"endTimeInMills,omitempty"`
	// 分页大小
	MaxResults *int64 `json:"maxResults,omitempty" xml:"maxResults,omitempty"`
	// 分页起始值
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// 模板编码
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 开始时间
	StartTimeInMills *int64 `json:"startTimeInMills,omitempty" xml:"startTimeInMills,omitempty"`
}

func (s QueryAllProcessInstancesRequest) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesRequest) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesRequest) SetAppUuid(v string) *QueryAllProcessInstancesRequest {
	s.AppUuid = &v
	return s
}

func (s *QueryAllProcessInstancesRequest) SetEndTimeInMills(v int64) *QueryAllProcessInstancesRequest {
	s.EndTimeInMills = &v
	return s
}

func (s *QueryAllProcessInstancesRequest) SetMaxResults(v int64) *QueryAllProcessInstancesRequest {
	s.MaxResults = &v
	return s
}

func (s *QueryAllProcessInstancesRequest) SetNextToken(v string) *QueryAllProcessInstancesRequest {
	s.NextToken = &v
	return s
}

func (s *QueryAllProcessInstancesRequest) SetProcessCode(v string) *QueryAllProcessInstancesRequest {
	s.ProcessCode = &v
	return s
}

func (s *QueryAllProcessInstancesRequest) SetStartTimeInMills(v int64) *QueryAllProcessInstancesRequest {
	s.StartTimeInMills = &v
	return s
}

type QueryAllProcessInstancesResponseBody struct {
	// result
	Result *QueryAllProcessInstancesResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s QueryAllProcessInstancesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesResponseBody) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesResponseBody) SetResult(v *QueryAllProcessInstancesResponseBodyResult) *QueryAllProcessInstancesResponseBody {
	s.Result = v
	return s
}

type QueryAllProcessInstancesResponseBodyResult struct {
	// 是否有更多数据
	HasMore *bool                                             `json:"hasMore,omitempty" xml:"hasMore,omitempty"`
	List    []*QueryAllProcessInstancesResponseBodyResultList `json:"list,omitempty" xml:"list,omitempty" type:"Repeated"`
	// 总数
	MaxResults *int64 `json:"maxResults,omitempty" xml:"maxResults,omitempty"`
	// 下次获取数据的游标
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s QueryAllProcessInstancesResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesResponseBodyResult) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesResponseBodyResult) SetHasMore(v bool) *QueryAllProcessInstancesResponseBodyResult {
	s.HasMore = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResult) SetList(v []*QueryAllProcessInstancesResponseBodyResultList) *QueryAllProcessInstancesResponseBodyResult {
	s.List = v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResult) SetMaxResults(v int64) *QueryAllProcessInstancesResponseBodyResult {
	s.MaxResults = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResult) SetNextToken(v string) *QueryAllProcessInstancesResponseBodyResult {
	s.NextToken = &v
	return s
}

type QueryAllProcessInstancesResponseBodyResultList struct {
	// 附属单信息
	AttachedProcessInstanceIds *string `json:"attachedProcessInstanceIds,omitempty" xml:"attachedProcessInstanceIds,omitempty"`
	// 审批单编号
	BusinessId *string `json:"businessId,omitempty" xml:"businessId,omitempty"`
	// 审批单创建时间
	CreateTime *int64 `json:"createTime,omitempty" xml:"createTime,omitempty"`
	// 审批结束时间
	FinishTime          *int64                                                               `json:"finishTime,omitempty" xml:"finishTime,omitempty"`
	FormComponentValues []*QueryAllProcessInstancesResponseBodyResultListFormComponentValues `json:"formComponentValues,omitempty" xml:"formComponentValues,omitempty" type:"Repeated"`
	// 主单实例Id
	MainProcessInstanceId *string `json:"mainProcessInstanceId,omitempty" xml:"mainProcessInstanceId,omitempty"`
	// 审批单操作记录
	OperationRecords []*QueryAllProcessInstancesResponseBodyResultListOperationRecords `json:"operationRecords,omitempty" xml:"operationRecords,omitempty" type:"Repeated"`
	// 发起人部门id
	OriginatorDeptId *string `json:"originatorDeptId,omitempty" xml:"originatorDeptId,omitempty"`
	// 发起者userId
	OriginatorUserid *string `json:"originatorUserid,omitempty" xml:"originatorUserid,omitempty"`
	// 流程实例ID
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 审批结果
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 审批单状态
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// 任务列表
	Tasks []*QueryAllProcessInstancesResponseBodyResultListTasks `json:"tasks,omitempty" xml:"tasks,omitempty" type:"Repeated"`
	// 审批单标题
	Title *string `json:"title,omitempty" xml:"title,omitempty"`
}

func (s QueryAllProcessInstancesResponseBodyResultList) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesResponseBodyResultList) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetAttachedProcessInstanceIds(v string) *QueryAllProcessInstancesResponseBodyResultList {
	s.AttachedProcessInstanceIds = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetBusinessId(v string) *QueryAllProcessInstancesResponseBodyResultList {
	s.BusinessId = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetCreateTime(v int64) *QueryAllProcessInstancesResponseBodyResultList {
	s.CreateTime = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetFinishTime(v int64) *QueryAllProcessInstancesResponseBodyResultList {
	s.FinishTime = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetFormComponentValues(v []*QueryAllProcessInstancesResponseBodyResultListFormComponentValues) *QueryAllProcessInstancesResponseBodyResultList {
	s.FormComponentValues = v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetMainProcessInstanceId(v string) *QueryAllProcessInstancesResponseBodyResultList {
	s.MainProcessInstanceId = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetOperationRecords(v []*QueryAllProcessInstancesResponseBodyResultListOperationRecords) *QueryAllProcessInstancesResponseBodyResultList {
	s.OperationRecords = v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetOriginatorDeptId(v string) *QueryAllProcessInstancesResponseBodyResultList {
	s.OriginatorDeptId = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetOriginatorUserid(v string) *QueryAllProcessInstancesResponseBodyResultList {
	s.OriginatorUserid = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetProcessInstanceId(v string) *QueryAllProcessInstancesResponseBodyResultList {
	s.ProcessInstanceId = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetResult(v string) *QueryAllProcessInstancesResponseBodyResultList {
	s.Result = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetStatus(v string) *QueryAllProcessInstancesResponseBodyResultList {
	s.Status = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetTasks(v []*QueryAllProcessInstancesResponseBodyResultListTasks) *QueryAllProcessInstancesResponseBodyResultList {
	s.Tasks = v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultList) SetTitle(v string) *QueryAllProcessInstancesResponseBodyResultList {
	s.Title = &v
	return s
}

type QueryAllProcessInstancesResponseBodyResultListFormComponentValues struct {
	// 控件扩展数据
	ExtValue *string `json:"extValue,omitempty" xml:"extValue,omitempty"`
	// 控件id
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 控件名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 控件数据
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s QueryAllProcessInstancesResponseBodyResultListFormComponentValues) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesResponseBodyResultListFormComponentValues) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesResponseBodyResultListFormComponentValues) SetExtValue(v string) *QueryAllProcessInstancesResponseBodyResultListFormComponentValues {
	s.ExtValue = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListFormComponentValues) SetId(v string) *QueryAllProcessInstancesResponseBodyResultListFormComponentValues {
	s.Id = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListFormComponentValues) SetName(v string) *QueryAllProcessInstancesResponseBodyResultListFormComponentValues {
	s.Name = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListFormComponentValues) SetValue(v string) *QueryAllProcessInstancesResponseBodyResultListFormComponentValues {
	s.Value = &v
	return s
}

type QueryAllProcessInstancesResponseBodyResultListOperationRecords struct {
	// 评论附件
	Attachments []*QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments `json:"attachments,omitempty" xml:"attachments,omitempty" type:"Repeated"`
	// 操作类型
	OperationType *string `json:"operationType,omitempty" xml:"operationType,omitempty"`
	// 评论
	Remark *string `json:"remark,omitempty" xml:"remark,omitempty"`
	// 操作结果
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 操作时间戳
	Timestamp *int64 `json:"timestamp,omitempty" xml:"timestamp,omitempty"`
	// 操作人staffId
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s QueryAllProcessInstancesResponseBodyResultListOperationRecords) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesResponseBodyResultListOperationRecords) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecords) SetAttachments(v []*QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments) *QueryAllProcessInstancesResponseBodyResultListOperationRecords {
	s.Attachments = v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecords) SetOperationType(v string) *QueryAllProcessInstancesResponseBodyResultListOperationRecords {
	s.OperationType = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecords) SetRemark(v string) *QueryAllProcessInstancesResponseBodyResultListOperationRecords {
	s.Remark = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecords) SetResult(v string) *QueryAllProcessInstancesResponseBodyResultListOperationRecords {
	s.Result = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecords) SetTimestamp(v int64) *QueryAllProcessInstancesResponseBodyResultListOperationRecords {
	s.Timestamp = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecords) SetUserId(v string) *QueryAllProcessInstancesResponseBodyResultListOperationRecords {
	s.UserId = &v
	return s
}

type QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments struct {
	// 附件钉盘id
	FileId *string `json:"fileId,omitempty" xml:"fileId,omitempty"`
	// 附件名称
	FileName *string `json:"fileName,omitempty" xml:"fileName,omitempty"`
	// 文件大小
	FileSize *string `json:"fileSize,omitempty" xml:"fileSize,omitempty"`
	// 文件类型
	FileType *string `json:"fileType,omitempty" xml:"fileType,omitempty"`
}

func (s QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments) SetFileId(v string) *QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments {
	s.FileId = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments) SetFileName(v string) *QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments {
	s.FileName = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments) SetFileSize(v string) *QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments {
	s.FileSize = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments) SetFileType(v string) *QueryAllProcessInstancesResponseBodyResultListOperationRecordsAttachments {
	s.FileType = &v
	return s
}

type QueryAllProcessInstancesResponseBodyResultListTasks struct {
	// 节点id
	ActivityId *string `json:"activityId,omitempty" xml:"activityId,omitempty"`
	// 任务创建时间戳
	CreateTimestamp *int64 `json:"createTimestamp,omitempty" xml:"createTimestamp,omitempty"`
	// 任务结束时间戳
	FinishTimestamp *int64 `json:"finishTimestamp,omitempty" xml:"finishTimestamp,omitempty"`
	// 任务结果
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 任务状态
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// 任务Id
	TaskId *int64 `json:"taskId,omitempty" xml:"taskId,omitempty"`
	// 任务处理人
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s QueryAllProcessInstancesResponseBodyResultListTasks) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesResponseBodyResultListTasks) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesResponseBodyResultListTasks) SetActivityId(v string) *QueryAllProcessInstancesResponseBodyResultListTasks {
	s.ActivityId = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListTasks) SetCreateTimestamp(v int64) *QueryAllProcessInstancesResponseBodyResultListTasks {
	s.CreateTimestamp = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListTasks) SetFinishTimestamp(v int64) *QueryAllProcessInstancesResponseBodyResultListTasks {
	s.FinishTimestamp = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListTasks) SetResult(v string) *QueryAllProcessInstancesResponseBodyResultListTasks {
	s.Result = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListTasks) SetStatus(v string) *QueryAllProcessInstancesResponseBodyResultListTasks {
	s.Status = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListTasks) SetTaskId(v int64) *QueryAllProcessInstancesResponseBodyResultListTasks {
	s.TaskId = &v
	return s
}

func (s *QueryAllProcessInstancesResponseBodyResultListTasks) SetUserId(v string) *QueryAllProcessInstancesResponseBodyResultListTasks {
	s.UserId = &v
	return s
}

type QueryAllProcessInstancesResponse struct {
	Headers map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *QueryAllProcessInstancesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s QueryAllProcessInstancesResponse) String() string {
	return tea.Prettify(s)
}

func (s QueryAllProcessInstancesResponse) GoString() string {
	return s.String()
}

func (s *QueryAllProcessInstancesResponse) SetHeaders(v map[string]*string) *QueryAllProcessInstancesResponse {
	s.Headers = v
	return s
}

func (s *QueryAllProcessInstancesResponse) SetBody(v *QueryAllProcessInstancesResponseBody) *QueryAllProcessInstancesResponse {
	s.Body = v
	return s
}

type QueryFormByBizTypeHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s QueryFormByBizTypeHeaders) String() string {
	return tea.Prettify(s)
}

func (s QueryFormByBizTypeHeaders) GoString() string {
	return s.String()
}

func (s *QueryFormByBizTypeHeaders) SetCommonHeaders(v map[string]*string) *QueryFormByBizTypeHeaders {
	s.CommonHeaders = v
	return s
}

func (s *QueryFormByBizTypeHeaders) SetXAcsDingtalkAccessToken(v string) *QueryFormByBizTypeHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type QueryFormByBizTypeRequest struct {
	// 应用搭建id
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 表单业务标识
	BizTypes []*string `json:"bizTypes,omitempty" xml:"bizTypes,omitempty" type:"Repeated"`
}

func (s QueryFormByBizTypeRequest) String() string {
	return tea.Prettify(s)
}

func (s QueryFormByBizTypeRequest) GoString() string {
	return s.String()
}

func (s *QueryFormByBizTypeRequest) SetAppUuid(v string) *QueryFormByBizTypeRequest {
	s.AppUuid = &v
	return s
}

func (s *QueryFormByBizTypeRequest) SetBizTypes(v []*string) *QueryFormByBizTypeRequest {
	s.BizTypes = v
	return s
}

type QueryFormByBizTypeResponseBody struct {
	// 模板列表
	Result []*QueryFormByBizTypeResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Repeated"`
}

func (s QueryFormByBizTypeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s QueryFormByBizTypeResponseBody) GoString() string {
	return s.String()
}

func (s *QueryFormByBizTypeResponseBody) SetResult(v []*QueryFormByBizTypeResponseBodyResult) *QueryFormByBizTypeResponseBody {
	s.Result = v
	return s
}

type QueryFormByBizTypeResponseBodyResult struct {
	// 表单类型，0为流程表单，1为数据表单
	AppType *int32 `json:"appType,omitempty" xml:"appType,omitempty"`
	// 应用搭建id
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 业务标识
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 表单控件描述
	Content *string `json:"content,omitempty" xml:"content,omitempty"`
	// 创建时间
	CreateTime *int64 `json:"createTime,omitempty" xml:"createTime,omitempty"`
	// 创建人
	Creator *string `json:"creator,omitempty" xml:"creator,omitempty"`
	// 模板code
	FormCode *string `json:"formCode,omitempty" xml:"formCode,omitempty"`
	// 表单uuid
	FormUuid *string `json:"formUuid,omitempty" xml:"formUuid,omitempty"`
	// 模板描述
	Memo *string `json:"memo,omitempty" xml:"memo,omitempty"`
	// 修改时间
	ModifedTime *int64 `json:"modifedTime,omitempty" xml:"modifedTime,omitempty"`
	// 模板名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 数据归属id
	OwnerId *string `json:"ownerId,omitempty" xml:"ownerId,omitempty"`
	// 模板状态
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s QueryFormByBizTypeResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s QueryFormByBizTypeResponseBodyResult) GoString() string {
	return s.String()
}

func (s *QueryFormByBizTypeResponseBodyResult) SetAppType(v int32) *QueryFormByBizTypeResponseBodyResult {
	s.AppType = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetAppUuid(v string) *QueryFormByBizTypeResponseBodyResult {
	s.AppUuid = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetBizType(v string) *QueryFormByBizTypeResponseBodyResult {
	s.BizType = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetContent(v string) *QueryFormByBizTypeResponseBodyResult {
	s.Content = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetCreateTime(v int64) *QueryFormByBizTypeResponseBodyResult {
	s.CreateTime = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetCreator(v string) *QueryFormByBizTypeResponseBodyResult {
	s.Creator = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetFormCode(v string) *QueryFormByBizTypeResponseBodyResult {
	s.FormCode = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetFormUuid(v string) *QueryFormByBizTypeResponseBodyResult {
	s.FormUuid = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetMemo(v string) *QueryFormByBizTypeResponseBodyResult {
	s.Memo = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetModifedTime(v int64) *QueryFormByBizTypeResponseBodyResult {
	s.ModifedTime = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetName(v string) *QueryFormByBizTypeResponseBodyResult {
	s.Name = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetOwnerId(v string) *QueryFormByBizTypeResponseBodyResult {
	s.OwnerId = &v
	return s
}

func (s *QueryFormByBizTypeResponseBodyResult) SetStatus(v string) *QueryFormByBizTypeResponseBodyResult {
	s.Status = &v
	return s
}

type QueryFormByBizTypeResponse struct {
	Headers map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *QueryFormByBizTypeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s QueryFormByBizTypeResponse) String() string {
	return tea.Prettify(s)
}

func (s QueryFormByBizTypeResponse) GoString() string {
	return s.String()
}

func (s *QueryFormByBizTypeResponse) SetHeaders(v map[string]*string) *QueryFormByBizTypeResponse {
	s.Headers = v
	return s
}

func (s *QueryFormByBizTypeResponse) SetBody(v *QueryFormByBizTypeResponseBody) *QueryFormByBizTypeResponse {
	s.Body = v
	return s
}

type QueryFormInstanceHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s QueryFormInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s QueryFormInstanceHeaders) GoString() string {
	return s.String()
}

func (s *QueryFormInstanceHeaders) SetCommonHeaders(v map[string]*string) *QueryFormInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *QueryFormInstanceHeaders) SetXAcsDingtalkAccessToken(v string) *QueryFormInstanceHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type QueryFormInstanceRequest struct {
	// 应用搭建id
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 表单模板Code
	FormCode *string `json:"formCode,omitempty" xml:"formCode,omitempty"`
	// 表单实例id
	FormInstanceId *string `json:"formInstanceId,omitempty" xml:"formInstanceId,omitempty"`
}

func (s QueryFormInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s QueryFormInstanceRequest) GoString() string {
	return s.String()
}

func (s *QueryFormInstanceRequest) SetAppUuid(v string) *QueryFormInstanceRequest {
	s.AppUuid = &v
	return s
}

func (s *QueryFormInstanceRequest) SetFormCode(v string) *QueryFormInstanceRequest {
	s.FormCode = &v
	return s
}

func (s *QueryFormInstanceRequest) SetFormInstanceId(v string) *QueryFormInstanceRequest {
	s.FormInstanceId = &v
	return s
}

type QueryFormInstanceResponseBody struct {
	// 应用搭建id
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 扩展信息
	Attributes map[string]interface{} `json:"attributes,omitempty" xml:"attributes,omitempty"`
	// 实例创建时间戳
	CreateTimestamp *int64 `json:"createTimestamp,omitempty" xml:"createTimestamp,omitempty"`
	// 创建人
	Creator *string `json:"creator,omitempty" xml:"creator,omitempty"`
	// 表单模板id
	FormCode *string `json:"formCode,omitempty" xml:"formCode,omitempty"`
	// 表单数据
	FormInstDataList []*QueryFormInstanceResponseBodyFormInstDataList `json:"formInstDataList,omitempty" xml:"formInstDataList,omitempty" type:"Repeated"`
	// 实例id
	FormInstanceId *string `json:"formInstanceId,omitempty" xml:"formInstanceId,omitempty"`
	// 修改人
	Modifier *string `json:"modifier,omitempty" xml:"modifier,omitempty"`
	// 实例最近修改时间戳
	ModifyTimestamp *int64 `json:"modifyTimestamp,omitempty" xml:"modifyTimestamp,omitempty"`
	// 外联业务code
	OutBizCode *string `json:"outBizCode,omitempty" xml:"outBizCode,omitempty"`
	// 外联业务实例id
	OutInstanceId *string `json:"outInstanceId,omitempty" xml:"outInstanceId,omitempty"`
	// 表单标题
	Title *string `json:"title,omitempty" xml:"title,omitempty"`
}

func (s QueryFormInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s QueryFormInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *QueryFormInstanceResponseBody) SetAppUuid(v string) *QueryFormInstanceResponseBody {
	s.AppUuid = &v
	return s
}

func (s *QueryFormInstanceResponseBody) SetAttributes(v map[string]interface{}) *QueryFormInstanceResponseBody {
	s.Attributes = v
	return s
}

func (s *QueryFormInstanceResponseBody) SetCreateTimestamp(v int64) *QueryFormInstanceResponseBody {
	s.CreateTimestamp = &v
	return s
}

func (s *QueryFormInstanceResponseBody) SetCreator(v string) *QueryFormInstanceResponseBody {
	s.Creator = &v
	return s
}

func (s *QueryFormInstanceResponseBody) SetFormCode(v string) *QueryFormInstanceResponseBody {
	s.FormCode = &v
	return s
}

func (s *QueryFormInstanceResponseBody) SetFormInstDataList(v []*QueryFormInstanceResponseBodyFormInstDataList) *QueryFormInstanceResponseBody {
	s.FormInstDataList = v
	return s
}

func (s *QueryFormInstanceResponseBody) SetFormInstanceId(v string) *QueryFormInstanceResponseBody {
	s.FormInstanceId = &v
	return s
}

func (s *QueryFormInstanceResponseBody) SetModifier(v string) *QueryFormInstanceResponseBody {
	s.Modifier = &v
	return s
}

func (s *QueryFormInstanceResponseBody) SetModifyTimestamp(v int64) *QueryFormInstanceResponseBody {
	s.ModifyTimestamp = &v
	return s
}

func (s *QueryFormInstanceResponseBody) SetOutBizCode(v string) *QueryFormInstanceResponseBody {
	s.OutBizCode = &v
	return s
}

func (s *QueryFormInstanceResponseBody) SetOutInstanceId(v string) *QueryFormInstanceResponseBody {
	s.OutInstanceId = &v
	return s
}

func (s *QueryFormInstanceResponseBody) SetTitle(v string) *QueryFormInstanceResponseBody {
	s.Title = &v
	return s
}

type QueryFormInstanceResponseBodyFormInstDataList struct {
	BizAlias      *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	ComponentType *string `json:"componentType,omitempty" xml:"componentType,omitempty"`
	ExtendValue   *string `json:"extendValue,omitempty" xml:"extendValue,omitempty"`
	Key           *string `json:"key,omitempty" xml:"key,omitempty"`
	Label         *string `json:"label,omitempty" xml:"label,omitempty"`
	Value         *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s QueryFormInstanceResponseBodyFormInstDataList) String() string {
	return tea.Prettify(s)
}

func (s QueryFormInstanceResponseBodyFormInstDataList) GoString() string {
	return s.String()
}

func (s *QueryFormInstanceResponseBodyFormInstDataList) SetBizAlias(v string) *QueryFormInstanceResponseBodyFormInstDataList {
	s.BizAlias = &v
	return s
}

func (s *QueryFormInstanceResponseBodyFormInstDataList) SetComponentType(v string) *QueryFormInstanceResponseBodyFormInstDataList {
	s.ComponentType = &v
	return s
}

func (s *QueryFormInstanceResponseBodyFormInstDataList) SetExtendValue(v string) *QueryFormInstanceResponseBodyFormInstDataList {
	s.ExtendValue = &v
	return s
}

func (s *QueryFormInstanceResponseBodyFormInstDataList) SetKey(v string) *QueryFormInstanceResponseBodyFormInstDataList {
	s.Key = &v
	return s
}

func (s *QueryFormInstanceResponseBodyFormInstDataList) SetLabel(v string) *QueryFormInstanceResponseBodyFormInstDataList {
	s.Label = &v
	return s
}

func (s *QueryFormInstanceResponseBodyFormInstDataList) SetValue(v string) *QueryFormInstanceResponseBodyFormInstDataList {
	s.Value = &v
	return s
}

type QueryFormInstanceResponse struct {
	Headers map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *QueryFormInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s QueryFormInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s QueryFormInstanceResponse) GoString() string {
	return s.String()
}

func (s *QueryFormInstanceResponse) SetHeaders(v map[string]*string) *QueryFormInstanceResponse {
	s.Headers = v
	return s
}

func (s *QueryFormInstanceResponse) SetBody(v *QueryFormInstanceResponseBody) *QueryFormInstanceResponse {
	s.Body = v
	return s
}

type QueryIntegratedTodoTaskHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s QueryIntegratedTodoTaskHeaders) String() string {
	return tea.Prettify(s)
}

func (s QueryIntegratedTodoTaskHeaders) GoString() string {
	return s.String()
}

func (s *QueryIntegratedTodoTaskHeaders) SetCommonHeaders(v map[string]*string) *QueryIntegratedTodoTaskHeaders {
	s.CommonHeaders = v
	return s
}

func (s *QueryIntegratedTodoTaskHeaders) SetXAcsDingtalkAccessToken(v string) *QueryIntegratedTodoTaskHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type QueryIntegratedTodoTaskRequest struct {
	// 在此时间戳之前创建的
	CreateBefore *int64 `json:"createBefore,omitempty" xml:"createBefore,omitempty"`
	// 第几页，取值范围为 1 ≤ page ≤ 1000
	PageNumber *int32 `json:"pageNumber,omitempty" xml:"pageNumber,omitempty"`
	// 分页大小，取值范围为 1 ≤ pageSize ≤ 40
	PageSize *int32  `json:"pageSize,omitempty" xml:"pageSize,omitempty"`
	UserId   *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s QueryIntegratedTodoTaskRequest) String() string {
	return tea.Prettify(s)
}

func (s QueryIntegratedTodoTaskRequest) GoString() string {
	return s.String()
}

func (s *QueryIntegratedTodoTaskRequest) SetCreateBefore(v int64) *QueryIntegratedTodoTaskRequest {
	s.CreateBefore = &v
	return s
}

func (s *QueryIntegratedTodoTaskRequest) SetPageNumber(v int32) *QueryIntegratedTodoTaskRequest {
	s.PageNumber = &v
	return s
}

func (s *QueryIntegratedTodoTaskRequest) SetPageSize(v int32) *QueryIntegratedTodoTaskRequest {
	s.PageSize = &v
	return s
}

func (s *QueryIntegratedTodoTaskRequest) SetUserId(v string) *QueryIntegratedTodoTaskRequest {
	s.UserId = &v
	return s
}

type QueryIntegratedTodoTaskResponseBody struct {
	Result *QueryIntegratedTodoTaskResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s QueryIntegratedTodoTaskResponseBody) String() string {
	return tea.Prettify(s)
}

func (s QueryIntegratedTodoTaskResponseBody) GoString() string {
	return s.String()
}

func (s *QueryIntegratedTodoTaskResponseBody) SetResult(v *QueryIntegratedTodoTaskResponseBodyResult) *QueryIntegratedTodoTaskResponseBody {
	s.Result = v
	return s
}

type QueryIntegratedTodoTaskResponseBodyResult struct {
	// 是否还有下一页
	HasMore *bool                                            `json:"hasMore,omitempty" xml:"hasMore,omitempty"`
	List    []*QueryIntegratedTodoTaskResponseBodyResultList `json:"list,omitempty" xml:"list,omitempty" type:"Repeated"`
}

func (s QueryIntegratedTodoTaskResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s QueryIntegratedTodoTaskResponseBodyResult) GoString() string {
	return s.String()
}

func (s *QueryIntegratedTodoTaskResponseBodyResult) SetHasMore(v bool) *QueryIntegratedTodoTaskResponseBodyResult {
	s.HasMore = &v
	return s
}

func (s *QueryIntegratedTodoTaskResponseBodyResult) SetList(v []*QueryIntegratedTodoTaskResponseBodyResultList) *QueryIntegratedTodoTaskResponseBodyResult {
	s.List = v
	return s
}

type QueryIntegratedTodoTaskResponseBodyResultList struct {
	// 待办组ID，需要在调用创建流程中心集成任务接口时，主动设置该值。
	ActivityId *string `json:"activityId,omitempty" xml:"activityId,omitempty"`
	// OA审批任务创建时间
	CreateTime *int64 `json:"createTime,omitempty" xml:"createTime,omitempty"`
	// OA审批任务完成时间
	FinishTime *string `json:"finishTime,omitempty" xml:"finishTime,omitempty"`
	// 流程实例ID
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 任务处理结果：agree 或 refuse
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 任务状态
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// OA审批任务ID
	TaskId *int64 `json:"taskId,omitempty" xml:"taskId,omitempty"`
	// OA审批任务执行人的用户ID
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s QueryIntegratedTodoTaskResponseBodyResultList) String() string {
	return tea.Prettify(s)
}

func (s QueryIntegratedTodoTaskResponseBodyResultList) GoString() string {
	return s.String()
}

func (s *QueryIntegratedTodoTaskResponseBodyResultList) SetActivityId(v string) *QueryIntegratedTodoTaskResponseBodyResultList {
	s.ActivityId = &v
	return s
}

func (s *QueryIntegratedTodoTaskResponseBodyResultList) SetCreateTime(v int64) *QueryIntegratedTodoTaskResponseBodyResultList {
	s.CreateTime = &v
	return s
}

func (s *QueryIntegratedTodoTaskResponseBodyResultList) SetFinishTime(v string) *QueryIntegratedTodoTaskResponseBodyResultList {
	s.FinishTime = &v
	return s
}

func (s *QueryIntegratedTodoTaskResponseBodyResultList) SetProcessInstanceId(v string) *QueryIntegratedTodoTaskResponseBodyResultList {
	s.ProcessInstanceId = &v
	return s
}

func (s *QueryIntegratedTodoTaskResponseBodyResultList) SetResult(v string) *QueryIntegratedTodoTaskResponseBodyResultList {
	s.Result = &v
	return s
}

func (s *QueryIntegratedTodoTaskResponseBodyResultList) SetStatus(v string) *QueryIntegratedTodoTaskResponseBodyResultList {
	s.Status = &v
	return s
}

func (s *QueryIntegratedTodoTaskResponseBodyResultList) SetTaskId(v int64) *QueryIntegratedTodoTaskResponseBodyResultList {
	s.TaskId = &v
	return s
}

func (s *QueryIntegratedTodoTaskResponseBodyResultList) SetUserId(v string) *QueryIntegratedTodoTaskResponseBodyResultList {
	s.UserId = &v
	return s
}

type QueryIntegratedTodoTaskResponse struct {
	Headers map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *QueryIntegratedTodoTaskResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s QueryIntegratedTodoTaskResponse) String() string {
	return tea.Prettify(s)
}

func (s QueryIntegratedTodoTaskResponse) GoString() string {
	return s.String()
}

func (s *QueryIntegratedTodoTaskResponse) SetHeaders(v map[string]*string) *QueryIntegratedTodoTaskResponse {
	s.Headers = v
	return s
}

func (s *QueryIntegratedTodoTaskResponse) SetBody(v *QueryIntegratedTodoTaskResponseBody) *QueryIntegratedTodoTaskResponse {
	s.Body = v
	return s
}

type QueryProcessByBizCategoryIdHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s QueryProcessByBizCategoryIdHeaders) String() string {
	return tea.Prettify(s)
}

func (s QueryProcessByBizCategoryIdHeaders) GoString() string {
	return s.String()
}

func (s *QueryProcessByBizCategoryIdHeaders) SetCommonHeaders(v map[string]*string) *QueryProcessByBizCategoryIdHeaders {
	s.CommonHeaders = v
	return s
}

func (s *QueryProcessByBizCategoryIdHeaders) SetXAcsDingtalkAccessToken(v string) *QueryProcessByBizCategoryIdHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type QueryProcessByBizCategoryIdRequest struct {
	// 业务标识
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 用户id
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s QueryProcessByBizCategoryIdRequest) String() string {
	return tea.Prettify(s)
}

func (s QueryProcessByBizCategoryIdRequest) GoString() string {
	return s.String()
}

func (s *QueryProcessByBizCategoryIdRequest) SetBizType(v string) *QueryProcessByBizCategoryIdRequest {
	s.BizType = &v
	return s
}

func (s *QueryProcessByBizCategoryIdRequest) SetUserId(v string) *QueryProcessByBizCategoryIdRequest {
	s.UserId = &v
	return s
}

type QueryProcessByBizCategoryIdResponseBody struct {
	Result []*QueryProcessByBizCategoryIdResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Repeated"`
}

func (s QueryProcessByBizCategoryIdResponseBody) String() string {
	return tea.Prettify(s)
}

func (s QueryProcessByBizCategoryIdResponseBody) GoString() string {
	return s.String()
}

func (s *QueryProcessByBizCategoryIdResponseBody) SetResult(v []*QueryProcessByBizCategoryIdResponseBodyResult) *QueryProcessByBizCategoryIdResponseBody {
	s.Result = v
	return s
}

type QueryProcessByBizCategoryIdResponseBodyResult struct {
	// 模板描述
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// 模板名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 模板code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 模版发布状态。
	//
	// - PUBLISHED：已启用
	//
	// - INVALID：停用
	//
	// - SAVED：已保存
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s QueryProcessByBizCategoryIdResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s QueryProcessByBizCategoryIdResponseBodyResult) GoString() string {
	return s.String()
}

func (s *QueryProcessByBizCategoryIdResponseBodyResult) SetDescription(v string) *QueryProcessByBizCategoryIdResponseBodyResult {
	s.Description = &v
	return s
}

func (s *QueryProcessByBizCategoryIdResponseBodyResult) SetName(v string) *QueryProcessByBizCategoryIdResponseBodyResult {
	s.Name = &v
	return s
}

func (s *QueryProcessByBizCategoryIdResponseBodyResult) SetProcessCode(v string) *QueryProcessByBizCategoryIdResponseBodyResult {
	s.ProcessCode = &v
	return s
}

func (s *QueryProcessByBizCategoryIdResponseBodyResult) SetStatus(v string) *QueryProcessByBizCategoryIdResponseBodyResult {
	s.Status = &v
	return s
}

type QueryProcessByBizCategoryIdResponse struct {
	Headers map[string]*string                       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *QueryProcessByBizCategoryIdResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s QueryProcessByBizCategoryIdResponse) String() string {
	return tea.Prettify(s)
}

func (s QueryProcessByBizCategoryIdResponse) GoString() string {
	return s.String()
}

func (s *QueryProcessByBizCategoryIdResponse) SetHeaders(v map[string]*string) *QueryProcessByBizCategoryIdResponse {
	s.Headers = v
	return s
}

func (s *QueryProcessByBizCategoryIdResponse) SetBody(v *QueryProcessByBizCategoryIdResponseBody) *QueryProcessByBizCategoryIdResponse {
	s.Body = v
	return s
}

type QuerySchemaByProcessCodeHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s QuerySchemaByProcessCodeHeaders) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeHeaders) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeHeaders) SetCommonHeaders(v map[string]*string) *QuerySchemaByProcessCodeHeaders {
	s.CommonHeaders = v
	return s
}

func (s *QuerySchemaByProcessCodeHeaders) SetXAcsDingtalkAccessToken(v string) *QuerySchemaByProcessCodeHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type QuerySchemaByProcessCodeRequest struct {
	// 应用搭建隔离信息
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 表单的唯一码
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s QuerySchemaByProcessCodeRequest) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeRequest) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeRequest) SetAppUuid(v string) *QuerySchemaByProcessCodeRequest {
	s.AppUuid = &v
	return s
}

func (s *QuerySchemaByProcessCodeRequest) SetProcessCode(v string) *QuerySchemaByProcessCodeRequest {
	s.ProcessCode = &v
	return s
}

type QuerySchemaByProcessCodeResponseBody struct {
	// 返回结果详情。
	Result *QuerySchemaByProcessCodeResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s QuerySchemaByProcessCodeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBody) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBody) SetResult(v *QuerySchemaByProcessCodeResponseBodyResult) *QuerySchemaByProcessCodeResponseBody {
	s.Result = v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResult struct {
	// 表单类型。
	AppType *int32 `json:"appType,omitempty" xml:"appType,omitempty"`
	// 表单应用 uuid 或者 corpId。
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 代表表单业务含义的类型。
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 创建人 userId。
	CreatorUserId *string `json:"creatorUserId,omitempty" xml:"creatorUserId,omitempty"`
	// 业务自定义设置数据。
	CustomSetting *string `json:"customSetting,omitempty" xml:"customSetting,omitempty"`
	// 引擎类型，表单：0，页面：1
	EngineType *int32 `json:"engineType,omitempty" xml:"engineType,omitempty"`
	// 表单的唯一码。
	FormCode *string `json:"formCode,omitempty" xml:"formCode,omitempty"`
	// 表单 uuid。
	FormUuid *string `json:"formUuid,omitempty" xml:"formUuid,omitempty"`
	// 创建时间的时间戳。
	GmtCreate *string `json:"gmtCreate,omitempty" xml:"gmtCreate,omitempty"`
	// 修改时间的时间戳。
	GmtModified *string `json:"gmtModified,omitempty" xml:"gmtModified,omitempty"`
	// 图标。
	Icon *string `json:"icon,omitempty" xml:"icon,omitempty"`
	// 排序 id。
	ListOrder *int32 `json:"listOrder,omitempty" xml:"listOrder,omitempty"`
	// 说明文案。
	Memo *string `json:"memo,omitempty" xml:"memo,omitempty"`
	// 表单名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 数据归属者的 id 类型。企业(orgId), 群(cid), 人(uid)。
	OwnerIdType *string `json:"ownerIdType,omitempty" xml:"ownerIdType,omitempty"`
	// 目标类型: inner, outer, customer。
	ProcType *string `json:"procType,omitempty" xml:"procType,omitempty"`
	// 表单 schema 详情。
	SchemaContent *QuerySchemaByProcessCodeResponseBodyResultSchemaContent `json:"schemaContent,omitempty" xml:"schemaContent,omitempty" type:"Struct"`
	// 状态, PUBLISHED(启用), INVALID(停用), SAVED(草稿)
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// 可见范围类型。
	VisibleRange *string `json:"visibleRange,omitempty" xml:"visibleRange,omitempty"`
}

func (s QuerySchemaByProcessCodeResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResult) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetAppType(v int32) *QuerySchemaByProcessCodeResponseBodyResult {
	s.AppType = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetAppUuid(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.AppUuid = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetBizType(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.BizType = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetCreatorUserId(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.CreatorUserId = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetCustomSetting(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.CustomSetting = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetEngineType(v int32) *QuerySchemaByProcessCodeResponseBodyResult {
	s.EngineType = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetFormCode(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.FormCode = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetFormUuid(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.FormUuid = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetGmtCreate(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.GmtCreate = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetGmtModified(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.GmtModified = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetIcon(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.Icon = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetListOrder(v int32) *QuerySchemaByProcessCodeResponseBodyResult {
	s.ListOrder = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetMemo(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.Memo = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetName(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.Name = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetOwnerIdType(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.OwnerIdType = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetProcType(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.ProcType = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetSchemaContent(v *QuerySchemaByProcessCodeResponseBodyResultSchemaContent) *QuerySchemaByProcessCodeResponseBodyResult {
	s.SchemaContent = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetStatus(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.Status = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResult) SetVisibleRange(v string) *QuerySchemaByProcessCodeResponseBodyResult {
	s.VisibleRange = &v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContent struct {
	// 图标
	Icon *string `json:"icon,omitempty" xml:"icon,omitempty"`
	// 控件列表
	Items []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems `json:"items,omitempty" xml:"items,omitempty" type:"Repeated"`
	// 表单名称。
	Title *string `json:"title,omitempty" xml:"title,omitempty"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContent) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContent) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContent) SetIcon(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContent {
	s.Icon = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContent) SetItems(v []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems) *QuerySchemaByProcessCodeResponseBodyResultSchemaContent {
	s.Items = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContent) SetTitle(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContent {
	s.Title = &v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems struct {
	// 子控件列表
	Children []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildren `json:"children,omitempty" xml:"children,omitempty" type:"Repeated"`
	// 控件类型，取值：
	ComponentName *string `json:"componentName,omitempty" xml:"componentName,omitempty"`
	// 控件属性。
	Props *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps `json:"props,omitempty" xml:"props,omitempty" type:"Struct"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems) SetChildren(v []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildren) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems {
	s.Children = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems) SetComponentName(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems {
	s.ComponentName = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems) SetProps(v *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItems {
	s.Props = v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildren struct {
	// 控件类型
	ComponentName *string `json:"componentName,omitempty" xml:"componentName,omitempty"`
	// 子控件属性
	Props *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps `json:"props,omitempty" xml:"props,omitempty" type:"Struct"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildren) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildren) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildren) SetComponentName(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildren {
	s.ComponentName = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildren) SetProps(v *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildren {
	s.Props = v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps struct {
	// 控件业务别名
	BizAlias *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	// 控件id
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 控件名称
	Label *string `json:"label,omitempty" xml:"label,omitempty"`
	// 是否必填
	Required *bool `json:"required,omitempty" xml:"required,omitempty"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps) SetBizAlias(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps {
	s.BizAlias = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps) SetId(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps {
	s.Id = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps) SetLabel(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps {
	s.Label = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps) SetRequired(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsChildrenProps {
	s.Required = &v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps struct {
	// 加班套件4.0新增 加班明细名称。
	ActionName *string `json:"actionName,omitempty" xml:"actionName,omitempty"`
	// textnote的样式，top|middle|bottom。
	Align *string `json:"align,omitempty" xml:"align,omitempty"`
	// ISV 微应用 appId，用于ISV身份权限识别，ISV可获得相应数据。
	AppId *int64 `json:"appId,omitempty" xml:"appId,omitempty"`
	// 套件是否开启异步获取分条件规则，true：开启；false：不开启。
	AsyncCondition *bool `json:"asyncCondition,omitempty" xml:"asyncCondition,omitempty"`
	// 请假、出差、外出、加班类型标签。
	AttendTypeLabel *string `json:"attendTypeLabel,omitempty" xml:"attendTypeLabel,omitempty"`
	// 表单关联控件列表。
	BehaviorLinkage []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkage `json:"behaviorLinkage,omitempty" xml:"behaviorLinkage,omitempty" type:"Repeated"`
	// 控件业务自定义别名。
	BizAlias *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	// 业务套件类型。
	BizType *string `json:"bizType,omitempty" xml:"bizType,omitempty"`
	// 套件内子组件可见性
	ChildFieldVisible map[string]*bool `json:"childFieldVisible,omitempty" xml:"childFieldVisible,omitempty"`
	// 内部联系人choice，1表示多选，0表示单选。
	Choice *int32 `json:"choice,omitempty" xml:"choice,omitempty"`
	// common field的commonBizType。
	CommonBizType *string `json:"commonBizType,omitempty" xml:"commonBizType,omitempty"`
	// 是否可编辑。
	Disabled *bool `json:"disabled,omitempty" xml:"disabled,omitempty"`
	// 是否自动计算时长。
	Duration *bool `json:"duration,omitempty" xml:"duration,omitempty"`
	// 兼容字段。
	DurationLabel *string `json:"durationLabel,omitempty" xml:"durationLabel,omitempty"`
	// e签宝专用标识。
	ESign *bool `json:"eSign,omitempty" xml:"eSign,omitempty"`
	// 套件值是否打平
	Extract *bool `json:"extract,omitempty" xml:"extract,omitempty"`
	// 关联表单中的fields存储
	FieldsInfo *string `json:"fieldsInfo,omitempty" xml:"fieldsInfo,omitempty"`
	// 时间格式(DDDateField和DDDateRangeField)。
	Format *string `json:"format,omitempty" xml:"format,omitempty"`
	// 公式。
	Formula *string `json:"formula,omitempty" xml:"formula,omitempty"`
	// 加班套件4.0新增 加班明细是否隐藏。
	Hidden *bool `json:"hidden,omitempty" xml:"hidden,omitempty"`
	// textnote在详情页是否隐藏，true隐藏， false不隐藏
	HiddenInApprovalDetail *bool `json:"hiddenInApprovalDetail,omitempty" xml:"hiddenInApprovalDetail,omitempty"`
	// 加班套件4.0新增 加班明细是否隐藏标签。
	HideLabel *bool `json:"hideLabel,omitempty" xml:"hideLabel,omitempty"`
	// 兼容出勤套件类型。
	HolidayOptions []map[string]*string `json:"holidayOptions,omitempty" xml:"holidayOptions,omitempty" type:"Repeated"`
	// 控件 id。
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 控件名称。
	Label *string `json:"label,omitempty" xml:"label,omitempty"`
	// label是否可修改 true：不可修改。
	LabelEditableFreeze *bool `json:"labelEditableFreeze,omitempty" xml:"labelEditableFreeze,omitempty"`
	// 说明文案的链接地址。
	Link *string `json:"link,omitempty" xml:"link,omitempty"`
	// 加班套件4.0新增 加班明细描述。
	MainTitle *string `json:"mainTitle,omitempty" xml:"mainTitle,omitempty"`
	// 是否参与打印(1表示不打印, 0表示打印)。
	NotPrint *string `json:"notPrint,omitempty" xml:"notPrint,omitempty"`
	// 是否需要大写 默认是需要; 1:不需要大写, 空或者0:需要大写。
	NotUpper *string `json:"notUpper,omitempty" xml:"notUpper,omitempty"`
	// 选项内容列表，提供给业务方更多的选择器操作。
	ObjOptions []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsObjOptions `json:"objOptions,omitempty" xml:"objOptions,omitempty" type:"Repeated"`
	// 单选框选项列表。
	Options []*string `json:"options,omitempty" xml:"options,omitempty" type:"Repeated"`
	// 是否有支付属性。
	PayEnable *bool `json:"payEnable,omitempty" xml:"payEnable,omitempty"`
	// 占位符。
	Placeholder *string `json:"placeholder,omitempty" xml:"placeholder,omitempty"`
	// 同步到考勤, 表示是否设置为员工状态。
	Push *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush `json:"push,omitempty" xml:"push,omitempty" type:"Struct"`
	// 推送到考勤, 子类型(DDSelectField)。
	PushToAttendance *bool `json:"pushToAttendance,omitempty" xml:"pushToAttendance,omitempty"`
	// 是否推送管理日历(DDDateRangeField, 1表示推送, 0表示不推送, 该属性为兼容保留)。
	PushToCalendar *int32 `json:"pushToCalendar,omitempty" xml:"pushToCalendar,omitempty"`
	// 是否必填。
	Required *bool `json:"required,omitempty" xml:"required,omitempty"`
	// 必填是否可修改 true：不可修改。
	RequiredEditableFreeze *bool `json:"requiredEditableFreeze,omitempty" xml:"requiredEditableFreeze,omitempty"`
	// 兼容出勤套件类型。
	ShowAttendOptions *bool `json:"showAttendOptions,omitempty" xml:"showAttendOptions,omitempty"`
	// 是否开启员工状态。
	StaffStatusEnabled *bool `json:"staffStatusEnabled,omitempty" xml:"staffStatusEnabled,omitempty"`
	// 需要计算总和的明细组件
	StatField []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField `json:"statField,omitempty" xml:"statField,omitempty" type:"Repeated"`
	// 数字组件/日期区间组件单位属性。
	Unit *string `json:"unit,omitempty" xml:"unit,omitempty"`
	// 是否使用考勤日历。
	UseCalendar *bool `json:"useCalendar,omitempty" xml:"useCalendar,omitempty"`
	// 明细打印排版方式 false：横向 true：纵向。
	VerticalPrint *bool `json:"verticalPrint,omitempty" xml:"verticalPrint,omitempty"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetActionName(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.ActionName = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetAlign(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Align = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetAppId(v int64) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.AppId = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetAsyncCondition(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.AsyncCondition = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetAttendTypeLabel(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.AttendTypeLabel = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetBehaviorLinkage(v []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkage) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.BehaviorLinkage = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetBizAlias(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.BizAlias = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetBizType(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.BizType = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetChildFieldVisible(v map[string]*bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.ChildFieldVisible = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetChoice(v int32) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Choice = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetCommonBizType(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.CommonBizType = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetDisabled(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Disabled = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetDuration(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Duration = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetDurationLabel(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.DurationLabel = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetESign(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.ESign = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetExtract(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Extract = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetFieldsInfo(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.FieldsInfo = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetFormat(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Format = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetFormula(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Formula = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetHidden(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Hidden = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetHiddenInApprovalDetail(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.HiddenInApprovalDetail = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetHideLabel(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.HideLabel = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetHolidayOptions(v []map[string]*string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.HolidayOptions = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetId(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Id = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetLabel(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Label = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetLabelEditableFreeze(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.LabelEditableFreeze = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetLink(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Link = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetMainTitle(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.MainTitle = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetNotPrint(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.NotPrint = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetNotUpper(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.NotUpper = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetObjOptions(v []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsObjOptions) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.ObjOptions = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetOptions(v []*string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Options = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetPayEnable(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.PayEnable = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetPlaceholder(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Placeholder = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetPush(v *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Push = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetPushToAttendance(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.PushToAttendance = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetPushToCalendar(v int32) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.PushToCalendar = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetRequired(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Required = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetRequiredEditableFreeze(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.RequiredEditableFreeze = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetShowAttendOptions(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.ShowAttendOptions = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetStaffStatusEnabled(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.StaffStatusEnabled = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetStatField(v []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.StatField = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetUnit(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.Unit = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetUseCalendar(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.UseCalendar = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps) SetVerticalPrint(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsProps {
	s.VerticalPrint = &v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkage struct {
	// 关联控件列表。
	Targets []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkageTargets `json:"targets,omitempty" xml:"targets,omitempty" type:"Repeated"`
	// 控件值。
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkage) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkage) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkage) SetTargets(v []*QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkageTargets) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkage {
	s.Targets = v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkage) SetValue(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkage {
	s.Value = &v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkageTargets struct {
	// 行为。
	Behavior *string `json:"behavior,omitempty" xml:"behavior,omitempty"`
	// 字段 id。
	FieldId *string `json:"fieldId,omitempty" xml:"fieldId,omitempty"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkageTargets) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkageTargets) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkageTargets) SetBehavior(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkageTargets {
	s.Behavior = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkageTargets) SetFieldId(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsBehaviorLinkageTargets {
	s.FieldId = &v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsObjOptions struct {
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsObjOptions) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsObjOptions) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsObjOptions) SetValue(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsObjOptions {
	s.Value = &v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush struct {
	// 考勤类型(1表示请假, 2表示出差, 3表示加班, 4表示外出)
	AttendanceRule *int32 `json:"attendanceRule,omitempty" xml:"attendanceRule,omitempty"`
	// 开启状态(1表示开启, 0表示关闭)
	PushSwitch *int32 `json:"pushSwitch,omitempty" xml:"pushSwitch,omitempty"`
	// 状态显示名称
	PushTag *string `json:"pushTag,omitempty" xml:"pushTag,omitempty"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush) SetAttendanceRule(v int32) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush {
	s.AttendanceRule = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush) SetPushSwitch(v int32) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush {
	s.PushSwitch = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush) SetPushTag(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsPush {
	s.PushTag = &v
	return s
}

type QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField struct {
	// id 值。
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 名称。
	Label *string `json:"label,omitempty" xml:"label,omitempty"`
	// 单位。
	Unit *string `json:"unit,omitempty" xml:"unit,omitempty"`
	// 大写。
	Upper *bool `json:"upper,omitempty" xml:"upper,omitempty"`
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField) SetId(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField {
	s.Id = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField) SetLabel(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField {
	s.Label = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField) SetUnit(v string) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField {
	s.Unit = &v
	return s
}

func (s *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField) SetUpper(v bool) *QuerySchemaByProcessCodeResponseBodyResultSchemaContentItemsPropsStatField {
	s.Upper = &v
	return s
}

type QuerySchemaByProcessCodeResponse struct {
	Headers map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *QuerySchemaByProcessCodeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s QuerySchemaByProcessCodeResponse) String() string {
	return tea.Prettify(s)
}

func (s QuerySchemaByProcessCodeResponse) GoString() string {
	return s.String()
}

func (s *QuerySchemaByProcessCodeResponse) SetHeaders(v map[string]*string) *QuerySchemaByProcessCodeResponse {
	s.Headers = v
	return s
}

func (s *QuerySchemaByProcessCodeResponse) SetBody(v *QuerySchemaByProcessCodeResponseBody) *QuerySchemaByProcessCodeResponse {
	s.Body = v
	return s
}

type RedirectWorkflowTaskHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s RedirectWorkflowTaskHeaders) String() string {
	return tea.Prettify(s)
}

func (s RedirectWorkflowTaskHeaders) GoString() string {
	return s.String()
}

func (s *RedirectWorkflowTaskHeaders) SetCommonHeaders(v map[string]*string) *RedirectWorkflowTaskHeaders {
	s.CommonHeaders = v
	return s
}

func (s *RedirectWorkflowTaskHeaders) SetXAcsDingtalkAccessToken(v string) *RedirectWorkflowTaskHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type RedirectWorkflowTaskRequest struct {
	// 操作节点名
	ActionName *string `json:"actionName,omitempty" xml:"actionName,omitempty"`
	// 操作人的用户ID，需要跟任务的当前执行人保持一致，否则无法通过校验
	OperateUserId *string `json:"operateUserId,omitempty" xml:"operateUserId,omitempty"`
	// 转交备注信息
	Remark *string `json:"remark,omitempty" xml:"remark,omitempty"`
	// OA审批任务ID
	TaskId *int64 `json:"taskId,omitempty" xml:"taskId,omitempty"`
	// OA审批任务被转交对象的用户ID
	ToUserId *string `json:"toUserId,omitempty" xml:"toUserId,omitempty"`
}

func (s RedirectWorkflowTaskRequest) String() string {
	return tea.Prettify(s)
}

func (s RedirectWorkflowTaskRequest) GoString() string {
	return s.String()
}

func (s *RedirectWorkflowTaskRequest) SetActionName(v string) *RedirectWorkflowTaskRequest {
	s.ActionName = &v
	return s
}

func (s *RedirectWorkflowTaskRequest) SetOperateUserId(v string) *RedirectWorkflowTaskRequest {
	s.OperateUserId = &v
	return s
}

func (s *RedirectWorkflowTaskRequest) SetRemark(v string) *RedirectWorkflowTaskRequest {
	s.Remark = &v
	return s
}

func (s *RedirectWorkflowTaskRequest) SetTaskId(v int64) *RedirectWorkflowTaskRequest {
	s.TaskId = &v
	return s
}

func (s *RedirectWorkflowTaskRequest) SetToUserId(v string) *RedirectWorkflowTaskRequest {
	s.ToUserId = &v
	return s
}

type RedirectWorkflowTaskResponseBody struct {
	// 是否转交成功
	Result *bool `json:"result,omitempty" xml:"result,omitempty"`
}

func (s RedirectWorkflowTaskResponseBody) String() string {
	return tea.Prettify(s)
}

func (s RedirectWorkflowTaskResponseBody) GoString() string {
	return s.String()
}

func (s *RedirectWorkflowTaskResponseBody) SetResult(v bool) *RedirectWorkflowTaskResponseBody {
	s.Result = &v
	return s
}

type RedirectWorkflowTaskResponse struct {
	Headers map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *RedirectWorkflowTaskResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RedirectWorkflowTaskResponse) String() string {
	return tea.Prettify(s)
}

func (s RedirectWorkflowTaskResponse) GoString() string {
	return s.String()
}

func (s *RedirectWorkflowTaskResponse) SetHeaders(v map[string]*string) *RedirectWorkflowTaskResponse {
	s.Headers = v
	return s
}

func (s *RedirectWorkflowTaskResponse) SetBody(v *RedirectWorkflowTaskResponseBody) *RedirectWorkflowTaskResponse {
	s.Body = v
	return s
}

type SaveIntegratedInstanceHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s SaveIntegratedInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s SaveIntegratedInstanceHeaders) GoString() string {
	return s.String()
}

func (s *SaveIntegratedInstanceHeaders) SetCommonHeaders(v map[string]*string) *SaveIntegratedInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *SaveIntegratedInstanceHeaders) SetXAcsDingtalkAccessToken(v string) *SaveIntegratedInstanceHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type SaveIntegratedInstanceRequest struct {
	FormComponentValueList []*SaveIntegratedInstanceRequestFormComponentValueList `json:"formComponentValueList,omitempty" xml:"formComponentValueList,omitempty" type:"Repeated"`
	Notifiers              []*SaveIntegratedInstanceRequestNotifiers              `json:"notifiers,omitempty" xml:"notifiers,omitempty" type:"Repeated"`
	// 审批实例接收人的userId。
	OriginatorUserId *string `json:"originatorUserId,omitempty" xml:"originatorUserId,omitempty"`
	// 审批模板唯一码
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 实例标题
	Title *string `json:"title,omitempty" xml:"title,omitempty"`
	// 三方审批系统中审批单的详情页地址
	Url *string `json:"url,omitempty" xml:"url,omitempty"`
}

func (s SaveIntegratedInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s SaveIntegratedInstanceRequest) GoString() string {
	return s.String()
}

func (s *SaveIntegratedInstanceRequest) SetFormComponentValueList(v []*SaveIntegratedInstanceRequestFormComponentValueList) *SaveIntegratedInstanceRequest {
	s.FormComponentValueList = v
	return s
}

func (s *SaveIntegratedInstanceRequest) SetNotifiers(v []*SaveIntegratedInstanceRequestNotifiers) *SaveIntegratedInstanceRequest {
	s.Notifiers = v
	return s
}

func (s *SaveIntegratedInstanceRequest) SetOriginatorUserId(v string) *SaveIntegratedInstanceRequest {
	s.OriginatorUserId = &v
	return s
}

func (s *SaveIntegratedInstanceRequest) SetProcessCode(v string) *SaveIntegratedInstanceRequest {
	s.ProcessCode = &v
	return s
}

func (s *SaveIntegratedInstanceRequest) SetTitle(v string) *SaveIntegratedInstanceRequest {
	s.Title = &v
	return s
}

func (s *SaveIntegratedInstanceRequest) SetUrl(v string) *SaveIntegratedInstanceRequest {
	s.Url = &v
	return s
}

type SaveIntegratedInstanceRequestFormComponentValueList struct {
	// 控件别名
	BizAlias *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	// 控件类型，取值：
	//
	// TextField：单行输入框
	//
	// TextareaField：多行输入框
	//
	// NumberField：数字输入框
	//
	// DDSelectField：单选框
	//
	// DDMultiSelectField：多选框
	//
	// DDDateField：日期控件
	//
	// DDDateRangeField：时间区间控件
	//
	// TextNote：文字说明控件
	//
	// PhoneField：电话控件
	//
	// DDPhotoField：图片控件
	//
	// MoneyField：金额控件
	//
	// TableField：明细控件
	//
	// DDAttachment：附件
	//
	// InnerContactField：联系人控件
	//
	// RelateField：关联审批单
	//
	// FormRelateField：关联控件
	//
	// AddressField：省市区控件
	//
	// StarRatingField：评分控件
	//
	// DepartmentField：部门控件
	ComponentType *string `json:"componentType,omitempty" xml:"componentType,omitempty"`
	// 表单扩展值
	ExtValue *string `json:"extValue,omitempty" xml:"extValue,omitempty"`
	// 控件id
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 表单名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 表单值
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s SaveIntegratedInstanceRequestFormComponentValueList) String() string {
	return tea.Prettify(s)
}

func (s SaveIntegratedInstanceRequestFormComponentValueList) GoString() string {
	return s.String()
}

func (s *SaveIntegratedInstanceRequestFormComponentValueList) SetBizAlias(v string) *SaveIntegratedInstanceRequestFormComponentValueList {
	s.BizAlias = &v
	return s
}

func (s *SaveIntegratedInstanceRequestFormComponentValueList) SetComponentType(v string) *SaveIntegratedInstanceRequestFormComponentValueList {
	s.ComponentType = &v
	return s
}

func (s *SaveIntegratedInstanceRequestFormComponentValueList) SetExtValue(v string) *SaveIntegratedInstanceRequestFormComponentValueList {
	s.ExtValue = &v
	return s
}

func (s *SaveIntegratedInstanceRequestFormComponentValueList) SetId(v string) *SaveIntegratedInstanceRequestFormComponentValueList {
	s.Id = &v
	return s
}

func (s *SaveIntegratedInstanceRequestFormComponentValueList) SetName(v string) *SaveIntegratedInstanceRequestFormComponentValueList {
	s.Name = &v
	return s
}

func (s *SaveIntegratedInstanceRequestFormComponentValueList) SetValue(v string) *SaveIntegratedInstanceRequestFormComponentValueList {
	s.Value = &v
	return s
}

type SaveIntegratedInstanceRequestNotifiers struct {
	// 抄送位置，可以值有：
	// start - 审批发起时，通知抄送人
	// finish - 审批通过后，通知抄送人
	// start_finish - 审批发起时和审批通过后，都通知抄送人
	Position *string `json:"position,omitempty" xml:"position,omitempty"`
	// 抄送接收人用户ID
	Userid *string `json:"userid,omitempty" xml:"userid,omitempty"`
}

func (s SaveIntegratedInstanceRequestNotifiers) String() string {
	return tea.Prettify(s)
}

func (s SaveIntegratedInstanceRequestNotifiers) GoString() string {
	return s.String()
}

func (s *SaveIntegratedInstanceRequestNotifiers) SetPosition(v string) *SaveIntegratedInstanceRequestNotifiers {
	s.Position = &v
	return s
}

func (s *SaveIntegratedInstanceRequestNotifiers) SetUserid(v string) *SaveIntegratedInstanceRequestNotifiers {
	s.Userid = &v
	return s
}

type SaveIntegratedInstanceResponseBody struct {
	Result *SaveIntegratedInstanceResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s SaveIntegratedInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s SaveIntegratedInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *SaveIntegratedInstanceResponseBody) SetResult(v *SaveIntegratedInstanceResponseBodyResult) *SaveIntegratedInstanceResponseBody {
	s.Result = v
	return s
}

type SaveIntegratedInstanceResponseBodyResult struct {
	// 实例id
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
}

func (s SaveIntegratedInstanceResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s SaveIntegratedInstanceResponseBodyResult) GoString() string {
	return s.String()
}

func (s *SaveIntegratedInstanceResponseBodyResult) SetProcessInstanceId(v string) *SaveIntegratedInstanceResponseBodyResult {
	s.ProcessInstanceId = &v
	return s
}

type SaveIntegratedInstanceResponse struct {
	Headers map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *SaveIntegratedInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s SaveIntegratedInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s SaveIntegratedInstanceResponse) GoString() string {
	return s.String()
}

func (s *SaveIntegratedInstanceResponse) SetHeaders(v map[string]*string) *SaveIntegratedInstanceResponse {
	s.Headers = v
	return s
}

func (s *SaveIntegratedInstanceResponse) SetBody(v *SaveIntegratedInstanceResponseBody) *SaveIntegratedInstanceResponse {
	s.Body = v
	return s
}

type SaveProcessHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s SaveProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s SaveProcessHeaders) GoString() string {
	return s.String()
}

func (s *SaveProcessHeaders) SetCommonHeaders(v map[string]*string) *SaveProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *SaveProcessHeaders) SetXAcsDingtalkAccessToken(v string) *SaveProcessHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type SaveProcessRequest struct {
	// 表单模板描述
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// 表单控件列表
	FormComponents []*FormComponent `json:"formComponents,omitempty" xml:"formComponents,omitempty" type:"Repeated"`
	// 表单模板名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 模板code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 流程中心集成配置
	ProcessFeatureConfig *SaveProcessRequestProcessFeatureConfig `json:"processFeatureConfig,omitempty" xml:"processFeatureConfig,omitempty" type:"Struct"`
}

func (s SaveProcessRequest) String() string {
	return tea.Prettify(s)
}

func (s SaveProcessRequest) GoString() string {
	return s.String()
}

func (s *SaveProcessRequest) SetDescription(v string) *SaveProcessRequest {
	s.Description = &v
	return s
}

func (s *SaveProcessRequest) SetFormComponents(v []*FormComponent) *SaveProcessRequest {
	s.FormComponents = v
	return s
}

func (s *SaveProcessRequest) SetName(v string) *SaveProcessRequest {
	s.Name = &v
	return s
}

func (s *SaveProcessRequest) SetProcessCode(v string) *SaveProcessRequest {
	s.ProcessCode = &v
	return s
}

func (s *SaveProcessRequest) SetProcessFeatureConfig(v *SaveProcessRequestProcessFeatureConfig) *SaveProcessRequest {
	s.ProcessFeatureConfig = v
	return s
}

type SaveProcessRequestProcessFeatureConfig struct {
	// 配置列表
	Features []*SaveProcessRequestProcessFeatureConfigFeatures `json:"features,omitempty" xml:"features,omitempty" type:"Repeated"`
}

func (s SaveProcessRequestProcessFeatureConfig) String() string {
	return tea.Prettify(s)
}

func (s SaveProcessRequestProcessFeatureConfig) GoString() string {
	return s.String()
}

func (s *SaveProcessRequestProcessFeatureConfig) SetFeatures(v []*SaveProcessRequestProcessFeatureConfigFeatures) *SaveProcessRequestProcessFeatureConfig {
	s.Features = v
	return s
}

type SaveProcessRequestProcessFeatureConfigFeatures struct {
	Callback *SaveProcessRequestProcessFeatureConfigFeaturesCallback `json:"callback,omitempty" xml:"callback,omitempty" type:"Struct"`
	// 手机端链接
	MobileUrl *string `json:"mobileUrl,omitempty" xml:"mobileUrl,omitempty"`
	// 名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// pc端链接
	PcUrl *string `json:"pcUrl,omitempty" xml:"pcUrl,omitempty"`
	// 运行方式：
	// ORIGIN：原生运行，即在官方审批内运行对应功能；
	// REDIRECT：外部跳转运行，需要跳转到三方地址运行对应功能
	RunType *string `json:"runType,omitempty" xml:"runType,omitempty"`
}

func (s SaveProcessRequestProcessFeatureConfigFeatures) String() string {
	return tea.Prettify(s)
}

func (s SaveProcessRequestProcessFeatureConfigFeatures) GoString() string {
	return s.String()
}

func (s *SaveProcessRequestProcessFeatureConfigFeatures) SetCallback(v *SaveProcessRequestProcessFeatureConfigFeaturesCallback) *SaveProcessRequestProcessFeatureConfigFeatures {
	s.Callback = v
	return s
}

func (s *SaveProcessRequestProcessFeatureConfigFeatures) SetMobileUrl(v string) *SaveProcessRequestProcessFeatureConfigFeatures {
	s.MobileUrl = &v
	return s
}

func (s *SaveProcessRequestProcessFeatureConfigFeatures) SetName(v string) *SaveProcessRequestProcessFeatureConfigFeatures {
	s.Name = &v
	return s
}

func (s *SaveProcessRequestProcessFeatureConfigFeatures) SetPcUrl(v string) *SaveProcessRequestProcessFeatureConfigFeatures {
	s.PcUrl = &v
	return s
}

func (s *SaveProcessRequestProcessFeatureConfigFeatures) SetRunType(v string) *SaveProcessRequestProcessFeatureConfigFeatures {
	s.RunType = &v
	return s
}

type SaveProcessRequestProcessFeatureConfigFeaturesCallback struct {
	// 网关接口标识
	ApiKey *string `json:"apiKey,omitempty" xml:"apiKey,omitempty"`
	// 网关接口对应应用的uuid
	AppUuid *string `json:"appUuid,omitempty" xml:"appUuid,omitempty"`
	// 网关接口版本
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s SaveProcessRequestProcessFeatureConfigFeaturesCallback) String() string {
	return tea.Prettify(s)
}

func (s SaveProcessRequestProcessFeatureConfigFeaturesCallback) GoString() string {
	return s.String()
}

func (s *SaveProcessRequestProcessFeatureConfigFeaturesCallback) SetApiKey(v string) *SaveProcessRequestProcessFeatureConfigFeaturesCallback {
	s.ApiKey = &v
	return s
}

func (s *SaveProcessRequestProcessFeatureConfigFeaturesCallback) SetAppUuid(v string) *SaveProcessRequestProcessFeatureConfigFeaturesCallback {
	s.AppUuid = &v
	return s
}

func (s *SaveProcessRequestProcessFeatureConfigFeaturesCallback) SetVersion(v string) *SaveProcessRequestProcessFeatureConfigFeaturesCallback {
	s.Version = &v
	return s
}

type SaveProcessResponseBody struct {
	// 表单模板信息
	Result *SaveProcessResponseBodyResult `json:"result,omitempty" xml:"result,omitempty" type:"Struct"`
}

func (s SaveProcessResponseBody) String() string {
	return tea.Prettify(s)
}

func (s SaveProcessResponseBody) GoString() string {
	return s.String()
}

func (s *SaveProcessResponseBody) SetResult(v *SaveProcessResponseBodyResult) *SaveProcessResponseBody {
	s.Result = v
	return s
}

type SaveProcessResponseBodyResult struct {
	// 保存或更新的表单code
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
}

func (s SaveProcessResponseBodyResult) String() string {
	return tea.Prettify(s)
}

func (s SaveProcessResponseBodyResult) GoString() string {
	return s.String()
}

func (s *SaveProcessResponseBodyResult) SetProcessCode(v string) *SaveProcessResponseBodyResult {
	s.ProcessCode = &v
	return s
}

type SaveProcessResponse struct {
	Headers map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *SaveProcessResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s SaveProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s SaveProcessResponse) GoString() string {
	return s.String()
}

func (s *SaveProcessResponse) SetHeaders(v map[string]*string) *SaveProcessResponse {
	s.Headers = v
	return s
}

func (s *SaveProcessResponse) SetBody(v *SaveProcessResponseBody) *SaveProcessResponse {
	s.Body = v
	return s
}

type StartProcessInstanceHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s StartProcessInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s StartProcessInstanceHeaders) GoString() string {
	return s.String()
}

func (s *StartProcessInstanceHeaders) SetCommonHeaders(v map[string]*string) *StartProcessInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *StartProcessInstanceHeaders) SetXAcsDingtalkAccessToken(v string) *StartProcessInstanceHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type StartProcessInstanceRequest struct {
	// 不使用审批流模板时，直接指定审批人列表
	Approvers []*StartProcessInstanceRequestApprovers `json:"approvers,omitempty" xml:"approvers,omitempty" type:"Repeated"`
	// 抄送人userId列表
	CcList []*string `json:"ccList,omitempty" xml:"ccList,omitempty" type:"Repeated"`
	// 抄送时间
	CcPosition *string `json:"ccPosition,omitempty" xml:"ccPosition,omitempty"`
	// 部门ID
	DeptId *int64 `json:"deptId,omitempty" xml:"deptId,omitempty"`
	// 表单数据内容，控件列表
	FormComponentValues []*StartProcessInstanceRequestFormComponentValues `json:"formComponentValues,omitempty" xml:"formComponentValues,omitempty" type:"Repeated"`
	// 企业微应用标识
	MicroappAgentId *int64 `json:"microappAgentId,omitempty" xml:"microappAgentId,omitempty"`
	// 审批发起人的userId
	OriginatorUserId *string `json:"originatorUserId,omitempty" xml:"originatorUserId,omitempty"`
	// 审批流的唯一码
	ProcessCode *string `json:"processCode,omitempty" xml:"processCode,omitempty"`
	// 使用审批流模板时，模板上的自选操作人列表
	TargetSelectActioners []*StartProcessInstanceRequestTargetSelectActioners `json:"targetSelectActioners,omitempty" xml:"targetSelectActioners,omitempty" type:"Repeated"`
}

func (s StartProcessInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s StartProcessInstanceRequest) GoString() string {
	return s.String()
}

func (s *StartProcessInstanceRequest) SetApprovers(v []*StartProcessInstanceRequestApprovers) *StartProcessInstanceRequest {
	s.Approvers = v
	return s
}

func (s *StartProcessInstanceRequest) SetCcList(v []*string) *StartProcessInstanceRequest {
	s.CcList = v
	return s
}

func (s *StartProcessInstanceRequest) SetCcPosition(v string) *StartProcessInstanceRequest {
	s.CcPosition = &v
	return s
}

func (s *StartProcessInstanceRequest) SetDeptId(v int64) *StartProcessInstanceRequest {
	s.DeptId = &v
	return s
}

func (s *StartProcessInstanceRequest) SetFormComponentValues(v []*StartProcessInstanceRequestFormComponentValues) *StartProcessInstanceRequest {
	s.FormComponentValues = v
	return s
}

func (s *StartProcessInstanceRequest) SetMicroappAgentId(v int64) *StartProcessInstanceRequest {
	s.MicroappAgentId = &v
	return s
}

func (s *StartProcessInstanceRequest) SetOriginatorUserId(v string) *StartProcessInstanceRequest {
	s.OriginatorUserId = &v
	return s
}

func (s *StartProcessInstanceRequest) SetProcessCode(v string) *StartProcessInstanceRequest {
	s.ProcessCode = &v
	return s
}

func (s *StartProcessInstanceRequest) SetTargetSelectActioners(v []*StartProcessInstanceRequestTargetSelectActioners) *StartProcessInstanceRequest {
	s.TargetSelectActioners = v
	return s
}

type StartProcessInstanceRequestApprovers struct {
	// 审批类型
	ActionType *string `json:"actionType,omitempty" xml:"actionType,omitempty"`
	// 审批人列表
	UserIds []*string `json:"userIds,omitempty" xml:"userIds,omitempty" type:"Repeated"`
}

func (s StartProcessInstanceRequestApprovers) String() string {
	return tea.Prettify(s)
}

func (s StartProcessInstanceRequestApprovers) GoString() string {
	return s.String()
}

func (s *StartProcessInstanceRequestApprovers) SetActionType(v string) *StartProcessInstanceRequestApprovers {
	s.ActionType = &v
	return s
}

func (s *StartProcessInstanceRequestApprovers) SetUserIds(v []*string) *StartProcessInstanceRequestApprovers {
	s.UserIds = v
	return s
}

type StartProcessInstanceRequestFormComponentValues struct {
	// 控件别名
	BizAlias *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	// 控件类型
	ComponentType *string                                                  `json:"componentType,omitempty" xml:"componentType,omitempty"`
	Details       []*StartProcessInstanceRequestFormComponentValuesDetails `json:"details,omitempty" xml:"details,omitempty" type:"Repeated"`
	// 控件扩展值
	ExtValue *string `json:"extValue,omitempty" xml:"extValue,omitempty"`
	// 控件id
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 控件名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 控件值
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s StartProcessInstanceRequestFormComponentValues) String() string {
	return tea.Prettify(s)
}

func (s StartProcessInstanceRequestFormComponentValues) GoString() string {
	return s.String()
}

func (s *StartProcessInstanceRequestFormComponentValues) SetBizAlias(v string) *StartProcessInstanceRequestFormComponentValues {
	s.BizAlias = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValues) SetComponentType(v string) *StartProcessInstanceRequestFormComponentValues {
	s.ComponentType = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValues) SetDetails(v []*StartProcessInstanceRequestFormComponentValuesDetails) *StartProcessInstanceRequestFormComponentValues {
	s.Details = v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValues) SetExtValue(v string) *StartProcessInstanceRequestFormComponentValues {
	s.ExtValue = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValues) SetId(v string) *StartProcessInstanceRequestFormComponentValues {
	s.Id = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValues) SetName(v string) *StartProcessInstanceRequestFormComponentValues {
	s.Name = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValues) SetValue(v string) *StartProcessInstanceRequestFormComponentValues {
	s.Value = &v
	return s
}

type StartProcessInstanceRequestFormComponentValuesDetails struct {
	// 控件别名
	BizAlias *string                                                         `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	Details  []*StartProcessInstanceRequestFormComponentValuesDetailsDetails `json:"details,omitempty" xml:"details,omitempty" type:"Repeated"`
	// 控件扩展值
	ExtValue *string `json:"extValue,omitempty" xml:"extValue,omitempty"`
	// 控件id
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 控件名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 控件值
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s StartProcessInstanceRequestFormComponentValuesDetails) String() string {
	return tea.Prettify(s)
}

func (s StartProcessInstanceRequestFormComponentValuesDetails) GoString() string {
	return s.String()
}

func (s *StartProcessInstanceRequestFormComponentValuesDetails) SetBizAlias(v string) *StartProcessInstanceRequestFormComponentValuesDetails {
	s.BizAlias = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetails) SetDetails(v []*StartProcessInstanceRequestFormComponentValuesDetailsDetails) *StartProcessInstanceRequestFormComponentValuesDetails {
	s.Details = v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetails) SetExtValue(v string) *StartProcessInstanceRequestFormComponentValuesDetails {
	s.ExtValue = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetails) SetId(v string) *StartProcessInstanceRequestFormComponentValuesDetails {
	s.Id = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetails) SetName(v string) *StartProcessInstanceRequestFormComponentValuesDetails {
	s.Name = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetails) SetValue(v string) *StartProcessInstanceRequestFormComponentValuesDetails {
	s.Value = &v
	return s
}

type StartProcessInstanceRequestFormComponentValuesDetailsDetails struct {
	// 控件别名
	BizAlias *string `json:"bizAlias,omitempty" xml:"bizAlias,omitempty"`
	// 控件类型
	ComponentType *string `json:"componentType,omitempty" xml:"componentType,omitempty"`
	// 控件扩展值
	ExtValue *string `json:"extValue,omitempty" xml:"extValue,omitempty"`
	// 控件id
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 控件名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 控件值
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s StartProcessInstanceRequestFormComponentValuesDetailsDetails) String() string {
	return tea.Prettify(s)
}

func (s StartProcessInstanceRequestFormComponentValuesDetailsDetails) GoString() string {
	return s.String()
}

func (s *StartProcessInstanceRequestFormComponentValuesDetailsDetails) SetBizAlias(v string) *StartProcessInstanceRequestFormComponentValuesDetailsDetails {
	s.BizAlias = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetailsDetails) SetComponentType(v string) *StartProcessInstanceRequestFormComponentValuesDetailsDetails {
	s.ComponentType = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetailsDetails) SetExtValue(v string) *StartProcessInstanceRequestFormComponentValuesDetailsDetails {
	s.ExtValue = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetailsDetails) SetId(v string) *StartProcessInstanceRequestFormComponentValuesDetailsDetails {
	s.Id = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetailsDetails) SetName(v string) *StartProcessInstanceRequestFormComponentValuesDetailsDetails {
	s.Name = &v
	return s
}

func (s *StartProcessInstanceRequestFormComponentValuesDetailsDetails) SetValue(v string) *StartProcessInstanceRequestFormComponentValuesDetailsDetails {
	s.Value = &v
	return s
}

type StartProcessInstanceRequestTargetSelectActioners struct {
	// 自选节点的规则key
	ActionerKey *string `json:"actionerKey,omitempty" xml:"actionerKey,omitempty"`
	// 操作人userId列表
	ActionerUserIds []*string `json:"actionerUserIds,omitempty" xml:"actionerUserIds,omitempty" type:"Repeated"`
}

func (s StartProcessInstanceRequestTargetSelectActioners) String() string {
	return tea.Prettify(s)
}

func (s StartProcessInstanceRequestTargetSelectActioners) GoString() string {
	return s.String()
}

func (s *StartProcessInstanceRequestTargetSelectActioners) SetActionerKey(v string) *StartProcessInstanceRequestTargetSelectActioners {
	s.ActionerKey = &v
	return s
}

func (s *StartProcessInstanceRequestTargetSelectActioners) SetActionerUserIds(v []*string) *StartProcessInstanceRequestTargetSelectActioners {
	s.ActionerUserIds = v
	return s
}

type StartProcessInstanceResponseBody struct {
	// 审批实例id
	InstanceId *string `json:"instanceId,omitempty" xml:"instanceId,omitempty"`
}

func (s StartProcessInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s StartProcessInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *StartProcessInstanceResponseBody) SetInstanceId(v string) *StartProcessInstanceResponseBody {
	s.InstanceId = &v
	return s
}

type StartProcessInstanceResponse struct {
	Headers map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *StartProcessInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s StartProcessInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s StartProcessInstanceResponse) GoString() string {
	return s.String()
}

func (s *StartProcessInstanceResponse) SetHeaders(v map[string]*string) *StartProcessInstanceResponse {
	s.Headers = v
	return s
}

func (s *StartProcessInstanceResponse) SetBody(v *StartProcessInstanceResponseBody) *StartProcessInstanceResponse {
	s.Body = v
	return s
}

type TerminateProcessInstanceHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s TerminateProcessInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s TerminateProcessInstanceHeaders) GoString() string {
	return s.String()
}

func (s *TerminateProcessInstanceHeaders) SetCommonHeaders(v map[string]*string) *TerminateProcessInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *TerminateProcessInstanceHeaders) SetXAcsDingtalkAccessToken(v string) *TerminateProcessInstanceHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type TerminateProcessInstanceRequest struct {
	// 是否通过系统操作：
	//
	// true：由系统直接终止
	//
	// false：由指定的操作者终止
	IsSystem *bool `json:"isSystem,omitempty" xml:"isSystem,omitempty"`
	// 操作人的userid。
	//
	// 当is_system为false时，该参数必传。
	OperatingUserId *string `json:"operatingUserId,omitempty" xml:"operatingUserId,omitempty"`
	// 审批实例ID，可通过调用获取审批实例ID列表接口获取。
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 终止说明。
	Remark *string `json:"remark,omitempty" xml:"remark,omitempty"`
}

func (s TerminateProcessInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s TerminateProcessInstanceRequest) GoString() string {
	return s.String()
}

func (s *TerminateProcessInstanceRequest) SetIsSystem(v bool) *TerminateProcessInstanceRequest {
	s.IsSystem = &v
	return s
}

func (s *TerminateProcessInstanceRequest) SetOperatingUserId(v string) *TerminateProcessInstanceRequest {
	s.OperatingUserId = &v
	return s
}

func (s *TerminateProcessInstanceRequest) SetProcessInstanceId(v string) *TerminateProcessInstanceRequest {
	s.ProcessInstanceId = &v
	return s
}

func (s *TerminateProcessInstanceRequest) SetRemark(v string) *TerminateProcessInstanceRequest {
	s.Remark = &v
	return s
}

type TerminateProcessInstanceResponseBody struct {
	// 撤销结果。
	Result *bool `json:"result,omitempty" xml:"result,omitempty"`
	// 接口调用是否成功。
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s TerminateProcessInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s TerminateProcessInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *TerminateProcessInstanceResponseBody) SetResult(v bool) *TerminateProcessInstanceResponseBody {
	s.Result = &v
	return s
}

func (s *TerminateProcessInstanceResponseBody) SetSuccess(v bool) *TerminateProcessInstanceResponseBody {
	s.Success = &v
	return s
}

type TerminateProcessInstanceResponse struct {
	Headers map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *TerminateProcessInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s TerminateProcessInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s TerminateProcessInstanceResponse) GoString() string {
	return s.String()
}

func (s *TerminateProcessInstanceResponse) SetHeaders(v map[string]*string) *TerminateProcessInstanceResponse {
	s.Headers = v
	return s
}

func (s *TerminateProcessInstanceResponse) SetBody(v *TerminateProcessInstanceResponseBody) *TerminateProcessInstanceResponse {
	s.Body = v
	return s
}

type UpdateIntegratedTaskHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s UpdateIntegratedTaskHeaders) String() string {
	return tea.Prettify(s)
}

func (s UpdateIntegratedTaskHeaders) GoString() string {
	return s.String()
}

func (s *UpdateIntegratedTaskHeaders) SetCommonHeaders(v map[string]*string) *UpdateIntegratedTaskHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UpdateIntegratedTaskHeaders) SetXAcsDingtalkAccessToken(v string) *UpdateIntegratedTaskHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type UpdateIntegratedTaskRequest struct {
	// OA审批流程实例ID，过创建实例接口获取
	ProcessInstanceId *string                             `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	Tasks             []*UpdateIntegratedTaskRequestTasks `json:"tasks,omitempty" xml:"tasks,omitempty" type:"Repeated"`
}

func (s UpdateIntegratedTaskRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateIntegratedTaskRequest) GoString() string {
	return s.String()
}

func (s *UpdateIntegratedTaskRequest) SetProcessInstanceId(v string) *UpdateIntegratedTaskRequest {
	s.ProcessInstanceId = &v
	return s
}

func (s *UpdateIntegratedTaskRequest) SetTasks(v []*UpdateIntegratedTaskRequestTasks) *UpdateIntegratedTaskRequest {
	s.Tasks = v
	return s
}

type UpdateIntegratedTaskRequestTasks struct {
	// 当status为COMPLETED时，必须指定任务结果：
	// AGREE：同意
	// REFUSE：拒绝
	//
	// 说明 当status为CANCELED时，不需要传result。
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 更新为目标任务状态，可选值 CANCELED、COMPLETED
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// OA审批任务ID
	TaskId *int64 `json:"taskId,omitempty" xml:"taskId,omitempty"`
}

func (s UpdateIntegratedTaskRequestTasks) String() string {
	return tea.Prettify(s)
}

func (s UpdateIntegratedTaskRequestTasks) GoString() string {
	return s.String()
}

func (s *UpdateIntegratedTaskRequestTasks) SetResult(v string) *UpdateIntegratedTaskRequestTasks {
	s.Result = &v
	return s
}

func (s *UpdateIntegratedTaskRequestTasks) SetStatus(v string) *UpdateIntegratedTaskRequestTasks {
	s.Status = &v
	return s
}

func (s *UpdateIntegratedTaskRequestTasks) SetTaskId(v int64) *UpdateIntegratedTaskRequestTasks {
	s.TaskId = &v
	return s
}

type UpdateIntegratedTaskResponseBody struct {
	// 是否更新成功
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s UpdateIntegratedTaskResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpdateIntegratedTaskResponseBody) GoString() string {
	return s.String()
}

func (s *UpdateIntegratedTaskResponseBody) SetSuccess(v bool) *UpdateIntegratedTaskResponseBody {
	s.Success = &v
	return s
}

type UpdateIntegratedTaskResponse struct {
	Headers map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *UpdateIntegratedTaskResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpdateIntegratedTaskResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateIntegratedTaskResponse) GoString() string {
	return s.String()
}

func (s *UpdateIntegratedTaskResponse) SetHeaders(v map[string]*string) *UpdateIntegratedTaskResponse {
	s.Headers = v
	return s
}

func (s *UpdateIntegratedTaskResponse) SetBody(v *UpdateIntegratedTaskResponseBody) *UpdateIntegratedTaskResponse {
	s.Body = v
	return s
}

type UpdateProcessInstanceHeaders struct {
	CommonHeaders           map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XAcsDingtalkAccessToken *string            `json:"x-acs-dingtalk-access-token,omitempty" xml:"x-acs-dingtalk-access-token,omitempty"`
}

func (s UpdateProcessInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s UpdateProcessInstanceHeaders) GoString() string {
	return s.String()
}

func (s *UpdateProcessInstanceHeaders) SetCommonHeaders(v map[string]*string) *UpdateProcessInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UpdateProcessInstanceHeaders) SetXAcsDingtalkAccessToken(v string) *UpdateProcessInstanceHeaders {
	s.XAcsDingtalkAccessToken = &v
	return s
}

type UpdateProcessInstanceRequest struct {
	// 抄送列表，注意最大抄送人数量不能超过30。
	Notifiers []*UpdateProcessInstanceRequestNotifiers `json:"notifiers,omitempty" xml:"notifiers,omitempty" type:"Repeated"`
	// 审批实例ID。
	ProcessInstanceId *string `json:"processInstanceId,omitempty" xml:"processInstanceId,omitempty"`
	// 实例结果：
	// 实例状态是COMPLETED，必须设置代表以下含义。
	// agree：同意
	// refuse：拒绝
	// 实例状态为TERMINATED，必须设置代表含义，result取值agree和refuse均代表撤销审批流。
	Result *string `json:"result,omitempty" xml:"result,omitempty"`
	// 实例状态：
	// COMPLETED：结束审批流
	// TERMINATED：终止审批流
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s UpdateProcessInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateProcessInstanceRequest) GoString() string {
	return s.String()
}

func (s *UpdateProcessInstanceRequest) SetNotifiers(v []*UpdateProcessInstanceRequestNotifiers) *UpdateProcessInstanceRequest {
	s.Notifiers = v
	return s
}

func (s *UpdateProcessInstanceRequest) SetProcessInstanceId(v string) *UpdateProcessInstanceRequest {
	s.ProcessInstanceId = &v
	return s
}

func (s *UpdateProcessInstanceRequest) SetResult(v string) *UpdateProcessInstanceRequest {
	s.Result = &v
	return s
}

func (s *UpdateProcessInstanceRequest) SetStatus(v string) *UpdateProcessInstanceRequest {
	s.Status = &v
	return s
}

type UpdateProcessInstanceRequestNotifiers struct {
	// 抄送接收人用户userId。
	UserId *string `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s UpdateProcessInstanceRequestNotifiers) String() string {
	return tea.Prettify(s)
}

func (s UpdateProcessInstanceRequestNotifiers) GoString() string {
	return s.String()
}

func (s *UpdateProcessInstanceRequestNotifiers) SetUserId(v string) *UpdateProcessInstanceRequestNotifiers {
	s.UserId = &v
	return s
}

type UpdateProcessInstanceResponseBody struct {
	Success *bool `json:"success,omitempty" xml:"success,omitempty"`
}

func (s UpdateProcessInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpdateProcessInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *UpdateProcessInstanceResponseBody) SetSuccess(v bool) *UpdateProcessInstanceResponseBody {
	s.Success = &v
	return s
}

type UpdateProcessInstanceResponse struct {
	Headers map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *UpdateProcessInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpdateProcessInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateProcessInstanceResponse) GoString() string {
	return s.String()
}

func (s *UpdateProcessInstanceResponse) SetHeaders(v map[string]*string) *UpdateProcessInstanceResponse {
	s.Headers = v
	return s
}

func (s *UpdateProcessInstanceResponse) SetBody(v *UpdateProcessInstanceResponseBody) *UpdateProcessInstanceResponse {
	s.Body = v
	return s
}

type Client struct {
	openapi.Client
}

func NewClient(config *openapi.Config) (*Client, error) {
	client := new(Client)
	err := client.Init(config)
	return client, err
}

func (client *Client) Init(config *openapi.Config) (_err error) {
	_err = client.Client.Init(config)
	if _err != nil {
		return _err
	}
	client.EndpointRule = tea.String("")
	if tea.BoolValue(util.Empty(client.Endpoint)) {
		client.Endpoint = tea.String("api.dingtalk.com")
	}

	return nil
}

func (client *Client) AddApproveDentryAuth(request *AddApproveDentryAuthRequest) (_result *AddApproveDentryAuthResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &AddApproveDentryAuthHeaders{}
	_result = &AddApproveDentryAuthResponse{}
	_body, _err := client.AddApproveDentryAuthWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) AddApproveDentryAuthWithOptions(request *AddApproveDentryAuthRequest, headers *AddApproveDentryAuthHeaders, runtime *util.RuntimeOptions) (_result *AddApproveDentryAuthResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.FileInfos)) {
		body["fileInfos"] = request.FileInfos
	}

	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		body["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &AddApproveDentryAuthResponse{}
	_body, _err := client.DoROARequest(tea.String("AddApproveDentryAuth"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processInstances/spaces/files/authDownload"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) AddProcessInstanceComment(request *AddProcessInstanceCommentRequest) (_result *AddProcessInstanceCommentResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &AddProcessInstanceCommentHeaders{}
	_result = &AddProcessInstanceCommentResponse{}
	_body, _err := client.AddProcessInstanceCommentWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) AddProcessInstanceCommentWithOptions(request *AddProcessInstanceCommentRequest, headers *AddProcessInstanceCommentHeaders, runtime *util.RuntimeOptions) (_result *AddProcessInstanceCommentResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CommentUserId)) {
		body["commentUserId"] = request.CommentUserId
	}

	if !tea.BoolValue(util.IsUnset(request.File)) {
		body["file"] = request.File
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		body["processInstanceId"] = request.ProcessInstanceId
	}

	if !tea.BoolValue(util.IsUnset(request.Text)) {
		body["text"] = request.Text
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &AddProcessInstanceCommentResponse{}
	_body, _err := client.DoROARequest(tea.String("AddProcessInstanceComment"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processInstances/comments"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) BatchUpdateProcessInstance(request *BatchUpdateProcessInstanceRequest) (_result *BatchUpdateProcessInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &BatchUpdateProcessInstanceHeaders{}
	_result = &BatchUpdateProcessInstanceResponse{}
	_body, _err := client.BatchUpdateProcessInstanceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) BatchUpdateProcessInstanceWithOptions(request *BatchUpdateProcessInstanceRequest, headers *BatchUpdateProcessInstanceHeaders, runtime *util.RuntimeOptions) (_result *BatchUpdateProcessInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.UpdateProcessInstanceRequests)) {
		body["updateProcessInstanceRequests"] = request.UpdateProcessInstanceRequests
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &BatchUpdateProcessInstanceResponse{}
	_body, _err := client.DoROARequest(tea.String("BatchUpdateProcessInstance"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("PUT"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/instances/batch"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CancelIntegratedTask(request *CancelIntegratedTaskRequest) (_result *CancelIntegratedTaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CancelIntegratedTaskHeaders{}
	_result = &CancelIntegratedTaskResponse{}
	_body, _err := client.CancelIntegratedTaskWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CancelIntegratedTaskWithOptions(request *CancelIntegratedTaskRequest, headers *CancelIntegratedTaskHeaders, runtime *util.RuntimeOptions) (_result *CancelIntegratedTaskResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ActivityId)) {
		body["activityId"] = request.ActivityId
	}

	if !tea.BoolValue(util.IsUnset(request.ActivityIds)) {
		body["activityIds"] = request.ActivityIds
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		body["processInstanceId"] = request.ProcessInstanceId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CancelIntegratedTaskResponse{}
	_body, _err := client.DoROARequest(tea.String("CancelIntegratedTask"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/tasks/cancel"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CleanProcessData(request *CleanProcessDataRequest) (_result *CleanProcessDataResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CleanProcessDataHeaders{}
	_result = &CleanProcessDataResponse{}
	_body, _err := client.CleanProcessDataWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CleanProcessDataWithOptions(request *CleanProcessDataRequest, headers *CleanProcessDataHeaders, runtime *util.RuntimeOptions) (_result *CleanProcessDataResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CorpId)) {
		body["corpId"] = request.CorpId
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		body["processCode"] = request.ProcessCode
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CleanProcessDataResponse{}
	_body, _err := client.DoROARequest(tea.String("CleanProcessData"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processes/clean"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CopyProcess(request *CopyProcessRequest) (_result *CopyProcessResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CopyProcessHeaders{}
	_result = &CopyProcessResponse{}
	_body, _err := client.CopyProcessWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CopyProcessWithOptions(request *CopyProcessRequest, headers *CopyProcessHeaders, runtime *util.RuntimeOptions) (_result *CopyProcessResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CopyOptions)) {
		body["copyOptions"] = request.CopyOptions
	}

	if !tea.BoolValue(util.IsUnset(request.SourceCorpId)) {
		body["sourceCorpId"] = request.SourceCorpId
	}

	if !tea.BoolValue(util.IsUnset(request.SourceProcessVOList)) {
		body["sourceProcessVOList"] = request.SourceProcessVOList
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CopyProcessResponse{}
	_body, _err := client.DoROARequest(tea.String("CopyProcess"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processes/copy"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateIntegratedTask(request *CreateIntegratedTaskRequest) (_result *CreateIntegratedTaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CreateIntegratedTaskHeaders{}
	_result = &CreateIntegratedTaskResponse{}
	_body, _err := client.CreateIntegratedTaskWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateIntegratedTaskWithOptions(request *CreateIntegratedTaskRequest, headers *CreateIntegratedTaskHeaders, runtime *util.RuntimeOptions) (_result *CreateIntegratedTaskResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ActivityId)) {
		body["activityId"] = request.ActivityId
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		body["processInstanceId"] = request.ProcessInstanceId
	}

	if !tea.BoolValue(util.IsUnset(request.Tasks)) {
		body["tasks"] = request.Tasks
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CreateIntegratedTaskResponse{}
	_body, _err := client.DoROARequest(tea.String("CreateIntegratedTask"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/tasks"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteProcess(request *DeleteProcessRequest) (_result *DeleteProcessResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteProcessHeaders{}
	_result = &DeleteProcessResponse{}
	_body, _err := client.DeleteProcessWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteProcessWithOptions(request *DeleteProcessRequest, headers *DeleteProcessHeaders, runtime *util.RuntimeOptions) (_result *DeleteProcessResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CleanRunningTask)) {
		query["cleanRunningTask"] = request.CleanRunningTask
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		query["processCode"] = request.ProcessCode
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &DeleteProcessResponse{}
	_body, _err := client.DoROARequest(tea.String("DeleteProcess"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("DELETE"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/schemas"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ExecuteProcessInstance(request *ExecuteProcessInstanceRequest) (_result *ExecuteProcessInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ExecuteProcessInstanceHeaders{}
	_result = &ExecuteProcessInstanceResponse{}
	_body, _err := client.ExecuteProcessInstanceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ExecuteProcessInstanceWithOptions(request *ExecuteProcessInstanceRequest, headers *ExecuteProcessInstanceHeaders, runtime *util.RuntimeOptions) (_result *ExecuteProcessInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ActionerUserId)) {
		body["actionerUserId"] = request.ActionerUserId
	}

	if !tea.BoolValue(util.IsUnset(request.File)) {
		body["file"] = request.File
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		body["processInstanceId"] = request.ProcessInstanceId
	}

	if !tea.BoolValue(util.IsUnset(request.Remark)) {
		body["remark"] = request.Remark
	}

	if !tea.BoolValue(util.IsUnset(request.Result)) {
		body["result"] = request.Result
	}

	if !tea.BoolValue(util.IsUnset(request.TaskId)) {
		body["taskId"] = request.TaskId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &ExecuteProcessInstanceResponse{}
	_body, _err := client.DoROARequest(tea.String("ExecuteProcessInstance"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processInstances/execute"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) FormCreate(request *FormCreateRequest) (_result *FormCreateResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &FormCreateHeaders{}
	_result = &FormCreateResponse{}
	_body, _err := client.FormCreateWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) FormCreateWithOptions(request *FormCreateRequest, headers *FormCreateHeaders, runtime *util.RuntimeOptions) (_result *FormCreateResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.FormComponents)) {
		body["formComponents"] = request.FormComponents
	}

	if !tea.BoolValue(util.IsUnset(request.Name)) {
		body["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		body["processCode"] = request.ProcessCode
	}

	if !tea.BoolValue(util.IsUnset(request.TemplateConfig)) {
		body["templateConfig"] = request.TemplateConfig
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &FormCreateResponse{}
	_body, _err := client.DoROARequest(tea.String("FormCreate"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/forms"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetAttachmentSpace(request *GetAttachmentSpaceRequest) (_result *GetAttachmentSpaceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetAttachmentSpaceHeaders{}
	_result = &GetAttachmentSpaceResponse{}
	_body, _err := client.GetAttachmentSpaceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetAttachmentSpaceWithOptions(request *GetAttachmentSpaceRequest, headers *GetAttachmentSpaceHeaders, runtime *util.RuntimeOptions) (_result *GetAttachmentSpaceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AgentId)) {
		body["agentId"] = request.AgentId
	}

	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		body["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &GetAttachmentSpaceResponse{}
	_body, _err := client.DoROARequest(tea.String("GetAttachmentSpace"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processInstances/spaces/infos/query"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetConditionFormComponent(request *GetConditionFormComponentRequest) (_result *GetConditionFormComponentResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetConditionFormComponentHeaders{}
	_result = &GetConditionFormComponentResponse{}
	_body, _err := client.GetConditionFormComponentWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetConditionFormComponentWithOptions(request *GetConditionFormComponentRequest, headers *GetConditionFormComponentHeaders, runtime *util.RuntimeOptions) (_result *GetConditionFormComponentResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AgentId)) {
		query["agentId"] = request.AgentId
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		query["processCode"] = request.ProcessCode
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &GetConditionFormComponentResponse{}
	_body, _err := client.DoROARequest(tea.String("GetConditionFormComponent"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/processes/conditions/components"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetCrmProcCodes() (_result *GetCrmProcCodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetCrmProcCodesHeaders{}
	_result = &GetCrmProcCodesResponse{}
	_body, _err := client.GetCrmProcCodesWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetCrmProcCodesWithOptions(headers *GetCrmProcCodesHeaders, runtime *util.RuntimeOptions) (_result *GetCrmProcCodesResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	_result = &GetCrmProcCodesResponse{}
	_body, _err := client.DoROARequest(tea.String("GetCrmProcCodes"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/crm/processes"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetManageProcessByStaffId(request *GetManageProcessByStaffIdRequest) (_result *GetManageProcessByStaffIdResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetManageProcessByStaffIdHeaders{}
	_result = &GetManageProcessByStaffIdResponse{}
	_body, _err := client.GetManageProcessByStaffIdWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetManageProcessByStaffIdWithOptions(request *GetManageProcessByStaffIdRequest, headers *GetManageProcessByStaffIdHeaders, runtime *util.RuntimeOptions) (_result *GetManageProcessByStaffIdResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		query["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &GetManageProcessByStaffIdResponse{}
	_body, _err := client.DoROARequest(tea.String("GetManageProcessByStaffId"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/processes/managements/templates"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetProcessCodeByName(request *GetProcessCodeByNameRequest) (_result *GetProcessCodeByNameResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetProcessCodeByNameHeaders{}
	_result = &GetProcessCodeByNameResponse{}
	_body, _err := client.GetProcessCodeByNameWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetProcessCodeByNameWithOptions(request *GetProcessCodeByNameRequest, headers *GetProcessCodeByNameHeaders, runtime *util.RuntimeOptions) (_result *GetProcessCodeByNameResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Name)) {
		query["name"] = request.Name
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &GetProcessCodeByNameResponse{}
	_body, _err := client.DoROARequest(tea.String("GetProcessCodeByName"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/schemaNames/processCodes"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetProcessConfig(request *GetProcessConfigRequest) (_result *GetProcessConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetProcessConfigHeaders{}
	_result = &GetProcessConfigResponse{}
	_body, _err := client.GetProcessConfigWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetProcessConfigWithOptions(request *GetProcessConfigRequest, headers *GetProcessConfigHeaders, runtime *util.RuntimeOptions) (_result *GetProcessConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ProcCode)) {
		query["procCode"] = request.ProcCode
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &GetProcessConfigResponse{}
	_body, _err := client.DoROARequest(tea.String("GetProcessConfig"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/crm/processes/configurations"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetProcessInstance(request *GetProcessInstanceRequest) (_result *GetProcessInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetProcessInstanceHeaders{}
	_result = &GetProcessInstanceResponse{}
	_body, _err := client.GetProcessInstanceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetProcessInstanceWithOptions(request *GetProcessInstanceRequest, headers *GetProcessInstanceHeaders, runtime *util.RuntimeOptions) (_result *GetProcessInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		query["processInstanceId"] = request.ProcessInstanceId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &GetProcessInstanceResponse{}
	_body, _err := client.DoROARequest(tea.String("GetProcessInstance"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/processInstances"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetSpaceWithDownloadAuth(request *GetSpaceWithDownloadAuthRequest) (_result *GetSpaceWithDownloadAuthResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetSpaceWithDownloadAuthHeaders{}
	_result = &GetSpaceWithDownloadAuthResponse{}
	_body, _err := client.GetSpaceWithDownloadAuthWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetSpaceWithDownloadAuthWithOptions(request *GetSpaceWithDownloadAuthRequest, headers *GetSpaceWithDownloadAuthHeaders, runtime *util.RuntimeOptions) (_result *GetSpaceWithDownloadAuthResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AgentId)) {
		body["agentId"] = request.AgentId
	}

	if !tea.BoolValue(util.IsUnset(request.FileId)) {
		body["fileId"] = request.FileId
	}

	if !tea.BoolValue(util.IsUnset(request.FileIdList)) {
		body["fileIdList"] = request.FileIdList
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		body["processInstanceId"] = request.ProcessInstanceId
	}

	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		body["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &GetSpaceWithDownloadAuthResponse{}
	_body, _err := client.DoROARequest(tea.String("GetSpaceWithDownloadAuth"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processInstances/spaces/authPreview"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetUserTodoTaskSum(request *GetUserTodoTaskSumRequest) (_result *GetUserTodoTaskSumResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetUserTodoTaskSumHeaders{}
	_result = &GetUserTodoTaskSumResponse{}
	_body, _err := client.GetUserTodoTaskSumWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetUserTodoTaskSumWithOptions(request *GetUserTodoTaskSumRequest, headers *GetUserTodoTaskSumHeaders, runtime *util.RuntimeOptions) (_result *GetUserTodoTaskSumResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		query["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &GetUserTodoTaskSumResponse{}
	_body, _err := client.DoROARequest(tea.String("GetUserTodoTaskSum"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/processes/todoTasks/numbers"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GrantCspaceAuthorization(request *GrantCspaceAuthorizationRequest) (_result *GrantCspaceAuthorizationResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GrantCspaceAuthorizationHeaders{}
	_result = &GrantCspaceAuthorizationResponse{}
	_body, _err := client.GrantCspaceAuthorizationWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GrantCspaceAuthorizationWithOptions(request *GrantCspaceAuthorizationRequest, headers *GrantCspaceAuthorizationHeaders, runtime *util.RuntimeOptions) (_result *GrantCspaceAuthorizationResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DurationSeconds)) {
		body["durationSeconds"] = request.DurationSeconds
	}

	if !tea.BoolValue(util.IsUnset(request.SpaceId)) {
		body["spaceId"] = request.SpaceId
	}

	if !tea.BoolValue(util.IsUnset(request.Type)) {
		body["type"] = request.Type
	}

	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		body["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &GrantCspaceAuthorizationResponse{}
	_body, _err := client.DoROARequest(tea.String("GrantCspaceAuthorization"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/spaces/authorize"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GrantProcessInstanceForDownloadFile(request *GrantProcessInstanceForDownloadFileRequest) (_result *GrantProcessInstanceForDownloadFileResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GrantProcessInstanceForDownloadFileHeaders{}
	_result = &GrantProcessInstanceForDownloadFileResponse{}
	_body, _err := client.GrantProcessInstanceForDownloadFileWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GrantProcessInstanceForDownloadFileWithOptions(request *GrantProcessInstanceForDownloadFileRequest, headers *GrantProcessInstanceForDownloadFileHeaders, runtime *util.RuntimeOptions) (_result *GrantProcessInstanceForDownloadFileResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.FileId)) {
		body["fileId"] = request.FileId
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		body["processInstanceId"] = request.ProcessInstanceId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &GrantProcessInstanceForDownloadFileResponse{}
	_body, _err := client.DoROARequest(tea.String("GrantProcessInstanceForDownloadFile"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processInstances/spaces/files/urls/download"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) InstallApp(request *InstallAppRequest) (_result *InstallAppResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &InstallAppHeaders{}
	_result = &InstallAppResponse{}
	_body, _err := client.InstallAppWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) InstallAppWithOptions(request *InstallAppRequest, headers *InstallAppHeaders, runtime *util.RuntimeOptions) (_result *InstallAppResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BizGroup)) {
		body["bizGroup"] = request.BizGroup
	}

	if !tea.BoolValue(util.IsUnset(request.InstallOption)) {
		body["installOption"] = request.InstallOption
	}

	if !tea.BoolValue(util.IsUnset(request.SourceDirName)) {
		body["sourceDirName"] = request.SourceDirName
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &InstallAppResponse{}
	_body, _err := client.DoROARequest(tea.String("InstallApp"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processes/apps/install"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListProcessInstanceIds(request *ListProcessInstanceIdsRequest) (_result *ListProcessInstanceIdsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListProcessInstanceIdsHeaders{}
	_result = &ListProcessInstanceIdsResponse{}
	_body, _err := client.ListProcessInstanceIdsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListProcessInstanceIdsWithOptions(request *ListProcessInstanceIdsRequest, headers *ListProcessInstanceIdsHeaders, runtime *util.RuntimeOptions) (_result *ListProcessInstanceIdsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		body["endTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.MaxResults)) {
		body["maxResults"] = request.MaxResults
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		body["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		body["processCode"] = request.ProcessCode
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		body["startTime"] = request.StartTime
	}

	if !tea.BoolValue(util.IsUnset(request.UserIds)) {
		body["userIds"] = request.UserIds
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &ListProcessInstanceIdsResponse{}
	_body, _err := client.DoROARequest(tea.String("ListProcessInstanceIds"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processes/instanceIds/query"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListTodoWorkRecords(request *ListTodoWorkRecordsRequest) (_result *ListTodoWorkRecordsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListTodoWorkRecordsHeaders{}
	_result = &ListTodoWorkRecordsResponse{}
	_body, _err := client.ListTodoWorkRecordsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListTodoWorkRecordsWithOptions(request *ListTodoWorkRecordsRequest, headers *ListTodoWorkRecordsHeaders, runtime *util.RuntimeOptions) (_result *ListTodoWorkRecordsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.MaxResults)) {
		query["maxResults"] = request.MaxResults
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Status)) {
		query["status"] = request.Status
	}

	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		query["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &ListTodoWorkRecordsResponse{}
	_body, _err := client.DoROARequest(tea.String("ListTodoWorkRecords"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/workRecords/todoTasks"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListUserVisibleBpmsProcesses(request *ListUserVisibleBpmsProcessesRequest) (_result *ListUserVisibleBpmsProcessesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListUserVisibleBpmsProcessesHeaders{}
	_result = &ListUserVisibleBpmsProcessesResponse{}
	_body, _err := client.ListUserVisibleBpmsProcessesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListUserVisibleBpmsProcessesWithOptions(request *ListUserVisibleBpmsProcessesRequest, headers *ListUserVisibleBpmsProcessesHeaders, runtime *util.RuntimeOptions) (_result *ListUserVisibleBpmsProcessesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.MaxResults)) {
		query["maxResults"] = request.MaxResults
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		query["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &ListUserVisibleBpmsProcessesResponse{}
	_body, _err := client.DoROARequest(tea.String("ListUserVisibleBpmsProcesses"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/processes/userVisibilities/templates"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ProcessForecast(request *ProcessForecastRequest) (_result *ProcessForecastResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ProcessForecastHeaders{}
	_result = &ProcessForecastResponse{}
	_body, _err := client.ProcessForecastWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ProcessForecastWithOptions(request *ProcessForecastRequest, headers *ProcessForecastHeaders, runtime *util.RuntimeOptions) (_result *ProcessForecastResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DeptId)) {
		body["deptId"] = request.DeptId
	}

	if !tea.BoolValue(util.IsUnset(request.FormComponentValues)) {
		body["formComponentValues"] = request.FormComponentValues
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		body["processCode"] = request.ProcessCode
	}

	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		body["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &ProcessForecastResponse{}
	_body, _err := client.DoROARequest(tea.String("ProcessForecast"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processes/forecast"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) QueryAllFormInstances(request *QueryAllFormInstancesRequest) (_result *QueryAllFormInstancesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &QueryAllFormInstancesHeaders{}
	_result = &QueryAllFormInstancesResponse{}
	_body, _err := client.QueryAllFormInstancesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) QueryAllFormInstancesWithOptions(request *QueryAllFormInstancesRequest, headers *QueryAllFormInstancesHeaders, runtime *util.RuntimeOptions) (_result *QueryAllFormInstancesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AppUuid)) {
		query["appUuid"] = request.AppUuid
	}

	if !tea.BoolValue(util.IsUnset(request.FormCode)) {
		query["formCode"] = request.FormCode
	}

	if !tea.BoolValue(util.IsUnset(request.MaxResults)) {
		query["maxResults"] = request.MaxResults
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &QueryAllFormInstancesResponse{}
	_body, _err := client.DoROARequest(tea.String("QueryAllFormInstances"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/forms/pages/instances"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) QueryAllProcessInstances(request *QueryAllProcessInstancesRequest) (_result *QueryAllProcessInstancesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &QueryAllProcessInstancesHeaders{}
	_result = &QueryAllProcessInstancesResponse{}
	_body, _err := client.QueryAllProcessInstancesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) QueryAllProcessInstancesWithOptions(request *QueryAllProcessInstancesRequest, headers *QueryAllProcessInstancesHeaders, runtime *util.RuntimeOptions) (_result *QueryAllProcessInstancesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AppUuid)) {
		query["appUuid"] = request.AppUuid
	}

	if !tea.BoolValue(util.IsUnset(request.EndTimeInMills)) {
		query["endTimeInMills"] = request.EndTimeInMills
	}

	if !tea.BoolValue(util.IsUnset(request.MaxResults)) {
		query["maxResults"] = request.MaxResults
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		query["processCode"] = request.ProcessCode
	}

	if !tea.BoolValue(util.IsUnset(request.StartTimeInMills)) {
		query["startTimeInMills"] = request.StartTimeInMills
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &QueryAllProcessInstancesResponse{}
	_body, _err := client.DoROARequest(tea.String("QueryAllProcessInstances"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/processes/pages/instances"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) QueryFormByBizType(request *QueryFormByBizTypeRequest) (_result *QueryFormByBizTypeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &QueryFormByBizTypeHeaders{}
	_result = &QueryFormByBizTypeResponse{}
	_body, _err := client.QueryFormByBizTypeWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) QueryFormByBizTypeWithOptions(request *QueryFormByBizTypeRequest, headers *QueryFormByBizTypeHeaders, runtime *util.RuntimeOptions) (_result *QueryFormByBizTypeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AppUuid)) {
		body["appUuid"] = request.AppUuid
	}

	if !tea.BoolValue(util.IsUnset(request.BizTypes)) {
		body["bizTypes"] = request.BizTypes
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &QueryFormByBizTypeResponse{}
	_body, _err := client.DoROARequest(tea.String("QueryFormByBizType"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/forms/forminfos/query"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) QueryFormInstance(request *QueryFormInstanceRequest) (_result *QueryFormInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &QueryFormInstanceHeaders{}
	_result = &QueryFormInstanceResponse{}
	_body, _err := client.QueryFormInstanceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) QueryFormInstanceWithOptions(request *QueryFormInstanceRequest, headers *QueryFormInstanceHeaders, runtime *util.RuntimeOptions) (_result *QueryFormInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AppUuid)) {
		query["appUuid"] = request.AppUuid
	}

	if !tea.BoolValue(util.IsUnset(request.FormCode)) {
		query["formCode"] = request.FormCode
	}

	if !tea.BoolValue(util.IsUnset(request.FormInstanceId)) {
		query["formInstanceId"] = request.FormInstanceId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &QueryFormInstanceResponse{}
	_body, _err := client.DoROARequest(tea.String("QueryFormInstance"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/forms/instances"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) QueryIntegratedTodoTask(request *QueryIntegratedTodoTaskRequest) (_result *QueryIntegratedTodoTaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &QueryIntegratedTodoTaskHeaders{}
	_result = &QueryIntegratedTodoTaskResponse{}
	_body, _err := client.QueryIntegratedTodoTaskWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) QueryIntegratedTodoTaskWithOptions(request *QueryIntegratedTodoTaskRequest, headers *QueryIntegratedTodoTaskHeaders, runtime *util.RuntimeOptions) (_result *QueryIntegratedTodoTaskResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CreateBefore)) {
		query["createBefore"] = request.CreateBefore
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["pageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["pageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		query["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &QueryIntegratedTodoTaskResponse{}
	_body, _err := client.DoROARequest(tea.String("QueryIntegratedTodoTask"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/todoTasks"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) QueryProcessByBizCategoryId(request *QueryProcessByBizCategoryIdRequest) (_result *QueryProcessByBizCategoryIdResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &QueryProcessByBizCategoryIdHeaders{}
	_result = &QueryProcessByBizCategoryIdResponse{}
	_body, _err := client.QueryProcessByBizCategoryIdWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) QueryProcessByBizCategoryIdWithOptions(request *QueryProcessByBizCategoryIdRequest, headers *QueryProcessByBizCategoryIdHeaders, runtime *util.RuntimeOptions) (_result *QueryProcessByBizCategoryIdResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BizType)) {
		query["bizType"] = request.BizType
	}

	if !tea.BoolValue(util.IsUnset(request.UserId)) {
		query["userId"] = request.UserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &QueryProcessByBizCategoryIdResponse{}
	_body, _err := client.DoROARequest(tea.String("QueryProcessByBizCategoryId"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/processes/categories/templates"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) QuerySchemaByProcessCode(request *QuerySchemaByProcessCodeRequest) (_result *QuerySchemaByProcessCodeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &QuerySchemaByProcessCodeHeaders{}
	_result = &QuerySchemaByProcessCodeResponse{}
	_body, _err := client.QuerySchemaByProcessCodeWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) QuerySchemaByProcessCodeWithOptions(request *QuerySchemaByProcessCodeRequest, headers *QuerySchemaByProcessCodeHeaders, runtime *util.RuntimeOptions) (_result *QuerySchemaByProcessCodeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AppUuid)) {
		query["appUuid"] = request.AppUuid
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		query["processCode"] = request.ProcessCode
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	_result = &QuerySchemaByProcessCodeResponse{}
	_body, _err := client.DoROARequest(tea.String("QuerySchemaByProcessCode"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("GET"), tea.String("AK"), tea.String("/v1.0/workflow/forms/schemas/processCodes"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RedirectWorkflowTask(request *RedirectWorkflowTaskRequest) (_result *RedirectWorkflowTaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &RedirectWorkflowTaskHeaders{}
	_result = &RedirectWorkflowTaskResponse{}
	_body, _err := client.RedirectWorkflowTaskWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RedirectWorkflowTaskWithOptions(request *RedirectWorkflowTaskRequest, headers *RedirectWorkflowTaskHeaders, runtime *util.RuntimeOptions) (_result *RedirectWorkflowTaskResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ActionName)) {
		body["actionName"] = request.ActionName
	}

	if !tea.BoolValue(util.IsUnset(request.OperateUserId)) {
		body["operateUserId"] = request.OperateUserId
	}

	if !tea.BoolValue(util.IsUnset(request.Remark)) {
		body["remark"] = request.Remark
	}

	if !tea.BoolValue(util.IsUnset(request.TaskId)) {
		body["taskId"] = request.TaskId
	}

	if !tea.BoolValue(util.IsUnset(request.ToUserId)) {
		body["toUserId"] = request.ToUserId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &RedirectWorkflowTaskResponse{}
	_body, _err := client.DoROARequest(tea.String("RedirectWorkflowTask"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/tasks/redirect"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) SaveIntegratedInstance(request *SaveIntegratedInstanceRequest) (_result *SaveIntegratedInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &SaveIntegratedInstanceHeaders{}
	_result = &SaveIntegratedInstanceResponse{}
	_body, _err := client.SaveIntegratedInstanceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) SaveIntegratedInstanceWithOptions(request *SaveIntegratedInstanceRequest, headers *SaveIntegratedInstanceHeaders, runtime *util.RuntimeOptions) (_result *SaveIntegratedInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.FormComponentValueList)) {
		body["formComponentValueList"] = request.FormComponentValueList
	}

	if !tea.BoolValue(util.IsUnset(request.Notifiers)) {
		body["notifiers"] = request.Notifiers
	}

	if !tea.BoolValue(util.IsUnset(request.OriginatorUserId)) {
		body["originatorUserId"] = request.OriginatorUserId
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		body["processCode"] = request.ProcessCode
	}

	if !tea.BoolValue(util.IsUnset(request.Title)) {
		body["title"] = request.Title
	}

	if !tea.BoolValue(util.IsUnset(request.Url)) {
		body["url"] = request.Url
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &SaveIntegratedInstanceResponse{}
	_body, _err := client.DoROARequest(tea.String("SaveIntegratedInstance"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/instances"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) SaveProcess(request *SaveProcessRequest) (_result *SaveProcessResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &SaveProcessHeaders{}
	_result = &SaveProcessResponse{}
	_body, _err := client.SaveProcessWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) SaveProcessWithOptions(request *SaveProcessRequest, headers *SaveProcessHeaders, runtime *util.RuntimeOptions) (_result *SaveProcessResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.FormComponents)) {
		body["formComponents"] = request.FormComponents
	}

	if !tea.BoolValue(util.IsUnset(request.Name)) {
		body["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		body["processCode"] = request.ProcessCode
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessFeatureConfig)) {
		body["processFeatureConfig"] = request.ProcessFeatureConfig
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &SaveProcessResponse{}
	_body, _err := client.DoROARequest(tea.String("SaveProcess"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/schemas"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) StartProcessInstance(request *StartProcessInstanceRequest) (_result *StartProcessInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &StartProcessInstanceHeaders{}
	_result = &StartProcessInstanceResponse{}
	_body, _err := client.StartProcessInstanceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) StartProcessInstanceWithOptions(request *StartProcessInstanceRequest, headers *StartProcessInstanceHeaders, runtime *util.RuntimeOptions) (_result *StartProcessInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Approvers)) {
		body["approvers"] = request.Approvers
	}

	if !tea.BoolValue(util.IsUnset(request.CcList)) {
		body["ccList"] = request.CcList
	}

	if !tea.BoolValue(util.IsUnset(request.CcPosition)) {
		body["ccPosition"] = request.CcPosition
	}

	if !tea.BoolValue(util.IsUnset(request.DeptId)) {
		body["deptId"] = request.DeptId
	}

	if !tea.BoolValue(util.IsUnset(request.FormComponentValues)) {
		body["formComponentValues"] = request.FormComponentValues
	}

	if !tea.BoolValue(util.IsUnset(request.MicroappAgentId)) {
		body["microappAgentId"] = request.MicroappAgentId
	}

	if !tea.BoolValue(util.IsUnset(request.OriginatorUserId)) {
		body["originatorUserId"] = request.OriginatorUserId
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessCode)) {
		body["processCode"] = request.ProcessCode
	}

	if !tea.BoolValue(util.IsUnset(request.TargetSelectActioners)) {
		body["targetSelectActioners"] = request.TargetSelectActioners
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &StartProcessInstanceResponse{}
	_body, _err := client.DoROARequest(tea.String("StartProcessInstance"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processInstances"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) TerminateProcessInstance(request *TerminateProcessInstanceRequest) (_result *TerminateProcessInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &TerminateProcessInstanceHeaders{}
	_result = &TerminateProcessInstanceResponse{}
	_body, _err := client.TerminateProcessInstanceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) TerminateProcessInstanceWithOptions(request *TerminateProcessInstanceRequest, headers *TerminateProcessInstanceHeaders, runtime *util.RuntimeOptions) (_result *TerminateProcessInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.IsSystem)) {
		body["isSystem"] = request.IsSystem
	}

	if !tea.BoolValue(util.IsUnset(request.OperatingUserId)) {
		body["operatingUserId"] = request.OperatingUserId
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		body["processInstanceId"] = request.ProcessInstanceId
	}

	if !tea.BoolValue(util.IsUnset(request.Remark)) {
		body["remark"] = request.Remark
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &TerminateProcessInstanceResponse{}
	_body, _err := client.DoROARequest(tea.String("TerminateProcessInstance"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("POST"), tea.String("AK"), tea.String("/v1.0/workflow/processInstances/terminate"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateIntegratedTask(request *UpdateIntegratedTaskRequest) (_result *UpdateIntegratedTaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &UpdateIntegratedTaskHeaders{}
	_result = &UpdateIntegratedTaskResponse{}
	_body, _err := client.UpdateIntegratedTaskWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateIntegratedTaskWithOptions(request *UpdateIntegratedTaskRequest, headers *UpdateIntegratedTaskHeaders, runtime *util.RuntimeOptions) (_result *UpdateIntegratedTaskResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		body["processInstanceId"] = request.ProcessInstanceId
	}

	if !tea.BoolValue(util.IsUnset(request.Tasks)) {
		body["tasks"] = request.Tasks
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &UpdateIntegratedTaskResponse{}
	_body, _err := client.DoROARequest(tea.String("UpdateIntegratedTask"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("PUT"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/tasks"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateProcessInstance(request *UpdateProcessInstanceRequest) (_result *UpdateProcessInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &UpdateProcessInstanceHeaders{}
	_result = &UpdateProcessInstanceResponse{}
	_body, _err := client.UpdateProcessInstanceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateProcessInstanceWithOptions(request *UpdateProcessInstanceRequest, headers *UpdateProcessInstanceHeaders, runtime *util.RuntimeOptions) (_result *UpdateProcessInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Notifiers)) {
		body["notifiers"] = request.Notifiers
	}

	if !tea.BoolValue(util.IsUnset(request.ProcessInstanceId)) {
		body["processInstanceId"] = request.ProcessInstanceId
	}

	if !tea.BoolValue(util.IsUnset(request.Result)) {
		body["result"] = request.Result
	}

	if !tea.BoolValue(util.IsUnset(request.Status)) {
		body["status"] = request.Status
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XAcsDingtalkAccessToken)) {
		realHeaders["x-acs-dingtalk-access-token"] = util.ToJSONString(headers.XAcsDingtalkAccessToken)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &UpdateProcessInstanceResponse{}
	_body, _err := client.DoROARequest(tea.String("UpdateProcessInstance"), tea.String("workflow_1.0"), tea.String("HTTP"), tea.String("PUT"), tea.String("AK"), tea.String("/v1.0/workflow/processCentres/instances"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}
