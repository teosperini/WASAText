package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.wrap(rt.doLoginAPI))
	rt.router.GET("/users", rt.wrap(rt.getUsersAPI))
	rt.router.PUT("/username", rt.wrap(rt.putUsernameAPI))
	rt.router.POST("/conversations", rt.wrap(rt.postConversationAPI))
	rt.router.GET("/conversations", rt.wrap(rt.getConversationsAPI))
	rt.router.GET("/conversations/:convId", rt.wrap(rt.getConversationAPI))
	rt.router.POST("/conversations/:convId", rt.wrap(rt.postMessageAPI))
	rt.router.GET("/image", rt.wrap(rt.getImageAPI))
	rt.router.PUT("/image", rt.wrap(rt.putImageAPI))
	rt.router.PUT("/conversations/:convId/name", rt.wrap(rt.putNameGroupAPI))
	rt.router.PUT("/conversations/:convId/image", rt.wrap(rt.putImageGroupAPI))
	rt.router.PUT("/conversations/:convId/members", rt.wrap(rt.addMemberToGroupAPI))
	rt.router.DELETE("/conversations/:convId/members", rt.wrap(rt.removeMemberFromGroupAPI))
	rt.router.POST("/upload", rt.wrap(rt.postMediaAPI))
	rt.router.DELETE("/conversations/:convId/messages/:messId", rt.wrap(rt.deleteMessageAPI))
	rt.router.POST("/conversations/:convId/messages/:messId", rt.wrap(rt.forwardMessageAPI))
	rt.router.PUT("/conversations/:convId/messages/:messId/emoji", rt.wrap(rt.putEmojiAPI))
	rt.router.DELETE("/conversations/:convId/messages/:messId/emoji", rt.wrap(rt.deleteEmojiAPI))

	rt.router.ServeFiles("/media/*filepath", http.Dir("./uploads/images"))
	rt.router.ServeFiles("/assets/*filepath", http.Dir("./assets"))

	return rt.router
}
