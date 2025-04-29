package mapper

import (
	"testing"

	"github.com/Junkes887/MMS/internal/database/entity"
	"github.com/Junkes887/MMS/internal/service/dto"
	"github.com/stretchr/testify/assert"
)

func TestMMSEntityToMSSResponse(t *testing.T) {
	tests := []struct {
		name     string
		entities []entity.MMSEntity
		mmsRange int
		expected []dto.MMSResponse
	}{
		{
			name: "Range 20",
			entities: []entity.MMSEntity{
				{Timestamp: 1000, MMS20: 10.5, MMS50: 20.5, MMS200: 30.5},
				{Timestamp: 2000, MMS20: 11.5, MMS50: 21.5, MMS200: 31.5},
			},
			mmsRange: 20,
			expected: []dto.MMSResponse{
				{Timestamp: 1000, MMS: 10.5},
				{Timestamp: 2000, MMS: 11.5},
			},
		},
		{
			name: "Range 50",
			entities: []entity.MMSEntity{
				{Timestamp: 1000, MMS20: 10.5, MMS50: 20.5, MMS200: 30.5},
				{Timestamp: 2000, MMS20: 11.5, MMS50: 21.5, MMS200: 31.5},
			},
			mmsRange: 50,
			expected: []dto.MMSResponse{
				{Timestamp: 1000, MMS: 20.5},
				{Timestamp: 2000, MMS: 21.5},
			},
		},
		{
			name: "Range 200",
			entities: []entity.MMSEntity{
				{Timestamp: 1000, MMS20: 10.5, MMS50: 20.5, MMS200: 30.5},
				{Timestamp: 2000, MMS20: 11.5, MMS50: 21.5, MMS200: 31.5},
			},
			mmsRange: 200,
			expected: []dto.MMSResponse{
				{Timestamp: 1000, MMS: 30.5},
				{Timestamp: 2000, MMS: 31.5},
			},
		},
		{
			name:     "Empty list",
			entities: []entity.MMSEntity{},
			mmsRange: 20,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MMSEntityToMSSResponse(tt.entities, tt.mmsRange)
			assert.Equal(t, tt.expected, result)
		})
	}
}
