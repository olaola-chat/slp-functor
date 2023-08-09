package dao

import (
	"context"
)

type voiceLoverAlbumDao struct {
}

var VoiceLoverAlbumDao = &voiceLoverAlbumDao{}

func (v *voiceLoverAlbumDao) Post(ctx context.Context) error {
	return nil
}
