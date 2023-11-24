package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// GetJobInfoResponseBodyJob 任务信息。
type GetJobInfoResponseBodyJob struct {

	// 任务ID。
	Id string `json:"id"`

	// 任务名称。
	Name string `json:"name"`

	// 任务执行状态。  取值： - 值为“Running”，表示任务正在执行。 - 值为“Completed”，表示任务执行成功。 - 值为“Failed”，表示任务执行失败。
	Status GetJobInfoResponseBodyJobStatus `json:"status"`

	// 创建时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	Created string `json:"created"`

	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	Ended *string `json:"ended,omitempty"`

	// 任务执行进度。执行中状态才返回执行进度，例如60%，否则返回“”。
	Process *string `json:"process,omitempty"`

	Instance *GetTaskDetailListRspJobsInstance `json:"instance"`

	// 根据不同的任务，显示不同的内容。
	Entities *interface{} `json:"entities,omitempty"`

	// 任务执行失败时的错误信息。
	FailReason *string `json:"fail_reason,omitempty"`
}

func (o GetJobInfoResponseBodyJob) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetJobInfoResponseBodyJob struct{}"
	}

	return strings.Join([]string{"GetJobInfoResponseBodyJob", string(data)}, " ")
}

type GetJobInfoResponseBodyJobStatus struct {
	value string
}

type GetJobInfoResponseBodyJobStatusEnum struct {
	RUNNING   GetJobInfoResponseBodyJobStatus
	COMPLETED GetJobInfoResponseBodyJobStatus
	FAILED    GetJobInfoResponseBodyJobStatus
}

func GetGetJobInfoResponseBodyJobStatusEnum() GetJobInfoResponseBodyJobStatusEnum {
	return GetJobInfoResponseBodyJobStatusEnum{
		RUNNING: GetJobInfoResponseBodyJobStatus{
			value: "Running",
		},
		COMPLETED: GetJobInfoResponseBodyJobStatus{
			value: "Completed",
		},
		FAILED: GetJobInfoResponseBodyJobStatus{
			value: "Failed",
		},
	}
}

func (c GetJobInfoResponseBodyJobStatus) Value() string {
	return c.value
}

func (c GetJobInfoResponseBodyJobStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *GetJobInfoResponseBodyJobStatus) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
