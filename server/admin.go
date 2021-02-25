package server

import (
	"context"
	"encoding/json"
	"github.com/kabukky/journey/dao"
	"github.com/kabukky/journey/dao/scheme"
	"github.com/kabukky/journey/logger"
	"github.com/kabukky/journey/repositories/post"
	"github.com/kabukky/journey/repositories/setting"
	"github.com/kabukky/journey/repositories/user"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kabukky/journey/authentication"
	"github.com/kabukky/journey/conversion"
	"github.com/kabukky/journey/date"
	"github.com/kabukky/journey/filenames"

	"github.com/kabukky/journey/structure"
	"github.com/kabukky/journey/templates"
	"github.com/satori/go.uuid"
)

type JsonPost struct {
	Id              int64
	Title           string
	Slug            string
	Markdown        string
	Html            string
	IsFeatured      bool
	IsPage          bool
	IsPublished     bool
	Image           string
	MetaDescription string
	Date            *time.Time
	Tags            string
}

type JsonBlog struct {
	Url             string
	Title           string
	Description     string
	Logo            string
	Cover           string
	Themes          []string
	ActiveTheme     string
	PostsPerPage    int64
	NavigationItems []structure.Navigation
}

type JsonUser struct {
	Id               int64
	Name             string
	Slug             string
	Email            string
	Image            string
	Cover            string
	Bio              string
	Website          string
	Location         string
	Password         string
	PasswordRepeated string
}

type JsonUserId struct {
	Id uint
}

type JsonImage struct {
	Filename string
}

// Function to serve the login page
func getLoginHandler(c echo.Context) (err error) {
	ctx := context.Background()
	db := dao.DB.WithContext(ctx)
	count, err := user.UsersCount(db)
	if err != nil {
		return
	}

	if count == 0 {
		err = c.Redirect(http.StatusFound, "/admin/register/")
		return
	}
	err = c.File(filepath.Join(filenames.AdminFilepath, "login.html"))
	return
}

// Function to receive a login form
func postLoginHandler(c echo.Context) (err error) {
	name := c.FormValue("name")
	password := c.FormValue("password")
	if name != "" && password != "" {
		if authentication.LoginIsCorrect(name, password) {
			logInUser(c, name)
		} else {
			logger.Info("Failed login attempt for user " + name)
		}
	}
	err = c.Redirect(http.StatusFound, "/admin/")
	return
}

// Function to serve the registration form
func getRegistrationHandler(c echo.Context) (err error) {
	ctx := context.Background()
	db := dao.DB.WithContext(ctx)
	count, err := user.UsersCount(db)
	if err != nil {
		return
	}

	if count == 0 {
		err = c.File(filepath.Join(filenames.AdminFilepath, "registration.html"))
		return
	}
	err = c.Redirect(http.StatusFound, "/admin/")
	return
}

// Function to recieve a registration form.
func postRegistrationHandler(c echo.Context) (err error) {
	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	count, err := user.UsersCount(db)
	if err != nil {
		return
	}

	if count == 0 { // TODO: Or check if authenticated user is admin when adding users from inside the admin area
		name := c.FormValue("name")
		email := c.FormValue("email")
		password := c.FormValue("password")

		otherData := make(map[string]interface{})
		otherData["Email"] = email
		_, err = user.Create(db, name, password, otherData)
		if err != nil {
			logger.Error("Create User Error %v", err)
		}

		err = c.Redirect(http.StatusFound, "/admin/")
		return
	} else {
		// TODO: Handle creation of other users (not just the first one)
		http.Error(c.Response(), "Not implemented yet.", http.StatusNotImplemented)
		return
	}
}

// Function to log out the user. Not used at the moment.
func logoutHandler(c echo.Context) (err error) {
	authentication.ClearSession(c)
	err = c.Redirect(http.StatusFound, "/admin/")
	return
}

// Function to route the /admin/ url accordingly. (Is user logged in? Is at least one user registered?)
func adminHandler(c echo.Context) (err error) {
	ctx := context.Background()
	db := dao.DB.WithContext(ctx)
	count, err := user.UsersCount(db)
	if err != nil {
		return
	}

	if count == 0 {
		err = c.Redirect(http.StatusFound, "/admin/register/")
		return
	}
	userName := authentication.GetUserName(c)
	if userName != "" {
		err = c.File(filepath.Join(filenames.AdminFilepath, "admin.html"))
		return
	}
	err = c.Redirect(http.StatusFound, "/admin/login/")
	return
}

// Function to serve files belonging to the admin interface.
func adminFileHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName != "" {
		// Get arguments (files)
		err = c.File(filepath.Join(filenames.AdminFilepath, c.Param("filepath")))
		return
	} else {
		http.NotFound(c.Response(), c.Request())
		return
	}
}

// API function to get all posts by pages
func apiPostsHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		err = c.String(http.StatusUnauthorized, "Not logged in!")
		return
	}

	var page int64
	number := c.Param("number")
	page, err = strconv.ParseInt(number, 10, 64)
	if err != nil || page < 1 {
		return
	}

	postsPerPage := int64(15)

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	posts, _, err := post.GetPostBySearch(db, nil, ((page - 1) * postsPerPage), postsPerPage, "")

	if err != nil {
		return
	}

	err = c.JSON(http.StatusOK, posts)
	return

}

// API function to get a post by id
func getApiPostHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		err = c.String(http.StatusUnauthorized, "Not logged in!")
		return
	}
	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	id := c.Param("id")
	// Get post
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil || postId < 1 {
		return
	}
	//post, err := database.RetrievePostById(postId)
	postObj, err := post.GetPostByID(db, uint(postId))
	if err != nil {
		return
	}

	err = c.JSON(http.StatusOK, postObj)
	return

}

// API function to create a post
func postApiPostHandler(c echo.Context) (err error) {
	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	userObj, err := user.GetUserByName(db, userName)
	if err != nil {
		return
	}

	data := &scheme.Post{}
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&data)
	if err != nil {
		return
	}

	savePostData := make(map[string]interface{})
	savePostData["Title"] = data.Title
	savePostData["Slug"] = data.Slug
	savePostData["Markdown"] = data.Markdown
	savePostData["HTML"] = conversion.GenerateHtmlFromMarkdown([]byte(data.Markdown))
	savePostData["Featured"] = data.Featured
	savePostData["Page"] = data.Page
	savePostData["PublishedAt"] = data.PublishedAt
	if data.PublishedAt != nil {
		savePostData["PublishedBy"] = userObj.ID
	}

	savePostData["MetaDescription"] = data.MetaDescription
	savePostData["MetaTitle"] = data.MetaTitle
	savePostData["Image"] = data.Image
	savePostData["tags_str"] = data.TagsStr
	savePostData["AuthorID"] = userObj.ID

	_, err = post.Create(db, savePostData)
	if err != nil {
		return
	}

	err = c.String(http.StatusOK, "Post created!")
	return

}

// API function to update a post.
func patchApiPostHandler(c echo.Context) (err error) {
	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	userObj, err := user.GetUserByName(db, userName)
	if err != nil {
		return
	}

	data := &scheme.Post{}
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&data)
	if err != nil {
		return
	}

	postObj, err := post.GetPostByID(db, data.ID)
	if err != nil {
		return
	}

	savePostData := make(map[string]interface{})
	savePostData["Title"] = data.Title
	savePostData["Slug"] = data.Slug
	savePostData["Markdown"] = data.Markdown
	savePostData["HTML"] = conversion.GenerateHtmlFromMarkdown([]byte(data.Markdown))
	savePostData["Featured"] = data.Featured
	savePostData["Page"] = data.Page
	savePostData["PublishedAt"] = data.PublishedAt
	if data.PublishedAt != nil {
		savePostData["PublishedBy"] = userObj.ID
	}

	savePostData["MetaDescription"] = data.MetaDescription
	savePostData["MetaTitle"] = data.MetaTitle
	savePostData["Image"] = data.Image
	savePostData["tags_str"] = data.TagsStr
	savePostData["UpdatedBy"] = userObj.ID
	savePostData["UpdatedAt"] = time.Now()

	_, err = post.Update(db, postObj, savePostData)
	if err != nil {
		return
	}

	err = c.String(http.StatusOK, "Post updated!")
	return

}

// API function to delete a post by id.
func deleteApiPostHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	id := c.Param("id")
	// Delete post
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil || postId < 1 {

		return
	}
	err = post.Delete(db, uint(postId))
	if err != nil {

		return
	}
	err = c.String(http.StatusOK, "Post deleted!")

	return

}

// API function to upload images
func apiUploadHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	// Create multipart reader
	reader, err := c.Request().MultipartReader()
	if err != nil {
		return
	}
	// Slice to hold all paths to the files
	allFilePaths := make([]string, 0)
	// Copy each part to destination.
	for {
		var part *multipart.Part
		part, err = reader.NextPart()
		if err == io.EOF {
			break
		}
		// If part.FileName() is empty, skip this iteration.
		if part.FileName() == "" {
			continue
		}
		// Folder structure: year/month/randomname
		currentDate := date.GetCurrentTime()
		filePath := filepath.Join(filenames.ImagesFilepath, currentDate.Format("2006"), currentDate.Format("01"))
		if os.MkdirAll(filePath, 0777) != nil {
			return
		}

		var dst *os.File
		dst, err = os.Create(filepath.Join(filePath, strconv.FormatInt(currentDate.Unix(), 10)+"_"+uuid.NewV4().String()+filepath.Ext(part.FileName())))
		defer dst.Close()
		if err != nil {
			return
		}
		if _, err = io.Copy(dst, part); err != nil {

			return
		}
		// Rewrite to file path on server
		filePath = strings.Replace(dst.Name(), filenames.ImagesFilepath, "/images", 1)
		// Make sure to always use "/" as path separator (to make a valid url that we can use on the blog)
		filePath = filepath.ToSlash(filePath)
		allFilePaths = append(allFilePaths, filePath)
	}

	err = c.JSON(http.StatusOK, allFilePaths)

	return

}

// API function to get all images by pages
func apiImagesHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	number := c.Param("number")
	page, err := strconv.Atoi(number)
	if err != nil || page < 1 {
		err = c.String(http.StatusInternalServerError, "Not a valid api function!")
		return
	}
	images := make([]string, 0)
	// Walk all files in images folder
	err = filepath.Walk(filenames.ImagesFilepath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && (strings.EqualFold(filepath.Ext(filePath), ".jpg") || strings.EqualFold(filepath.Ext(filePath), ".jpeg") || strings.EqualFold(filepath.Ext(filePath), ".gif") || strings.EqualFold(filepath.Ext(filePath), ".png") || strings.EqualFold(filepath.Ext(filePath), ".svg")) {
			// Rewrite to file path on server
			filePath = strings.Replace(filePath, filenames.ImagesFilepath, "/images", 1)
			// Make sure to always use "/" as path separator (to make a valid url that we can use on the blog)
			filePath = filepath.ToSlash(filePath)
			// Prepend file to slice (thus reversing the order)
			images = append([]string{filePath}, images...)
		}
		return nil
	})
	if len(images) == 0 {
		// Write empty json array
		err = c.JSON(http.StatusOK, images)
		return
	}
	imagesPerPage := 15
	start := (page * imagesPerPage) - imagesPerPage
	end := page * imagesPerPage
	if start > (len(images) - 1) {
		// Write empty json array
		err = c.JSON(http.StatusOK, []string{})
		return
	}
	if end > len(images) {
		end = len(images)
	}

	err = c.JSON(http.StatusOK, images[start:end])
	return

}

// API function to delete an image by its filename.
func deleteApiImageHandler(c echo.Context) (err error) {
	// TODO: Check if the user has permissions to delete the image
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	// Get the file name from the json data
	decoder := json.NewDecoder(c.Request().Body)
	var image JsonImage
	err = decoder.Decode(&image)
	if err != nil {
		return
	}
	err = filepath.Walk(filenames.ImagesFilepath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Base(filePath) == filepath.Base(image.Filename) {
			err = os.Remove(filePath)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {

		return
	}

	err = c.String(http.StatusOK, "Image deleted!")
	return

}

// API function to get blog settings
func getApiBlogHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	blog, err := setting.RetrieveBlog(db)
	if err != nil {
		return
	}
	// Read lock the global blog
	blogJson := blogToJson(blog)

	err = c.JSON(http.StatusOK, blogJson)
	return

}

// API function to update blog settings
func patchApiBlogHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	userObj, err := user.GetUserByName(db, userName)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(c.Request().Body)
	var blogData JsonBlog
	err = decoder.Decode(&blogData)
	if err != nil {
		return
	}

	// Make sure postPerPage is over 0
	if blogData.PostsPerPage < 1 {
		blogData.PostsPerPage = 1
	}
	// Remove blog url in front of navigation urls
	for index, _ := range blogData.NavigationItems {
		if strings.HasPrefix(blogData.NavigationItems[index].Url, blogData.Url) {
			blogData.NavigationItems[index].Url = strings.Replace(blogData.NavigationItems[index].Url, blogData.Url, "", 1)
			// If we removed the blog url, there should be a / in front of the url
			if !strings.HasPrefix(blogData.NavigationItems[index].Url, "/") {
				blogData.NavigationItems[index].Url = "/" + blogData.NavigationItems[index].Url
			}
		}
	}

	// Retrieve old blog settings for comparison

	saveData := make(map[string]interface{})
	saveData["title"] = conversion.XssFilter(blogData.Title)
	saveData["description"] = conversion.XssFilter(blogData.Description)
	saveData["logo"] = conversion.XssFilter(blogData.Logo)
	saveData["cover"] = conversion.XssFilter(blogData.Cover)
	saveData["postsPerPage"] = blogData.PostsPerPage
	saveData["activeTheme"] = blogData.ActiveTheme
	saveData["navigation"] = blogData.NavigationItems

	for k, v := range saveData {
		err = setting.Set(db, "blog", k, v, userObj)
		if err != nil {
			return
		}
	}

	err = c.String(http.StatusOK, "Blog settings updated!")
	return

}

// API function to get user settings
func getApiUserHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	userObj, err := user.GetUserByName(db, userName)
	if err != nil {
		return
	}

	id := c.Param("id")
	userIdToGet, err := strconv.ParseInt(id, 10, 64)
	if err != nil || userIdToGet < 1 {
		return
	} else if uint(userIdToGet) != userObj.ID { // Make sure the authenticated user is only accessing his/her own data. TODO: Make sure the user is admin when multiple users have been introduced
		err = c.String(http.StatusForbidden, "You don't have permission to access this data.")
		return
	}

	err = c.JSON(http.StatusOK, userObj)
	return

}

// API function to patch user settings
func patchApiUserHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	userObj, err := user.GetUserByName(db, userName)
	if err != nil {
		return
	}

	data := &scheme.User{}
	decoder := json.NewDecoder(c.Request().Body)
	if err = decoder.Decode(&data); err != nil {
		return
	}

	if data.ID < 1 {
		err = c.String(http.StatusBadRequest, "Wrong user id.")
		return
	}

	if data.ID != userObj.ID {
		// Make sure the authenticated user is only changing his/her own data. TODO: Make sure the user is admin when multiple users have been introduced
		err = c.String(http.StatusUnauthorized, "You don't have permission to change this data.")
		return
	}

	updateData := make(map[string]interface{})
	var checkUser *scheme.User
	// Check if new name is already taken
	if data.Name != "" {
		if data.Name != userObj.Name {
			checkUser, err = user.GetUserByName(db, data.Name)
			if err == nil && !errors.Is(gorm.ErrRecordNotFound, err) {
				return
			}
			if checkUser != nil && checkUser.ID != data.ID {
				err = c.String(http.StatusBadRequest, "User name already token.")
				return
			}
		}
		updateData["Name"] = data.Name
	}

	// Check if new slug is already taken
	if data.Slug != "" {
		if data.Slug != userObj.Slug {
			checkUser, err = user.GetUserBySlug(db, data.Slug)
			if err == nil && !errors.Is(gorm.ErrRecordNotFound, err) {
				return
			}
			if checkUser != nil && checkUser.ID != data.ID {
				err = c.String(http.StatusBadRequest, "User slug already token.")
				return
			}
		}
		updateData["Slug"] = data.Slug
	}

	if data.Password != "" {
		updateData["Password"] = data.Password

	}

	updateData["Email"] = data.Email
	updateData["Image"] = data.Image
	updateData["Cover"] = conversion.StripTagsFromHTML(data.BIO)
	updateData["Website"] = conversion.StripTagsFromHTML(data.Website)
	updateData["Location"] = conversion.StripTagsFromHTML(data.Location)
	updateData["UpdatedBy"] = userObj.ID

	newUser, err := user.Update(db, data, updateData)
	if err != nil {
		return
	}

	// Check if the user name was changed. If so, update the session cookie to the new user name.
	if newUser.Name != userObj.Name && newUser.ID == userObj.ID {
		logInUser(c, newUser.Name)
	}

	err = c.String(http.StatusOK, "User settings updated!")
	return

}

// API function to get the id of the currently authenticated user
func getApiUserIdHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	userObj, err := user.GetUserByName(db, userName)
	if err != nil {
		return
	}

	jsonUserId := JsonUserId{Id: userObj.ID}
	err = c.JSON(http.StatusOK, jsonUserId)
	return

}

func logInUser(c echo.Context, name string) {
	authentication.SetSession(c, name)

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	userObj, err := user.GetUserByName(db, name)
	if err != nil {
		logger.Error("Couldn't get id of logged in user:", err)
		return
	}

	err = userObj.UpdateLastLogin(db, c.RealIP(), time.Now())
	if err != nil {
		logger.Error("Couldn't update last login date of a user:", err)
	}
}

func blogToJson(blog *structure.Blog) *JsonBlog {
	var jsonBlog JsonBlog
	jsonBlog.Url = blog.Url
	jsonBlog.Title = blog.Title
	jsonBlog.Description = blog.Description
	jsonBlog.Logo = blog.Logo
	jsonBlog.Cover = blog.Cover
	jsonBlog.PostsPerPage = blog.PostsPerPage
	jsonBlog.Themes = templates.GetAllThemes()
	jsonBlog.ActiveTheme = blog.ActiveTheme
	jsonBlog.NavigationItems = blog.NavigationItems
	return &jsonBlog
}

func InitializeAdmin(router *echo.Echo) {
	// For admin panel
	router.GET("/admin/", adminHandler)
	router.GET("/admin/login/", getLoginHandler)
	router.POST("/admin/login/", postLoginHandler)
	router.GET("/admin/register/", getRegistrationHandler)
	router.POST("/admin/register/", postRegistrationHandler)
	router.GET("/admin/logout/", logoutHandler)
	router.GET("/admin/*filepath", adminFileHandler)

	// For admin API (no trailing slash)
	// Posts
	router.GET("/admin/api/posts/:number", apiPostsHandler)
	// Post
	router.GET("/admin/api/post/:id", getApiPostHandler)
	router.POST("/admin/api/post", postApiPostHandler)
	router.PATCH("/admin/api/post", patchApiPostHandler)
	router.DELETE("/admin/api/post/:id", deleteApiPostHandler)
	// Upload
	router.POST("/admin/api/upload", apiUploadHandler)
	// Images
	router.GET("/admin/api/images/:number", apiImagesHandler)
	router.DELETE("/admin/api/image", deleteApiImageHandler)
	// Blog
	router.GET("/admin/api/blog", getApiBlogHandler)
	router.PATCH("/admin/api/blog", patchApiBlogHandler)
	// User
	router.GET("/admin/api/user/:id", getApiUserHandler)
	router.PATCH("/admin/api/user", patchApiUserHandler)
	// User id
	router.GET("/admin/api/userid", getApiUserIdHandler)

	router.Static("/admin/", filenames.AdminFilepath)
}
