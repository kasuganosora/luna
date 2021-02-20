package post

import (
	"github.com/kabukky/journey/structure"
	"time"
)

func InsertPost(title []byte, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, meta_description []byte, image []byte, created_at time.Time, created_by int64) (int64, error) {
}
func InsertPostTag(post_id int64, tag_id int64) error                                        {}
func DeletePostTagsForPostId(post_id int64) error                                            {}
func DeletePostById(id int64) error                                                          {}
func RetrievePostById(id int64) (*structure.Post, error)                                     {}
func RetrievePostBySlug(slug string) (*structure.Post, error)                                {}
func RetrievePostsByUser(user_id int64, limit int64, offset int64) ([]structure.Post, error) {}
func RetrievePostsByTag(tag_id int64, limit int64, offset int64) ([]structure.Post, error)   {}
func RetrievePostsForIndex(limit int64, offset int64) ([]structure.Post, error)              {}
func RetrievePostsForApi(limit int64, offset int64) ([]structure.Post, error)                {}
func UpdatePost(id int64, title []byte, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, meta_description []byte, image []byte, updated_at time.Time, updated_by int64) error {
}

func CreatePost(data map[string]interface{}) {

}

func UpdatePost_(postID int64, data map[string]interface{}) {}
