package bonds

import (
	"encoding/json"
	"os"
	"time"
)

/* Struct for describe one bonds */
type BondsData struct {
	Name              string        `json:"name"` // Bond name
	CouponCount       int           `json:"couponCount"` // Count of remaining coupon payments
	CouponPeriod      int           `json:"couponPeriod"` // Period between coupon payments (Calc as NextDate - NearDate)
	CouponNearPayDate time.Time     `json:"nearPayDate"` // Near date of payment
	PayDates          []time.Time   `json:"-"` // Calculated dates of coupon payments
}

/* Struct for store multiply bonds */
type Bonds struct {
	Bonds []*BondsData
}

var (
	DefaultLocation, _ = time.LoadLocation("Asia/Yekaterinburg")
)

func BondsNew() *Bonds {
	obj := new(Bonds)
	obj.Bonds = make([]*BondsData, 0)

	return obj
}

/* Save all appended bonds into file as json */
func (self *Bonds) SaveToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(self.Bonds)

	if err != nil {
		return err
	}

	return nil
}

/*
Save all appended bonds into file as json
Overwrite current Bonds array
*/
func (self *Bonds) LoadFromFile(filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer file.Close()
	self.Bonds = make([]*BondsData, 0)

	encoder := json.NewDecoder(file)
	err = encoder.Decode(&self.Bonds)

	if err != nil {
		return err
	}

	for _, obj := range self.Bonds {
		obj.CalcAll()
	}

	return nil
}

/* Append new bonds, also call CalcCouponDates before append */
func (self *Bonds) Append(obj *BondsData) {
	obj.CalcAll()
	self.Bonds = append(self.Bonds, obj)
}

/* Count a payments by given year and month */
func (self *Bonds) PayCountByYearMonth(year, month int) int {
	var result int
	validMonth := time.Month(month)

	for _, obj := range self.Bonds {
		for _, date := range obj.PayDates {
			if date.Year() == year && date.Month() == validMonth {
				result++
			}
		}
	}

	return result
}

/*
Check is any bond has been expired:
    - Near pay date is in the past and coupon count is 1
*/
// func (self *Bonds) ExpiredList() []int {
//
// }

func BondsDataNew() *BondsData {
	obj := new(BondsData)
	obj.PayDates = make([]time.Time, 0)
	return obj
}

/** Create a pay day date with only year, month and day */
func CouponPayDay(year, month, day int) time.Time {
	obj := time.Date(year, time.Month(month), day, 0, 0, 0, 0, DefaultLocation)
	return obj
}

/* Create a payment duration with 2 next pay dates (near and next one) */
func CouponPeriodCreate(nearest, next time.Time) int {
	return int(next.Sub(nearest).Hours()) / 24
}

/*
Caclulate all related data:
    - Next coupon pay dates
    - Remove dates if they are in the past (date < time.Now)
*/
func (self *BondsData) CalcAll() {
    self.calcCouponDates()
    self.removePastDates()
}

/* Check is near pay date is in the past and set near as next, if next available */
func (self *BondsData) removePastDates() {
    timeNow := time.Now()

    for id, val := range self.PayDates {
        if val.After(timeNow) {
            self.PayDates = self.PayDates[id:]
            self.CouponNearPayDate = self.PayDates[0]
            self.CouponCount -= id
            break
        }
    }
}

/* Calculate all next pay dates */
func (self *BondsData) calcCouponDates() {
	self.PayDates = make([]time.Time, 0)
	tmp := self.CouponNearPayDate

	for i := 0; i < self.CouponCount; i++ {
		self.PayDates = append(self.PayDates, tmp)
		tmp = tmp.AddDate(0, 0, self.CouponPeriod)
	}
}

