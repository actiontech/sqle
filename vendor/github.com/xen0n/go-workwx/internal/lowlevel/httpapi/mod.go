package httpapi

import (
	"net/http"

	"github.com/xen0n/go-workwx/internal/lowlevel/encryptor"
	"github.com/xen0n/go-workwx/internal/lowlevel/envelope"
)

type LowlevelHandler struct {
	token     string
	encryptor *encryptor.WorkwxEncryptor
	ep        *envelope.Processor
	eh        EnvelopeHandler
}

var _ http.Handler = (*LowlevelHandler)(nil)

func NewLowlevelHandler(
	token string,
	encodingAESKey string,
	eh EnvelopeHandler,
) (*LowlevelHandler, error) {
	enc, err := encryptor.NewWorkwxEncryptor(encodingAESKey)
	if err != nil {
		return nil, err
	}

	ep, err := envelope.NewProcessor(token, encodingAESKey)
	if err != nil {
		return nil, err
	}

	return &LowlevelHandler{
		token:     token,
		encryptor: enc,
		ep:        ep,
		eh:        eh,
	}, nil
}

func (h *LowlevelHandler) ServeHTTP(
	rw http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		// 测试回调模式请求
		h.echoTestHandler(rw, r)

	case http.MethodPost:
		// 回调事件
		h.eventHandler(rw, r)

	default:
		// unhandled request method
		rw.WriteHeader(http.StatusNotImplemented)
	}
}
