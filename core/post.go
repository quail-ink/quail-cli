package core

import (
	"fmt"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type QuailPostFrontMatter struct {
	Slug          string     `yaml:"slug"`
	CoverImageUrl string     `yaml:"cover_image_url"`
	Title         string     `yaml:"title"`
	Summary       string     `yaml:"summary"`
	Theme         string     `yaml:"theme"`
	Tags          string     `yaml:"tags"`
	Datetime      *time.Time `yaml:"datetime"`
}

var datetimeFormats = []string{
	// add more datetime formats here
	time.RFC1123,
	time.RFC3339, // 2006-01-02T15:04:05Z07:00
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"02 Jan 2006",
	"02 Jan 2006 15:04",
	"02 Jan 2006 15:04:05",
	"January 2, 2006",
	"January 2, 2006 15:04",
	"January 2, 2006 15:04:05",
}

func (q *QuailPostFrontMatter) LoadFromYAML(data string, convertMap map[string]string) error {
	var frontMatterMap map[string]any
	err := yaml.Unmarshal([]byte(data), &frontMatterMap)
	if err != nil {
		return fmt.Errorf("could not parse frontmatter: %w", err)
	}

	// convert the frontMatterMap name to standard name by using convertMap
	for key, value := range convertMap {
		if val, ok := frontMatterMap[value]; ok {
			frontMatterMap[key] = val
			delete(frontMatterMap, value)
		}
	}

	if err := q.ConvertMapToFrontMatter(frontMatterMap); err != nil {
		return fmt.Errorf("could not convert map to front matter: %w", err)
	}

	return nil
}

func (q *QuailPostFrontMatter) ConvertMapToFrontMatter(frontMatterMap map[string]any) error {
	// handle the datetime field and tags field
	if rawDatetime, ok := frontMatterMap["datetime"]; ok {
		if datetimeStr, ok := rawDatetime.(string); ok {
			parsedTime, err := parseDateTime(datetimeStr)
			if err != nil {
				return err
			}
			frontMatterMap["datetime"] = parsedTime
		}
	}
	if rawTags, ok := frontMatterMap["tags"]; ok {
		frontMatterMap["tags"] = parseTags(rawTags)
	}

	// Marshal the map back into YAML
	yamlData, err := yaml.Marshal(frontMatterMap)
	if err != nil {
		return fmt.Errorf("could not marshal map to YAML: %w", err)
	}

	err = yaml.Unmarshal(yamlData, q)
	if err != nil {
		return fmt.Errorf("could not unmarshal YAML to struct: %w", err)
	}

	return nil
}

func parseDateTime(datetimeStr string) (*time.Time, error) {
	for _, layout := range datetimeFormats {
		parsedTime, err := time.Parse(layout, datetimeStr)
		if err == nil {
			return &parsedTime, nil
		}
	}
	return nil, fmt.Errorf("could not parse datetime: %s", datetimeStr)
}

func parseTags(rawTags any) string {
	tags := []string{}
	switch v := rawTags.(type) {
	case string:
		tags = strings.Split(v, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	case []any:
		for _, tag := range v {
			if tagStr, ok := tag.(string); ok {
				tags = append(tags, tagStr)
			}
		}
	}
	return strings.Join(tags, ",")
}
