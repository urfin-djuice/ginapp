package controller

import "github.com/gin-gonic/gin"

type HandlerList []gin.HandlerFunc

type Act struct {
	Method   string
	Route    string
	Handlers HandlerList
}

type Ctrl struct {
	Name     string
	Handlers HandlerList
	Acts     []Act
}
