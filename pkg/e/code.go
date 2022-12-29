package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
	HANDLE_FAIL    = 600

	ACL_NO_AUTH           = 405
	ACCOUNT_NEED_AUDIT    = 406
	REFRER_ILLEGAL        = 407
	LOGIN_USER_IS_DESTORY = 408
	LOGIN_USER_IS_REFUSE  = 409

	EXIST_TAG         = 10001
	NOT_EXIST_TAG     = 10002
	NOT_EXIST_ARTICLE = 10003

	AUTH_CHECK_TOKEN_FAIL    = 20001
	AUTH_CHECK_TOKEN_TIMEOUT = 20002
	AUTH_TOKEN               = 20003
	ERROR_AUTH               = 20004
)
