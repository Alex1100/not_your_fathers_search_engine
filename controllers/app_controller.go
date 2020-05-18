package controllers

import (
	links "not_your_fathers_search_engine/controllers/links"
)

// WORK IN PROGRESS

type AppController struct {
	LinkResource *links.LinkResource
}

func NewApiController() *AppController {
	return &AppController{
		LinkResource: links.NewLinkResource(),
	}
}