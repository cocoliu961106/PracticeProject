package service

import (
	"PracticeProject/define"
	"PracticeProject/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetCategoryList
// @Tags  公共方法
// @Summary 分类列表
// @Param page query int false "当前页，默认第一页"
// @Param size query int false "size，默认20"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-list [get]
func GetCategoryList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	keyword := c.Query("keyword")
	if err != nil {
		log.Println("GetProblemList Page strconv Error:", err)
		return
	}
	startPoint := (page - 1) * size
	categoryList, count, err := models.GetCategoryList(startPoint, size, keyword)
	if err != nil {
		log.Println("GetCategoryList Error", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "获取分类列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  categoryList,
			"count": count,
		},
	})
	return
}

// CategoryCreate
// @Tags 管理员私有方法
// @Summary 分类创建
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parentId formData int false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-create [post]
func CategoryCreate(c *gin.Context) {
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parentId"))
	err := models.CategoryCreate(name, parentId)
	if err != nil {
		log.Println("CategoryCreate Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建分类失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// CategoryModify
// @Tags 管理员私有方法
// @Summary 分类修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param name formData string true "name"
// @Param parentId formData int false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-modify [put]
func CategoryModify(c *gin.Context) {
	identity := c.PostForm("identity")
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parentId"))
	if name == "" || identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}

	cnt, _ := models.GetCategoryCount(identity)
	if cnt == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类不存在",
		})
		return
	}
	err := models.CategoryModify(identity, name, parentId)
	if err != nil {
		log.Println("CategoryModify Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "修改分类失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

// CategoryDelete
// @Tags 管理员私有方法
// @Summary 分类删除
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-delete [delete]
func CategoryDelete(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	ccnt, _ := models.GetCategoryCount(identity)
	if ccnt == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类不存在",
		})
		return
	}
	cnt, err := models.CountProblemByCategory(identity)
	if err != nil {
		log.Println("Get ProblemCategory Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类关联的问题失败",
		})
		return
	}
	if cnt > 0 { // 若欲删除分类的id下已有关联问题则不允许删除
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类下面已存在问题，不可删除",
		})
		return
	}
	if cnt > 0 { // 若欲删除分类的id下已有关联问题则不允许删除
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类下面已存在问题，不可删除",
		})
		return
	}
	err = models.CategoryDelete(identity)
	if err != nil {
		log.Println("Delete CategoryBasic Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
