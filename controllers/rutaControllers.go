package controllers

import (
	"backend/initializers"
	"backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRuta(c *gin.Context) {
	var body models.RutaModel

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "INSERT INTO ruta (activa, destino, cuota_destino) VALUES (?, ?, ?)"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(body.Activa, body.Destino, body.CuotaDestino)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error getting last insert ID", err)
		return
	}

	ruta := models.RutaModel{
		Id:           int(id),
		Activa:       body.Activa,
		Destino:      body.Destino,
		CuotaDestino: body.CuotaDestino,
	}

	sqlStatement2 := "INSERT INTO ruta_punto_control (ruta, punto_control) VALUES (?, ?)"

	stmt, err = initializers.DB.Prepare(sqlStatement2)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing second statement", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(ruta.Id, ruta.Destino)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ruta": ruta})
}

func AddPuntoControl(c *gin.Context) {
	var body struct {
		RutaId         int
		PuntoControlId int
	}

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "INSERT INTO ruta_punto_control (ruta, punto_control) VALUES (?, ?)"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(body.RutaId, body.PuntoControlId)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Punto control added successfully"})
}

func IsPedidoInRutaTableEmpty() (bool, error) {
	sqlStatement := "SELECT p.* FROM pedido p INNER JOIN ruta r ON p.ruta = r.id"

	rows, err := initializers.DB.Query(sqlStatement)
	if err != nil {
		fmt.Println("Error with the query: ", err)
		return false, err
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
	}

	return rowCount == 0, err
}

func UpdateRuta(c *gin.Context) {
	var body struct {
		Id           int
		Activa       bool
		CuotaDestino float32
	}

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	isPedidoInRutaTableEmpty, err := IsPedidoInRutaTableEmpty()
	if !isPedidoInRutaTableEmpty {
		handleError(c, http.StatusConflict, "The ruta has pedidos in queue", err)
		return
	}

	sqlStatement := "UPDATE ruta SET activa = ?, cuota_destino = ? WHERE id = ?"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.Activa, body.CuotaDestino, body.Id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "No ruta found with the provided id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ruta updated successfully"})

}

func GetRutas(c *gin.Context) {

	sqlStatement := "SELECT * FROM ruta"

	rows, err := initializers.DB.Query(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing query", err)
		return
	}
	defer rows.Close()

	var rutas []models.RutaModel

	for rows.Next() {
		var ruta models.RutaModel
		err := rows.Scan(&ruta.Id, &ruta.Activa, &ruta.Destino, &ruta.CuotaDestino)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Error scanning row", err)
			return
		}
		rutas = append(rutas, ruta)
	}

	if err := rows.Err(); err != nil {
		handleError(c, http.StatusInternalServerError, "Error iterating over rows", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"rutas": rutas})
}

func GetPuntosControlRuta(c *gin.Context) {
	var body struct {
		RutaId int
	}
	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "SELECT p.* FROM punto_control p INNER JOIN ruta_punto_control r ON p.id = r.punto_control AND r.ruta = ?"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(body.RutaId)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

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

func RemovePuntoControl(c *gin.Context) {
	var body struct {
		PuntoControlId int
		RutaId         int
	}

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "DELETE FROM ruta_punto_control WHERE ruta = ? AND punto_control = ?"

	// Prepare the SQL statement
	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(body.RutaId, body.PuntoControlId)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Punto control removed successfully"})
}
