package transifex_app_go_client

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type SearchStringsResponse struct {
	Data []string `json:"data"`
	Meta struct {
		Total    int    `json:"total"`
		Relation string `json:"relation"`
	} `json:"meta"`
}

type SlotSegment struct {
	ID                string `json:"id"`
	SourceID          int    `json:"source_id"`
	Status            int    `json:"status"`
	SourceString      string `json:"source_string"`
	TranslationString string `json:"translation_string"`
	SourceLanguage    struct {
		Code      string `json:"code"`
		Name      string `json:"name"`
		CharBased bool   `json:"char_based"`
		IsRtl     bool   `json:"is_rtl"`
	} `json:"source_language"`
	TargetLanguage struct {
		Code      string `json:"code"`
		Name      string `json:"name"`
		CharBased bool   `json:"char_based"`
		IsRtl     bool   `json:"is_rtl"`
	} `json:"target_language"`
	Resource struct {
		ID       int    `json:"id"`
		Slug     string `json:"slug"`
		Name     string `json:"name"`
		I18NType string `json:"i18n_type"`
	} `json:"resource"`
	Project struct {
		ID   int    `json:"id"`
		Slug string `json:"slug"`
		Name string `json:"name"`
	} `json:"project"`
	ChecksStatus int `json:"checks_status"`
}

type SlotSegments struct {
	Data []SlotSegment `json:"data"`
}

func (t *TransifexAppClient) GetResourseWebIDByTranslation(strToSearch, srcLang, dstLang string) (int, error) {
	var ssResponse SearchStringsResponse

	processedStrToSearch := url.QueryEscape(removeServiceSymbols(strToSearch))

	searchURL := strings.Join([]string{
		"https://app.transifex.com/_/search/ajax/carpentries-i18n/ids/?",
		// fmt.Sprintf("translation_text=\"%s\"", processedStrToSearch),
		fmt.Sprintf("translation_text=%s", processedStrToSearch),
		fmt.Sprintf("source_language=%s", srcLang),
		fmt.Sprintf("target_language=%s", dstLang),
		"ordering=relevance__desc",
	}, "&")

	// searchURL = url.QueryEscape(searchURL)

	_, _, errs := t.bot.Requester("GET", searchURL).
		Client.
		EndStruct(&ssResponse)
	if errs != nil {

		t.l.WithFields(logrus.Fields{
			"string-to-search": strToSearch,
			"searchURL":        searchURL,
		}).Error(errs)

		return -1, errors.New("errors")
	}

	var slSegms SlotSegments

	for _, translationApplicant := range ssResponse.Data {

		data := fmt.Sprintf("{\"ids\":[\"%s\"]}", translationApplicant)

		_, _, errs := t.bot.Requester("POST",
			"https://app.transifex.com/_/search/ajax/carpentries-i18n/slot_segments/",
		).
			Client.
			Send(data).
			EndStruct(&slSegms)
		if errs != nil {

			t.l.WithFields(logrus.Fields{
				"function":         "GetResourseWebIDByTranslation",
				"string-to-search": strToSearch,
				"lang-from":        srcLang,
				"lang-to":          dstLang,
			}).Error(errs)

			return -1, errors.New("errors")
		}

		for _, ss := range slSegms.Data {
			if strToSearch == ss.TranslationString {
				return ss.SourceID, nil
			}
		}
	}

	t.l.WithFields(logrus.Fields{
		"function":         "GetResourseWebIDByTranslation",
		"string-to-search": strToSearch,
		"lang-from":        srcLang,
		"lang-to":          dstLang,
	}).Debug("the translation is not found")

	return -1, fmt.Errorf("the translation is not found")
}

func removeServiceSymbols(s string) string {
	return strings.TrimLeft(strings.Replace(s, "\n", " ", -1), "#* ")
}
