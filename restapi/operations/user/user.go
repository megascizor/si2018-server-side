package user

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

// GetUsers get users list
func GetUsers(p si.GetUsersParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserImage := repositories.NewUserImageRepository()
	repoUserLike := repositories.NewUserLikeRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	if p.Limit <= 0 {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request: limit in query must be more than 1",
			})
	}
	if p.Offset < 0 {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request: offset in query must be not less than 1",
			})
	}

	entUserToken, err := repoUserToken.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserToken == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get user
	entUser, err := repoUser.GetByUserID(entUserToken.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUser == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request: 'GetByUserID' failed",
			})
	}

	// Get except user IDs (Users I liked)
	exceptIDs, err := repoUserLike.FindIDsILiked(entUser.ID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	exceptIDs = append(exceptIDs, entUser.ID) // Add me

	// Get users list
	users, err := repoUser.FindWithCondition(int(p.Limit), int(p.Offset), entUser.GetOppositeGender(), exceptIDs)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if users == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request: 'GetBuToken' failed",
			})
	}

	// Get users' IDs
	var userIDs []int64
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	// Get image entities
	entUserImages, err := repoUserImage.GetByUserIDs(userIDs)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var entUsers entities.Users
	entUsers = entities.Users(users)
	sEnt := entUsers.Build()

	// Put image URI to user.image_URI
	for _, u := range sEnt {
		for _, entUserImage := range entUserImages {
			if u.ID == entUserImage.UserID {
				u.ImageURI = entUserImage.Path
			}
		}
	}

	return si.NewGetUsersOK().WithPayload(sEnt)
}

// GetProfileByUserID gets user profile
func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserImage := repositories.NewUserImageRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	entUserToken, err := repoUserToken.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserToken == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get users (look and looked)
	lookUser, err := repoUser.GetByUserID(entUserToken.UserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if lookUser == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found",
			})
	}

	lookedUser, err := repoUser.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if lookedUser == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found",
			})
	}

	// Check whether look and looked user is the same-gender
	// You can see your profile.
	if lookUser.Gender == lookedUser.Gender && lookUser.ID != lookedUser.ID {
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code:    "400",
				Message: "Bad Request: the same-gender",
			})
	}

	// Get user image
	entUserImage, err := repoUserImage.GetByUserID(lookedUser.ID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserImage == nil {
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code:    "400",
				Message: "Bad Request: 'GetByUserID' failed",
			})
	}

	sEnt := lookedUser.Build()
	// Input image uri
	sEnt.ImageURI = entUserImage.Path
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

// PutProfile update user profile
func PutProfile(p si.PutProfileParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserImage := repositories.NewUserImageRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	entUserToken, err := repoUserToken.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserToken == nil || p.Params.Token == "" {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get user profile
	entUser, err := repoUser.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUser == nil {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message: "Bad Request: 'GetByUserID' failed",
			})
	}
	if p.UserID != entUserToken.UserID {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:    "403",
				Message: "Forbidden: must not update others profile",
			})
	}

	// Update profile
	entUser.ApplyParams(p.Params)

	err = repoUser.Update(entUser)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Get updated profile
	entUpdatedUser, err := repoUser.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUpdatedUser == nil {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message: "Bad Request: 'GetByUserID' (after update) failed",
			})
	}

	// Get image URI
	entUserImage, err := repoUserImage.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserImage == nil {
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code:    "400",
				Message: "Bad Request: 'GetByUserID' failed",
			})
	}

	sEnt := entUpdatedUser.Build()
	// Input image uri
	sEnt.ImageURI = entUserImage.Path
	return si.NewPutProfileOK().WithPayload(&sEnt)
}
