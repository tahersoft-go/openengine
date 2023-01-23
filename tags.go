package openengine

import (
	"errors"

	"gitlab.hoitek.fi/openapi/openengine/engine"
)

func (p *openEngine) AddTag(tag engine.Tag) OpenEngine {
	p.Tags = append(p.Tags, tag)
	return p
}

func (p *openEngine) ParseTags(modelsDirPath string, allIgnoredDirs ...[]string) OpenEngine {
	if modelsDirPath == "" {
		p.err = errors.New("path is not set to parse tags from subdirectories")
		return p
	}

	tagNames, err := engine.ExtractDirNames(modelsDirPath)
	if err != nil {
		p.err = err
		return p
	}

	if len(tagNames) == 0 {
		p.err = errors.New("no tags found in path")
		return p
	}
	ignoredDirs := engine.TerIf(len(allIgnoredDirs) == 0, []string{}, allIgnoredDirs[0])
	approvedTags := []string{}
	if len(ignoredDirs) > 0 {
		for _, tagName := range tagNames {
			if !engine.StringInSlice(tagName, &ignoredDirs) {
				approvedTags = append(approvedTags, tagName)
			}
		}
		tagNames = approvedTags
	}

	var tags []engine.Tag
	for _, name := range tagNames {
		tags = append(tags, engine.Tag{
			Name: name,
		})
	}

	p.Tags = tags
	return p
}
