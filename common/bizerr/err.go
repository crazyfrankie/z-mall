package bizerr

import (
	"github.com/crazyfrankie/gem/gerrors"
	
	"zmall/common/constant"
)

type BizError struct {
	gerrors.BizErrorIface
}

func NewBizError(errCode constant.ErrCode) *BizError {
	return &BizError{gerrors.NewBizError(errCode.Code, errCode.Msg)}
}
