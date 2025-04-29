package mapper

import (
	"github.com/Junkes887/MMS/internal/database/entity"
	"github.com/Junkes887/MMS/internal/service/dto"
)

func MMSEntityToMSSResponse(list []entity.MMSEntity, mmsRange int) []dto.MMSResponse {
	var dtos []dto.MMSResponse

	for _, v := range list {
		dto := dto.MMSResponse{
			Timestamp: v.Timestamp,
		}

		switch mmsRange {
		case 20:
			dto.MMS = v.MMS20
		case 50:
			dto.MMS = v.MMS50
		case 200:
			dto.MMS = v.MMS200
		}

		dtos = append(dtos, dto)
	}

	return dtos
}
