package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Add_company(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}

	dbbb := dbb.(*db.DB)
	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)

	if err != nil {
		fmt.Printf("hello2 %s", err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(comp.Password.String), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	comp.Password.String = string(hash)

	err = db.InsertCompany(dbbb, comp)

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
	})
}

func Get_company(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	comp, _ := db.GetCompanies(dbbb)

	c.JSON(200, comp)
}

func Update_company(c *gin.Context) {
	if !Is_authorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	var newComp db.Company
	newComp, err = db.UpdateCompany(dbbb, comp)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	c.JSON(200, newComp)
}

func Add_slots(c *gin.Context) {
	if !Is_authorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var slots []db.Slot
	err := json.NewDecoder(c.Request.Body).Decode(&slots)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, slot := range slots {
		db.AddSlot(dbbb, slot)
	}
}

func Get_slots(c *gin.Context) {
	if !Is_authorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)
	if err != nil {
		fmt.Println(err)
		return
	}

	var slots []db.Slot
	slots, err = db.GetSlotsByCompany(dbbb, int(comp.ID.Int64))
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, slots)
}

func Update_slot(c *gin.Context) {
	if !Is_authorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var slot db.Slot
	err := json.NewDecoder(c.Request.Body).Decode(&slot)
	if err != nil {
		fmt.Println(err)
		return
	}

	var newSlot db.Slot
	newSlot, err = db.UpdateSlot(dbbb, slot)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, newSlot)
}

func Get_slot(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var slot db.Slot
	err := json.NewDecoder(c.Request.Body).Decode(&slot)
	if err != nil {
		fmt.Println(err)
		return
	}

	var newSlot db.Slot
	newSlot, err = db.GetSlot(dbbb, int(slot.ID.Int64))
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, newSlot)
}

/*
func Get_code(c *gin.Context) {

}
*/
