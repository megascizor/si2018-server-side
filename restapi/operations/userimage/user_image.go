package userimage

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"bytes"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	repoUserImage := repositories.NewUserImageRepository()
	repoUserToken := repositories.NewUserTokenRepository()

	// Validation
	entUserToken, err := repoUserToken.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUserToken == nil {
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	// Check image format
	ioLeader := bytes.NewBuffer(p.Params.Image)
	_, format, err := image.DecodeConfig(ioLeader)
	if err != nil {
		si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: 'image.DecodeConfig' failed",
			})
	}
	if format != "jpeg" && format != "png" && format != "gif" {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request: image format is invalid",
			})
	}

	// Set path for image file
	assetsPath := os.Getenv("ASSETS_PATH")
	userID := entUserToken.UserID
	strUserID := strconv.Itoa(int(userID))
	// updatedTime := time.Now()
	// strUpdatedTime := updatedTime.String()
	fileName := strings.Join([]string{"icon", strUserID}, "_") + "." + format
	filePath := path.Join(assetsPath, fileName)

	// Create file
	f, err := os.Create(filePath)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: 'os.Create' failed",
			})
	}
	defer f.Close()

	// Write image file to created file
	f.Write(p.Params.Image)

	// Create UserImage for update
	var entUserImage entities.UserImage
	entUserImage.UserID = userID
	entUserImage.Path = filePath
	entUserImage.UpdatedAt = strfmt.DateTime(time.Now())

	// Update image
	errImage := repoUserImage.Update(entUserImage)
	if errImage != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: 'Update' failed: " + errImage.Error(),
			})
	}

	// entUpdatedUserImage, errUpdatedImage := repoUserImage.GetByUserID(userID)
	// if errUpdatedImage != nil {
	// 	return si.NewPostImagesInternalServerError().WithPayload(
	// 		&si.PostImagesInternalServerErrorBody{
	// 			Code:    "500",
	// 			Message: "Internal Server Error: 'GetByUserID' failed: " + errUpdatedImage.Error(),
	// 		})
	// }

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(entUserImage.Path),
		})
}
