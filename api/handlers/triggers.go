package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/headers"
	"github.com/pikami/cosmium/internal/constants"
	"github.com/pikami/cosmium/internal/datastore"
)

func (h *Handlers) GetAllTriggers(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	triggers, status := h.dataStore.GetAllTriggers(databaseId, collectionId)

	if status == datastore.StatusOk {
		c.Header(headers.ItemCount, fmt.Sprintf("%d", len(triggers)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "Triggers": triggers, "_count": len(triggers)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) GetTrigger(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	triggerId := c.Param("triggerId")

	trigger, status := h.dataStore.GetTrigger(databaseId, collectionId, triggerId)

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, trigger)
		return
	}

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) DeleteTrigger(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	triggerId := c.Param("triggerId")

	status := h.dataStore.DeleteTrigger(databaseId, collectionId, triggerId)
	if status == datastore.StatusOk {
		c.Status(http.StatusNoContent)
		return
	}

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) ReplaceTrigger(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	triggerId := c.Param("triggerId")

	var trigger datastore.Trigger
	if err := c.BindJSON(&trigger); err != nil {
		c.IndentedJSON(http.StatusBadRequest, constants.BadRequestResponse)
		return
	}

	status := h.dataStore.DeleteTrigger(databaseId, collectionId, triggerId)
	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	createdTrigger, status := h.dataStore.CreateTrigger(databaseId, collectionId, trigger)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, createdTrigger)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) CreateTrigger(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	var trigger datastore.Trigger
	if err := c.BindJSON(&trigger); err != nil {
		c.IndentedJSON(http.StatusBadRequest, constants.BadRequestResponse)
		return
	}

	createdTrigger, status := h.dataStore.CreateTrigger(databaseId, collectionId, trigger)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdTrigger)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}
