package models

import "google.golang.org/grpc/codes"

type CustomErr struct {
	Code codes.Code
	Msg  string
}
