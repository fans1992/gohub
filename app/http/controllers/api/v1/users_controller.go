package v1

import (
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/auth"
	"gohub/pkg/config"
	"gohub/pkg/file"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	BaseAPIController
}

// CurrentUser 当前登录用户信息
func (ctrl *UsersController) CurrentUser(c *gin.Context) {
	userModel := auth.CurrentUser(c)
	response.Data(c, userModel)
}

// Index 所有用户
func (ctrl *UsersController) Index(c *gin.Context) {
	request := requests.PaginationRequest{}
	if ok := requests.Validate(c, &request, requests.Pagination); !ok {
		return
	}

	data, paper := user.Paginate(c, 10)
	response.JSON(c, gin.H{
		"data":  data,
		"paper": paper,
	})
}

func (ctrl *UsersController) UpdateProfile(c *gin.Context) {
	request := requests.UserUpdateProfileRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdateProfile); !ok {
		return
	}

	currentUser := auth.CurrentUser(c)
	currentUser.Name = request.Name
	currentUser.Introduction = request.Introduction
	currentUser.City = request.City
	rowsAffected := currentUser.Save()

	if rowsAffected > 0 {
		response.Data(c, currentUser)
		return
	}

	response.Abort500(c, "修改个人信息失败， 请稍后再试~")
}

func (ctrl *UsersController) UpdateEmail(c *gin.Context) {
	request := requests.UserUpdateEmailRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdateEmail); !ok {
		return
	}

	currentUser := auth.CurrentUser(c)
	currentUser.Email = request.Email
	rowsAffected := currentUser.Save()

	if rowsAffected > 0 {
		response.Success(c)
		return
	}

	// 失败，显示错误提示
	response.Abort500(c, "更新失败，请稍后重试~")
}

func (ctrl *UsersController) UpdatePhone(c *gin.Context) {

	request := requests.UserUpdatePhoneRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdatePhone); !ok {
		return
	}

	currentUser := auth.CurrentUser(c)
	currentUser.Phone = request.Phone
	rowsAffected := currentUser.Save()

	if rowsAffected > 0 {
		response.Success(c)
		return
	}

	response.Abort500(c, "手机号更新失败， 请稍候再试")
}

func (ctrl *UsersController) UpdatePassword(c *gin.Context) {

	request := requests.UserUpdatePasswordRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdatePassword); !ok {
		return
	}

	currentUser := auth.CurrentUser(c)
	// 验证原始密码是否正确
	if _, err := auth.Attempt(currentUser.Name, request.Password); err != nil {
		// 失败，显示错误提示
		response.Unauthorized(c, "原密码不正确")
		return
	}

	// 更新密码为新密码
	currentUser.Password = request.NewPassword
	currentUser.Save()

	response.Success(c)
}

func (ctrl *UsersController) UpdateAvatar(c *gin.Context) {

	request := requests.UserUpdateAvatarRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdateAvatar); !ok {
		return
	}

	avatar, err := file.SaveUploadAvatar(c, request.Avatar)
	if err != nil {
		response.Abort500(c, "上传头像失败，请稍后尝试~")
		return
	}

	currentUser := auth.CurrentUser(c)
	currentUser.Avatar = config.GetString("app.url") + avatar
	currentUser.Save()

	response.Data(c, currentUser)

}
