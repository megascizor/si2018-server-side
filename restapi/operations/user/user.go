package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	// ur := repositories.NewUserRepository()
	// tr := repositories.NewUserTokenRepository()
	//
	// token_user, err := tr.GetByToken(p.Token)
	//
	// user, err := ur.GetByUserID(token.UserID)
	//
	// p.Limit = 20
	// p.Offset = 0
	//
	// var ids []int64
	// users, err := ur.FindWithCondition(p.Limit, p.Offset, user.Gender, ids)
	//
	return si.NewGetUsersOK()
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	r := repositories.NewUserRepository()

	ent, err := r.GetByUserID(p.UserID)

	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if ent == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Profile Not Found",
			})
	}
	// if .Token != token.Token {
	// 	return si.NewGetProfileByUserIDUnauthorized().WithPayload(
	// 		&si.GetProfileByUserIDUnauthorizedBody{
	// 			Code:    "401",
	// 			Message: "Unauthorized",
	// 		})
	// }

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}
