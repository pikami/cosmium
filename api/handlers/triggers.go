package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (h *Handlers) GetAllTriggers(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	triggers, status := h.repository.GetAllTriggers(databaseId, collectionId)

	if status == repositorymodels.StatusOk {
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(triggers)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "Triggers": triggers, "_count": len(triggers)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) GetTrigger(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	triggerId := c.Param("triggerId")

	trigger, status := h.repository.GetTrigger(databaseId, collectionId, triggerId)

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, trigger)
		return
	}

	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) DeleteTrigger(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	triggerId := c.Param("triggerId")

	status := h.repository.DeleteTrigger(databaseId, collectionId, triggerId)
	if status == repositorymodels.StatusOk {
		c.Status(http.StatusNoContent)
		return
	}

	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) ReplaceTrigger(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	triggerId := c.Param("triggerId")

	var trigger repositorymodels.Trigger
	if err := c.BindJSON(&trigger); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
		return
	}

	status := h.repository.DeleteTrigger(databaseId, collectionId, triggerId)
	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	createdTrigger, status := h.repository.CreateTrigger(databaseId, collectionId, trigger)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, createdTrigger)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) CreateTrigger(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	var trigger repositorymodels.Trigger
	if err := c.BindJSON(&trigger); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
		return
	}

	createdTrigger, status := h.repository.CreateTrigger(databaseId, collectionId, trigger)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdTrigger)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
