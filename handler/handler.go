package handler

import (
	"em/controls"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPersonHandler(c *gin.Context) {
	var req controls.Person
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "JSON error"})
		return
	}

	err := controls.InsertPerson(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error in inserting "})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Messsage": "Inserted","person":req})
}

func DeletePersonHandler(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("JSON Binding Error:", err) 
		c.JSON(http.StatusBadRequest, gin.H{"message": "JSON problem"})
		return
	}

	err := controls.DeletePerson(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Deleting error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Deleted"})
}

func UpdatePersonHandler(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
		Name string `json:"name"`
		Surname string `json:"surname"`
		Patronymic string `json:"patronymic"`
		Gender string `json:"gender"`
		Age int `json:"age"`
		Nationality string `json:"nationality"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("JSON Binding Error:", err) 
		c.JSON(http.StatusBadRequest, gin.H{"Message": "JSON problem"})
		return
	}

	err := controls.UpdatePerson(req.ID, req.Age, req.Name,req.Surname,req.Patronymic,req.Gender,req.Nationality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Update problem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Updated"})
}

func ViewPerson(c * gin.Context){
	name := c.Query("name")
	surname := c.Query("surname")
	gender := c.Query("gender")
	nationality := c.Query("nationality")
	age := c.Query("age")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	persons, err := controls.GetFilteredPersons(name, surname, gender, nationality, age, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, persons)
}