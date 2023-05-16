package logic

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/jinzhu/copier"
	"github.com/three-body/hertz-scaffold/biz/dal/model"
	"github.com/three-body/hertz-scaffold/biz/dal/query"
	"github.com/three-body/hertz-scaffold/biz/hmodel/user"
)

type UserLogic struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserLogic(ctx context.Context, c *app.RequestContext) *UserLogic {
	return &UserLogic{
		ctx: ctx,
		c:   c,
	}
}

func (l *UserLogic) Signup(req *user.SignupRequest) (*user.User, error) {
	newUser := &model.User{
		UID:      "tt",
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Nickname: req.GetNickname(),
	}

	err := query.Q.Transaction(func(tx *query.Query) error {
		u := tx.User
		oldUser, err := u.WithContext(l.ctx).Where(u.Email.Eq(req.GetEmail())).Take()
		if err != nil {
			return err
		}
		if oldUser != nil {
			return fmt.Errorf("user account already exists")
		}
		err = u.WithContext(l.ctx).Create(newUser)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.NewPrivate(err.Error()).SetMeta("login failed")
	}

	resp := new(user.User)
	if err = copier.Copy(resp, newUser); err != nil {
		return nil, errors.NewPrivate(err.Error()).SetMeta("login failed")
	}
	return resp, nil
}

func (l *UserLogic) Login(req *user.LoginRequest) (*user.User, error) {
	u := query.User
	record, err := u.WithContext(l.ctx).Where(u.Email.Eq(req.GetEmail()), u.Password.Eq(req.GetPassword())).Take()
	if err != nil {
		return nil, errors.NewPrivate(err.Error()).SetMeta("login failed")
	}
	resp := new(user.User)
	if err = copier.Copy(resp, record); err != nil {
		return nil, errors.NewPrivate(err.Error()).SetMeta("login failed")
	}
	return resp, nil
}
