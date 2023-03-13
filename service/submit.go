package service

import (
	"PracticeProject/define"
	"PracticeProject/helper"
	"PracticeProject/models"
	codeExecutor "PracticeProject/service/code_executor"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// GetSubmitList
// @Tags  公共方法
// @Summary 提交列表
// @Param page query int false "当前页，默认第一页"
// @Param size query int false "size，默认20"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query string false "status"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	status, _ := strconv.Atoi(c.Query("status"))
	if err != nil {
		log.Println("GetProblemList Page strconv Error:", err)
		return
	}
	startPoint := (page - 1) * size

	var count int64
	list := make([]models.SubmitBasic, 0)

	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")

	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Offset(startPoint).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Submit List Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/submit [post]
func Submit(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	code, err := ioutil.ReadAll(c.Request.Body) // 获取用户提交的代码内容
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Read Code Error" + err.Error(),
		})
		return
	}

	// 核心功能 - 代码判断
	problemBasic, err := models.GetProblemTestCase(problemIdentity) // 获取提交问题的测试用例
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Problem Error:" + err.Error(),
		})
		return
	}
	// 保存代码到本地
	path, dir, err := helper.CodeSave(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Code Save Error:" + err.Error(),
		})
		return
	}
	u, _ := c.Get("user")
	userClaim := u.(*helper.UserClaims)
	submitBasic := &models.SubmitBasic{
		Identity:        helper.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaim.Identity,
		Path:            path,
	}

	msg := codeExecutor.GoExecuteTestCases(problemBasic, submitBasic, path)

	// 代码执行完后，删除存在本地的代码文件
	// TODO 后续区分环境后，只有dev环境要删除
	go func() {
		err := helper.CodeDelete(dir)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	// 数据保存到数据库
	if err = models.Submit(submitBasic, userClaim, problemIdentity); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Submit Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"status": submitBasic.Status,
			"msg":    msg,
		},
	})
}
