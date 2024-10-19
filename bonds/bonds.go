package bonds

import (
	"encoding/json"
	"os"
	"time"
)

/* Struct for describe one bonds */
type BondsData struct {
	Name              string        `json:"name"`
	CouponCount       int           `json:"couponCount"`
	CouponPeriod      time.Duration `json:"couponPeriod"`
	CouponNearPayDate time.Time     `json:"nearPayDate"`
	PayDates          []time.Time   `json:"-"`
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
		obj.CalcCouponDates()
	}

	return nil
}

/* Append new bonds, also call CalcCouponDates before append */
func (self *Bonds) Append(obj *BondsData) {
	obj.CalcCouponDates()
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

func BondsDataNew() *BondsData {
	obj := new(BondsData)
	obj.PayDates = make([]time.Time, 0)
	return obj
}

func (self *BondsData) CalcCouponDates() {
	self.PayDates = make([]time.Time, 0)
	tmp := self.CouponNearPayDate

	for i := 0; i < self.CouponCount; i++ {
		self.PayDates = append(self.PayDates, tmp)
		tmp = tmp.Add(self.CouponPeriod)
	}
}

/** Create a pay day date with only valid year, month and day */
func CouponPayDay(year, month, day int) time.Time {
	obj := time.Date(year, time.Month(month), day, 0, 0, 0, 0, DefaultLocation)
	return obj
}

/* Create a payment duration by 2 pay dates */
func CouponPeriodCreate(nearest, next time.Time) time.Duration {
	return next.Sub(nearest)
}
