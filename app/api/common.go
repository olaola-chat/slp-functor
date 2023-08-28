package api

//func OutputCustomError(r *ghttp.Request, err consts.CommonError) {
//	g.Log().Ctx(r.GetCtx()).Info(err)
//	response.Output(r, &pb.CommonResp{
//		Success: false,
//		Code:    err.Code(),
//		Msg:     err.Msg(),
//		Data:    nil,
//	})
//}
//
//func OutputCustomData(r *ghttp.Request, data proto.Message) {
//	var anyData *anypb.Any
//	if data != nil {
//		any1, err := anypb.New(data)
//		if err != nil {
//			g.Log().Ctx(r.GetCtx()).Errorf("AnyPB Error %+v", err)
//			OutputCustomError(r, consts.ERROR_SYSTEM)
//			return
//		}
//		anyData = any1
//	}
//	response.Output(r, &pb.CommonResp{
//		Success: true,
//		Code:    consts.ERROR_SUCCESS.Code(),
//		Msg:     consts.ERROR_SUCCESS.Msg(),
//		Data:    anyData,
//	})
//}
