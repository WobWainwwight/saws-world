package db_test

import (
	"os"
	"slices"
	"testing"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wobwainwwight/sa-photos/db"
)

func TestDB(t *testing.T) {

	t.Run("should setup saws.db if no arg", func(t *testing.T) {
		require.NoFileExists(t, "saws.db")
		_, err := db.NewImageTable("")
		require.NoError(t, err)

		stat, err := os.Stat("saws.db")
		require.NoError(t, err)
		assert.False(t, stat.IsDir())

		require.NoError(t, os.Remove("saws.db"))
	})

	t.Run("should add image row", func(t *testing.T) {
		table := newTestTable(t)
		defer table.Close()

		err := table.Save(db.Image{
			ID:         "image123",
			MimeType:   "jpg",
			UploadedAt: time.Now(),
			CreatedAt:  time.Now(),
		})
		require.NoError(t, err)

		image, err := table.GetByID("image123")
		require.NoError(t, err)

		assert.Equal(t, image.ID, "image123")
	})

	t.Run("should get list of image rows", func(t *testing.T) {
		table := newTestTable(t)
		defer table.Close()

		err := table.Save(db.Image{
			ID:         "image123",
			MimeType:   "jpg",
			UploadedAt: time.Now(),
			CreatedAt:  time.Now(),
		})
		require.NoError(t, err)

		err = table.Save(db.Image{
			ID:         "image456",
			MimeType:   "jpg",
			UploadedAt: time.Now(),
			CreatedAt:  time.Now(),
		})
		require.NoError(t, err)

		rows, err := table.Get()
		require.NoError(t, err)

		assertContainsRowWithID(t, rows, "image123")
		assertContainsRowWithID(t, rows, "image456")
	})

	t.Run("should get rows sorted by created order", func(t *testing.T) {
		table := newTestTable(t)
		defer table.Close()

		err := table.Save(db.Image{
			ID:        "image123",
			MimeType:  "jpg",
			CreatedAt: time.Now().Add(-time.Hour),
		})
		require.NoError(t, err)

		err = table.Save(db.Image{
			ID:        "image456",
			MimeType:  "jpg",
			CreatedAt: time.Now(),
		})
		require.NoError(t, err)

		rows, err := table.Get(db.GetOpts{OrderDirection: db.DESC})
		require.NoError(t, err)

		require.Len(t, rows, 2)
		assert.Equal(t, "image456", rows[0].ID)
		assert.Equal(t, "image123", rows[1].ID)

		rows, err = table.Get()
		require.NoError(t, err)

		require.Len(t, rows, 2)
		assert.Equal(t, "image123", rows[0].ID)
		assert.Equal(t, "image456", rows[1].ID)
	})

	t.Run("should get all localities", func(t *testing.T) {
		table := newTestTable(t)
		defer table.Close()

		countryLocalities := map[string][]string{
			"United States": {"New York", "Washington DC", "Los Angeles"},
			"Wales":         {"Cardiff", "Swansea", "Newport"},
			"Argentina":     {"San Carlos de Bariloche", "Mendoza"},
		}

		givenCountriesAndLocalities(t, table, countryLocalities)

		l, err := table.GetLocalities()
		require.NoError(t, err)

		for k, v := range countryLocalities {
			assertContainsLocality(t, l, k, v)
		}

	})
}

func givenCountriesAndLocalities(t *testing.T, it testTable, countryLocalities map[string][]string) {
	for country, localities := range countryLocalities {
		for _, l := range localities {
			err := it.Save(givenImageInLocale(country, l))
			require.NoError(t, err)
		}
	}
}

func givenImageInLocale(country, locality string) db.Image {
	return db.Image{
		ID:       gonanoid.Must(),
		Locality: locality,
		Country:  country,
	}
}

func assertContainsLocality(t *testing.T, ll []db.Locality, country string, localities []string) {
	found := false
	for _, l := range ll {
		if l.Country == country {
			found = true

			for _, lc := range localities {
				assert.Contains(t, l.Localities, lc)
			}
		}

	}
	assert.Truef(t, found, "country %s not found", country)

}

type testTable struct {
	t *testing.T
	*db.ImageTable
}

func newTestTable(t *testing.T) testTable {
	require.NoFileExists(t, "saws.db")
	table, err := db.NewImageTable("")
	require.NoError(t, err)
	return testTable{t, table}
}

func (t *testTable) Close() error {
	_ = t.ImageTable.Close()
	return os.Remove("saws.db")
}

func assertContainsRowWithID(t *testing.T, imgs []db.Image, id string) {
	assert.Truef(t, slices.ContainsFunc(imgs, func(img db.Image) bool {
		return img.ID == id
	}), "does not contain row with id %s", id)
}
