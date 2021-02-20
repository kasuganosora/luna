package server

import (
	"encoding/json"
	"github.com/kabukky/journey/logger"
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kabukky/journey/authentication"
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/conversion"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/date"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/slug"
	"github.com/kabukky/journey/structure"
	"github.com/kabukky/journey/structure/methods"
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
	Id int64
}

type JsonImage struct {
	Filename string
}

// Function to serve the login page
func getLoginHandler(c echo.Context) (err error) {
	if database.RetrieveUsersCount() == 0 {
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
	if database.RetrieveUsersCount() == 0 {
		err = c.File(filepath.Join(filenames.AdminFilepath, "registration.html"))
		return
	}
	err = c.Redirect(http.StatusFound, "/admin/")
	return
}

// Function to recieve a registration form.
func postRegistrationHandler(c echo.Context) (err error) {
	if database.RetrieveUsersCount() == 0 { // TODO: Or check if authenticated user is admin when adding users from inside the admin area
		name := c.FormValue("name")
		email := c.FormValue("email")
		password := c.FormValue("password")

		var hashedPassword string
		if name != "" && password != "" {
			hashedPassword, err = authentication.EncryptPassword(password)
			if err != nil {
				return
			}
			user := structure.User{Name: []byte(name), Slug: slug.Generate(name, "users"), Email: []byte(email), Image: []byte(filenames.DefaultUserImageFilename), Cover: []byte(filenames.DefaultUserCoverFilename), Role: 4}
			err = methods.SaveUser(&user, hashedPassword, 1)
			if err != nil {
				return
			}
			err = c.Redirect(http.StatusFound, "/admin/")
			return
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
	if database.RetrieveUsersCount() == 0 {
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
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	var page int64
	number := c.Param("number")
	page, err = strconv.ParseInt(number, 10, 64)
	if err != nil || page < 1 {
		return
	}

	postsPerPage := int64(15)
	posts, err := database.RetrievePostsForApi(postsPerPage, ((page - 1) * postsPerPage))
	if err != nil {
		return
	}

	err = c.JSON(http.StatusOK, postsToJson(posts))
	return

}

// API function to get a post by id
func getApiPostHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	id := c.Param("id")
	// Get post
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil || postId < 1 {
		return
	}
	post, err := database.RetrievePostById(postId)
	if err != nil {
		return
	}

	err = c.JSON(http.StatusOK, postToJson(post))
	return

}

// API function to create a post
func postApiPostHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	userId, err := getUserId(userName)
	if err != nil {
		return
	}
	// Create post
	decoder := json.NewDecoder(c.Request().Body)
	var postJSON JsonPost
	err = decoder.Decode(&postJSON)
	if err != nil {
		return
	}
	var postSlug string
	if postJSON.Slug != "" { // Ceck if user has submitted a custom slug
		postSlug = slug.Generate(postJSON.Slug, "posts")
	} else {
		postSlug = slug.Generate(postJSON.Title, "posts")
	}
	currentTime := date.GetCurrentTime()
	post := structure.Post{Title: []byte(postJSON.Title), Slug: postSlug, Markdown: []byte(postJSON.Markdown), Html: conversion.GenerateHtmlFromMarkdown([]byte(postJSON.Markdown)), IsFeatured: postJSON.IsFeatured, IsPage: postJSON.IsPage, IsPublished: postJSON.IsPublished, MetaDescription: []byte(postJSON.MetaDescription), Image: []byte(postJSON.Image), Date: &currentTime, Tags: methods.GenerateTagsFromCommaString(postJSON.Tags), Author: &structure.User{Id: userId}}
	err = methods.SavePost(&post)
	if err != nil {
		return
	}

	err = c.String(http.StatusOK, "Post created!")
	return

}

// API function to update a post.
func patchApiPostHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	userId, err := getUserId(userName)
	if err != nil {
		return
	}
	// Update post
	decoder := json.NewDecoder(c.Request().Body)
	var postJSON JsonPost
	err = decoder.Decode(&postJSON)
	if err != nil {
		return
	}
	var postSlug string
	// Get current slug of post
	post, err := database.RetrievePostById(postJSON.Id)
	if err != nil {
		return
	}
	if postJSON.Slug != post.Slug { // Check if user has submitted a custom slug
		postSlug = slug.Generate(postJSON.Slug, "posts")
	} else {
		postSlug = post.Slug
	}
	currentTime := date.GetCurrentTime()
	*post = structure.Post{Id: postJSON.Id, Title: []byte(postJSON.Title), Slug: postSlug, Markdown: []byte(postJSON.Markdown), Html: conversion.GenerateHtmlFromMarkdown([]byte(postJSON.Markdown)), IsFeatured: postJSON.IsFeatured, IsPage: postJSON.IsPage, IsPublished: postJSON.IsPublished, MetaDescription: []byte(postJSON.MetaDescription), Image: []byte(postJSON.Image), Date: &currentTime, Tags: methods.GenerateTagsFromCommaString(postJSON.Tags), Author: &structure.User{Id: userId}}
	err = methods.UpdatePost(post)
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

	id := c.Param("id")
	// Delete post
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil || postId < 1 {

		return
	}
	err = methods.DeletePost(postId)
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
	// Read lock the global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	blogJson := blogToJson(methods.Blog)

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

	userId, err := getUserId(userName)
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
	blog, err := database.RetrieveBlog()
	if err != nil {
		return
	}
	tempBlog := structure.Blog{Url: []byte(configuration.Config.Url), Title: []byte(blogData.Title), Description: []byte(blogData.Description), Logo: []byte(blogData.Logo), Cover: []byte(blogData.Cover), AssetPath: []byte("/assets/"), PostCount: blog.PostCount, PostsPerPage: blogData.PostsPerPage, ActiveTheme: blogData.ActiveTheme, NavigationItems: blogData.NavigationItems}
	err = methods.UpdateBlog(&tempBlog, userId)
	// Check if active theme setting has been changed, if so, generate templates from new theme
	if tempBlog.ActiveTheme != blog.ActiveTheme {
		err = templates.Generate()
		if err != nil {
			// If there's an error while generating the new templates, the whole program must be stopped.
			logger.Fatal("Fatal error: Template data couldn't be generated from theme files: " + err.Error())
			return
		}
	}
	if err != nil {
		return
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

	userId, err := getUserId(userName)
	if err != nil {
		return
	}
	id := c.Param("id")
	userIdToGet, err := strconv.ParseInt(id, 10, 64)
	if err != nil || userIdToGet < 1 {
		return
	} else if userIdToGet != userId { // Make sure the authenticated user is only accessing his/her own data. TODO: Make sure the user is admin when multiple users have been introduced
		err = c.String(http.StatusForbidden, "You don't have permission to access this data.")
		return
	}
	user, err := database.RetrieveUser(userIdToGet)
	if err != nil {

		return
	}
	userJson := userToJson(user)
	err = c.JSON(http.StatusOK, userJson)
	return

}

// API function to patch user settings
func patchApiUserHandler(c echo.Context) (err error) {
	userName := authentication.GetUserName(c)
	if userName == "" {
		http.Error(c.Response(), "Not logged in!", http.StatusUnauthorized)
		return
	}

	userId, err := getUserId(userName)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(c.Request().Body)
	var userData JsonUser
	err = decoder.Decode(&userData)
	if err != nil {
		return
	}
	// Make sure user id is over 0
	if userData.Id < 1 {
		err = c.String(http.StatusBadRequest, "Wrong user id.")
		return
	} else if userId != userData.Id { // Make sure the authenticated user is only changing his/her own data. TODO: Make sure the user is admin when multiple users have been introduced
		err = c.String(http.StatusUnauthorized, "You don't have permission to change this data.")
		return
	}
	// Get old user data to compare
	tempUser, err := database.RetrieveUser(userData.Id)
	if err != nil {
		return
	}
	// Make sure user email is provided
	if userData.Email == "" {
		userData.Email = string(tempUser.Email)
	}
	// Make sure user name is provided
	if userData.Name == "" {
		userData.Name = string(tempUser.Name)
	}
	// Make sure user slug is provided
	if userData.Slug == "" {
		userData.Slug = tempUser.Slug
	}
	// Check if new name is already taken
	if userData.Name != string(tempUser.Name) {
		_, err = database.RetrieveUserByName([]byte(userData.Name))
		if err == nil {
			// The new user name is already taken. Assign the old name.
			// TODO: Return error that will be displayed in the admin interface.
			userData.Name = string(tempUser.Name)
		}
	}
	// Check if new slug is already taken
	if userData.Slug != tempUser.Slug {
		_, err = database.RetrieveUserBySlug(userData.Slug)
		if err == nil {
			// The new user slug is already taken. Assign the old slug.
			// TODO: Return error that will be displayed in the admin interface.
			userData.Slug = tempUser.Slug
		}
	}
	user := structure.User{Id: userData.Id, Name: []byte(userData.Name), Slug: userData.Slug, Email: []byte(userData.Email), Image: []byte(userData.Image), Cover: []byte(userData.Cover), Bio: []byte(userData.Bio), Website: []byte(userData.Website), Location: []byte(userData.Location)}
	err = methods.UpdateUser(&user, userId)
	if err != nil {
		return
	}
	if userData.Password != "" && (userData.Password == userData.PasswordRepeated) { // Update password if a new one was submitted
		var encryptedPassword string
		encryptedPassword, err = authentication.EncryptPassword(userData.Password)
		if err != nil {
			return
		}
		err = database.UpdateUserPassword(user.Id, encryptedPassword, date.GetCurrentTime(), userData.Id)
		if err != nil {
			return
		}
	}
	// Check if the user name was changed. If so, update the session cookie to the new user name.
	if userData.Name != string(tempUser.Name) {
		logInUser(c, userData.Name)
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

	userId, err := getUserId(userName)
	if err != nil {
		return
	}
	jsonUserId := JsonUserId{Id: userId}
	err = c.JSON(http.StatusOK, jsonUserId)
	return

}

func getUserId(userName string) (int64, error) {
	user, err := database.RetrieveUserByName([]byte(userName))
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func logInUser(c echo.Context, name string) {
	authentication.SetSession(c, name)
	userId, err := getUserId(name)
	if err != nil {
		logger.Error("Couldn't get id of logged in user:", err)
	}
	err = database.UpdateLastLogin(date.GetCurrentTime(), userId)
	if err != nil {
		logger.Error("Couldn't update last login date of a user:", err)
	}
}

func postsToJson(posts []structure.Post) *[]JsonPost {
	jsonPosts := make([]JsonPost, len(posts))
	for index, _ := range posts {
		jsonPosts[index] = *postToJson(&posts[index])
	}
	return &jsonPosts
}

func postToJson(post *structure.Post) *JsonPost {
	var jsonPost JsonPost
	jsonPost.Id = post.Id
	jsonPost.Title = string(post.Title)
	jsonPost.Slug = post.Slug
	jsonPost.Markdown = string(post.Markdown)
	jsonPost.Html = string(post.Html)
	jsonPost.IsFeatured = post.IsFeatured
	jsonPost.IsPage = post.IsPage
	jsonPost.IsPublished = post.IsPublished
	jsonPost.MetaDescription = string(post.MetaDescription)
	jsonPost.Image = string(post.Image)
	jsonPost.Date = post.Date
	tags := make([]string, len(post.Tags))
	for index, _ := range post.Tags {
		tags[index] = string(post.Tags[index].Name)
	}
	jsonPost.Tags = strings.Join(tags, ",")
	return &jsonPost
}

func blogToJson(blog *structure.Blog) *JsonBlog {
	var jsonBlog JsonBlog
	jsonBlog.Url = string(blog.Url)
	jsonBlog.Title = string(blog.Title)
	jsonBlog.Description = string(blog.Description)
	jsonBlog.Logo = string(blog.Logo)
	jsonBlog.Cover = string(blog.Cover)
	jsonBlog.PostsPerPage = blog.PostsPerPage
	jsonBlog.Themes = templates.GetAllThemes()
	jsonBlog.ActiveTheme = blog.ActiveTheme
	jsonBlog.NavigationItems = blog.NavigationItems
	return &jsonBlog
}

func userToJson(user *structure.User) *JsonUser {
	var jsonUser JsonUser
	jsonUser.Id = user.Id
	jsonUser.Name = string(user.Name)
	jsonUser.Slug = user.Slug
	jsonUser.Email = string(user.Email)
	jsonUser.Image = string(user.Image)
	jsonUser.Cover = string(user.Cover)
	jsonUser.Bio = string(user.Bio)
	jsonUser.Website = string(user.Website)
	jsonUser.Location = string(user.Location)
	return &jsonUser
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
