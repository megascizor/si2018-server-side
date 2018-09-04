package message

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	r := repositories.NewUserMessageRepository()
	ur := repositories.NewUserRepository()

	sendUser, _ := ur.GetByToken(p.Params.Token)
	/*
	* Add error handling
	 */

	var ent = entities.UserMessage{}
	ent.UserID = sendUser.ID
	ent.PartnerID = p.UserID

	err := r.Create(ent)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	return si.NewPostMessageOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	r := repositories.NewUserMessageRepository()
	mr := repositories.NewUserMatchRepository()
	ur := repositories.NewUserRepository()

	uEnt, _ := ur.GetByToken(p.Token)
	/*
	* Add error handling
	 */

	parentIDs, _ := mr.FindAllByUserID(uEnt.ID)
	/*
	* Add error handling
	 */

	var ents entities.UserMessages
	for _, parentID := range parentIDs {
		messages, _ := r.GetMessages(uEnt.ID, parentID, int(*p.Limit), p.Latest, p.Oldest)
		for _, message := range messages {
			ents = append(ents, message)
		}
	}

	sEnt := ents.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}
