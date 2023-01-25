package router

//
//import (
//	"BackendSimple/auth"
//	"BackendSimple/db"
//	"github.com/gofiber/fiber/v2"
//	"github.com/golang-jwt/jwt/v4"
//)
//
//type Middleware struct {
//	Auth auth.Auth
//	DB   db.DB
//}
//
//func (m Middleware) SilentRefresh(refrsh, client string) (auth.TokenPair, error) {
//	token := db.RefreshToken{
//		Refresh: refrsh,
//	}
//	if !m.DB.CheckAndGetTocken(&token) {
//		return auth.TokenPair{}, RefreshIsObsentError
//	}
//
//	err := m.DB.DelToken(&token)
//	if err != nil {
//		return auth.TokenPair{}, err
//	}
//
//	newTockenPair, err := m.Auth.GenerateToken(token.UserID, token.UserLevel, client)
//	if err != nil {
//		return auth.TokenPair{}, err
//	}
//
//	return newTockenPair, nil
//}
//
//func silentRefreshEntry(c *fiber.Ctx, m Middleware, refresh string) error {
//	tocks, err := m.SilentRefresh(refresh, c.GetReqHeaders()["User-Agent"])
//	if err != nil {
//		return Drop401Error(c)
//	}
//	c.ClearCookie("refresh", "access")
//	sendKeys(c, tocks)
//	return nil
//}
//
//func (m Middleware) GetAccessCheck(only []int) func(c *fiber.Ctx) error {
//	return func(c *fiber.Ctx) error {
//		refresh := c.Cookies("refresh")
//		access := c.Cookies("access")
//
//		if refresh == "" && access == "" {
//			return Drop401Error(c)
//		}
//		if access == "" {
//			err := silentRefreshEntry(c, m, refresh)
//			if err != nil {
//				return Drop500Error(c, err)
//			}
//		}
//
//		cl, err := m.Auth.GetTokenClaims(access)
//		if err == jwt.ErrTokenExpired {
//			err := silentRefreshEntry(c, m, refresh)
//			if err != nil {
//				return Drop500Error(c, err)
//			}
//		} else if err != nil {
//			return Drop401Error(c)
//		}
//
//		acclelev, accok := cl["alv"]
//		userid, uidok := cl["sub"]
//		if !accok || !uidok {
//			return Drop401Error(c)
//		}
//
//		level := int(acclelev.(float64))
//
//		c.Locals("userid", uint(userid.(float64)))
//		c.Locals("level", level)
//
//		if !contains(only, level) {
//			return Drop401Error(c)
//		}
//
//		return c.Next()
//	}
//}
