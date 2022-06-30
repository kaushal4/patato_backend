package listings

import (
	"fmt"
	"net/http"
	"potato/backend/db"
	"potato/backend/user"
	"potato/backend/utils"

	"github.com/gin-gonic/gin"
)

type Listings struct {
	Item      string `json:"item" form:"item" db:"item"`
	Freshness string `json:"freshness" form:"freshness" db:"freshness"`
	Photo     string `db:"photo"`
	UserId    string `db:"userid"`
}

func AddListings(c *gin.Context) {
	var listings Listings
	claims, err := user.CheckClaims(c)
	if utils.CheckError(err, http.StatusBadRequest, c) {
		return
	}
	file, err := c.FormFile("photo")
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	filename, err := utils.UploadImage("listings/photos/", file, c)

	if err != nil {

		fmt.Println(err)
		return
	}
	err = c.Bind(&listings)
	if utils.CheckError(err, http.StatusInternalServerError, c) {
		fmt.Println("There's an error here bubbo")
		return
	}
	con, err := db.Connect()
	if utils.CheckError(err, http.StatusInternalServerError, c) {
		return
	}
	defer con.Close()
	listings.Photo = filename
	listings.UserId = claims.Id
	_, err = con.NamedExec("insert into listings(item,freshness,photo,userId) values(:item,:freshness,:photo,:userid)", &listings)
	if utils.CheckError(err, http.StatusInternalServerError, c) {
		return
	}
	fmt.Println(listings, filename, claims.Id)
	c.JSON(http.StatusOK, gin.H{"status": "inserted"})
}
