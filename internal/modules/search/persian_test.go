package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePersian_ArabicYeToPersian(t *testing.T) {
	assert.Equal(t, "ای", NormalizePersian("اي"))
}

func TestNormalizePersian_ArabicKeToPersian(t *testing.T) {
	assert.Equal(t, "ک", NormalizePersian("ك"))
}

func TestNormalizePersian_ArabicDigits(t *testing.T) {
	assert.Equal(t, "۱۲۳", NormalizePersian("١٢٣"))
}

func TestNormalizePersian_TaMarbuta(t *testing.T) {
	assert.Equal(t, "ه", NormalizePersian("ة"))
}

func TestNormalizePersian_RemoveDiacritics(t *testing.T) {
	assert.Equal(t, "ب", NormalizePersian("بِ"))
}

func TestExpandQuery_PreservesOriginal(t *testing.T) {
	expansions := ExpandQuery("test")
	assert.Contains(t, expansions, "test")
}

func TestExpandQuery_Deduplicates(t *testing.T) {
	expansions := ExpandQuery("test")
	seen := make(map[string]bool)
	for _, e := range expansions {
		if seen[e] {
			t.Errorf("duplicate expansion: %s", e)
		}
		seen[e] = true
	}
}

func TestExpandQuery_NonEmpty(t *testing.T) {
	expansions := ExpandQuery("گوگوش")
	assert.NotEmpty(t, expansions)
}

func TestIsStopWord_Empty(t *testing.T) {
	assert.True(t, IsStopWord(""))
}

func TestIsStopWord_PersianStopWord(t *testing.T) {
	assert.True(t, IsStopWord("و"))
	assert.True(t, IsStopWord("در"))
	assert.True(t, IsStopWord("به"))
}

func TestIsStopWord_NonStopWord(t *testing.T) {
	assert.False(t, IsStopWord("موسیقی"))
	assert.False(t, IsStopWord("آهنگ"))
}

func TestFilterStopWords_RemovesStopWords(t *testing.T) {
	result := FilterStopWords("آهنگ در موسیقی")
	assert.Equal(t, "آهنگ موسیقی", result)
}

func TestFilterStopWords_AllStopWords(t *testing.T) {
	result := FilterStopWords("و در به")
	assert.Equal(t, "", result)
}

func TestTransliterateEnglishToPersian_KnownName(t *testing.T) {
	result := TransliterateEnglishToPersian("googoosh")
	assert.Contains(t, result, "گوگوش")
}

func TestTransliteratePersianToEnglish_Basic(t *testing.T) {
	result := TransliteratePersianToEnglish("گوگوش")
	assert.Contains(t, result, "g")
}

func TestGetSynonyms_Pop(t *testing.T) {
	syns := GetSynonyms("pop")
	assert.Contains(t, syns, "پاپ")
}

func TestGetSynonyms_Unknown(t *testing.T) {
	syns := GetSynonyms("unknownword")
	assert.Nil(t, syns)
}

func TestStopWords_NotEmpty(t *testing.T) {
	words := StopWords()
	assert.NotEmpty(t, words)
	assert.Contains(t, words, "و")
	assert.Contains(t, words, "در")
}
