package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sms-dashboard/internal/config"
	"sms-dashboard/internal/database"
	"sms-dashboard/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SMSHandler struct{}

func NewSMSHandler() *SMSHandler {
	return &SMSHandler{}
}

func (h *SMSHandler) Receive(c *gin.Context) {
	var req struct {
		Content     string `json:"content" form:"content" binding:"required"`
		SendTime    string `json:"sendTime" form:"sendTime"`
		OrgContent  string `json:"org_content" form:"org_content"`
		Sign        string `json:"sign" form:"sign"`
		Timestamp   string `json:"timestamp" form:"timestamp"`
		ReceiveTime string `json:"receive_time" form:"receive_time"`
	}

	// For JSON requests, we can also use ShouldBindJSON
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	cfg := config.LoadConfig()
	//校验Sign
	// 1. 获取毫秒时间戳
	timestamp := req.Timestamp
	secret := cfg.Secret

	// 2. 拼接待签名字符串 (timestamp + \n + secret)
	stringToSign := timestamp + "\n" + secret

	// 3. HmacSHA256 加密
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(stringToSign))
	signData := hash.Sum(nil)

	// 4. Base64 编码
	base64Sign := base64.StdEncoding.EncodeToString(signData)

	// 5. URL Encode (对应 Java 的 URLEncoder.encode)
	sign := url.QueryEscape(base64Sign)

	// 打印结果
	fmt.Printf("Timestamp: %v\n", timestamp)
	fmt.Printf("Final Sign: %v\n", sign)

	//如果不对则返回错误
	if sign != req.Sign {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Sign"})
		return
	}
	sms := model.SMS{
		Content:  req.Content,
		SendTime: req.SendTime,
	}

	if sms.SendTime == "" {
		sms.SendTime = time.Now().Format("2006-01-02 15:04:05")
	}

	// 对 content 进行 URLDecode 处理
	if decodedContent, err := url.QueryUnescape(sms.Content); err == nil {
		sms.Content = decodedContent
	}

	if err := database.DB.Create(&sms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "id": sms.ID})
}

func (h *SMSHandler) List(c *gin.Context) {
	var smsList []model.SMS
	var total int64
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	database.DB.Model(&model.SMS{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&smsList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": smsList,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}

func (h *SMSHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := database.DB.Delete(&model.SMS{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
