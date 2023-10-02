package models

import (
	"cache/internal/handlers"
)

type Article struct {
	ID    string
	URL   string
	Title string
}

func MapArticleToResponse(a Article) handlers.ArticleResponse {
	return handlers.ArticleResponse{
		ID:    a.ID,
		URL:   a.URL,
		Title: a.Title,
	}
}

func MapReqToArticle(a handlers.ArticleRequest) Article {
	return Article{
		URL:   a.URL,
		Title: a.Title,
	}
}

func MapArticlesToResponse(articles []Article) []handlers.ArticleResponse {
	var resp []handlers.ArticleResponse
	for _, a := range articles {
		r := handlers.ArticleResponse{
			ID:    a.ID,
			URL:   a.URL,
			Title: a.Title,
		}
		resp = append(resp, r)
	}
	return resp
}
