package returncodes

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"sync"
)

// 状态码对应着httpStatus，业务异常请定义 >= 1000(4位数)
var (
	Success = of(http.StatusOK, http.StatusText(http.StatusOK))

	ErrSystem   = NewErrorCoder(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	ErrBadRequest = NewErrorCoder(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	ErrTimeout  = NewErrorCoder(http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout))

	codes []int
)

type ReturnCoder interface {
	Code() int
	Message() string
	Data() interface{}
}

type returnCode struct {
	C int         `json:"Code"`
	M string      `json:"Message,omitempty"`
	D interface{} `json:"Data,omitempty"`
}

func (rc returnCode) Code() int {
	return rc.C
}

func (rc returnCode) Message() string {
	return rc.M
}

func (rc returnCode) Data() interface{} {
	return rc.D
}

var _ ReturnCoder = &returnCode{}

type ErrorCoder interface {
	ReturnCoder
	error
}

type errorCode struct {
	returnCode
}

func (ec errorCode) Error() string {
	return ec.Message()
}

var _ ErrorCoder = &errorCode{}

func of(code int, message string) ReturnCoder {
	checkCode(code)
	return &returnCode{C: code, M: message}
}

func NewErrorCoder(code int, message string) ErrorCoder {
	checkCode(code)
	return &errorCode{
		returnCode: returnCode{C: code, M: message},
	}
}

func Fail(i interface{}) ErrorCoder {
	switch i.(type) {
	case error:
		err := errors.Cause(i.(error))
		if err, ok := err.(ErrorCoder); ok {
			return err
		}
		return &errorCode{
			returnCode: returnCode{C: ErrSystem.Code(), M: err.Error()},
		}
	default:
		return &errorCode{
			returnCode: returnCode{C: ErrSystem.Code(), M: fmt.Sprint(i)},
		}
	}
}

func Succ(data interface{}) ReturnCoder {
	return &returnCode{C: Success.Code(), D: data}
}

var lock sync.Mutex

// 检查code是否已经重复
func checkCode(code int) {
	lock.Lock()
	defer lock.Unlock()

	for _, c := range codes {
		if c == code {
			panic("Duplicate code = " + strconv.Itoa(code))
		}
	}

	codes = append(codes, code)
}
