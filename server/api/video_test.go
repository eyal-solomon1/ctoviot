package api

import (
	"bytes"

	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	mockdb "github.com/eyal-solomon1/ctoviot/db/mock"
	db "github.com/eyal-solomon1/ctoviot/db/sqlc"
	"github.com/eyal-solomon1/ctoviot/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateVideoAPI(t *testing.T) {
	user, _ := createRandomUser(t)

	var param gin.H = gin.H{
		"username": user.Username,
		"video_info": gin.H{
			"owner":             "John Doe",
			"video_name":        "My Video",
			"video_length":      int64(15),
			"video_remote_path": "/path/to/video",
			"video_decs":        "A sample video",
		},
	}

	testCases := []struct {
		name          string
		params        gin.H
		buildStubs    func(store *mockdb.MockStore, aws *mockdb.MockAWS, openai *mockdb.MockOpenAI, ffmpeg *mockdb.MockFFMPEG)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "Happy scenario",
			params: param,
			buildStubs: func(store *mockdb.MockStore, aws *mockdb.MockAWS, openai *mockdb.MockOpenAI, ffmpeg *mockdb.MockFFMPEG) {
				// db
				store.EXPECT().GetUsersVideosCount(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(int64(1), nil)
				store.EXPECT().VideoTx(gomock.Any(), gomock.Any()).Times(1).Return(db.VideoTxResult{}, nil)

				// services
				aws.EXPECT().CreateFile(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)
				ffmpeg.EXPECT().CreateAudioFile(gomock.Any()).Times(1).Return("../assets/test.mp3", nil)
				ffmpeg.EXPECT().GetFileDuration(gomock.Any()).Times(1).Return(10.0, nil)
				openai.EXPECT().GetAudioTranscription(gomock.Any(), gomock.Any()).Times(1).Return("test transcription", nil)
				ffmpeg.EXPECT().MatchTranscriptionsToVideo(gomock.Any(), gomock.Any()).Times(1).Return("../assets/ready.mp4", nil)

			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		// {
		// 	name:   "UserNotFound",
		// 	params: param,
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusNotFound, recorder.Code)
		// 	},
		// },
		// {
		// 	name:   "BadRequest",
		// 	params: gin.H{},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0).Return(user, nil)
		// 		store.EXPECT().VideoTx(gomock.Any(), gomock.Any()).Times(0).Return(db.VideoTxResult{}, nil)
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name:   "MissingAuthentication",
		// 	params: param,
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0).Return(db.User{}, nil)
		// 		store.EXPECT().VideoTx(gomock.Any(), gomock.Any()).Times(0).Return(db.VideoTxResult{}, nil)
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
		// {
		// 	name:   "UnAuthorized",
		// 	params: param,
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, nil)
		// 		store.EXPECT().VideoTx(gomock.Any(), gomock.Any()).Times(0).Return(db.VideoTxResult{}, nil)
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "none", time.Minute)

		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
	}

	for i := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		store := mockdb.NewMockStore(ctrl)
		aws := mockdb.NewMockAWS(ctrl)
		openai := mockdb.NewMockOpenAI(ctrl)
		ffmpeg := mockdb.NewMockFFMPEG(ctrl)

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			filePath := filepath.Join("../assets", "test.mp4")
			fileName := "test.mp4"

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)

			fileField, _ := writer.CreateFormFile("inputFile", fileName)
			file, err := os.Open(filePath)
			require.NoError(t, err)

			_, err = io.Copy(fileField, file)
			require.NoError(t, err)

			_ = file.Close()

			writer.Close()

			tc.buildStubs(store, aws, openai, ffmpeg)

			server := newTestServer(t, WithStore(store), WithAWSService(aws), WithFFMPEGService(ffmpeg), WithOpenAIService(openai))
			recorder := httptest.NewRecorder()

			urlS := "/new_video"

			request, err := http.NewRequest(http.MethodPost, urlS, body)
			request.PostForm = url.Values{
				"videoName":   []string{"Mock Video"},
				"description": []string{"Mock Description"},
			}
			request.Header.Add("Content-Type", writer.FormDataContentType())
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}
