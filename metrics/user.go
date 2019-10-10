package metrics

import (
	"github.com/xesina/golang-echo-realworld-example-app/model"
	"github.com/xesina/golang-echo-realworld-example-app/user"
)

type instrumentedUserStore struct {
	s user.Store
	m *StoreMetrics
}

func NewUserStore(s user.Store, m *StoreMetrics) user.Store {
	return &instrumentedUserStore{
		s: s,
		m: m,
	}
}

func (s *instrumentedUserStore) GetByID(id uint) (u *model.User, err error) {
	s.m.wrapRequest(opRead, func() bool {
		u, err = s.s.GetByID(id)
		return err == nil
	})
	return u, err
}

func (s *instrumentedUserStore) GetByEmail(e string) (u *model.User, err error) {
	s.m.wrapRequest(opRead, func() bool {
		u, err = s.s.GetByEmail(e)
		return err == nil
	})
	return u, err
}

func (s *instrumentedUserStore) GetByUsername(username string) (u *model.User, err error) {
	s.m.wrapRequest(opRead, func() bool {
		u, err = s.s.GetByUsername(username)
		return err == nil
	})
	return u, err
}

func (s *instrumentedUserStore) Create(u *model.User) (err error) {
	s.m.wrapRequest(opCreate, func() bool {
		err = s.s.Create(u)
		return err == nil
	})
	return err
}

func (s *instrumentedUserStore) Update(u *model.User) (err error) {
	s.m.wrapRequest(opUpdate, func() bool {
		err = s.s.Update(u)
		return err == nil
	})
	return err
}

func (s *instrumentedUserStore) AddFollower(u *model.User, followerID uint) (err error) {
	s.m.wrapRequest(opCreate, func() bool {
		err = s.s.AddFollower(u, followerID)
		return err == nil
	})
	return err
}

func (s *instrumentedUserStore) RemoveFollower(u *model.User, followerID uint) (err error) {
	s.m.wrapRequest(opDelete, func() bool {
		err = s.s.RemoveFollower(u, followerID)
		return err == nil
	})
	return err
}

func (s *instrumentedUserStore) IsFollower(userID, followerID uint) (b bool, err error) {
	s.m.wrapRequest(opRead, func() bool {
		b, err = s.s.IsFollower(userID, followerID)
		return err == nil
	})
	return b, err
}
