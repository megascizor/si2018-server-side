package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	// Repositories
	repoUser := repositories.NewUserRepository()
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
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get matches
	entUserMatches, errUserMatches := repoUserMatch.FindByUserIDWithLimitOffset(entUserToken.UserID, int(p.Limit), int(p.Offset))
	if errUserMatches != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error : Matches",
			})
	}

	var entMatchUserResponses entities.MatchUserResponses
	for _, entUserMatch := range entUserMatches {
		var res entities.MatchUserResponse
		// Get partner ID
		var partnerID int64
		if entUserToken.UserID == entUserMatch.UserID {
			partnerID = entUserMatch.PartnerID
		} else if entUserToken.UserID == entUserMatch.PartnerID {
			partnerID = entUserMatch.UserID
		}

		entUser, errUser := repoUser.GetByUserID(partnerID)
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
		entMatchUserResponses = append(entMatchUserResponses, res)
	}

	sEnt := entMatchUserResponses.Build()

	// uEnt, _ := ur.GetByToken(p.Token)
	//
	// mlEnts, _ := mr.FindByUserIDWithLimitOffset(uEnt.ID, int(p.Limit), int(p.Offset))
	// var ents entities.MatchUserResponses
	// for _, matchUser := range mlEnts {
	// 	var res = entities.MatchUserResponse{}
	// 	// user, err := ur.GetByUserID(matchUser.UserID)
	//
	// 	// Add error handler
	//
	// 	user, _ := ur.GetByUserID(matchUser.UserID)
	// 	res.ApplyUser(*user)
	// 	ents = append(ents, res)
	// }
	//
	// sEnt := ents.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
