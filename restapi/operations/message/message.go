package message

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/restapi/utils"
	"github.com/go-openapi/runtime/middleware"
	// strfmt "github.com/go-openapi/strfmt"
	"time"
)

// PostMessage post message
func PostMessage(p si.PostMessageParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserMatch := repositories.NewUserMatchRepository()
	repoUserMessage := repositories.NewUserMessageRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	entUserToken, err := repoUserToken.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserToken == nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get users (sender, receiver)
	sendUser, err := repoUser.GetByUserID(entUserToken.UserID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if sendUser == nil {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "404",
				Message: "Bad Request: 'GetByUserID' (sender) failed: " + err.Error(),
			})
	}

	// recvUser, err := repoUser.GetByUserID(p.UserID)
	// if err != nil {
	// 	return si.NewPostMessageInternalServerError().WithPayload(
	// 		&si.PostMessageInternalServerErrorBody{
	// 			Code:    "500",
	// 			Message: "Internal Server Error",
	// 		})
	// }
	// if recvUser == nil {
	// 	return si.NewPostMessageBadRequest().WithPayload(
	// 		&si.PostMessageBadRequestBody{
	// 			Code:    "404",
	// 			Message: "Bad Request: 'GetByUserID' (receiver) failed" + err.Error(),
	// 		})
	// }

	sendID := sendUser.ID
	recvID := p.UserID
	// recvID := recvUser.ID

	// Check whether to match
	matchedIDs, errMatch := repoUserMatch.FindAllByUserID(sendID)
	if errMatch != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if !utils.IsContained(recvID, matchedIDs) {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request: not mathced",
			})
	}

	// Apply parameters (message)
	var entMessage entities.UserMessage
	entMessage.UserID = sendID
	entMessage.PartnerID = recvID
	entMessage.Message = p.Params.Message
	// entMessage.CreatedAt = strfmt.DateTime(time.Now())
	// entMessage.UpdatedAt = strfmt.DateTime(time.Now())

	err = repoUserMessage.Create(entMessage)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

// GetMessages get messages
func GetMessages(p si.GetMessagesParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserMatch := repositories.NewUserMatchRepository()
	repoUserMessage := repositories.NewUserMessageRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	if *p.Limit <= 0 {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request: limit in query must be more than 1",
			})
	}
	if (p.Oldest != nil) && (p.Latest != nil) {
		// Latest
		if time.Time(*p.Latest).Before(time.Time(*p.Oldest)) {
			return si.NewGetMessagesBadRequest().WithPayload(
				&si.GetMessagesBadRequestBody{
					Code:    "400",
					Message: "Bad Request: latest must be more than oldest",
				})
		}
	}
	/* Todo: Add validater for p.Latest and p.Oldest */

	entUserToken, err := repoUserToken.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserToken == nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Get users' entities (sender, receiver)
	sendUser, err := repoUser.GetByUserID(entUserToken.UserID)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if sendUser == nil {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "404",
				Message: "Bad Request: 'GetByUserID' (sender) failed: " + err.Error(),
			})
	}

	// recvUser, err := repoUser.GetByUserID(p.UserID)
	// if err != nil {
	// 	return si.NewGetMessagesInternalServerError().WithPayload(
	// 		&si.GetMessagesInternalServerErrorBody{
	// 			Code:    "500",
	// 			Message: "Internal Server Error",
	// 		})
	// }
	// if recvUser == nil {
	// 	return si.NewGetUsersBadRequest().WithPayload(
	// 		&si.GetUsersBadRequestBody{
	// 			Code:    "404",
	// 			Message: "Bad Request: 'GetByUserID' (receiver) failed: " + err.Error(),
	// 		})
	// }

	sendID := sendUser.ID
	recvID := p.UserID
	// recvID := recvUser.ID

	// Check whether to match
	matchedIDs, err := repoUserMatch.FindAllByUserID(sendID)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if !utils.IsContained(recvID, matchedIDs) {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request: not matched",
			})
	}

	// Get messages
	var entMessages entities.UserMessages
	entMessages, err = repoUserMessage.GetMessages(sendID, recvID, int(*p.Limit), p.Latest, p.Oldest)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// for _, message := range messages {
	// 	entMessages = append(entMessages, message)
	// }

	sEnt := entMessages.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}
