package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"test/go_server/db_client"
	"time"

	"github.com/labstack/echo/v4"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Test struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ErrMsg struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
}

func CreatePost(c echo.Context) error {
	// reqBody := c.Request().Body
	var reqBody Test
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&reqBody)
	if err != nil {
		log.Printf("Failed processing CreatePost request: %s", err)
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   "parse request body error",
		})
	}

	res, err := db_client.DBClient.Exec("INSERT INTO test (name) values (?);",
		reqBody.Name,
	)
	if err != nil {
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   fmt.Sprint(err),
		})
	}
	id, _ := res.LastInsertId()
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"error": false,
		"id":    id,
	})
}

func GetPosts(c echo.Context) error {
	var tests []Test

	rows, err := db_client.DBClient.Query("SELECT id, name FROM test;")
	if err != nil {
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   fmt.Sprint(err),
		})
	}

	for rows.Next() {
		var singleTest Test
		if err := rows.Scan(&singleTest.ID, &singleTest.Name); err != nil {
			return c.JSON(500, ErrMsg{
				Error: true,
				Msg:   fmt.Sprint(err),
			})
		}
		tests = append(tests, singleTest)
	}

	return c.JSON(200, tests)
}

func GetPost(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	row := db_client.DBClient.QueryRow("SELECT id, name FROM test WHERE id = ?;", id)
	var test Test
	if err := row.Scan(&test.ID, &test.Name); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, ErrMsg{
				Error: true,
				Msg:   err.Error(),
			})
		}
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   err.Error(),
		})
	}

	return c.JSON(200, test)
}
