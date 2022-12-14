package middleware

import
(
	"net/http"
)

type CheckTokenMiddleware struct {
}

func NewCheckTokenMiddleware() *CheckTokenMiddleware {
	return &CheckTokenMiddleware{}
}

func (m *CheckTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//if _,ok := r.Context().Value("userId").(json.Number);!ok {
		//		httpx.WriteJson(w,http.StatusOK,errorx.NewCodeDefaultError("token 不存在"))
		//		return
		//}
		next(w, r)
	}
}
