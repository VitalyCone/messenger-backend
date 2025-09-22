package endpoints

import (
	"net/http"
	"strings"

	"github.com/VitalyCone/websocket-messenger/internal/app/apiserver/chat"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	// "github.com/sirupsen/logrus"
)

func (ep *Endpoints) ConnectUserToChats(g *gin.Context){
protocols := g.Request.Header["Sec-Websocket-Protocol"]
    var tokenString string
    logrus.Println(protocols)

    // Исправленный парсинг заголовка
    for _, p := range protocols {
        // Ищем как с пробелом после запятой, так и без
        if strings.HasPrefix(p, "token, ") {
            tokenString = p[7:]
            break
        } else if strings.HasPrefix(p, "token,") {
            tokenString = p[6:]
            break
        }
    }
	// tokenString := g.GetHeader("Sec-Websocket-Protocol")
	// if tokenString == ""{
	// 	newErrorResponse(g, http.StatusUnauthorized, "token nil")
	// 	return
	// }
	username, err := ep.services.User.GetUsernameFromToken(tokenString)
	if err != nil{
		newErrorResponse(g, http.StatusUnauthorized, err.Error())
		return
	}

	   g.Header("Sec-WebSocket-Protocol", "token")
	chat.ServeWS(g, username, ep.hub)
}