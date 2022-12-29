package access

// 日志记录到文件
//func LoginCheckMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if c.Request.URL.Path == "/api/v1/health" ||
//			c.Request.URL.Path == "/api/v1/login/sign_out" ||
//			c.Request.URL.Path == "/api/v1/login/simulation_test" {
//			c.Next()
//			return
//		}
//
//		appG := util.Gin{C: c}
//
//		user, err := model.GetLoginUserInfo(c)
//		if err != nil {
//			if err.Error() == e.GetMsg(e.ERROR_LOGIN_NO_COOKIE) {
//				util.Log.Warning(c, "GetLoginUserInfo fail "+err.Error())
//			} else {
//				util.Log.Error(c, "GetLoginUserInfo fail"+err.Error())
//			}
//
//			appG.Response(http.StatusOK, e.ERROR_LOGIN_NO_COOKIE, nil)
//			c.Abort()
//			return
//		}
//		if user.Status == -1 { //账户待审核
//			util.Log.Warning(c, e.GetMsg(e.ERROR_ACCOUNT_NEED_AUDIT)+" "+string(user.Uid))
//			appG.Response(http.StatusOK, e.ERROR_ACCOUNT_NEED_AUDIT, map[string]int{"login_src": user.LoginSrc})
//			c.Abort()
//			return
//		}
//
//		if user.Status == 1 && user.LoginSrc == 1 { //账户待审核，且非模拟登录
//			util.Log.Warning(c, e.GetMsg(e.ERROR_LOGIN_USER_IS_DESTORY)+" "+string(user.Uid))
//			appG.Response(http.StatusOK, e.ERROR_LOGIN_USER_IS_DESTORY, map[string]int{"login_src": user.LoginSrc})
//			c.Abort()
//			return
//		}
//		if user.Status == -2 && user.LoginSrc == 1 { //账户审核拒绝，且非模拟登录
//			util.Log.Warning(c, e.GetMsg(e.ERROR_LOGIN_USER_IS_REFUSE)+" "+string(user.Uid))
//			appG.Response(http.StatusOK, e.ERROR_LOGIN_USER_IS_REFUSE, map[string]int{"login_src": user.LoginSrc})
//			c.Abort()
//			return
//		}
//		//登录信息
//		c.Set(util.LoginUidKey, user.Uid)
//		c.Set(util.LoginUserNameKey, user.LoginName)
//		c.Set(util.LoginUserRoleKey, user.Role)
//		c.Set(util.LoginSrc, user.LoginSrc)
//
//		c.Next()
//	}
//}
