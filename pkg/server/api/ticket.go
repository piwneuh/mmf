package api

import (
	"context"

	"mmf/pkg/ws"
	"mmf/wires"

	"github.com/gin-gonic/gin"
)

func RegisterTicket(router *gin.Engine, ctx context.Context) {
	tickets := router.Group("/tickets")
	{
		// tickets.POST("/submit/:game", submitTicket)
		tickets.GET("/fetch/:queue", fetchTickets)
	}

	router.GET("/ws/:queue/:steamId", wsGet)

}

func wsGet(c *gin.Context) {
	queue := c.Param("queue")
	steamId := c.Param("steamId")

	ws.StartWebSocket(queue, steamId, c)
}

func fetchTickets(c *gin.Context) {
	queue := c.Param("queue")
	tickets := wires.Instance.TicketService.GetAllTickets(c, queue)
	c.JSON(200, tickets)
}
