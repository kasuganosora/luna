package tag

import (
	"github.com/kabukky/journey/structure"
	"time"
)

func InsertTag(name []byte, slug string, created_at time.Time, created_by int64) (int64, error) {}
func RetrieveTags(postId int64) ([]structure.Tag, error)                                        {}
func RetrieveTag(tagId int64) (*structure.Tag, error)                                           {}
func RetrieveTagBySlug(slug string) (*structure.Tag, error)                                     {}
func RetrieveTagIdBySlug(slug string) (int64, error)                                            {}
