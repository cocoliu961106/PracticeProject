package service

import (
	"PracticeProject/define"
	"PracticeProject/helper"
	"PracticeProject/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// 业务逻辑-获取问题列表

// GetProblemList
// @Tags  公共方法
// @Summary 问题列表
// @Param page query int false "当前页，默认第一页"
// @Param size query int false "size，默认20"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage)) // 字符串转换成int类型
	if err != nil {
		log.Println("GetProblemList Page strconv Error:", err)
		return
	}
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	startPoint := (page - 1) * size // 列表数据的起始值
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")

	list, count, err := models.GetProblemList(startPoint, size, keyword, categoryIdentity)
	if err != nil {
		log.Println("Get Problem List Error:", err)
	}
	// c.String(http.StatusOK, "Get Problem List")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"data":  list,
			"count": count,
		},
	})
}

// GetProblemDetail
// @Tags  公共方法
// @Summary 问题详情
// @Param  identity query string false "problem identity"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}
	tx, detail := models.GetProblemDetail(identity)
	err := tx.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get ProblemDetail Error:" + err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": detail,
	})

}

// ProblemCreate
// @Tags 管理员私有方法
// @Summary 问题创建
// @Accept json
// @Param authorization header string true "authorization"
// @Param data body define.ProblemBasic true "ProblemBasic"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
func ProblemCreate(c *gin.Context) {
	in := new(define.ProblemBasic)
	err := c.ShouldBindJSON(in) // 请求body中的数据能否符合ProblemBasic结构体的格式
	if err != nil {
		log.Println("[JsonBind Error] : ", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}

	if in.Title == "" || in.Content == "" || len(in.ProblemCategories) == 0 || len(in.TestCases) == 0 || in.MaxRuntime == 0 || in.MaxMem == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	identity := helper.GetUUID()
	data := &models.ProblemBasic{
		Identity:   identity,
		Title:      in.Title,
		Content:    in.Content,
		MaxRuntime: in.MaxRuntime,
		MaxMem:     in.MaxMem,
	}

	// 处理分类:对于新增问题的分类列表，其中每一个分类应该在ProblemCategory表里面对应一条数据
	categoryBasics := make([]*models.ProblemCategory, 0)
	for _, id := range in.ProblemCategories {
		categoryBasics = append(categoryBasics, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(id),
		})
	}
	data.ProblemCategories = categoryBasics

	// 处理测试用例
	testCaseBasics := make([]*models.TestCase, 0)
	for _, v := range in.TestCases {
		// 举个例子 {"input":"1 2\n","output":"3\n"}
		testCaseBasic := &models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: identity,
			Input:           v.Input,
			Output:          v.Output,
		}
		testCaseBasics = append(testCaseBasics, testCaseBasic)
	}
	data.TestCases = testCaseBasics // 测试用例插入到data中

	// 创建问题
	err = models.ProblemCreate(data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Create Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": data.Identity,
		},
	})
}

// ProblemModify
// @Tags 管理员私有方法
// @Summary 问题修改
// @Param authorization header string true "authorization"
// @Param data body define.ProblemBasic true "ProblemBasic"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-modify [put]
func ProblemModify(c *gin.Context) {
	in := new(define.ProblemBasic)
	err := c.ShouldBindJSON(in)
	if err != nil {
		log.Println("[JsonBind Error] : ", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}
	if in.Identity == "" || in.Title == "" || in.Content == "" || len(in.ProblemCategories) == 0 || len(in.TestCases) == 0 || in.MaxRuntime == 0 || in.MaxMem == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	if cnt, _ := models.GetProblemCount(in.Identity); cnt == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "不存在该问题",
		})
		return
	}
	if err := models.ProblemModify(in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Modify Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题修改成功",
	})
}
