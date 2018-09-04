package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	ulr := repositories.NewUserLikeRepository()
	ur := repositories.NewUserRepository()
	mr := repositories.NewUserMatchRepository()

	// var ents entities.UserLikes

	// uEnt, err := ur.GetByToken(p.Token)
	uEnt, _ := ur.GetByToken(p.Token)

	// Add error handler for "UserRepository"

	// exceptedIds, err := mr.FindAllByUserID(uEnd.ID)
	exceptedIds, _ := mr.FindAllByUserID(uEnt.ID)

	// Add error handler for "UserMatchRepository"

	// ents, err = lr.FindGotLikeWithLimitOffset(uEnt.ID, int(p.Limit), int(p.Offset), exceptedIds)
	ulEnts, _ := ulr.FindGotLikeWithLimitOffset(uEnt.ID, int(p.Limit), int(p.Offset), exceptedIds)

	// Add error handler for "UserLikeRepository"

	var ents entities.LikeUserResponses
	for _, likeUser := range ulEnts {
		var res = entities.LikeUserResponse{}
		// user, err := ur.GetByUserID(likeUser.UserID)

		// Add error handler

		user, _ := ur.GetByUserID(likeUser.UserID)
		res.ApplyUser(*user)
		ents = append(ents, res)
	}

	// ents, _ := lur.FindByIDs(ids)

	sEnt := ents.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	// r := repositories.NewUserLikeRepository()

	return si.NewPostLikeOK()
}
