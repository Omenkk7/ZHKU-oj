package impl

import (
	"zhku-oj/internal/repository/interfaces"
)

/*
@Author: omenkk7
@Date: 2025/9/14 16:13
@Description:
*/

type SubmitServiceImpl struct {
	repo interfaces.SubmitRepository
}

func NewSubmitServicfeImpl(repo interfaces.SubmitRepository) *SubmitServiceImpl {
	return &SubmitServiceImpl{repo: repo}
}
