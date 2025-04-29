package entity

type MMSEntity struct {
	Pair      string `gorm:"primaryKey"`
	Timestamp int64  `gorm:"primaryKey"`
	MMS20     float64
	MMS50     float64
	MMS200    float64
}

func (MMSEntity) TableName() string {
	return "medias_moveis_simples"
}
