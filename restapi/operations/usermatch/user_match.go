package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	ur := repositories.NewUserRepository()
	mr := repositories.NewUserMatchRepository()

	uEnt, _ := ur.GetByToken(p.Token)

	mlEnts, _ := mr.FindByUserIDWithLimitOffset(uEnt.ID, int(p.Limit), int(p.Offset))
	var ents entities.MatchUserResponses
	for _, matchUser := range mlEnts {
		var res = entities.MatchUserResponse{}
		// user, err := ur.GetByUserID(matchUser.UserID)

		// Add error handler

		user, _ := ur.GetByUserID(matchUser.UserID)
		res.ApplyUser(*user)
		ents = append(ents, res)
	}

	sEnt := ents.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
