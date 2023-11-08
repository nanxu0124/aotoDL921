package web

import (
	"autoDL921/internal/domain"
	"autoDL921/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler 定义和用户有关的所有路由
type UserHandler struct {
	svc         *service.UserService
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	return &UserHandler{
		svc:         svc,
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")

	ug.POST("/signup", u.SignUp)
	ug.POST("/signin", u.SignIn)
	ug.POST("/edit", u.Edit)
	ug.POST("/profile", u.Profile)

}

// SignUp 处理注册用户请求
func (u *UserHandler) SignUp(ctx *gin.Context) {
	type UserReq struct {
		UserId          string `json:"userId"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req UserReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码输入不一致")
		return
	}
	isPassword, err := u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须大于8位，包含特殊字符")
		return
	}

	err = u.svc.SignUp(ctx, domain.User{
		UserId:   req.UserId,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateUserId {
		ctx.String(http.StatusOK, "用户名已存在")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, "注册成功")
}
func (u *UserHandler) SignIn(ctx *gin.Context) {
	type SignReq struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
	}

	var req SignReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.SignIn(ctx, req.UserId, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "账号或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	sess := sessions.Default(ctx)
	sess.Set("userId", user.UserId)
	sess.Save()

	ctx.String(http.StatusOK, "登录成功")
	return
}
func (u *UserHandler) Edit(ctx *gin.Context) {

}
func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是你的profile")
	return
}
