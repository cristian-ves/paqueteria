package controllers

import (
	"backend/initializers"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDatosGlobales(c *gin.Context) {

	sqlStatement := "SELECT * FROM datos_globales WHERE id=1"

	row, err := initializers.DB.Query(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing query", err)
		return
	}
	defer row.Close()

	var datosGlobales models.DatosGlobales

	row.Next()
	newErr := row.Scan(&datosGlobales.Id, &datosGlobales.TarifaOperacion, &datosGlobales.PrecioLibra)
	if newErr != nil {
		handleError(c, http.StatusInternalServerError, "Error scanning row", newErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"datos globales": datosGlobales})

}

func UpdateDatosGlobales(c *gin.Context) {
	var body models.DatosGlobales

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "UPDATE datos_globales SET tarifa_operacion = ?, precio_libra = ? WHERE id = 1"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.TarifaOperacion, body.PrecioLibra)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error getting rows affected", err)
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No instace found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})
}
