package handler

import (
	"em/controls"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AddPersonResponse struct {
	Message string          `json:"message"`
	Person  controls.Person `json:"person"`
}

type DeleteRequest struct {
	ID int `json:"id"`
}

type UpdateRequest struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}

// AddPersonHandler godoc
// @Summary Add new person
// @Description Insert a new person and auto-enrich with age, gender, and nationality
// @Accept  json
// @Produce  json
// @Param person body controls.Person true "Person data"
// @Success 200 {object} AddPersonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /add [post]
func AddPersonHandler(c *gin.Context) {
	var req controls.Person
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("DEBUG: Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON format"})
		return
	}

	log.Printf("INFO: Adding person: %+v\n", req)

	err := controls.InsertPerson(&req)
	if err != nil {
		log.Printf("ERROR: Failed to insert person: %v\n", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database insert error"})
		return
	}
	log.Println("INFO: Person inserted successfully")
	c.JSON(http.StatusOK, AddPersonResponse{
		Message: "Person added successfully",
		Person:  req,
	})
}

// DeletePersonHandler godoc
// @Summary Delete a person
// @Description Delete person by ID
// @Accept  json
// @Produce  json
// @Param id body DeleteRequest true "ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /delete [post]
func DeletePersonHandler(c *gin.Context) {
	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("DEBUG: Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	log.Printf("INFO: Deleting person with ID: %d\n", req.ID)
	err := controls.DeletePerson(req.ID)
	if err != nil {
		log.Printf("ERROR: Failed to delete person: %v\n", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database delete error"})
		return
	}
	log.Println("INFO: Person deleted successfully")
	c.JSON(http.StatusOK, SuccessResponse{Message: "Person deleted successfully"})
}

// UpdatePersonHandler godoc
// @Summary Update a person
// @Description Update all person fields by ID
// @Accept  json
// @Produce  json
// @Param person body UpdateRequest true "Person data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /update [post]
func UpdatePersonHandler(c *gin.Context) {
	var req UpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("DEBUG: Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON format"})
		return
	}

	log.Printf("INFO: Updating person ID %d with data: %+v\n", req.ID, req)
	err := controls.UpdatePerson(
		req.ID,
		req.Age,
		req.Name,
		req.Surname,
		req.Patronymic,
		req.Gender,
		req.Nationality,
	)
	if err != nil {
		log.Printf("ERROR: Failed to update person: %v\n", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database update error"})
		return
	}
	log.Println("INFO: Person updated successfully")
	c.JSON(http.StatusOK, SuccessResponse{Message: "Person updated successfully"})
}

// ViewPerson godoc
// @Summary Get filtered list of persons
// @Description Filter by name, surname, gender, nationality, age with pagination
// @Produce  json
// @Param name query string false "Name"
// @Param surname query string false "Surname"
// @Param gender query string false "Gender"
// @Param nationality query string false "Nationality"
// @Param age query string false "Age"
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} controls.Person
// @Failure 500 {object} ErrorResponse
// @Router /person [get]
func ViewPerson(c *gin.Context) {
	name := c.Query("name")
	surname := c.Query("surname")
	gender := c.Query("gender")
	nationality := c.Query("nationality")
	age := c.Query("age")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	log.Printf("INFO: Getting people with filters: name=%s, surname=%s, gender=%s, nationality=%s, age=%s, limit=%s, offset=%s\n",
		name, surname, gender, nationality, age, limit, offset)

	persons, err := controls.GetFilteredPersons(name, surname, gender, nationality, age, limit, offset)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve persons: %v\n", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database query error"})
		return
	}

	log.Printf("INFO: Retrieved %d person(s)\n", len(persons))
	c.JSON(http.StatusOK, persons)
}