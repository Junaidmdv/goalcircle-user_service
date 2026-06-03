package dtos 


type BlockUserReq struct{
	UserId string `json:"user_id" validate:"required"` 
}


type UnBlockUserReq struct{
	UserId string `json:"user_id" validate:"required"` 
}