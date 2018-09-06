package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	"time"
)

// GetLikes get likes
func GetLikes(p si.GetLikesParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserLike := repositories.NewUserLikeRepository()
	repoUserMatch := repositories.NewUserMatchRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	if p.Limit < 0 {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "400",
				Message: "Bad Request: limit in query must be not less than 1",
			})
	}
	if p.Offset < 0 {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request: offset in query must be not less than 1",
			})
	}

	entUserToken, err := repoUserToken.GetByToken(p.Token)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserToken == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get excepted users (matched users)
	userID := entUserToken.UserID
	exceptIDs, err := repoUserMatch.FindAllByUserID(userID)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Get likes
	entUserLikes, err := repoUserLike.FindGotLikeWithLimitOffset(userID, int(p.Limit), int(p.Offset), exceptIDs)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Converte UserLikes -> LikeUserResponses
	var entLikeUserResponses entities.LikeUserResponses
	for _, entUserLike := range entUserLikes {
		var res entities.LikeUserResponse
		entUser, err := repoUser.GetByUserID(entUserLike.UserID)
		if err != nil {
			return si.NewGetLikesInternalServerError().WithPayload(
				&si.GetLikesInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
		if entUser == nil {
			return si.NewGetUsersBadRequest().WithPayload(
				&si.GetUsersBadRequestBody{
					Code:    "500",
					Message: "Bad Request: 'GetByUserID' failed: " + err.Error(),
				})
		}
		res.ApplyUser(*entUser)
		entLikeUserResponses = append(entLikeUserResponses, res)
	}

	sEnt := entLikeUserResponses.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

// PostLike post like
func PostLike(p si.PostLikeParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserLike := repositories.NewUserLikeRepository()
	repoUserMatch := repositories.NewUserMatchRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	entUserToken, err := repoUserToken.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserToken == nil {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "500",
				Message: "Unauthorized",
			})
	}

	// Get users (sender and receiver)
	sendUser, err := repoUser.GetByUserID(entUserToken.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if sendUser == nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request: 'GetByUserID' (sender) failed: " + err.Error(),
			})
	}

	recvUser, err := repoUser.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if recvUser == nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request: 'GetByUserID' (receiver) failed: " + err.Error(),
			})
	}

	// Check whether a sender do "like" to hisself
	if sendUser.ID == recvUser.ID {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request: must not do 'like' for myself",
			})
	}

	// Check whether a sender do "like" to the same-gender person
	if sendUser.Gender == recvUser.Gender {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request: must not do 'like' to the same-gender",
			})
	}

	// Check whether a sender do two times "like" to the same person
	exceptIDs, err := repoUserLike.FindLikedIDs(sendUser.ID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	for _, exceptID := range exceptIDs {
		if recvUser.ID == exceptID {
			return si.NewPostLikeBadRequest().WithPayload(
				&si.PostLikeBadRequestBody{
					Code:    "404",
					Message: "Bad Request : Already did 'like'",
				})
		}
	}

	// Apply parameters (for "like")
	var entUserLike entities.UserLike
	entUserLike.UserID = sendUser.ID
	entUserLike.PartnerID = recvUser.ID
	entUserLike.CreatedAt = strfmt.DateTime(time.Now())
	entUserLike.UpdatedAt = entUserLike.CreatedAt

	// Do "like"
	err = repoUserLike.Create(entUserLike)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Make matching when they have done "like" each other
	likedIDs, err := repoUserLike.FindLikedIDs(recvUser.ID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	for _, likedID := range likedIDs {
		if sendUser.ID == likedID {
			// Apply parameters (for matching)
			var entUserMatch entities.UserMatch
			entUserMatch.UserID = entUserLike.PartnerID
			entUserMatch.PartnerID = entUserLike.UserID
			entUserMatch.CreatedAt = strfmt.DateTime(time.Now())
			entUserMatch.UpdatedAt = entUserMatch.CreatedAt

			// Matching
			err = repoUserMatch.Create(entUserMatch)
			if err != nil {
				return si.NewPostLikeInternalServerError().WithPayload(
					&si.PostLikeInternalServerErrorBody{
						Code:    "500",
						Message: "Internal Server Error",
					})
			}

			return si.NewPostLikeOK().WithPayload(
				&si.PostLikeOKBody{
					Code:    "201",
					Message: "Matched",
				})
		}
	}

	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "201",
			Message: "Liked",
		})
}
