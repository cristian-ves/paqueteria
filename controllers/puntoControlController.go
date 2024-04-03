package controllers

import (
	"backend/initializers"
	"backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePuntoControl(c *gin.Context) {
	var body models.PuntoControl

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "INSERT INTO punto_control (localizacion, paquetes_maximos, tarifa, operador) VALUES (?, ?, ?, ?)"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(body.Localizacion, body.PaquetesMaximos, body.Tarifa, body.Operador)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error getting last insert ID", err)
		return
	}

	newPuntoControl := models.PuntoControl{
		Id:              int(id),
		Localizacion:    body.Localizacion,
		PaquetesMaximos: body.PaquetesMaximos,
		Tarifa:          body.Tarifa,
		Operador:        body.Operador,
	}

	c.JSON(http.StatusOK, gin.H{"usuario": newPuntoControl})
}

func GetPuntosControlByOperador(c *gin.Context) {
	var body struct {
		Operador string
	}

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "SELECT * FROM punto_control WHERE operador = '" + body.Operador + "'"
	rows, err := initializers.DB.Query(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing query", err)
		return
	}
	defer rows.Close()

	var puntosControl []models.PuntoControl

	for rows.Next() {
		var puntoControl models.PuntoControl
		err := rows.Scan(&puntoControl.Id, &puntoControl.Localizacion, &puntoControl.PaquetesMaximos, &puntoControl.Tarifa, &puntoControl.Operador)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Error scanning row", err)
			return
		}
		puntosControl = append(puntosControl, puntoControl)
	}

	if err := rows.Err(); err != nil {
		handleError(c, http.StatusInternalServerError, "Error iterating over rows", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"puntosControl": puntosControl})
}

func UpdatePuntoControl(c *gin.Context) {

	var body models.PuntoControl

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "SELECT * FROM proceso WHERE punto_control = ?"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}

	result, err := stmt.Exec(body.Id)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error getting rows affected", err)
		return
	}

	if rowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "This punto de control has paquetes in queue"})
		return
	}

	sqlStatement = "UPDATE punto_control SET localizacion = ?, paquetes_maximos = ?, tarifa = ?, operador = ? WHERE id = ?"

	stmt, err = initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(body.Localizacion, body.PaquetesMaximos, body.Tarifa, body.Operador, body.Id)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PuntoControl updated successfully"})

}

func GetPuntosControl(c *gin.Context) {
	sqlStatement := "SELECT * FROM punto_control"

	rows, err := initializers.DB.Query(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing query", err)
		return
	}
	defer rows.Close()

	var puntosControl []models.PuntoControl

	for rows.Next() {
		var puntoControl models.PuntoControl
		fmt.Println(puntoControl)
		err := rows.Scan(&puntoControl.Id, &puntoControl.Localizacion, &puntoControl.PaquetesMaximos, &puntoControl.Tarifa, &puntoControl.Operador)
		fmt.Println(puntoControl)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Error scanning row", err)
			return
		}
		puntosControl = append(puntosControl, puntoControl)
	}

	if err := rows.Err(); err != nil {
		handleError(c, http.StatusInternalServerError, "Error iterating over rows", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"puntos control": puntosControl})
}
