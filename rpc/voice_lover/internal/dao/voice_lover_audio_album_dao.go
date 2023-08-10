package dao

import (
	"context"

	"github.com/gogf/gf/frame/g"
	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
)

type voiceLoverAudioAlbumDao struct {
}

var VoiceLoverAudioAlbumDao = &voiceLoverAudioAlbumDao{}

func (v *voiceLoverAudioAlbumDao) GetCountByAlbumId(ctx context.Context, albumId uint64) (int, error) {
	total, err := functor2.VoiceLoverAudioAlbum.Ctx(ctx).Where("album_id", albumId).Count()
	if err != nil {
		g.Log().Errorf("voiceLoverAudioAlbumDao GetAudioCountByAlbumId error=%v", err)
		return 0, err
	}
	return total, nil
}
