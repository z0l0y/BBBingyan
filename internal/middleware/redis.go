package middleware

import (
	utils "BBBingyan/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"
)

var redisClient *redis.Client
var cfg, _ = utils.LoadConfig("config/application.yaml")

func RedisInit() {
	config := &redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
	}
	redisClient = redis.NewClient(config)
}

func SendVerificationCode(c *gin.Context) {
	email := c.PostForm("email")

	println("email:", email)
	// 检查用户是否在一分钟内发送过验证码
	if err := checkSendFrequency(c, email); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 生成新的验证码
	verificationCode := generateVerificationCode()
	println("verificationCode:", verificationCode)
	// 将验证码保存到 Redis 中,并设置过期时间为 5 分钟
	if err := storeVerificationCode(c, email, verificationCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 发送验证码到用户邮箱
	if err := sendVerificationCodeByEmail(email, verificationCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification code sent successfully",
	})
}

func VerifyVerificationCode(c *gin.Context) {
	email := c.PostForm("email")
	code := c.PostForm("code")

	// 从 Redis 中获取用户的验证码
	storedCode, err := getVerificationCode(c, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid verification code",
		})
		return
	}

	// 比较用户输入的验证码是否正确
	if storedCode != code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid verification code",
		})
		return
	}

	// 验证码验证成功,删除 Redis 中的验证码
	if err := deleteVerificationCode(c, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回登录成功的响应
	utils.Success(c, "success", "Login successful")
}

func checkSendFrequency(c *gin.Context, email string) error {
	// 检查用户在一分钟内是否发送过验证码
	lastSentTime, err := getLastSentTime(c, email)
	if err == nil && time.Since(lastSentTime) < time.Minute {
		return fmt.Errorf("you can only request a verification code once per minute")
	}
	return nil
}

func generateVerificationCode() string {
	// 生成 6 位数的随机验证码
	return fmt.Sprintf("%06d", rand.Intn(900000)+100000)
}

func storeVerificationCode(c *gin.Context, email, verificationCode string) error {
	// 将验证码存储在 Redis 中,并设置 5 分钟的过期时间
	key := fmt.Sprintf("verification_code:%s", email)
	if err := redisClient.Set(c.Request.Context(), key, verificationCode, 5*time.Minute).Err(); err != nil {
		return err
	}
	// 记录最后一次发送验证码的时间
	if err := setLastSentTime(c, email, time.Now()); err != nil {
		return err
	}
	return nil
}

func getLastSentTime(c *gin.Context, email string) (time.Time, error) {
	key := fmt.Sprintf("last_sent_time:%s", email)
	lastSentTimeStr, err := redisClient.Get(c.Request.Context(), key).Result()
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, lastSentTimeStr)
}

func setLastSentTime(c *gin.Context, email string, t time.Time) error {
	key := fmt.Sprintf("last_sent_time:%s", email)
	return redisClient.Set(c.Request.Context(), key, t.Format(time.RFC3339), 0).Err()
}

func getVerificationCode(c *gin.Context, email string) (string, error) {
	key := fmt.Sprintf("verification_code:%s", email)
	return redisClient.Get(c.Request.Context(), key).Result()
}

func deleteVerificationCode(c *gin.Context, email string) error {
	key := fmt.Sprintf("verification_code:%s", email)
	return redisClient.Del(c.Request.Context(), key).Err()
}

func sendVerificationCodeByEmail(to, code string) error {
	from := cfg.Email.Value
	password := cfg.Email.Password
	body := fmt.Sprintf("Your verification code is: %s", code)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\n \r\n\r\n%s",
		from, to, body,
	))
	println(string(msg))
	// 连接 SMTP 服务器
	println("正在连接服务器")
	auth := smtp.PlainAuth("", from, password, cfg.Email.Host)
	err := smtp.SendMail(cfg.Email.Addr, auth, from, []string{to}, msg)
	if err != nil {
		println("连接服务器失败")
		return err
	}
	println("连接成功")
	return nil
}
