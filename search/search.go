package search

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type SearchResults struct {
	ready   bool
	Query   string
	Results []Result
}

type Result struct {
	Name        string
	Description string
	URL         string
}

func (sr *SearchResults) UnmarshalJSON(bs []byte) error {
	array := []interface{}{}
	if err := json.Unmarshal(bs, &array); err != nil {
		return err
	}
	sr.Query = array[0].(string)
	for i := range array[1].([]interface{}) {
		sr.Results = append(sr.Results, Result{
			array[1].([]interface{})[i].(string),
			array[2].([]interface{})[i].(string),
			array[3].([]interface{})[i].(string),
		})
	}
	return nil
}

func WikipediaAPI(request string) (answer []string) {
	s := make([]string, 3)

	if response, err := http.Get(request); err != nil {
		s[0] = "Wikipedia не отвечает"
	} else {
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			s[0] = "нет контента"
		}

		sr := &SearchResults{}

		if err = json.Unmarshal([]byte(contents), sr); err != nil {
			s[0] = "Что-то пошло не так, измените запрос"
		}

		if !sr.ready {
			s[0] = "Что-то пошло не так, измените запрос"
		}

		for i := range sr.Results {
			s[i] = sr.Results[i].URL
		}
	}

	return s
}
