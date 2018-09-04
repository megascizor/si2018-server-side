package user

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	lr := repositories.NewUserLikeRepository()
	ur := repositories.NewUserRepository()

	var ents entities.Users

	// uEnt, err := ur.GetByToken(p.Token)
	uEnt, _ := ur.GetByToken(p.Token)

	// Add error handling for "UserRepository"

	oppositeGender := uEnt.GetOppositeGender()
	// exceptedIds, err := lr.FindLikeAll(uEnt.ID)
	exceptedIds, _ := lr.FindLikeAll(uEnt.ID)

	// Add error handling for "UserLikeRepository"

	// ents, err = ur.FindWithCondition(int(p.Limit), int(p.Offset), oppositeGender, exceptedIds)
	ents, _ = ur.FindWithCondition(int(p.Limit), int(p.Offset), oppositeGender, exceptedIds)

	sEnt := ents.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
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
	r := repositories.NewUserRepository()

	// ent, err := r.GetByUserID(p.UserID)
	ent, _ := r.GetByUserID(p.UserID)

	// Add error handling

	ent.ApplyParams(p.Params)

	err := r.Update(ent)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	sEnt := ent.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}
