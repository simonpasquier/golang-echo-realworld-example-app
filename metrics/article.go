package metrics

import (
	"github.com/xesina/golang-echo-realworld-example-app/article"
	"github.com/xesina/golang-echo-realworld-example-app/model"
)

type instrumentedArticleStore struct {
	s article.Store
	m *StoreMetrics
}

func NewArticleStore(s article.Store, m *StoreMetrics) article.Store {
	return &instrumentedArticleStore{
		s: s,
		m: m,
	}
}

func (as *instrumentedArticleStore) GetBySlug(s string) (a *model.Article, err error) {
	as.m.wrapRequest(opRead, func() bool {
		a, err = as.s.GetBySlug(s)
		return err != nil
	})
	return a, err
}

func (as *instrumentedArticleStore) GetUserArticleBySlug(userID uint, slug string) (a *model.Article, err error) {
	as.m.wrapRequest(opRead, func() bool {
		a, err = as.s.GetUserArticleBySlug(userID, slug)
		return err != nil
	})
	return a, err
}

func (as *instrumentedArticleStore) CreateArticle(a *model.Article) (err error) {
	as.m.wrapRequest(opCreate, func() bool {
		err = as.s.CreateArticle(a)
		return err != nil
	})
	return err
}

func (as *instrumentedArticleStore) UpdateArticle(a *model.Article, tagList []string) (err error) {
	as.m.wrapRequest(opUpdate, func() bool {
		err = as.s.UpdateArticle(a, tagList)
		return err != nil
	})
	return err
}

func (as *instrumentedArticleStore) DeleteArticle(a *model.Article) (err error) {
	as.m.wrapRequest(opDelete, func() bool {
		err = as.s.DeleteArticle(a)
		return err != nil
	})
	return err
}

func (as *instrumentedArticleStore) List(offset, limit int) (a []model.Article, n int, err error) {
	as.m.wrapRequest(opRead, func() bool {
		a, n, err = as.s.List(offset, limit)
		return err != nil
	})
	return a, n, err
}

func (as *instrumentedArticleStore) ListByTag(tag string, offset, limit int) (a []model.Article, n int, err error) {
	as.m.wrapRequest(opRead, func() bool {
		a, n, err = as.s.ListByTag(tag, offset, limit)
		return err != nil
	})
	return a, n, err
}

func (as *instrumentedArticleStore) ListByAuthor(username string, offset, limit int) (a []model.Article, n int, err error) {
	as.m.wrapRequest(opRead, func() bool {
		a, n, err = as.s.ListByAuthor(username, offset, limit)
		return err != nil
	})
	return a, n, err
}

func (as *instrumentedArticleStore) ListByWhoFavorited(username string, offset, limit int) (a []model.Article, n int, err error) {
	as.m.wrapRequest(opRead, func() bool {
		a, n, err = as.s.ListByWhoFavorited(username, offset, limit)
		return err != nil
	})
	return a, n, err
}

func (as *instrumentedArticleStore) ListFeed(userID uint, offset, limit int) (a []model.Article, n int, err error) {
	as.m.wrapRequest(opRead, func() bool {
		a, n, err = as.s.ListFeed(userID, offset, limit)
		return err != nil
	})
	return a, n, err
}

func (as *instrumentedArticleStore) AddComment(a *model.Article, c *model.Comment) (err error) {
	as.m.wrapRequest(opCreate, func() bool {
		err = as.s.AddComment(a, c)
		return err != nil
	})
	return err
}

func (as *instrumentedArticleStore) GetCommentsBySlug(slug string) (c []model.Comment, err error) {
	as.m.wrapRequest(opRead, func() bool {
		c, err = as.s.GetCommentsBySlug(slug)
		return err != nil
	})
	return c, err
}

func (as *instrumentedArticleStore) GetCommentByID(id uint) (c *model.Comment, err error) {
	as.m.wrapRequest(opRead, func() bool {
		c, err = as.s.GetCommentByID(id)
		return err != nil
	})
	return c, err
}

func (as *instrumentedArticleStore) DeleteComment(c *model.Comment) (err error) {
	as.m.wrapRequest(opDelete, func() bool {
		err = as.s.DeleteComment(c)
		return err != nil
	})
	return err
}

func (as *instrumentedArticleStore) AddFavorite(a *model.Article, userID uint) (err error) {
	as.m.wrapRequest(opDelete, func() bool {
		err = as.s.AddFavorite(a, userID)
		return err != nil
	})
	return err
}

func (as *instrumentedArticleStore) RemoveFavorite(a *model.Article, userID uint) (err error) {
	as.m.wrapRequest(opDelete, func() bool {
		err = as.s.RemoveFavorite(a, userID)
		return err != nil
	})
	return err
}

func (as *instrumentedArticleStore) ListTags() (t []model.Tag, err error) {
	as.m.wrapRequest(opRead, func() bool {
		t, err = as.s.ListTags()
		return err != nil
	})
	return t, err
}
