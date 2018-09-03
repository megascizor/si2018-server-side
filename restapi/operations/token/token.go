package token

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetTokenByUserID(p si.GetTokenByUserIDParams) middleware.Responder {
	r := repositories.NewUserTokenRepository()

	ent, err := r.GetByUserID(p.UserID)

	if err != nil {
		return si.NewGetTokenByUserIDInternalServerError().WithPayload(
			&si.GetTokenByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if ent == nil {
		return si.NewGetTokenByUserIDNotFound().WithPayload(
			&si.GetTokenByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Token Not Found",
			})
	}

	sEnt := ent.Build()
	return si.NewGetTokenByUserIDOK().WithPayload(&sEnt)
}
