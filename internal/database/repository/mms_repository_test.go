package repository

import (
	"testing"
	"time"

	"github.com/Junkes887/MMS/internal/database"
	"github.com/Junkes887/MMS/internal/database/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *database.DBConnection {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to sqlite memory: %v", err)
	}

	err = db.AutoMigrate(&entity.MMSEntity{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return &database.DBConnection{DB: db}
}

func TestFindMSS_Success(t *testing.T) {
	dbConn := setupTestDB(t)
	logger := zap.NewNop()
	repo := NewRepository(dbConn, logger)

	from := time.Now().AddDate(0, 0, -3).Unix()
	to := time.Now().Unix()

	testData := []entity.MMSEntity{
		{Pair: "BRLBTC", Timestamp: from, MMS20: 50000.0},
		{Pair: "BRLBTC", Timestamp: from + 86400, MMS20: 51000.0},
		{Pair: "BRLBTC", Timestamp: from + 2*86400, MMS20: 52000.0},
	}

	for _, data := range testData {
		err := repo.SaveMSS(data)
		assert.NoError(t, err)
	}

	result, err := repo.FindMSS("BRLBTC", from, to)

	assert.NoError(t, err)
	assert.Len(t, result, 3)

	assert.Equal(t, 50000.0, result[0].MMS20)
	assert.Equal(t, 51000.0, result[1].MMS20)
	assert.Equal(t, 52000.0, result[2].MMS20)
}

func TestFindMSS_NoResults(t *testing.T) {
	dbConn := setupTestDB(t)
	logger := zap.NewNop()
	repo := NewRepository(dbConn, logger)

	from := time.Now().AddDate(0, 0, -3).Unix()
	to := time.Now().Unix()

	result, err := repo.FindMSS("BRLBTC", from, to)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestFindMSS_Error(t *testing.T) {
	dbConn := setupTestDB(t)
	logger := zap.NewNop()
	repo := NewRepository(dbConn, logger)

	from := time.Now().AddDate(0, 0, -3).Unix()
	to := time.Now().Unix()

	err := dbConn.DB.Migrator().DropTable(&entity.MMSEntity{})
	assert.NoError(t, err)

	result, err := repo.FindMSS("BRLBTC", from, to)

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestSaveMSS_Success(t *testing.T) {
	dbConn := setupTestDB(t)
	logger := zap.NewNop()
	repo := NewRepository(dbConn, logger)

	mmsEntity := entity.MMSEntity{
		Pair:      "BRLBTC",
		Timestamp: time.Now().Unix(),
		MMS20:     50000.0,
	}

	err := repo.SaveMSS(mmsEntity)

	assert.NoError(t, err)

	var saved entity.MMSEntity
	result := dbConn.DB.First(&saved, "pair = ? AND timestamp = ?", mmsEntity.Pair, mmsEntity.Timestamp)
	assert.NoError(t, result.Error)

	assert.Equal(t, mmsEntity.Pair, saved.Pair)
	assert.Equal(t, mmsEntity.Timestamp, saved.Timestamp)
	assert.Equal(t, mmsEntity.MMS20, saved.MMS20)
}

func TestSaveMSS_Error(t *testing.T) {
	dbConn := setupTestDB(t)
	logger := zap.NewNop()
	repo := NewRepository(dbConn, logger)

	err := dbConn.DB.Migrator().DropTable(&entity.MMSEntity{})
	assert.NoError(t, err)

	mmsEntity := entity.MMSEntity{
		Pair:      "BRLBTC",
		Timestamp: time.Now().Unix(),
		MMS20:     50000.0,
	}

	err = repo.SaveMSS(mmsEntity)

	assert.Error(t, err)
}

func TestBuscarDiasFaltantes_Success(t *testing.T) {
	dbConn := setupTestDB(t)
	logger := zap.NewNop()
	repo := NewRepository(dbConn, logger)

	from := time.Now().AddDate(0, 0, -3).Unix()
	to := time.Now().Unix()

	testData := []entity.MMSEntity{
		{Pair: "BRLBTC", Timestamp: from + 0*86400, MMS20: 50000.0},
		{Pair: "BRLBTC", Timestamp: from + 1*86400, MMS20: 51000.0},
		{Pair: "BRLBTC", Timestamp: from + 2*86400, MMS20: 52000.0},
	}

	for _, data := range testData {
		err := repo.SaveMSS(data)
		assert.NoError(t, err)
	}

	existentes := repo.BuscarDiasFaltantes("BRLBTC", from, to)

	assert.Len(t, existentes, 3)
	assert.Contains(t, existentes, from)
	assert.Contains(t, existentes, from+86400)
	assert.Contains(t, existentes, from+2*86400)
}

func TestBuscarDiasFaltantes_Error(t *testing.T) {
	dbConn := setupTestDB(t)
	logger := zap.NewNop()
	repo := NewRepository(dbConn, logger)

	from := time.Now().AddDate(0, 0, -3).Unix()
	to := time.Now().Unix()

	err := dbConn.DB.Migrator().DropTable(&entity.MMSEntity{})
	assert.NoError(t, err)

	existentes := repo.BuscarDiasFaltantes("BRLBTC", from, to)

	assert.Nil(t, existentes)
}
