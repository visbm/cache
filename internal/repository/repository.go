package repository

import (	
	"cache/internal/repository/models"
	"cache/internal/repository/pgsqldb"
	"cache/internal/repository/redisdb"
	"cache/pkg/logger"
	)

type ArticleRepository interface {
	GetArticle(id string) (models.Article, error)
	GetArticles() ([]models.Article, error)
	CreateArticle(article models.Article) (models.Article, error)
}

type MainRepository struct {
	cache  *redisdb.RedisDatabase
	db     *pgsqldb.PgSqlRepository
	logger logger.Logger
}

func NewMainRepository(cache *redisdb.RedisDatabase, db *pgsqldb.PgSqlRepository, logger logger.Logger) ArticleRepository {
	return &MainRepository{
		cache:  cache,
		db:     db,
		logger: logger,
	}
}

func (r *MainRepository) GetArticles() ([]models.Article, error) {

	cacheArticles, err := r.cache.GetArticles()
	if err == nil {
		return cacheArticles, nil
	}
	r.logger.Errorf("No value in cache %s", err)

	dbArticles, err := r.db.GetArticles()
	if err != nil {
		r.logger.Errorf("error on get articles %s", err)
		return nil, err
	}
	r.logger.Infof("len articles from db %s", len(dbArticles))

	err = r.cache.SaveAll(dbArticles)
	if err != nil {
		r.logger.Errorf("error on save articles in cache %s", err)
		return nil, err
	}
	r.logger.Info("articles saved in cache")

	return dbArticles, nil
}

func (r *MainRepository) GetArticle(id string) (models.Article, error) {
	cacheArticle, err := r.cache.GetArticle(id)
	if err == nil {
		return cacheArticle, nil
	}
	r.logger.Errorf("No value in cache %s", err)

	dbArticle, err := r.db.GetArticle(id)
	if err != nil {
		r.logger.Errorf("%s", err)
		return models.Article{}, err
	}
	r.logger.Infof("article's id from db %s", dbArticle.ID)

	_, err = r.cache.CreateArticle(dbArticle)
	if err != nil {
		r.logger.Errorf("%s", err)
		return models.Article{}, err
	}
	r.logger.Info("article saved in cache db")
	return dbArticle, nil
}

func (r *MainRepository) CreateArticle(article models.Article) (models.Article, error) {
	article, err := r.db.CreateArticle(article)
	if err != nil {
		r.logger.Errorf("%s", err)
		return article, err
	}	
	r.logger.Info("article saved in db")

	_, err = r.cache.CreateArticle(article)
	if err != nil {
		r.logger.Errorf("%s", err)
		return article, err
	}
	r.logger.Info("article saved in cache db")
	return article, nil
}
