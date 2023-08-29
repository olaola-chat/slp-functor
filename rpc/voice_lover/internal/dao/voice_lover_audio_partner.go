package dao

import (
	"context"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAudioPartnerDao struct {
}

var VoiceLoverAudioPartnerDao = &voiceLoverAudioPartnerDao{}

func (v *voiceLoverAudioPartnerDao) GetAudioPartnerByAudioId(ctx context.Context, id uint64) ([]*functor.EntityVoiceLoverAudioPartner, error) {
	res, err := functor2.VoiceLoverAudioPartner.Ctx(ctx).Where("audio_id", id).FindAll()
	if err != nil {
		return nil, err
	}
	return res, nil
}
