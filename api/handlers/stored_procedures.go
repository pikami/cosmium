package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (h *Handlers) GetAllStoredProcedures(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	sps, status := h.repository.GetAllStoredProcedures(databaseId, collectionId)

	if status == repositorymodels.StatusOk {
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(sps)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "StoredProcedures": sps, "_count": len(sps)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) GetStoredProcedure(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	spId := c.Param("spId")

	sp, status := h.repository.GetStoredProcedure(databaseId, collectionId, spId)

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, sp)
		return
	}

	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) DeleteStoredProcedure(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	spId := c.Param("spId")

	status := h.repository.DeleteStoredProcedure(databaseId, collectionId, spId)
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

func (h *Handlers) ReplaceStoredProcedure(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	spId := c.Param("spId")

	var sp repositorymodels.StoredProcedure
	if err := c.BindJSON(&sp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
		return
	}

	status := h.repository.DeleteStoredProcedure(databaseId, collectionId, spId)
	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	createdSP, status := h.repository.CreateStoredProcedure(databaseId, collectionId, sp)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, createdSP)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) CreateStoredProcedure(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	var sp repositorymodels.StoredProcedure
	if err := c.BindJSON(&sp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
		return
	}

	createdSP, status := h.repository.CreateStoredProcedure(databaseId, collectionId, sp)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdSP)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
