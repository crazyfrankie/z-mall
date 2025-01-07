package web

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"zmall/common/bizerr"
	"zmall/common/constant"
	"zmall/common/response"
	"zmall/server/middleware"
	"zmall/server/user/domain"
	"zmall/server/user/service"
)

type UserHandler struct {
	svc      *service.UserService
	jwt      *middleware.JWTHandler
	stateKey []byte
}

type StateClaims struct {
	State string
	jwt.StandardClaims
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	wsConnection = make(map[string]*websocket.Conn)
)

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (ctl *UserHandler) RegisterRoute(r *gin.Engine) {
	userGroup := r.Group("api/user")
	{
		userGroup.POST("login", ctl.Login())
		userGroup.Any("callback", ctl.Callback())
	}
}

func (ctl *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		state := uuid.New().String()

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			response.Error(c, bizerr.NewBizError(constant.Success))
			return
		}

		wsConnection[state] = conn

		var url string
		url, err = ctl.svc.AuthUrl(c.Request.Context(), state)
		if err != nil {
			response.Error(c, bizerr.NewBizError(constant.Success))
			return
		}

		if err := ctl.SetCookie(c, state); err != nil {
			response.Error(c, bizerr.NewBizError(constant.Success))
			return
		}

		if err := conn.WriteJSON(map[string]string{
			"auth_url":   url,
			"session_id": state,
		}); err != nil {
			response.Error(c, bizerr.NewBizError(constant.Success))
		}

		defer func() {
			err := conn.Close()
			delete(wsConnection, state)
			if err != nil {
				panic(err)
			}
		}()
	}
}

func (ctl *UserHandler) Callback() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")

		state, err := ctl.VerifyState(c)
		if err != nil {
			response.Error(c, bizerr.NewBizError(constant.Success))
			return
		}

		info, er := ctl.svc.VerifyCode(c.Request.Context(), code)
		if er != nil {
			response.Error(c, bizerr.NewBizError(constant.Success))
			return
		}

		var user domain.User
		user, err = ctl.svc.FindOrCreateUser(c.Request.Context(), info)
		if err != nil {
			response.Error(c, bizerr.NewBizError(constant.Success))
			return
		}

		conn := wsConnection[state]

		var token string
		token, err = ctl.jwt.SetToken(user.ID, user.Role)
		if err != nil {
			conn.WriteJSON(map[string]string{
				"status":  "error",
				"message": "failed to generate JWT",
			})
			conn.Close()
		}

		conn.WriteJSON(map[string]string{
			"status": "success",
			"token":  token,
		})

		// 关闭连接
		conn.Close()
		delete(wsConnection, state)
	}
}

func (ctl *UserHandler) SetCookie(c *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	})
	tokenStr, err := token.SignedString(ctl.stateKey)
	if err != nil {
		return err
	}
	c.SetCookie("jwt-state", tokenStr, 600, "/api/user/callback", "", false, true)

	return nil
}

func (ctl *UserHandler) VerifyState(c *gin.Context) (string, error) {
	state := c.Query("state")
	jwtState, err := c.Cookie("jwt-state")
	if err != nil {
		return "", fmt.Errorf("拿不到 state 的 cookie, %w", err)
	}

	var sc StateClaims
	token, err := jwt.ParseWithClaims(jwtState, &sc, func(token *jwt.Token) (interface{}, error) {
		return ctl.stateKey, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("token 已经过期, %w", err)
	}

	if sc.State != state {
		return "", errors.New("state 不相等")
	}

	return state, nil
}
