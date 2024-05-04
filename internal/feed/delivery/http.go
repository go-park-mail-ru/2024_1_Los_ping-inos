package delivery

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	. "main.go/config"
	gen "main.go/internal/auth/proto"
	"main.go/internal/feed"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	"net/http"
)

type FeedHandler struct {
	usecase     feed.UseCase
	AuthManager gen.AuthHandlClient
}

func NewFeedDelivery(uc feed.UseCase, am gen.AuthHandlClient) *FeedHandler {
	return &FeedHandler{
		usecase:     uc,
		AuthManager: am,
	}
}

// GetCardsHandler godoc
// @Summary Получить ленту
// @Tags    Продукт
// @Router  /cards [get]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Success 200		  {array}  models.PersonWithInterests
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
// @Failure 500       {string} string
func (deliver *FeedHandler) GetCardsHandler() func(http.ResponseWriter, *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)

		cards, err := deliver.usecase.GetCards(request.Context().Value(RequestUserID).(types.UserID), request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, cards)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent cards")
	}
}

// CreateLike godoc
// @Summary Создать лайк
// @Tags    Лайк
// @Router  /like [post]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Param   profile2   body   string false "profile id to like"
// @Success 200
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
func (deliver *FeedHandler) CreateLike() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)

		body, err := io.ReadAll(request.Body)
		if err != nil { // TODO эти два блока вынести в отдельную функцию и напихать её во все ручки
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		var requestBody requests.CreateLikeRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		err = deliver.usecase.CreateLike(request.Context().Value(RequestUserID).(types.UserID), requestBody.Profile2, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't update profile: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, nil)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("create like sent response")
	}
}

func (deliver *FeedHandler) CreateDislike() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("not implemented dislike")
		requests.SendResponse(respWriter, request, http.StatusOK, nil)
	}
}

func (deliver *FeedHandler) GetChat() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		body, err := io.ReadAll(request.Body)
		if err != nil { // TODO эти два блока вынести в отдельную функцию и напихать её во все ручки
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		var requestBody feed.GetChatRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		messages, err := deliver.usecase.GetChat(request.Context(), request.Context().Value(RequestUserID).(types.UserID), requestBody.Person)
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, messages)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("sent chat")
	}
}

func upgradeConnection() websocket.Upgrader {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Пропускаем любой запрос
		},
	}
	return upgrader
}

func (deliver *FeedHandler) ServeMessages() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		upgrader := upgradeConnection()
		connection, err := upgrader.Upgrade(respWriter, request, nil)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Error("can't open connection: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}
		deliver.handleWebsocket(request.Context(), connection)
	}
}

func (deliver *FeedHandler) handleWebsocket(ctx context.Context, connection *websocket.Conn) {
	logger := ctx.Value(Logg).(Log)
	sender := ctx.Value(RequestUserID).(types.UserID)
	deliver.usecase.AddConnection(ctx, connection, sender)

	for {
		mt, mess, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		var message feed.Message
		err = json.Unmarshal(mess, &message)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Error unmarshalling message: ", err.Error())
			continue // TODO отправлять ошибку на фронт?
		}
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("got ws message")

		message.Sender = sender
		_, err = deliver.usecase.SaveMessage(ctx, message)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Error saving message: ", err.Error())
			continue // TODO отправлять ошибку на фронт?
		}

		// отправляем сообщение получателю
		resConnection, ok := deliver.usecase.GetConnection(ctx, message.Receiver)
		if !ok {
			continue
		}
		err = resConnection.WriteMessage(mt, mess)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Error sending message: ", err.Error())
		}
	}
}

func (deliver *FeedHandler) GetAllChats() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		chats, err := deliver.AuthManager.GetMatches(request.Context(), &gen.GetMatchesRequest{
			UserID:    int64(request.Context().Value(RequestUserID).(types.UserID)),
			RequestID: logger.RequestID,
		})
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Error("can't get chats: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}
		requests.SendResponse(respWriter, request, http.StatusOK, chats.Chats)
	}
}
