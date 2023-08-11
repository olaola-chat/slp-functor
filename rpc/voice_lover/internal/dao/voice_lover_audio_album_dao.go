package dao

import (
	"context"
	"time"

	"github.com/gogf/gf/frame/g"
	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
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

func (v *voiceLoverAudioAlbumDao) GetCountByAudioId(ctx context.Context, audioId uint64) (int, error) {
	count, err := functor2.VoiceLoverAudioAlbum.Ctx(ctx).Where("audio_id", audioId).Count()
	if err != nil {
		g.Log().Errorf("voiceLoverAudioAlbumDao GetCountByAudioId error=%v", err)
		return 0, err
	}
	return count, nil
}

func (v *voiceLoverAudioAlbumDao) GetAudioAlbumByAudioIdAlbumId(ctx context.Context, audioId uint64, albumId uint64) (*functor.EntityVoiceLoverAudioAlbum, error) {
	data, err := functor2.VoiceLoverAudioAlbum.Ctx(ctx).Where("audio_id", audioId).Where("album_id", albumId).One()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (v *voiceLoverAudioAlbumDao) Create(ctx context.Context, audioId uint64, albumId uint64) error {
	_, err := functor2.VoiceLoverAudioAlbum.Ctx(ctx).Insert(&functor.EntityVoiceLoverAudioAlbum{
		AudioId:    audioId,
		AlbumId:    albumId,
		CreateTime: uint64(time.Now().Unix()),
		UpdateTime: uint64(time.Now().Unix()),
	})
	return err
}

func (v *voiceLoverAudioAlbumDao) Del(ctx context.Context, audioId uint64, albumId uint64) error {
	_, err := functor2.VoiceLoverAudioAlbum.Ctx(ctx).Where("audio_id", audioId).Where("album_id", albumId).Delete()
	return err
}

func (v *voiceLoverAudioAlbumDao) GetAlbumIdsByAudioId(ctx context.Context, audioId uint64) ([]uint64, error) {
	res, err := functor2.VoiceLoverAudioAlbum.Ctx(ctx).Where("audio_id", audioId).FindAll()
	if err != nil {
		return nil, err
	}
	albumIds := make([]uint64, 0)
	for _, r := range res {
		albumIds = append(albumIds, r.AlbumId)
	}
	return albumIds, nil
}
