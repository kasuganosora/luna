package structure

// Blog: settings that are used for template execution
type Blog struct {
	Url             string
	Title           string
	Description     string
	Logo            string
	Cover           string
	AssetPath       string
	PostCount       int64
	PostsPerPage    int64
	ActiveTheme     string
	NavigationItems []Navigation
}
