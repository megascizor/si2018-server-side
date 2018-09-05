package user

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	// Repositories
	repoUser := repositories.NewUserRepository()
	repoUserLike := repositories.NewUserLikeRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	entUserToken, errToken := repoUserToken.GetByToken(p.Token)
	if errToken != nil {
		return si.NewGetTokenByUserIDInternalServerError().WithPayload(
			&si.GetTokenByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : Token",
			})
	}
	if entUserToken == nil {
		return si.NewGetTokenByUserIDNotFound().WithPayload(
			&si.GetTokenByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found : Token",
			})
	}

	if p.Token != entUserToken.Token {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get user's entitie
	entUser, errUser := repoUser.GetByUserID(entUserToken.UserID)
	if errUser != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : User",
			})
	}
	if entUser == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request : User",
			})
	}

	// Get excepted users
	oppositeGender := entUser.GetOppositeGender()
	exceptIDs, errLike := repoUserLike.FindLikeAll(entUser.ID)
	if errLike != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : Like",
			})
	}

	var entUsers entities.Users
	entUsers, errUsers := repoUser.FindWithCondition(int(p.Limit), int(p.Offset), oppositeGender, exceptIDs)
	if errUsers != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : Users",
			})
	}

	sEnt := entUsers.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	// Repositories
	repoUser := repositories.NewUserRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	entUserToken, errToken := repoUserToken.GetByUserID(p.UserID)
	if errToken != nil {
		return si.NewGetTokenByUserIDInternalServerError().WithPayload(
			&si.GetTokenByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : Token",
			})
	}
	if entUserToken == nil {
		return si.NewGetTokenByUserIDNotFound().WithPayload(
			&si.GetTokenByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found : Token",
			})
	}

	if p.Token != entUserToken.Token {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get profile
	entUser, errUser := repoUser.GetByUserID(p.UserID)
	if errUser != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : User",
			})
	}
	if entUser == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found : User",
			})
	}
	/* Need Bad Request? */
	if p.UserID != entUser.ID {
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code:    "400",
				Message: "Bad Request : User",
			})
	}

	sEnt := entUser.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	// Repositories
	repoUser := repositories.NewUserRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	entUserToken, errToken := repoUserToken.GetByUserID(p.UserID)
	if errToken != nil {
		return si.NewGetTokenByUserIDInternalServerError().WithPayload(
			&si.GetTokenByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : Token",
			})
	}
	if entUserToken == nil {
		return si.NewGetTokenByUserIDNotFound().WithPayload(
			&si.GetTokenByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found : Token",
			})
	}

	if p.Params.Token != entUserToken.Token {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get user profile
	entUser, errUser := repoUser.GetByUserID(p.UserID)
	if errUser != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : User",
			})
	}
	if p.UserID != entUserToken.UserID {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:    "403",
				Message: "Forbidden",
			})
	}
	if entUser == nil {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message: "Bad Request : User",
			})
	}

	// Apply parameters
	entUser.ApplyParams(p.Params)

	errUpdate := repoUser.Update(entUser)
	if errUpdate != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : put profile",
			})
	}

	sEnt := entUser.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}
