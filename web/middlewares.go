package web

import (
	"Backend/Errors"
	"Backend/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (r *Router) RequestCounter(c *gin.Context) {
	err := r.Metrics.Actions("any_request")
	if err != nil {
		r.Error.Http_MetricsError(err)
	}
	c.Next()
}

func (r *Router) UnseriousIDScanning(c *gin.Context) {
	access, err := c.Cookie("access")
	accessPresent := true
	if err != nil {
		accessPresent = false
	}

	refresh, err := c.Cookie("refresh")
	refreshPresent := true
	if err != nil {
		refreshPresent = false
	}

	action, err := c.Cookie("action")
	actionPresent := true
	if err != nil {
		actionPresent = false
	}

	c.Set("ACCESS_PRESENT", accessPresent)
	c.Set("REFRESH_PRESENT", refreshPresent)
	c.Set("ACTION_PRESENT", actionPresent)

	if refreshPresent {
		c.Set("REFRESH", refresh)
	}
	if accessPresent {
		c.Set("ACCESS", access)
		claims, err := r.Auth.GetAccessClaims(access)
		accessExpired := false
		accessInvalid := false
		userID := uint(0)
		userLevel := 0

		if err == jwt.ErrTokenExpired {
			accessExpired = true
		} else if err != nil {
			accessInvalid = true
		} else {
			userID = claims.Subject
			userLevel = claims.LevelOfAccess
		}

		c.Set("ID", userID)
		c.Set("LEVEL", userLevel)
		c.Set("ACCESS_EXPIRED", accessExpired)
		c.Set("ACCESS_INVALID", accessInvalid)
	}
	if actionPresent {
		c.Set("ACTION", action)
		actionClaims, err := r.Auth.GetActionClaims(action)
		actionExpired := false
		actionInvalid := false
		var actionOperations []string

		if err == jwt.ErrTokenExpired {
			actionExpired = true
		} else if err != nil {
			actionInvalid = true
		} else {
			actionOperations = actionClaims.Actions
		}
		c.Set("ACTION_EXPIRED", actionExpired)
		c.Set("ACTION_INVALID", actionInvalid)

		c.Set("ACTION_OPERATION", actionOperations)
	}

	c.Next()
}

func (r *Router) RequestLog(c *gin.Context) {
	accessPresent := c.GetBool("ACCESS_PRESENT")
	refreshPresent := c.GetBool("REFRESH_PRESENT")

	var accessExpired, accessInvalid bool
	var userID uint
	var userLevel int
	var flag []zap.Field
	if accessPresent {
		accessExpired = c.GetBool("ACCESS_EXPIRED")
		accessInvalid = c.GetBool("ACCESS_INVALID")

		userID = c.GetUint("ID")
		userLevel = c.GetInt("LEVEL")

		flag = []zap.Field{{
			Key:     "UserID",
			Type:    zapcore.Uint32Type,
			Integer: int64(userID),
		}, {
			Key:     "User level",
			Type:    zapcore.Int32Type,
			Integer: int64(userLevel),
		}}
	}

	flag = append(flag, zap.Field{
		Key:    "Path",
		Type:   zapcore.StringType,
		String: c.FullPath(),
	}, zap.Field{
		Key:    "Access token present",
		Type:   zapcore.BoolType,
		String: fmt.Sprint(accessPresent),
	}, zap.Field{
		Key:    "Refresh token present",
		Type:   zapcore.BoolType,
		String: fmt.Sprint(refreshPresent),
	}, zap.Field{
		Key:    "Access expired",
		Type:   zapcore.BoolType,
		String: fmt.Sprint(accessExpired),
	}, zap.Field{
		Key:    "Access invalid",
		Type:   zapcore.BoolType,
		String: fmt.Sprint(accessInvalid),
	})

	r.Log.Debug("New request", flag...)
}

func (r *Router) Refresh(c *gin.Context) {
	accessPresent := c.GetBool("ACCESS_PRESENT")
	refreshPresent := c.GetBool("REFRESH_PRESENT")

	if !refreshPresent {
		c.Next()
		return
	}

	var refresh string
	var accessExpired, accessInvalid bool
	if refreshPresent {
		refresh = c.GetString("REFRESH")
	}
	if accessPresent {
		accessExpired = c.GetBool("ACCESS_EXPIRED")
		accessInvalid = c.GetBool("ACCESS_INVALID")
	}

	if accessPresent && !accessExpired && !accessInvalid {
		c.Next()
		return
	}

	id, err := r.DataBase.GetUserDataFromRefresh(refresh)
	if err != nil {
		r.Error.DropError(c, Errors.TokenIsIvalidError)
		return
	}

	level := r.DataBase.GetUserLevel(id)
	if level == 0 {
		r.Error.DropError(c, Errors.TokenIsIvalidError)
		return
	}

	pair, err := r.Auth.GenerateToken(id, level)
	if err != nil {
		r.Error.DropError(c, err)
		return
	}

	err = r.DataBase.RefreshSession(auth.Token{Key: refresh}, pair.Refresh)
	if err != nil {
		r.Error.DropError(c, Errors.Unauthorizated401Error)
		return
	}

	c.SetCookie("access", pair.Access.Key, int(pair.Access.Expires.Unix()), "/", "*", false, true)
	c.SetCookie("refresh", pair.Refresh.Key, int(pair.Refresh.Expires.Unix()), "/", "*", false, true)

	c.Next()
}

func (r *Router) AccessCheck(c *gin.Context) {
	accessLevels, alOk := r.pathAccess[c.FullPath()]
	if !alOk {
		r.Error.DropError(c, Errors.PageNotExistError)
		return
	}
	if contains(accessLevels, 0) {
		c.Next()
		return
	}

	level := c.GetInt("LEVEL")

	if !contains(accessLevels, level) {
		r.Error.DropError(c, Errors.Unauthorizated401Error)
		return
	}
	c.Next()
}
