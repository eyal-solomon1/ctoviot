package api

import (
	"context"
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	db "github.com/eyal-solomon1/ctoviot/db/sqlc"
	"github.com/eyal-solomon1/ctoviot/token"
	"github.com/gin-gonic/gin"
)

type getVideosResponse struct {
	Videos []db.Video `json:"videos"`
}

type deleteVideoRequest struct {
	VideoIdentifier string `json:"video_identifier"`
}

func (server *Server) deleteVideo(ctx *gin.Context) {
	// * Get user from payload and verify it √
	// * Get request video based of requesting user √
	// * Delete video (entry + video + remote s3 file) √
	var json deleteVideoRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusBadRequest, errorResponse(invalidRequesParams))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUser(ctx, payload.Username)

	if err != nil {
		errMsg := errors.New("user is invalid")
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusBadRequest, errorResponse(errMsg))
		return
	}

	vid, err := server.store.GetVideo(ctx, db.GetVideoParams{
		VideoIdentifier: json.VideoIdentifier,
		Owner:           payload.Username,
	})

	if err != nil {
		errMsg := errors.New(fmt.Sprintf("couldn't find requested video for %s", payload.Username))
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusBadRequest, errorResponse(errMsg))
		return
	}

	if vid.Owner != payload.Username {
		errMsg := errors.New("user not authorized to delete this video")
		server.errLogger(VideoAPILogGroup, errMsg)
		ctx.JSON(http.StatusUnauthorized, errorResponse(errMsg))
		return
	}

	vid, err = server.store.DeleteVideo(ctx, json.VideoIdentifier)

	if err != nil {
		errMsg := errors.New("couldn't delete your video from DB")
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusBadRequest, errorResponse(errMsg))
		return
	}

	_, err = server.awsService.DeleteFile(ctx, s3.DeleteObjectInput{
		Bucket: &server.config.S3BucketName,
		Key:    &vid.VideoRemotePath,
	})

	if err != nil {
		errMsg := errors.New("couldn't delete your video from cloud backup")
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusBadRequest, errorResponse(errMsg))
		return
	}

	server.infoLogger(VideoAPILogGroup, fmt.Sprintf("deleted %s video for %s user", vid.VideoName, user.Username))
	ctx.JSON(http.StatusNoContent, gin.H{"ok": true})
}

func (server *Server) getVideos(ctx *gin.Context) {
	// TODO add a dynamic limit and offset
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUser(ctx, payload.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := errors.New(fmt.Sprintf("user %v not found", payload.Username))
			server.errLogger(VideoAPILogGroup, err)
			ctx.JSON(http.StatusNotFound, errorResponse(errMsg))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	videos, err := server.store.ListVideos(ctx, db.ListVideosParams{Owner: payload.Username, Limit: 5, Offset: 0})
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := errors.New(fmt.Sprintf("no videos found for user %v", payload.Username))
			server.errLogger(VideoAPILogGroup, err)
			ctx.JSON(http.StatusNotFound, errorResponse(errMsg))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := getVideosResponse{
		Videos: videos,
	}

	server.infoLogger(VideoAPILogGroup, fmt.Sprintf("found %d videos from %s user", len(videos), user.Username))
	ctx.JSON(http.StatusOK, okResponse(resp))

}

func (server *Server) createVideo(ctx *gin.Context) {

	// TODO
	// * checks if video name already exists √
	// * transform video to audio ussing ffmpeg √
	// * get audio transcription from audio using openai √
	// * save received transcription to a 'subtitles' file √
	// * transform audio + transcription to a transcribed video √
	// * save file to AWS for the requedted user √
	// * create a new video db transaction √
	// * cleanup (remove temp files) √
	// * ?

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	videoName := ctx.PostForm("videoName")
	description := ctx.PostForm("description")
	inputFile, err := ctx.FormFile("inputFile")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": errorResponse(err)})
		return
	}

	if videoName == "" || description == "" || inputFile == nil {
		errMsg := errors.New("all fields are required")
		server.errLogger(VideoAPILogGroup, errMsg)
		ctx.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": errorResponse(errMsg)})
		return
	}

	file, err := ctx.FormFile("inputFile")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": errorResponse(err)})
		return
	}
	videoFile, err := server.saveVideoInputToFileSystem(ctx, file)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": errorResponse(err)})
		return
	}
	defer videoFile.Close()

	videoFilePath := videoFile.Name()
	defer server.deleteUneededFiles(videoFilePath)

	validVideo, err := server.videoIsValid(videoFilePath, payload.Username)

	if !validVideo && err != nil {
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}

	audioFilePath, err := server.ffmpegService.CreateAudioFile(videoFilePath)
	defer server.deleteUneededFiles(audioFilePath)

	if err != nil {
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": errorResponse(err)})
		return
	}

	transcriptionFilePath, err := server.createTranscriptionFileFromAudio(audioFilePath, description)
	defer server.deleteUneededFiles(transcriptionFilePath)

	if err != nil {
		errMsg := errors.New("failed creating an audio file from the provided video")
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": errorResponse(errMsg)})
		return
	}

	transcribedFilePath, err := server.ffmpegService.MatchTranscriptionsToVideo(videoFilePath, transcriptionFilePath)
	defer server.deleteUneededFiles(transcribedFilePath)

	if err != nil {
		errMsg := errors.New("failed matching transcription and video")
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": errorResponse(errMsg)})
		return
	}

	fileToUpload, err := os.Open(transcribedFilePath)

	if err != nil {
		errMsg := errors.New("failed reading file from filesystem")
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": errorResponse(errMsg)})
		return
	}

	defer fileToUpload.Close()

	videoNameEncoded := b64.StdEncoding.EncodeToString([]byte(videoName))
	duration, err := server.ffmpegService.GetFileDuration(transcribedFilePath)

	if err != nil {
		server.errLogger(VideoAPILogGroup, err)
		errMsg := errors.New("couldn't get video duration")
		ctx.JSON(http.StatusInternalServerError, errorResponse(errMsg))
		return
	}

	objKey := fmt.Sprintf("videos/users/%s/%s.mp4", payload.Username, videoNameEncoded)
	_, err = server.awsService.CreateFile(ctx, s3.PutObjectInput{
		Bucket: &server.config.S3BucketName,
		Key:    &objKey,
		Body:   fileToUpload,
	})

	if err != nil {
		errMsg := errors.New("failed saving video to the cloud")
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": errorResponse(errMsg)})
		return
	}

	txResult, err := server.store.VideoTx(ctx, db.VideoTxParam{
		Username: payload.Username,
		Video: db.CreateVideoParams{
			Owner:           payload.Username,
			VideoName:       videoName,
			VideoIdentifier: videoNameEncoded,
			VideoLength:     int64(duration),
			VideoRemotePath: objKey,
			VideoDecs:       description,
		},
	})

	if err != nil {
		server.errLogger(VideoAPILogGroup, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": errorResponse(err)})
		return
	}

	server.infoLogger(VideoAPILogGroup, fmt.Sprintf("created a new video named %s for %s user", videoName, payload.Username))
	ctx.JSON(http.StatusOK, gin.H{"ok": true, "payload": txResult})
}

// videoIsValid checks if the video length is under the maximum allowed (currectly 15s) and if the users passed his videos quata
func (server *Server) videoIsValid(filePath, username string) (bool, error) {

	// TODO move video quata + video length limit to config somehow

	fileDuration, err := server.ffmpegService.GetFileDuration(filePath)

	if err != nil {
		return false, errors.New("couldn't get video duration")
	}

	if fileDuration > 15 {
		return false, errors.New("video duration is longer then the allowed duration")
	}

	count, err := server.store.GetUsersVideosCount(context.Background(), username)

	if err != nil {
		return false, errors.New("couldn't get current user video's")
	}

	if count >= 3 {
		return false, errors.New("user reached videos qutata")
	}

	return true, nil

}

// saveVideoInputToFileSystem recives 'formdata' input video file,and saves it to filesystem
// returning *os.File and an error
func (server *Server) saveVideoInputToFileSystem(ctx *gin.Context, file *multipart.FileHeader) (*os.File, error) {
	tempFile, err := os.CreateTemp("/tmp", "video_*.mp4")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": errorResponse(err)})
		return nil, err
	}
	defer tempFile.Close()

	tempFilePath := tempFile.Name()

	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": errorResponse(err)})
		return nil, err
	}
	return tempFile, nil
}

// createTranscriptionFileFromAudio creates a 'subtitles' file from the recived audio file
func (server *Server) createTranscriptionFileFromAudio(audioFilePath, description string) (string, error) {
	audioTranscription, err := server.openaiService.GetAudioTranscription(audioFilePath, description)

	if err != nil {
		return "", err
	}

	tempTranscriptionFile, err := os.CreateTemp("/tmp", "transcription_*.srt")
	if err != nil {
		return "", err
	}

	_, err = tempTranscriptionFile.Write([]byte(audioTranscription))
	if err != nil {
		return "", err
	}

	return tempTranscriptionFile.Name(), nil
}

// deleteUneededFiles removes all file which are uneeded, recives a string slice
func (server *Server) deleteUneededFiles(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
