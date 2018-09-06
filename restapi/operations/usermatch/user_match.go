package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserMatch := repositories.NewUserMatchRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	if p.Limit <= 0 {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code:    "400",
				Message: "Bad Request: limit in query must be not less than 1",
			})
	}
	if p.Offset < 0 {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code:    "400",
				Message: "Bad Request: offset in query must be not less than 1",
			})
	}

	entUserToken, err := repoUserToken.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserToken == nil {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get matches
	entUserMatches, err := repoUserMatch.FindByUserIDWithLimitOffset(entUserToken.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
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

		entUser, err := repoUser.GetByUserID(partnerID)
		if err != nil {
			return si.NewGetMatchesInternalServerError().WithPayload(
				&si.GetMatchesInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
		if entUser == nil {
			return si.NewGetMatchesBadRequest().WithPayload(
				&si.GetMatchesBadRequestBody{
					Code:    "400",
					Message: "Bad Request: 'GetByUserID' failed: " + err.Error(),
				})
		}
		res.ApplyUser(*entUser)
		entMatchUserResponses = append(entMatchUserResponses, res)
	}

	sEnt := entMatchUserResponses.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
