package delivery

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	gen "main.go/internal/auth/proto"
	"main.go/internal/feed"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
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
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, feed.CardsToSend{Cards: cards})
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
		UID := request.Context().Value(RequestUserID).(types.UserID)
		body, err := io.ReadAll(request.Body)
		if err != nil { // TODO эти два блока вынести в отдельную функцию и напихать её во все ручки
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		var requestBody feed.CreateLikeRequest
		err = easyjson.Unmarshal(body, &requestBody)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		err = deliver.usecase.CreateLike(UID, requestBody.Profile2, request.Context())

		if err != nil {
			logrus.Info(err.Error())
		}

		if err != nil && err.Error() != "sql: Rows are closed" {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't update profile: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		if err != nil && err.Error() == feed.NoMatchFoundErr.Error() {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("created like")
			requests.SendSimpleResponse(respWriter, request, http.StatusOK, "ok")
			return
		}
		if err != nil && err.Error() == feed.NoLikesLeftErr.Error() {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("no likes left")
			requests.SendSimpleResponse(respWriter, request, http.StatusConflict, err.Error())
			return
		}

		err = deliver.sendMatchNotice(request.Context(), requestBody.Profile2, UID)
		err = deliver.sendMatchNotice(request.Context(), UID, requestBody.Profile2)

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't send match: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
		}

		requests.SendSimpleResponse(respWriter, request, http.StatusOK, "ok")
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("create like sent response")
	}
}

func (deliver *FeedHandler) sendMatchNotice(ctx context.Context, id1, id2 types.UserID) error {
	connection, ok := deliver.usecase.GetConnection(ctx, id1)
	if !ok {
		return nil
	}
	resp := feed.Message{MsgType: "match", Properties: feed.MsgProperties{Sender: id1, Receiver: id2}}
	respCoded, err := easyjson.Marshal(resp)
	if err != nil {
		return err
	}
	err = connection.WriteMessage(1, respCoded)

	if err != nil && err.Error() == feed.WSClosedErr.Error() {
		deliver.usecase.DeleteConnection(ctx, id1)
		return nil
	}
	return err
}

func (deliver *FeedHandler) CreateDislike() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("not implemented dislike")
		requests.SendSimpleResponse(respWriter, request, http.StatusOK, "ok")
	}
}

func (deliver *FeedHandler) GetChat() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		body, err := io.ReadAll(request.Body)
		if err != nil { // TODO эти два блока вынести в отдельную функцию и напихать её во все ручки
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		var requestBody feed.GetChatRequest
		err = easyjson.Unmarshal(body, &requestBody)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		messages, images, persons, err := deliver.usecase.GetChat(request.Context(), request.Context().Value(RequestUserID).(types.UserID), requestBody.Person)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}
		if messages == nil {
			messages = []feed.Message{}
		}

		resp := feed.GetChatFull{
			Messages: messages,
			Images:   images,
			Person:   persons[0],
		}

		requests.SendResponse(respWriter, request, http.StatusOK, resp)
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
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
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
		var message feed.MessageToReceive
		err = easyjson.Unmarshal(mess, &message)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Error unmarshalling message: ", err.Error())
			continue
		}
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("got ws message")

		_ = connection.WriteMessage(mt, mess)

		message.Properties.Sender = sender
		_, err = deliver.usecase.SaveMessage(ctx, message)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Error saving message: ", err.Error())
			continue
		}

		// отправляем сообщение получателю
		resConnection, ok := deliver.usecase.GetConnection(ctx, message.Properties.Receiver)
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
		matches, err := deliver.AuthManager.GetMatches(request.Context(), &gen.GetMatchesRequest{
			UserID:    int64(request.Context().Value(RequestUserID).(types.UserID)),
			RequestID: logger.RequestID,
		})
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Error("can't get chats: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		messages, err := deliver.usecase.GetLastMessages(request.Context(),
			int64(request.Context().Value(RequestUserID).(types.UserID)), getIds(matches))

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Error("can't get chats: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		var allChats feed.AllChats
		allChats.Chats = make([]feed.ChatPreview, len(matches.Chats))
		for i := range matches.Chats {
			allChats.Chats[i] = feed.ChatPreview{
				PersonID: matches.Chats[i].PersonID,
				Name:     matches.Chats[i].Name,
				Photo:    matches.Chats[i].Photo,
				//Premuim:  matches.Chats[i].Premium,
			}

			for j := range messages {
				if int64(messages[j].Properties.Sender) == allChats.Chats[i].PersonID || int64(messages[j].Properties.Receiver) == allChats.Chats[i].PersonID {
					allChats.Chats[i].LastMessage = feed.Message{
						MsgType: "message",
						Properties: feed.MsgProperties{
							Id:       messages[j].Properties.Id,
							Data:     messages[j].Properties.Data,
							Sender:   messages[j].Properties.Sender,
							Receiver: messages[j].Properties.Receiver,
							Time:     messages[j].Properties.Time,
						},
					}
					break
				}
			}
		}

		requests.SendResponse(respWriter, request, http.StatusOK, allChats)
	}
}

func getIds(l *gen.GetMatchesResponse) []int64 {
	res := make([]int64, len(l.Chats))
	for i := range l.Chats {
		res[i] = l.Chats[i].PersonID
	}
	return res
}

func (deliver *FeedHandler) CreateClaim() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		body, err := io.ReadAll(request.Body)
		if err != nil { // TODO эти два блока вынести в отдельную функцию и напихать её во все ручки
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}
		var requestBody feed.CreateClaimRequest
		err = easyjson.Unmarshal(body, &requestBody)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("can't unmarshal body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		err = deliver.usecase.CreateClaim(request.Context(), requestBody.Type,
			int64(request.Context().Value(RequestUserID).(types.UserID)), requestBody.ReceiverID)

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't create claim: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		requests.SendSimpleResponse(respWriter, request, http.StatusOK, "ok")
	}
}

func (deliver *FeedHandler) GetAlClaims() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		claims, err := deliver.usecase.GetClaims(request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't get claim types: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}
		requests.SendResponse(respWriter, request, http.StatusOK, feed.ClaimsToSend{Claims: claims})
	}
}

func MetricTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next.ServeHTTP(respWriter, request)

		end := time.Since(start)
		path := request.URL.Path
		if path != "/metrics" {
			feed.TotalHits.WithLabelValues().Inc()
			feed.HitDuration.WithLabelValues(request.Method, path).Set(float64(end.Milliseconds()))
		}
	})
}
