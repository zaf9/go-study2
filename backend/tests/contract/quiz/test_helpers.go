package contract

import (
	"reflect"
	"unsafe"

	"go-study2/internal/app/http_server/handler"
	appquiz "go-study2/internal/app/quiz"
)

// setQuizService 使用反射注入测验服务
func setQuizService(h *handler.Handler, svc *appquiz.Service) {
	val := reflect.ValueOf(h).Elem().FieldByName("quizService")
	ptr := unsafe.Pointer(val.UnsafeAddr())
	reflect.NewAt(val.Type(), ptr).Elem().Set(reflect.ValueOf(svc))
}
