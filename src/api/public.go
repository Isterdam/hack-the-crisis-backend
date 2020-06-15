package api

import (
	"database/sql"
	"net/http"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/Isterdam/hack-the-crisis-backend/src/tz"
	"github.com/gin-gonic/gin"
	null "gopkg.in/guregu/null.v3"

	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Interval struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// GetStoreSlots godoc
// @Summary Gets all slots for a certain company on a certain day.
// @Produce json
// @Param day path string true "Day"
// @Param store path string true "Store"
// @Success 200 {array} db.Slot
// @Router /stores/{store}/day/{day}/slots [get]
func GetStoreSlots(c *gin.Context) {
	var interval Interval
	err := json.NewDecoder(c.Request.Body).Decode(&interval)
	isInter := false

	if err == nil {
		isInter = true
	}

	storeIDStr := c.Param("store")
	storeID, _ := strconv.Atoi(storeIDStr)

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var slots []db.Slot

	if !isInter {
		slots, err = db.GetSlotsByCompany(dbbb, storeID)

	} else {
		slots, err = db.GetSlotsByCompanyAndBetween(dbbb, storeID, interval.StartTime, interval.EndTime)
	}

	c.JSON(http.StatusOK, slots)
}

// BookTime godoc
// @Summary "Books" a certain time by creating a confirmation link that is sent to the user by text. Does NOT add booking to database.
// @Consume json
// @Param booking body db.Booking true "Booking"
// @Router /book [post]
func BookTime(c *gin.Context) {
	var booking db.Booking
	err := json.NewDecoder(c.Request.Body).Decode(&booking)
	// could not parse enough arguments
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Page not found",
		})
		return
	}

	ticketCode := generateTicketCode(booking)
	booking.Code = null.StringFrom(ticketCode)
	// whitelist ticked code - to be checked at confirmation if it is contained
	ConfirmedBookings[ticketCode] = booking

	url := "www.shopalone.se" + c.Request.URL.Path + "/confirm/" + ticketCode

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	timeSlot, err := db.GetSlot(dbbb, int(booking.SlotID.Int64))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
		})
	}
	store, err := db.GetCompanyByID(dbbb, int(timeSlot.CompanyID.Int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
		})
	}

	// only gets the zeroth element of zone list (because European countries only have single time zones)
	loc, err := time.LoadLocation(tz.GetCountry(store.Country.String).Zones[0].Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find the location for time zone!",
		})
	}
	timeStart := timeSlot.StartTime.Time.In(loc)
	timeStop := timeSlot.EndTime.Time.In(loc)

	confirmation := "Hello " + booking.FirstName.String + "!\n\n" + "Please confirm your booking at " + store.Name.String + " from " + timeStart.Format("15:04") + " to " + timeStop.Format("15:04") + " and get your ticket in the link below:\n\n" + url

	go Send_text(c, booking.PhoneNumber.String, confirmation)

	c.JSON(200, gin.H{
		"message": "Booking was successful",
	})
}

func generateTicketCode(booking db.Booking) string {
	// last 2 digits of current time + random num [10, 100) + booking name (where space is replaced by underscore)
	return strconv.FormatInt(time.Now().Unix(), 10)[8:] + strconv.Itoa(10+rand.Intn(90)) + strings.ReplaceAll(booking.FirstName.String, " ", "_")
}

// ConfirmBookAndGetTicket godoc
// @Summary Confirms a booking and adds it to the database if first time. Gets a ticket if it has already been added to database.
// @Produce json
// @Param code path string true "Code"
// @Router /book/confirm/{code} [post]
func ConfirmBookAndGetTicket(c *gin.Context) {
	var bookingExists bool

	code := c.Param("code")

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	booking, err := db.GetBooking(dbbb, code)

	if err != nil {
		if err == sql.ErrNoRows {
			bookingExists = false
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}
	} else {
		bookingExists = true
	}

	if ConfirmedBookings[code].PhoneNumber.String == "" && !bookingExists {
		// booking does not exist
		c.JSON(200, gin.H{
			"message": "This booking does not exist.",
		})
		return
	} else if ConfirmedBookings[code].PhoneNumber.String == "" && bookingExists {
		// booking exists and has been added to database
		c.JSON(200, gin.H{
			"message": "Ticket already confirmed!",
			"data":    booking,
		})
	} else {
		// booking exists but has not yet been added to database
		err := db.InsertBooking(dbbb, ConfirmedBookings[code])

		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		booking, err = db.GetBooking(dbbb, code)
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Ticket confirmed!",
			"data":    booking,
		})
		delete(ConfirmedBookings, code) // delete entry from map
	}
}

// Unbook godoc
// @Summary Unbooks a ticket by removing it from the database by code.
// @Param code path string true "Code"
// @Router /unbook [post]
func Unbook(c *gin.Context) {
	code := c.Param("code")

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	err := db.RemoveBooking(dbbb, code)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not remove booking!")
	}
}

// GetSlotLoad godoc
// @Summary Gets the load of a slot by returning maxAmount of customers and amount of booked customers as JSON.
// @Produce json
// @Param slotID path string true "slotID"
// @Success 200 "JSON with "maxAmount", "bookingsAmount""
// @Router /slot/{slotID}/load [get]
func GetSlotLoad(c *gin.Context) {
	slotIDStr := c.Param("slotID")
	slotID, _ := strconv.Atoi(slotIDStr)

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	slot, err := db.GetSlot(dbbb, slotID)
	if err != nil {
		fmt.Println(err)
		return
	}

	maxAmount := strconv.Itoa(int(slot.MaxAmount.Int64))

	bookings, err := db.GetBookingsBySlotID(dbbb, slotID)
	if err != nil {
		fmt.Println(err)
		return
	}

	bookingsAmount := strconv.Itoa(len(bookings))

	c.JSON(200, gin.H{
		"maxAmount":      maxAmount,
		"bookingsAmount": bookingsAmount,
	})
}

func GetCompanyAvailability(c *gin.Context) {
	var req struct {
		CompanyIDs []int     `json:"company_ids"`
		StartTime  time.Time `json:"start_time"`
		Days       int       `json:"days"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&req)

	if err != nil {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	av := make([]db.Availabilty, len(req.CompanyIDs))

	for i := range req.CompanyIDs {
		_, err := db.GetCompanyByID(dbbb, req.CompanyIDs[i])

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Company does not exist",
			})
			return
		}

		ret, err := db.GetCompanyAverageAvailability(dbbb, req.CompanyIDs[i], req.StartTime, req.Days)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Company does not exist",
			})
			return
		}

		av[i].DailyAvailable = ret
		av[i].CompanyID = req.CompanyIDs[i]

		rett, err := db.GetCompanySlotAvailability(dbbb, req.CompanyIDs[i], req.StartTime, req.Days)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Company does not exist",
			})
			return
		}

		av[i].AvailableSlots = rett
	}

	c.JSON(200, gin.H{
		"data": av,
	})
}
