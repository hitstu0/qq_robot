package user

type UserFuncHandler interface {
	ParseParm(parm *string) error
	HandlerRequest() (*string, error)
}