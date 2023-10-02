package service

import (
	"cache/internal/handlers"
	"cache/internal/repository"
	"cache/internal/repository/models"
	"cache/pkg/logger"
	"github.com/satori/go.uuid"
)

type ArticleService struct {
	articleRepository repository.ArticleRepository
	logger            logger.Logger
}

func NewArticleService(articleRepository repository.ArticleRepository, logger logger.Logger) *ArticleService {
	return &ArticleService{
		articleRepository: articleRepository,
		logger:            logger,
	}

}

func (s *ArticleService) GetArticles() ([]handlers.ArticleResponse, error) {
	articles, err := s.articleRepository.GetArticles()
	if err != nil {
		s.logger.Errorf("error while getting articles: %v", err)
		return nil, err
	}

	resp := models.MapArticlesToResponse(articles)

	return resp, nil
}

func (s *ArticleService) GetArticle(id string) (*handlers.ArticleResponse, error) {
	article, err := s.articleRepository.GetArticle(id)
	if err != nil {
		s.logger.Errorf("error while getting article: %v", err)
		return nil, err
	}
	resp := models.MapArticleToResponse(article)

	return &resp, nil
}

func (s *ArticleService) CreateArticle(req *handlers.ArticleRequest) (*handlers.ArticleResponse, error) {
	article := models.MapReqToArticle(*req)
	article.ID = uuid.NewV4().String()

	article, err := s.articleRepository.CreateArticle(article)
	if err != nil {
		s.logger.Errorf("error while creating article: %v", err)
		return nil, err
	}

	resp := models.MapArticleToResponse(article)

	return &resp, nil
}
