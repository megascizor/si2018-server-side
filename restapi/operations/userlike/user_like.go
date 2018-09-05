package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	strfmt "github.com/go-openapi/strfmt"
	"time"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserLike := repositories.NewUserLikeRepository()
	repoUserMatch := repositories.NewUserMatchRepository()
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
				Message: "User Token Not Found",
			})
	}

	if p.Token != entUserToken.Token {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get user's entitie
	// entUser, errUser := repoUser.GetByUserID(entUserToken.UserID)
	// if errUser != nil {
	// 	return si.NewGetUsersInternalServerError().WithPayload(
	// 		&si.GetUsersInternalServerErrorBody{
	// 			Code:    "500",
	// 			Message: "Internal Server Error : User",
	// 		})
	// }
	// if entUser == nil {
	// 	return si.NewGetUsersBadRequest().WithPayload(
	// 		&si.GetUsersBadRequestBody{
	// 			Code:    "400",
	// 			Message: "Bad Request : User",
	// 		})
	// }

	// Get excepted users
	exceptIDs, errMatch := repoUserMatch.FindAllByUserID(entUserToken.UserID)
	if errMatch != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : Match",
			})
	}

	entUserLikes, errLikes := repoUserLike.FindGotLikeWithLimitOffset(entUserToken.UserID, int(p.Limit), int(p.Offset), exceptIDs)
	if errLikes != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : Likes",
			})
	}

	var entLikeUserResponses entities.LikeUserResponses
	for _, entUserLike := range entUserLikes {
		var res entities.LikeUserResponse
		entUser, errUser := repoUser.GetByUserID(entUserLike.UserID)
		if errUser != nil {
			return si.NewGetUsersInternalServerError().WithPayload(
				&si.GetUsersInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error : User (for)",
				})
		}
		if entUser == nil {
			return si.NewGetUsersBadRequest().WithPayload(
				&si.GetUsersBadRequestBody{
					Code:    "500",
					Message: "Bad Request : User (for)",
				})
		}
		res.ApplyUser(*entUser)
		entLikeUserResponses = append(entLikeUserResponses, res)
	}

	sEnt := entLikeUserResponses.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	// Repositories
	repoUser := repositories.NewUserRepository()
	repoUserLike := repositories.NewUserLikeRepository()
	repoUserMatch := repositories.NewUserMatchRepository()

	/* Todo: Add error handling */
	sendUser, _ := repoUser.GetByToken(p.Params.Token)
	recvUser, _ := repoUser.GetByUserID(p.UserID)
	/* ------------------------ */

	sendID := sendUser.ID
	recvID := p.UserID

	// Check whether same-gender
	if sendUser.Gender == recvUser.Gender {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "404",
				Message: "Bad Request : Gender is the same",
			})
	}

	// Two times "like" error handling to the same person
	exceptIDs, _ := repoUserLike.FindLikedIDs(sendID)
	for _, exceptID := range exceptIDs {
		if recvUser.ID == exceptID {
			return si.NewPostLikeBadRequest().WithPayload(
				&si.PostLikeBadRequestBody{
					Code:    "404",
					Message: "Bad Request : Already did 'like'",
				})
		}
	}

	// Initialize
	var entUserLike entities.UserLike
	entUserLike.UserID = sendID
	entUserLike.PartnerID = recvID
	entUserLike.CreatedAt = strfmt.DateTime(time.Now())
	entUserLike.UpdatedAt = entUserLike.CreatedAt

	// Do "like"
	err := repoUserLike.Create(entUserLike)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Make matching when they have done "like" each other
	likedIDs, _ := repoUserLike.FindLikedIDs(recvID)
	for _, likedID := range likedIDs {
		if sendID == likedID {
			// Initialize
			var entUserMatch entities.UserMatch
			entUserMatch.UserID = entUserLike.PartnerID
			entUserMatch.PartnerID = entUserLike.UserID
			entUserMatch.CreatedAt = strfmt.DateTime(time.Now())
			entUserMatch.UpdatedAt = entUserMatch.CreatedAt

			// Matching
			err := repoUserMatch.Create(entUserMatch)
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
