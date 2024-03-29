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
	"html"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GetStoreSlots godoc
// @Summary Gets all slots for a certain company on a certain day.
// @Produce json
// @Param day path string true "Day"
// @Param store path string true "Store"
// @Success 200 {array} db.Slot
// @Router /stores/{store}/day/{day}/slots [get]
func GetStoreSlots(c *gin.Context) {
	dbbb := c.MustGet("db").(*db.DB)

	var req struct {
		CompanyID int `uri:"companyID"`
		//Includes start_time and end_time
		db.Filters
	}

	err := c.ShouldBindQuery(&req) //Binds with start_time and end_time in from query params

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Time Format.",
		})
		return
	}

	err = c.ShouldBindUri(&req) //Binds with :storeID in url

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid CompanyID.",
		})
		return
	}

	var slots []db.Slot

	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		slots, err = db.GetSlotsAvailableByCompany(dbbb, req.CompanyID)
	} else {
		slots, err = db.GetSlotsAvailableByCompanyAndBetween(dbbb, req.CompanyID, req.StartTime, req.EndTime)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Slots could not be fetched.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    slots,
	})
	return
}

// BookTime godoc
// @Summary "Books" a certain time by creating a confirmation link that is sent to the user by text. Does NOT add booking to database.
// @Consume json
// @Param booking body db.Booking true "Booking"
// @Router /book [post]
func BookTime(c *gin.Context) {
	dbb := c.MustGet("db").(*db.DB)

	var booking db.Booking
	err := json.NewDecoder(c.Request.Body).Decode(&booking)
	// could not parse enough arguments
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON.",
		})
		return
	}

	booking.Sanitize()

	/*
		if hasAlreadyBooked(booking.PhoneNumber.String, dbbb, c) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "This phone number has already booked a time!",
			})
			return
		}
	*/

	timeSlot, err := db.GetSlot(dbb, int(booking.SlotID.Int64))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get slot!",
		})
		return
	}

	if timeSlot.Booked.Int64 >= timeSlot.MaxAmount.Int64 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "This slot is full.",
		})
		return
	}

	ticketCode := generateTicketCode(booking)
	booking.Code = null.StringFrom(ticketCode)
	// whitelist ticked code - to be checked at confirmation if it is contained
	ConfirmedBookings[ticketCode] = booking

	url := "www.booklie.se" + c.Request.URL.Path + "/confirm/" + ticketCode

	store, err := db.GetCompanyByID(dbb, int(timeSlot.CompanyID.Int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get company by ID!",
		})
		return
	}

	// only gets the zeroth element of zone list (because European countries only have single time zones)
	loc, err := time.LoadLocation(tz.GetCountry(store.Country.String).Zones[0].Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find the location for time zone!",
		})
		return
	}
	timeStart := timeSlot.StartTime.Time.In(loc)
	timeStop := timeSlot.EndTime.Time.In(loc)

	confirmation := "Hej " + booking.FirstName.String + "!\n\n" + "Vänligen bekräfta din bokning på " + store.Name.String + " den " + timeStart.Format("2/1") + " klockan " + timeStart.Format("15:04") + "-" + timeStop.Format("15:04") + " i länken nedan:\n\n" + url

	go Send_text(c, booking.PhoneNumber.String, confirmation)

	c.JSON(200, gin.H{
		"message": "Booking was successful.",
	})
}

func generateTicketCode(booking db.Booking) string {
	// last 2 digits of current time + random num [10, 100) + booking name (where space is replaced by underscore)
	return strconv.FormatInt(time.Now().Unix(), 10)[8:] + strconv.Itoa(10+rand.Intn(90)) + strings.ReplaceAll(booking.FirstName.String, " ", "_")
}

func hasAlreadyBooked(phoneNum string, dbb *db.DB, c *gin.Context) bool {
	bookings, err := db.GetBookingsByPhoneNum(dbb, phoneNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong when trying to check for previous bookings with this phone number",
		})
	}

	currentTime := time.Now().UTC()

	if len(bookings) == 0 {
		return false
	} else { // also have to check if current bookings have already taken place
		for i := 0; i < len(bookings); i++ {
			slot, err := db.GetSlot(dbb, int(bookings[i].SlotID.Int64))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Could not get slot by ID in checking if phone number has already booked",
				})
			}
			if currentTime.Before(slot.StartTime.Time) {
				return true
			}
		}
		return false
	}
}

// ConfirmBookAndGetTicket godoc
// @Summary Confirms a booking and adds it to the database if first time. Gets a ticket if it has already been added to database.
// @Produce json
// @Param code path string true "Code"
// @Router /book/confirm/{code} [post]
func ConfirmBookAndGetTicket(c *gin.Context) {
	var bookingExists bool

	code := c.Param("code")
	code = html.EscapeString(code)

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
				"message": "Something went wrong when getting the booking",
			})
			return
		}
	} else {
		bookingExists = true
	}

	if ConfirmedBookings[code].PhoneNumber.String == "" && !bookingExists {
		// booking does not exist
		c.JSON(200, gin.H{
			"message": "This booking does not exist",
		})
		return
	} else if ConfirmedBookings[code].PhoneNumber.String == "" && bookingExists {
		// booking exists and has been added to database
		c.JSON(200, gin.H{
			"message": "Ticket already confirmed!",
			"data":    booking,
		})
		return
	} else {
		// booking exists but has not yet been added to database
		_, err := db.InsertBooking(dbbb, ConfirmedBookings[code])

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not insert the booking",
			})
			return
		}

		var booking db.BookingSlot
		booking.Booking, err = db.GetBooking(dbbb, code)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not get the booking",
			})
			return
		}
		booking.ID = booking.Booking.ID
		booking.Slot, err = db.GetSlot(dbbb, int(booking.Booking.SlotID.Int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not get the slot",
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
	code = html.EscapeString(code)

	dbb := c.MustGet("db").(*db.DB)

	updatedBooking, err := db.UpdateBookingStatusByCode(dbb, code, "cancelled")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not cancel booking.",
		})
		return
	}

	slot, err := db.GetSlot(dbb, int(updatedBooking.SlotID.Int64))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Booking has been cancelled but slot has not correct amount of bookings.",
		})
		return
	}

	slot, err = db.UpdateSlotBooked(dbb, int(slot.ID.Int64), int(slot.Booked.Int64)-1)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Booking has been cancelled but slot has not correct amount of bookings.",
		})
		return
	}

	var booking db.BookingSlot
	booking.Booking = updatedBooking
	booking.ID = updatedBooking.ID
	booking.Slot = slot

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    booking,
	})
}

// GetSlotLoad godoc
// @Summary Gets the load of a slot by returning maxAmount of customers and amount of booked customers as JSON.
// @Produce json
// @Param slotID path string true "slotID"
// @Success 200 "JSON with "maxAmount", "bookingsAmount""
// @Router /slot/{slotID}/load [get]
func GetSlotLoad(c *gin.Context) {
	slotIDStr := c.Param("slotID")
	slotID, err := strconv.Atoi(slotIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not parse slot ID into integer",
		})
		return
	}

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
