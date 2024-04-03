package controllers

import (
	"backend/initializers"
	"backend/models"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, statusCode int, message string, err error) {
	fmt.Println("Error:", err)
	c.JSON(statusCode, gin.H{"error": message})
}

func CreateUsuario(c *gin.Context) {
	var body models.Usuario

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "INSERT INTO usuario (username, password, nombre, activo, rol) VALUES (?, ?, ?, ?, ?)"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(body.Username, body.Password, body.Nombre, body.Activo, body.Rol)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing statement", err)
		return
	}

	user := models.Usuario{
		Username: body.Username,
		Password: body.Password,
		Nombre:   body.Nombre,
		Activo:   body.Activo,
		Rol:      body.Rol,
	}

	c.JSON(http.StatusOK, gin.H{"usuario": user})
}

func GetUsuarios(c *gin.Context) {
	type User struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nombre   string `json:"nombre"`
		Activo   bool   `json:"activo"`
		Rol      int    `json:"rol"`
	}

	sqlStatement := "SELECT * FROM usuario"

	rows, err := initializers.DB.Query(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error executing query", err)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.Password, &user.Nombre, &user.Activo, &user.Rol)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Error scanning row", err)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		handleError(c, http.StatusInternalServerError, "Error iterating over rows", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"usuarios": users})
}

func GetUsuarioByUsername(c *gin.Context) {
	var body struct {
		Username string
	}

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "SELECT * FROM usuario WHERE username = ?"
	row := initializers.DB.QueryRow(sqlStatement, body.Username)

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nombre   string `json:"nombre"`
		Activo   bool   `json:"activo"`
		Rol      int    `json:"rol"`
	}
	err := row.Scan(&user.Username, &user.Password, &user.Nombre, &user.Activo, &user.Rol)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		handleError(c, http.StatusInternalServerError, "Error retrieving user", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"usuario": user})
}

func UpdateUser(c *gin.Context) {
	var body struct {
		Username string
		Password string
		Nombre   string
		Activo   bool
	}

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "UPDATE usuario SET password = ?, nombre = ?, activo = ? WHERE username = ?"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.Password, body.Nombre, body.Activo, body.Username)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "No user found with the provided username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	var body struct {
		Username string
	}

	if err := c.BindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	sqlStatement := "DELETE FROM usuario WHERE username = ?"

	stmt, err := initializers.DB.Prepare(sqlStatement)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing statement", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.Username)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "No user found with the provided username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
