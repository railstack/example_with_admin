package controllers

import (
	"net/http"
	"strconv"

	m "../models"
	"gopkg.in/gin-gonic/gin.v1"
)

func IndexHandler(c *gin.Context) {
	posts, err := m.AllPosts()
	if err != nil {
		c.String(http.StatusNotFound, "Posts not found or some error occurred!")
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", posts)
}

func ShowHandler(c *gin.Context) {
	id, err := ToInt(c.Param("id"))
	post, err := m.FindPost(id)
	if err != nil {
		c.String(http.StatusNotFound, "Post not found or some error occurred!")
		return
	}
	c.HTML(http.StatusOK, "show.tmpl", post)
}

func ToInt(s string) (int64, error) {
	res, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		res = 0
	}
	return res, err
}
