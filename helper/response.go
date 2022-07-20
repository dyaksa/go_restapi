package helper

import (
	"golang_restapi/model/entity"
	"golang_restapi/model/web"
)

func ToCategoryResponse(category entity.Category) web.CategoryResponse {
	return web.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

func ToCategorySliceResponse(category []entity.Category) []web.CategoryResponse {
	var categories []web.CategoryResponse
	for _, cat := range category {
		category := web.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		}
		categories = append(categories, category)
	}
	return categories
}
