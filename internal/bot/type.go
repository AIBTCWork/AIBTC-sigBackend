package sms

import "context"

type botServiceInterface interface {
	Send(c context.Context, tpl string, args []string, numbers ...string) error
}
