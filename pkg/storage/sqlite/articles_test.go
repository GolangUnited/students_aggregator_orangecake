package sqlite

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var TEST_DB_DATA = []core.Article{
	{
		Title:       "1",
		Author:      "1",
		Link:        "1",
		PublishDate: time.Date(time.Now().Year(), time.Now().Month(), 9, 0, 0, 0, 0, time.UTC),
		Description: "1",
	},
	{
		Title:       "2",
		Author:      "2",
		Link:        "2",
		PublishDate: time.Date(time.Now().Year(), time.Now().Month(), 10, 0, 0, 0, 0, time.UTC),
		Description: "2",
	},
	{
		Title:       "3",
		Author:      "3",
		Link:        "3",
		PublishDate: time.Date(time.Now().Year(), time.Now().Month(), 11, 0, 0, 0, 0, time.UTC),
		Description: "3",
	},
	{
		Title:       "4",
		Author:      "4",
		Link:        "4",
		PublishDate: time.Date(time.Now().Year(), time.Now().Month(), 12, 0, 0, 0, 0, time.UTC),
		Description: "4",
	},
}

func TestNewArticleDB(t *testing.T) {
	var (
		lExpectedArticleDB = core.ArticleDB{
			ID: 0,
			Article: core.Article{
				Title:       "1",
				Author:      "1",
				Link:        "1",
				PublishDate: core.NormalizeDate(time.Now()),
				Description: "1",
			},
		}
		lTestArticle = core.Article{
			Title:       "1",
			Author:      "1",
			Link:        "1",
			PublishDate: core.NormalizeDate(time.Now()),
			Description: "1",
		}
	)

	lGotArticle := newArticleDB(&lTestArticle)
	assert.Equal(t, lExpectedArticleDB, *lGotArticle)
}

func TestCutField(t *testing.T) {

	lTestTable := []struct {
		field    string
		expected string
	}{
		{
			field:    "123456789101112",
			expected: "123456789",
		},
		{
			field:    "",
			expected: "",
		},
	}

	for _, lTestCase := range lTestTable {
		lResult := cut(lTestCase.field, 9)
		t.Logf("Calling cut(%s), result is %s\n", lTestCase.field, lResult)

		if lResult != lTestCase.expected {
			t.Errorf("Incorrect result. Expect %s, but got %s", lTestCase.expected, lResult)
		}
	}
}

func TestValidationFields(t *testing.T) {
	testCases := map[string]struct {
		testArticle     core.Article
		expectedArticle core.Article
	}{
		"valid article": {
			testArticle:     TEST_DB_DATA[0],
			expectedArticle: TEST_DB_DATA[0],
		},
		"invalid article": {
			testArticle: core.Article{
				Title:       string(make([]byte, 201)),
				Author:      string(make([]byte, 201)),
				Link:        "",
				PublishDate: time.Time{},
				Description: string(make([]byte, 6001)),
			},
			expectedArticle: core.Article{
				Title:       string(make([]byte, 200)),
				Author:      string(make([]byte, 200)),
				Link:        "",
				PublishDate: time.Time{},
				Description: string(make([]byte, 6000)),
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			validation(&testCase.testArticle)
			assert.Equal(t, testCase.expectedArticle, testCase.testArticle)
		})
	}
}

func TestSqliteStorage_WriteArticle(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)
	var lGotArticle core.ArticleDB

	testCases := map[string]struct {
		testData        core.Article
		expectedArticle core.Article
		expectedError   error
	}{
		"ok test case": {
			testData: core.Article{
				Title:       "1",
				Author:      "1",
				Link:        "1",
				PublishDate: core.NormalizeDate(time.Now()),
				Description: "1",
			},
			expectedArticle: core.Article{
				Title:       "1",
				Author:      "1",
				Link:        "1",
				PublishDate: core.NormalizeDate(time.Now()),
				Description: "1",
			},
			expectedError: nil,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			lDb, lErr := NewSqliteConnection("file::memory:", logger)
			if lErr != nil {
				t.Fail()
			}

			lErr = lDb.WriteArticle(testCase.testData)
			if lErr != nil {
				assert.Error(t, lErr)
			}

			lResult := lDb.db.Where("id = ?", 1).First(&lGotArticle)

			assert.Equal(t, testCase.expectedArticle, lGotArticle.Article)
			assert.Equal(t, testCase.expectedError, lResult.Error)
		})
	}
}

func TestSqliteStorage_WriteArticle_AlreadyExists(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)
	lDb, lErr := NewSqliteConnection("file::memory:", logger)
	if lErr != nil {
		assert.Fail(t, lErr.Error())
	}

	for _, dbDatum := range TEST_DB_DATA {
		err := lDb.db.Create(newArticleDB(&dbDatum)).Error
		if err != nil {
			assert.Fail(t, lErr.Error())
		}
	}

	lErr = lDb.WriteArticle(TEST_DB_DATA[0])
	assert.NoError(t, lErr)

	var gotArticles []core.ArticleDB
	lDb.db.Find(&gotArticles)
	assert.Equal(t, len(TEST_DB_DATA), len(gotArticles))
	for i, dbDatum := range TEST_DB_DATA {
		assert.Equal(t, dbDatum, gotArticles[i].Article)
	}
}

func TestSqliteStorage_WriteArticle_TransactionError(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)
	lDb, lErr := NewSqliteConnection("file::memory:", logger)
	if lErr != nil {
		assert.Fail(t, lErr.Error())
	}

	expectedError := "SQL logic error"

	for _, dbDatum := range TEST_DB_DATA {
		err := lDb.db.Create(newArticleDB(&dbDatum)).Error
		if err != nil {
			assert.Fail(t, lErr.Error())
		}
	}

	gotError := lDb.db.Transaction(func(tx *gorm.DB) error {
		lErr = lDb.WriteArticle(TEST_DB_DATA[0])
		if lErr == nil {
			assert.Fail(t, "error must be not nil")
		}
		return lErr
	})

	if !assert.ErrorContains(t, gotError, expectedError) {
		assert.Fail(t, fmt.Sprintf("mismatch between expected [%s] err and got error [%s]", expectedError, gotError))
	}
}

func TestSqliteStorage_WriteArticles(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)
	var lGotArticles []core.ArticleDB

	testCases := map[string]struct {
		testData         []core.Article
		expectedArticles []core.ArticleDB
		expectedError    error
	}{
		"ok test case": {
			testData: []core.Article{
				{
					Title:       "1",
					Author:      "1",
					Link:        "1",
					PublishDate: core.NormalizeDate(time.Now()),
					Description: "1",
				},
				{
					Title:       "2",
					Author:      "2",
					Link:        "2",
					PublishDate: core.NormalizeDate(time.Now()),
					Description: "2",
				},
			},
			expectedArticles: []core.ArticleDB{
				{
					ID: 1,
					Article: core.Article{
						Title:       "1",
						Author:      "1",
						Link:        "1",
						PublishDate: core.NormalizeDate(time.Now()),
						Description: "1",
					},
				},
				{
					ID: 2,
					Article: core.Article{
						Title:       "2",
						Author:      "2",
						Link:        "2",
						PublishDate: core.NormalizeDate(time.Now()),
						Description: "2",
					},
				},
			},
			expectedError: nil,
		},
	}

	for testName, testCase := range testCases {

		t.Run(testName, func(t *testing.T) {
			lDb, lErr := NewSqliteConnection("file::memory:", logger)
			if lErr != nil {
				t.Fail()
			}

			lErr = lDb.WriteArticles(testCase.testData)
			if lErr != nil {
				assert.Fail(t, lErr.Error())
			}

			lResult := lDb.db.Find(&lGotArticles)

			assert.Equal(t, testCase.expectedArticles, lGotArticles)
			assert.Equal(t, testCase.expectedError, lResult.Error)
		})
	}
}

func TestSqliteStorage_WriteArticles_AlreadyExists(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)
	lDb, lErr := NewSqliteConnection("file::memory:", logger)
	if lErr != nil {
		assert.Fail(t, lErr.Error())
	}

	for _, dbDatum := range TEST_DB_DATA {
		err := lDb.db.Create(newArticleDB(&dbDatum)).Error
		if err != nil {
			assert.Fail(t, err.Error())
		}
	}

	lErr = lDb.WriteArticles(TEST_DB_DATA)
	assert.NoError(t, lErr)

	var gotArticles []core.ArticleDB
	lDb.db.Find(&gotArticles)
	assert.Equal(t, len(TEST_DB_DATA), len(gotArticles))
	for i, dbDatum := range TEST_DB_DATA {
		assert.Equal(t, dbDatum, gotArticles[i].Article)
	}
}

func TestSqliteStorage_WriteArticles_TransactionError(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)
	lDb, lErr := NewSqliteConnection("file::memory:", logger)
	if lErr != nil {
		assert.Fail(t, lErr.Error())
	}

	expectedError := "SQL logic error"

	for _, dbDatum := range TEST_DB_DATA {
		err := lDb.db.Create(newArticleDB(&dbDatum)).Error
		if err != nil {
			assert.Fail(t, err.Error())
		}
	}

	gotError := lDb.db.Transaction(func(tx *gorm.DB) error {
		lErr = lDb.WriteArticles(TEST_DB_DATA)
		if lErr == nil {
			assert.Fail(t, "error must be not nil")
		}
		return lErr
	})

	assert.ErrorContains(t, gotError, expectedError)

	var gotArticles []core.ArticleDB
	lDb.db.Find(&gotArticles)
	assert.Equal(t, len(TEST_DB_DATA), len(gotArticles))
	for i, dbDatum := range TEST_DB_DATA {
		assert.Equal(t, dbDatum, gotArticles[i].Article)
	}
}

func TestSqliteStorage_UpdateArticle(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)
	var lUpArticle = core.Article{
		Title:       "3",
		Author:      "3",
		Link:        "3",
		PublishDate: core.NormalizeDate(time.Now()),
		Description: "3",
	}

	testCases := map[string]struct {
		testData        []core.Article
		expectedArticle core.Article
		articleId       uint
		expectedError   error
	}{
		"ok test case": {
			testData: []core.Article{
				{
					Title:       "1",
					Author:      "1",
					Link:        "1",
					PublishDate: core.NormalizeDate(time.Now()),
					Description: "1",
				},
				{
					Title:       "2",
					Author:      "2",
					Link:        "2",
					PublishDate: core.NormalizeDate(time.Now()),
					Description: "2",
				},
			},
			expectedArticle: core.Article{
				Title:       "3",
				Author:      "3",
				Link:        "3",
				PublishDate: core.NormalizeDate(time.Now()),
				Description: "3",
			},
			articleId:     1,
			expectedError: nil,
		},
	}

	for testName, testCase := range testCases {

		t.Run(testName, func(t *testing.T) {
			lDb, lErr := NewSqliteConnection("file::memory:", logger)
			if lErr != nil {
				t.Fail()
			}

			lErr = lDb.WriteArticles(testCase.testData)
			if lErr != nil {
				assert.Fail(t, lErr.Error())
			}

			lGotErr := lDb.UpdateArticle(testCase.articleId, lUpArticle)
			lGotArticle, _ := lDb.ReadArticleByID(testCase.articleId)

			assert.Equal(t, testCase.expectedArticle, lGotArticle.Article)
			assert.Equal(t, testCase.expectedError, lGotErr)
		})
	}
}

func TestSqliteStorage_UpdateArticle_NotExists(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)
	lDb, lErr := NewSqliteConnection("file::memory:", logger)
	if lErr != nil {
		assert.Fail(t, lErr.Error())
	}

	for _, dbDatum := range TEST_DB_DATA {
		err := lDb.db.Create(newArticleDB(&dbDatum)).Error
		if err != nil {
			assert.Fail(t, err.Error())
		}
	}

	lErr = lDb.UpdateArticle(99, core.Article{})
	assert.NoError(t, lErr)

	var gotArticles []core.ArticleDB
	lDb.db.Find(&gotArticles)
	assert.Equal(t, len(TEST_DB_DATA), len(gotArticles))
	for i, dbDatum := range TEST_DB_DATA {
		assert.Equal(t, dbDatum, gotArticles[i].Article)
	}
}

func TestSqliteStorage_UpdateArticle_Error(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)
	file, err := os.Create("test.db")
	if err != nil {
		assert.Fail(t, err.Error())
	}

	{
		lDb, lErr := NewSqliteConnection(fmt.Sprintf("file:%s", file.Name()), logger)
		if lErr != nil {
			assert.Fail(t, lErr.Error())
		}

		for _, dbDatum := range TEST_DB_DATA {
			err := lDb.db.Create(newArticleDB(&dbDatum)).Error
			if err != nil {
				assert.Fail(t, err.Error())
			}
		}

		sqlDB, _ := lDb.db.DB()
		err := sqlDB.Close()
		if err != nil {
			assert.Fail(t, err.Error())
		}
	}

	{
		roDB, err := NewSqliteConnection(fmt.Sprintf("file:%s?mode=ro", file.Name()), logger)
		assert.NoError(t, err)

		expectedError := fmt.Errorf("attempt to write a readonly database")
		gotError := roDB.UpdateArticle(1, core.Article{
			Title:       "n",
			Author:      "n",
			Link:        "n",
			PublishDate: time.Time{},
			Description: "n",
		})

		if !assert.ErrorContains(t, gotError, expectedError.Error()) {
			assert.Fail(t, fmt.Sprintf("mismatch between expected [%s] err and got error [%s]", expectedError, gotError))
		}

		var gotArticles []core.ArticleDB
		roDB.db.Find(&gotArticles)
		assert.Equal(t, len(TEST_DB_DATA), len(gotArticles))
		for i, dbDatum := range TEST_DB_DATA {
			assert.Equal(t, dbDatum, gotArticles[i].Article)
		}
		sqlDB, _ := roDB.db.DB()
		err = sqlDB.Close()
		if err != nil {
			assert.Fail(t, err.Error())
		}
	}

	_ = file.Close()
	err = os.Remove(file.Name())
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestSqliteStorage_ReadArticleByID(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)

	testCases := map[string]struct {
		testData        []core.Article
		articleId       uint
		expectedArticle core.Article
		expectedError   error
	}{
		"ok test case": {
			testData:        TEST_DB_DATA,
			articleId:       1,
			expectedArticle: TEST_DB_DATA[0],
			expectedError:   nil,
		},

		"empty db test case": {
			testData:        []core.Article{},
			articleId:       0,
			expectedArticle: core.Article{},
			expectedError:   fmt.Errorf("record not found"),
		},

		"wrong id test case": {
			testData:        TEST_DB_DATA,
			articleId:       100,
			expectedArticle: core.Article{},
			expectedError:   fmt.Errorf("record not found"),
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			lDb, err := NewSqliteConnection("file::memory:", logger)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			err = lDb.WriteArticles(testCase.testData)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			gotArticle, gotError := lDb.ReadArticleByID(testCase.articleId)
			assert.Equal(t, testCase.expectedArticle, gotArticle.Article)
			assert.Equal(t, testCase.expectedError, gotError)
		})
	}
}

func TestSqliteStorage_ReadArticlesByDateRange(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)

	testCases := map[string]struct {
		testData         []core.Article
		articlesMinDate  time.Time
		articlesMaxDate  time.Time
		expectedArticles []core.Article
		expectedError    error
	}{
		"ok test case": {
			testData:        TEST_DB_DATA,
			articlesMinDate: time.Date(time.Now().Year(), time.Now().Month(), 10, 0, 0, 0, 0, time.UTC),
			articlesMaxDate: time.Date(time.Now().Year(), time.Now().Month(), 11, 0, 0, 0, 0, time.UTC),
			expectedArticles: []core.Article{
				{
					Title:       "2",
					Author:      "2",
					Link:        "2",
					PublishDate: time.Date(time.Now().Year(), time.Now().Month(), 10, 0, 0, 0, 0, time.UTC),
					Description: "2",
				},
				{
					Title:       "3",
					Author:      "3",
					Link:        "3",
					PublishDate: time.Date(time.Now().Year(), time.Now().Month(), 11, 0, 0, 0, 0, time.UTC),
					Description: "3",
				},
			},
			expectedError: nil,
		},

		"empty db test case": {
			testData:         []core.Article{},
			articlesMinDate:  time.Date(time.Now().Year(), time.Now().Month(), 10, 0, 0, 0, 0, time.UTC),
			articlesMaxDate:  time.Date(time.Now().Year(), time.Now().Month(), 11, 0, 0, 0, 0, time.UTC),
			expectedArticles: []core.Article{},
			expectedError:    nil,
		},

		"wrong date arguments test case": {
			testData:         TEST_DB_DATA,
			articlesMinDate:  time.Date(time.Now().Year(), time.Now().Month(), 111, 0, 0, 0, 0, time.UTC),
			articlesMaxDate:  time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, time.UTC),
			expectedArticles: []core.Article{},
			expectedError:    nil,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			lDb, err := NewSqliteConnection("file::memory:", logger)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			err = lDb.WriteArticles(testCase.testData)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			gotArticles, gotError := lDb.ReadArticlesByDateRange(testCase.articlesMinDate, testCase.articlesMaxDate)
			assert.Equal(t, len(testCase.expectedArticles), len(gotArticles))
			for i, gotArticle := range gotArticles {
				assert.Equal(t, testCase.expectedArticles[i], gotArticle.Article)
			}

			assert.Equal(t, testCase.expectedError, gotError)
		})
	}
}

func TestSqliteStorage_ReadArticlesByDateRangeErrors(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)

	testCases := map[string]struct {
		testData         []core.Article
		expectedArticles []core.ArticleDB
		changeDB         func(t *testing.T, storage *SqliteStorage)
		expectedError    string
	}{
		"no such column": {
			testData:         TEST_DB_DATA,
			expectedArticles: nil,
			expectedError:    "no such column",
			changeDB: func(t *testing.T, storage *SqliteStorage) {
				err := storage.db.Migrator().DropColumn(&core.ArticleDB{}, "publish_date")
				if err != nil {
					assert.Fail(t, err.Error())
				}
			},
		},
		"no such table": {
			testData:         TEST_DB_DATA,
			expectedArticles: nil,
			expectedError:    "no such table",
			changeDB: func(t *testing.T, storage *SqliteStorage) {
				err := storage.db.Migrator().DropTable(&core.ArticleDB{})
				if err != nil {
					assert.Fail(t, err.Error())
				}
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			lDb, err := NewSqliteConnection("file::memory:", logger)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			err = lDb.WriteArticles(testCase.testData)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			testCase.changeDB(t, lDb)

			gotArticles, gotError := lDb.ReadArticlesByDateRange(time.Now(), time.Now())
			assert.Equal(t, testCase.expectedArticles, gotArticles)
			if !assert.ErrorContains(t, gotError, testCase.expectedError) {
				assert.Fail(t, fmt.Sprintf("mismatch between expected [%s] err and got error [%s]", testCase.expectedError, gotError))
			}
		})
	}
}
