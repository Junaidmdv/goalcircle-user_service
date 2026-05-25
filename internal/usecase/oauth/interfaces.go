package oauth

import (
	"context" 
	 uc_dtos "github.com/Junaidmdv/goalcircle-user_service/internal/usecase/dtos"

)

type OauthUsecase interface { 
	GoogleOauth(context.Context,*uc_dtos.GoogleOauthReq)(*uc_dtos.GoogleOauthRes,error)
	GoogleOauthCallback(context.Context,*uc_dtos.GoogleCallbackReq)(*uc_dtos.GoogleCallbackRes,error)
}