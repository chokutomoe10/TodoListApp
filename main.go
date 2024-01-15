package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todo_list/models"

	"github.com/labstack/echo/v4"
)

func main() {
	models.DatabaseConnect()

	e := echo.New()

	e.GET("/", GetAllLists)
	e.GET("/list/:id", GetList)
	e.POST("/list", CreateList)
	e.PUT("/list/:id", UpdateList)
	e.DELETE("/list/:id", DeleteList)

	e.GET("/subLists/:id", GetAllSubLists)
	e.GET("/subList/:id", GetOneSubList)
	e.POST("/subList/:id", CreateOneSubList)
	e.PUT("/subList/:id", UpdateSubList)
	e.DELETE("/subList/:id", DeleteSubList)

	e.Logger.Fatal(e.Start(":1323"))
}

func GetAllLists(c echo.Context) error {
	var lists []models.List

	models.DB.Find(&lists)
	// models.DB.Preload("SubLists").Find(&lists)

	return c.JSON(http.StatusOK, lists)
}

func GetList(c echo.Context) error {
	var list models.List

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	models.DB.Find(&list, "id = ?", id)
	// models.DB.Find(&list, "id = ?", c.Param("id"))
	// models.DB.Preload("SubLists").Find(&list, "id = ?", c.Param("id"))

	return c.JSON(http.StatusOK, list)
}

func CreateList(c echo.Context) error {
	var list models.List
	var fileType, fileName string
	var fileSlice []string

	form, err := c.MultipartForm()

	if err != nil {
		return err
	}

	files := form.File["files"]

	for _, file := range files {
		src, errSrc := file.Open()

		if errSrc != nil {
			return errSrc
		}

		defer src.Close()

		fileByte, _ := ioutil.ReadAll(src)
		fileType = http.DetectContentType(fileByte)

		randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		randomInt := randomizer.Int()
		randomString := fmt.Sprintf("%d", randomInt)

		if fileType == "application/pdf" {
			fileName = "uploads/List/" + randomString + ".pdf"

			errWrite := ioutil.WriteFile(fileName, fileByte, 0777)

			if errWrite != nil {
				return errWrite
			}

			fileSlice = append(fileSlice, fileName[13:])
		} else if fileType == "text/plain; charset=utf-8" {
			fileName = "uploads/List/" + randomString + ".txt"

			errWrite := ioutil.WriteFile(fileName, fileByte, 0777)

			if errWrite != nil {
				return errWrite
			}

			fileSlice = append(fileSlice, fileName[13:])
		} else {
			return c.JSON(http.StatusBadRequest, "cannot upload file")
		}
	}

	if errBody := c.Bind(&list); err != nil {
		return errBody
	}

	list.Files = strings.Join(fileSlice, ", ")

	models.DB.Create(&list)

	return c.JSON(http.StatusOK, list)
}

func UpdateList(c echo.Context) error {
	var list models.List
	var fileType, fileName string
	var fileSlice []string

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}

	form, errForm := c.MultipartForm()

	if errForm != nil {
		return errForm
	}

	files := form.File["files"]

	for _, file := range files {
		src, errSrc := file.Open()

		if errSrc != nil {
			return errSrc
		}

		defer src.Close()

		fileByte, _ := ioutil.ReadAll(src)
		fileType = http.DetectContentType(fileByte)

		randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		randomInt := randomizer.Int()
		randomString := fmt.Sprintf("%d", randomInt)

		if fileType == "application/pdf" {
			fileName = "uploads/List/" + randomString + ".pdf"

			errWrite := ioutil.WriteFile(fileName, fileByte, 0777)

			if errWrite != nil {
				return errWrite
			}

			fileSlice = append(fileSlice, fileName[13:])
		} else if fileType == "text/plain; charset=utf-8" {
			fileName = "uploads/List/" + randomString + ".txt"

			errWrite := ioutil.WriteFile(fileName, fileByte, 0777)

			if errWrite != nil {
				return errWrite
			}

			fileSlice = append(fileSlice, fileName[13:])
		} else {
			return c.JSON(http.StatusBadRequest, "cannot update file")
		}
	}

	if errBody := c.Bind(&list); err != nil {
		return errBody
	}

	list.Files = strings.Join(fileSlice, ", ")

	models.DB.Where("id = ?", c.Param("id")).Updates(&list)

	list.ID = uint(id)

	return c.JSON(http.StatusOK, list)
}

func DeleteList(c echo.Context) error {
	var list models.List

	models.DB.Delete(&list, "id = ?", c.Param("id"))

	return c.JSON(http.StatusOK, list)
}

func GetAllSubLists(c echo.Context) error {
	var subLists []models.SubList

	models.DB.Where("list_id = ?", c.Param("id")).Find(&subLists)

	return c.JSON(http.StatusOK, subLists)
}

func GetOneSubList(c echo.Context) error {
	var subList models.SubList

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	models.DB.Find(&subList, "id = ?", id)
	// models.DB.Find(&subList, "id = ?", c.Param("id"))

	return c.JSON(http.StatusOK, subList)
}

func CreateOneSubList(c echo.Context) error {
	var subList models.SubList
	var fileType, fileName string
	var fileSlice []string

	listId, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}

	subList.ListID = uint(listId)

	form, errForm := c.MultipartForm()

	if errForm != nil {
		return errForm
	}

	files := form.File["subList_files"]

	for _, file := range files {
		src, errSrc := file.Open()

		if errSrc != nil {
			return errSrc
		}

		defer src.Close()

		fileByte, _ := ioutil.ReadAll(src)
		fileType = http.DetectContentType(fileByte)

		randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		randomInt := randomizer.Int()
		randomString := fmt.Sprintf("%d", randomInt)

		if fileType == "application/pdf" {
			fileName = "uploads/SubList/" + randomString + ".pdf"

			errWrite := ioutil.WriteFile(fileName, fileByte, 0777)

			if errWrite != nil {
				return errWrite
			}

			fileSlice = append(fileSlice, fileName[16:])
		} else if fileType == "text/plain; charset=utf-8" {
			fileName = "uploads/SubList/" + randomString + ".txt"

			errWrite := ioutil.WriteFile(fileName, fileByte, 0777)

			if errWrite != nil {
				return errWrite
			}

			fileSlice = append(fileSlice, fileName[16:])
		} else {
			return c.JSON(http.StatusBadRequest, "cannot upload file")
		}
	}

	if errBody := c.Bind(&subList); err != nil {
		return errBody
	}

	subList.SubListFiles = strings.Join(fileSlice, ", ")

	models.DB.Create(&subList)

	return c.JSON(http.StatusOK, subList)
}

func UpdateSubList(c echo.Context) error {
	var subList models.SubList
	var getSubList models.SubList
	var fileType, fileName string
	var fileSlice []string

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}

	models.DB.Where("id = ?", id).Find(&getSubList)

	form, errForm := c.MultipartForm()

	if errForm != nil {
		return errForm
	}

	files := form.File["subList_files"]

	for _, file := range files {
		src, errSrc := file.Open()

		if errSrc != nil {
			return errSrc
		}

		defer src.Close()

		fileByte, _ := ioutil.ReadAll(src)
		fileType = http.DetectContentType(fileByte)

		randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		randomInt := randomizer.Int()
		randomString := fmt.Sprintf("%d", randomInt)

		if fileType == "application/pdf" {
			fileName = "uploads/SubList/" + randomString + ".pdf"

			errWrite := ioutil.WriteFile(fileName, fileByte, 0777)

			if errWrite != nil {
				return errWrite
			}

			fileSlice = append(fileSlice, fileName[16:])
		} else if fileType == "text/plain; charset=utf-8" {
			fileName = "uploads/SubList/" + randomString + ".txt"

			errWrite := ioutil.WriteFile(fileName, fileByte, 0777)

			if errWrite != nil {
				return errWrite
			}

			fileSlice = append(fileSlice, fileName[16:])
		} else {
			return c.JSON(http.StatusBadRequest, "cannot update file")
		}
	}

	if errBody := c.Bind(&subList); err != nil {
		return errBody
	}

	subList.SubListFiles = strings.Join(fileSlice, ", ")

	models.DB.Where("id = ?", c.Param("id")).Updates(&subList)

	subList.ID = uint(id)
	subList.ListID = getSubList.ListID

	return c.JSON(http.StatusOK, subList)
}

func DeleteSubList(c echo.Context) error {
	var subList models.SubList

	models.DB.Delete(&subList, "id = ?", c.Param("id"))

	return c.JSON(http.StatusOK, subList)
}
