package model

import (
	"time"

	"gorm.io/gorm"
)

type Periods []struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// func (p *Periods) Value() (driver.Value, error) {
// 	jsonBytes, err := json.Marshal(p)
// 	fmt.Println(strings.Repeat("*",50))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jsonBytes, nil
// }

// func (p *Periods) Scan(value interface{}) error {
// 	vToBytes, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("error while converting interface to array of bytes")
// 	}
// 	if err := json.Unmarshal(vToBytes, &p); err != nil {
// 		return err
// 	}
// 	return nil
// }

type Schedule struct {
	gorm.Model
	Periods      []byte `gorm:"column:periods;index;type:jsonb"`
	Each         int    `gorm:"column:each;type:integer"`
	ConsultantID int
	Consultant   Consultant
}
