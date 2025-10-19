package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"takeaway-go/internal/app/initializer/redis"
	"time"
)

const (
	// AccessTokenExpire 访问token过期时间
	AccessTokenExpire = 2 * time.Hour
	// RefreshTokenExpire 刷新token过期时间
	RefreshTokenExpire = 7 * 24 * time.Hour
)

// Claims 自定义 JWT 载荷结构体
type Claims struct {
	UserID     uint64 `json:"userID"`     // 用户 ID
	GrantScope string `json:"grantScope"` // 签发对象
	TokenType  string `json:"tokenType"`  // tokenType: access, refresh
	jwt.RegisteredClaims
}

// TokenResponse token 响应结构体
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // 过期时间, 单位秒
}

// GenerateTokenPair 生成一对 accessToken 和 refreshToken
func GenerateTokenPair(userID uint64, subject string, jwtSecret string) (*TokenResponse, error) {
	byteSecret := []byte(jwtSecret)
	// 生成 accessToken
	accessClaims := Claims{
		UserID:     userID,
		GrantScope: subject,
		TokenType:  "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                         // 签发者
			Subject:   subject,                                               // 签发对象
			Audience:  jwt.ClaimStrings{"PC", "Wechat_Program"},              // 签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpire)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                        // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                        // 生效时间
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(byteSecret)
	if err != nil {
		return nil, err
	}

	// 生成 refreshToken
	refreshClaims := Claims{
		UserID:     userID,
		GrantScope: subject,
		TokenType:  "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                          // 签发者
			Subject:   subject,                                                // 签发对象
			Audience:  jwt.ClaimStrings{"PC", "Wechat_Program"},               // 签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpire)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                         // 生效时间
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(byteSecret)
	if err != nil {
		return nil, err
	}

	// 将 token 缓存到 Redis
	if err := StoreTokenInRedis(userID, accessToken, refreshToken); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(AccessTokenExpire.Seconds()),
	}, nil
}

func StoreTokenInRedis(userID uint64, accessToken string, refreshToken string) error {
	ctx := context.Background()

	// 使用 hash 存储用户 token 信息
	userKey := fmt.Sprintf("user_tokens:%d", userID)
	err := redis.RedisClient.HSet(ctx, userKey, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"created_at":    time.Now().Unix(),
	}).Err()
	if err != nil {
		return err
	}

	// 设置过期时间为 refresh_token 的过期时间
	err = redis.RedisClient.Expire(ctx, userKey, RefreshTokenExpire).Err()
	if err != nil {
		return err
	}

	return nil
}

// GenerateToken 兼容旧接口, 返回 accessToken
func GenerateToken(userID uint64, subject string, jwtSecret string) (string, error) {
	tokenResponse, err := GenerateTokenPair(userID, subject, jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenResponse.AccessToken, nil
}

// ParseToken 解析和校验JWT Token
func ParseToken(tokenString string, jwtSecret string) (*Claims, error) {
	byteSecret := []byte(jwtSecret)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return byteSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 检查token是否在Redis中被撤销
		if !IsTokenValidInRedis(claims.UserID, tokenString, claims.TokenType) {
			return nil, errors.New("token已被撤销")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// IsTokenValidInRedis 检查token是否在Redis中有效
func IsTokenValidInRedis(userID uint64, token string, tokenType string) bool {
	ctx := context.Background()
	userKey := fmt.Sprintf("user_tokens:%d", userID)

	var redisToken string
	var err error

	if tokenType == "access" {
		redisToken, err = redis.RedisClient.HGet(ctx, userKey, "access_token").Result()
	} else {
		redisToken, err = redis.RedisClient.HGet(ctx, userKey, "refresh_token").Result()
	}

	if err != nil {
		return false
	}

	return redisToken == token
}

// RefreshAccessToken 使用刷新token生成新的访问token
func RefreshAccessToken(refreshToken string, jwtSecret string) (*TokenResponse, error) {
	// 解析刷新token
	claims, err := ParseToken(refreshToken, jwtSecret)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("无效的刷新token")
	}

	// 生成新的token对
	return GenerateTokenPair(claims.UserID, claims.GrantScope, jwtSecret)
}

// RevokeToken 撤销用户的所有token
func RevokeToken(userID uint64) error {
	ctx := context.Background()
	userKey := fmt.Sprintf("user_tokens:%d", userID)
	return redis.RedisClient.Del(ctx, userKey).Err()
}

// RevokeAllUserTokens 撤销所有用户的token（用于安全事件）
func RevokeAllUserTokens() error {
	ctx := context.Background()
	// 删除所有用户token
	keys, err := redis.RedisClient.Keys(ctx, "user_tokens:*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return redis.RedisClient.Del(ctx, keys...).Err()
	}
	return nil
}
