package search

import (
	"strings"
)

// Persian NLP utilities for search

// FinglishToPersian maps common English/Finglish spellings to Persian characters
var finglishMap = map[string]string{
	// Vowels
	"aa": "آ", "a": "ا", "e": "ه", "ee": "ی", "i": "ای",
	"o": "و", "u": "و", "ou": "و", "oo": "و",

	// Consonants
	"b": "ب", "p": "پ", "t": "ت", "s": "س", "j": "ج",
	"ch": "چ", "h": "ح", "kh": "خ", "d": "د", "z": "ز",
	"zh": "ژ", "r": "ر", "sh": "ش", "gh": "ق", "f": "ف",
	"k": "ک", "g": "گ", "l": "ل", "m": "م", "n": "ن",
	"v": "و", "y": "ی",

	// Special connected letters
	"al": "ال", "ol": "ال",

	// Common Persian word endings
	"eh": "ه", "ye": "یه", "in": "ین", "an": "ان", "ha": "ها",
}

// Common Persian → English transliterations (for reverse matching)
var persianToLatinMap = map[rune]string{
	'آ': "a", 'ا': "a", 'ب': "b", 'پ': "p", 'ت': "t",
	'ث': "s", 'ج': "j", 'چ': "ch", 'ح': "h", 'خ': "kh",
	'د': "d", 'ذ': "z", 'ر': "r", 'ز': "z", 'ژ': "zh",
	'س': "s", 'ش': "sh", 'ص': "s", 'ض': "z", 'ط': "t",
	'ظ': "z", 'ع': "a", 'غ': "gh", 'ف': "f", 'ق': "gh",
	'ک': "k", 'گ': "g", 'ل': "l", 'م': "m", 'ن': "n",
	'و': "v", 'ه': "h", 'ی': "y", 'ئ': "y",
}

// StopWords returns Persian stop words
func StopWords() []string {
	return []string{
		"و", "در", "به", "از", "که", "این", "آن", "با", "برای", "را",
		"است", "یا", "یک", "نه", "هست", "اما", "اگر", "تا", "شد", "شود",
		"می", "های", "ها", "ای", "ما", "شما", "او", "آنها", "خود", "هم",
		"بر", "باید", "مثل", "میتوان", "داشت", "دارد", "کرد", "کند", "گفت",
		"همه", "بعضی", "هیچ", "چند", "همین", "آنچه", "هر", "بین", "زیر",
		"روی", "بالا", "پایین", "داخل", "بیرون", "پیش", "بعد", "قبل",
		"چون", "زیرا", "البته", "حدود", "مثل", "دیگر", "دیگران",
		"اول", "آخر", "فقط", "حتی", "نیز", "همچنین", "خیلی", "بسیار",
		"کم", "زیاد", "چطور", "چگونه", "کجا", "چیست", "چه", "کی",
	}
}

// Common Persian synonyms/alternate spellings for search expansion
var synonymMap = map[string][]string{
	// Artist name variations
	"googoosh":   {"گوگوش"},
	"hayedeh":    {"هایده"},
	"mahasti":    {"مهستی"},
	"dariush":    {"داریوش"},
	"ebi":        {"ابی"},
	"homeyra":    {"حمیرا"},
	"farhad":     {"فرهاد"},
	"golpa":      {"گلپا"},
	"pooran":     {"پوران"},

	// Music genre terms
	"pop":           {"پاپ", "پاپ"},
	"rock":          {"راک", "راک"},
	"jazz":          {"جاز", "جز"},
	"classical":     {"کلاسیک"},
	"traditional":   {"سنتی", "سنتي"},
	"folk":          {"فولک", "محلی", "محلي"},
	"electronic":    {"الکترونیک"},
	"hip-hop":       {"هیپ هاپ", "هیپ-هاپ", "رپ"},
	"rap":           {"رپ"},
	"rnb":           {"آر اند بی"},
	"dance":         {"رقص", "دانسی"},
	"remix":         {"رمیکس", "ریمیکس"},

	// Common song/album terms
	"love":         {"عشق", "عشقي"},
	"song":         {"آهنگ", "ترانه", "تصنیف"},
	"music":        {"موسیقی", "موزیک"},
	"album":        {"آلبوم"},
	"singer":       {"خواننده", "آوازخوان"},
	"lyrics":       {"متن", "شعر", "ترانه"},
	"dastgah":      {"دستگاه", "دستگاه"},
	"radif":        {"ردیف", "رديف"},
	"avaz":         {"آواز", "آواز"},
	"tasnif":       {"تصنیف", "تصنيف"},
	"marg":         {"مرگ", "موت"},
	"gham":         {"غم"},
	"shadi":        {"شادی"},
	"bahar":        {"بهار"},
	"ghanari":      {"قناری"},
	"barkis":       {"برکش"},
	"torki":        {"ترکی", "تركي"},
	"kordi":        {"کردی", "كردي"},
	"farsi":        {"فارسی", "پارسی"},
	"irani":        {"ایرانی", "ايراني"},
	"persian":      {"پرشین"},

	// Numbers
	"0":  {"۰"},
	"1":  {"۱"},
	"2":  {"۲"},
	"3":  {"۳"},
	"4":  {"۴"},
	"5":  {"۵"},
	"6":  {"۶"},
	"7":  {"۷"},
	"8":  {"۸"},
	"9":  {"۹"},
	"10": {"۱۰"},
}

// NormalizePersian normalizes Persian/Arabic characters
func NormalizePersian(text string) string {
	replacer := strings.NewReplacer(
		"ي", "ی", // Arabic ye → Persian ye
		"ك", "ک", // Arabic ke → Persian ke
		"ة", "ه", // Ta marbuta → he
		"٠", "۰", "١", "۱", "٢", "۲", "٣", "۳", "٤", "۴",
		"٥", "۵", "٦", "۶", "٧", "۷", "٨", "۸", "٩", "۹", // Arabic → Persian digits
		"\u064B", "", // Fathatan
		"\u064C", "", // Dammatan
		"\u064D", "", // Kasratan
		"\u064E", "", // Fatha
		"\u064F", "", // Damma
		"\u0650", "", // Kasra
		"\u0651", "", // Shadda
		"\u0652", "", // Sukun
		"\u0653", "", // Maddah
	)
	return replacer.Replace(text)
}

// TransliterateEnglishToPersian attempts to convert Finglish (English spelling of Persian) to Persian script
func TransliterateEnglishToPersian(text string) string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	var result []string
	for _, word := range words {
		translated := transliterateWord(word)
		if translated != word {
			result = append(result, translated)
		}
	}
	return strings.Join(result, " ")
}

func transliterateWord(word string) string {
	// Exact match check first
	if p, ok := finglishMap[word]; ok {
		return p
	}

	// Common Persian name patterns
	knownNames := map[string]string{
		"googoosh": "گوگوش", "hayedeh": "هایده",
		"mahasti": "مهستی", "dariush": "داریوش",
		"farhad": "فرهاد", "homeyra": "حمیرا",
		"pooran": "پوران", "marjan": "مرجان",
		"shahram": "شهرام", "shohreh": "شهره",
		"leila": "لیلا", "golshifteh": "گلشیفته",
		"mohsen": "محسن", "ebrahim": "ابراهیم",
		"hamid": "حمید", "shahin": "شاهین",
		"mohammad": "محمد", "ali": "علی",
		"hosein": "حسین", "reza": "رضا",
		"ahmad": "احمد", "mehdi": "مهدی",
		"hadi": "هادی", "saeed": "سعید",
		"sina": "سینا", "arshia": "آرشیا",
		"navid": "نوید", "pedram": "پدرام",
		"saman": "سامان", "bahram": "بهرام",
		"behnam": "بهنام", "kambiz": "کامبیز",
		"parviz": "پرویز", "majid": "محمود",
		"ghorbani": "قربانی",
		"alizadeh": "علیزاده", "ahmadi": "احمدی",
		"mohammadi": "محمدی", "hosseini": "حسینی",
		"karimi": "کریمی", "moradi": "مرادی",
		"hejazi": "حجازی", "nazeri": "ناظری",
		"shajarian": "شجریان", "sharif": "شریف",
		"ebtehaj": "ابتهاج", "fardin": "فردین",
		"laleh": "لاله", "maryam": "مریم",
		"narges": "نرگس", "parastoo": "پرستو",
		"yazdani": "یزدانی", "mostafa": "مصطفی",
	}
	if p, ok := knownNames[word]; ok {
		return p
	}

	return word
}

// TransliteratePersianToEnglish converts Persian text to approximate English spelling
func TransliteratePersianToEnglish(text string) string {
	var result strings.Builder
	for _, r := range text {
		if lat, ok := persianToLatinMap[r]; ok {
			result.WriteString(lat)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ExpandQuery generates search expansions: original + transliterated forms
func ExpandQuery(query string) []string {
	seen := make(map[string]bool)
	var expansions []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" && !seen[s] {
			seen[s] = true
			expansions = append(expansions, s)
		}
	}

	// Original query
	add(query)

	// Normalized Persian
	add(NormalizePersian(query))

	// English → Persian transliteration
	persian := TransliterateEnglishToPersian(query)
	if persian != query {
		add(persian)
	}

	// Persian → English
	english := TransliteratePersianToEnglish(query)
	if english != query {
		add(english)
	}

	return expansions
}

// GetSynonyms returns synonyms for a term
func GetSynonyms(term string) []string {
	term = strings.ToLower(term)
	if syns, ok := synonymMap[term]; ok {
		return syns
	}
	return nil
}

// IsStopWord checks if a word is a Persian stop word
func IsStopWord(word string) bool {
	word = strings.TrimSpace(word)
	if word == "" {
		return true
	}
	for _, sw := range StopWords() {
		if word == sw {
			return true
		}
	}
	return false
}

// FilterStopWords removes stop words from a query
func FilterStopWords(query string) string {
	words := strings.Fields(query)
	var filtered []string
	for _, w := range words {
		if !IsStopWord(w) {
			filtered = append(filtered, w)
		}
	}
	return strings.Join(filtered, " ")
}
