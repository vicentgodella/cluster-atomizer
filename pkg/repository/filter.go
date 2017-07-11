package repository

import "fmt"

type FilterPatternBuilder interface {
	Build() string
}

func NewProjectFilterPatternBuilder(project string) *ProjectFilterPatternBuilder {
	return &ProjectFilterPatternBuilder{project}
}

type ProjectFilterPatternBuilder struct {
	project string
}

func (b *ProjectFilterPatternBuilder) Build() string {
	return fmt.Sprintf("^%s.*v([0-9]{3})$", b.project)
}

func NewProjectAndApplicationFilterPatternBuilder(project string, application string) *ProjectAndApplicationFilterPatternBuilder {
	return &ProjectAndApplicationFilterPatternBuilder{project, application}
}

type ProjectAndApplicationFilterPatternBuilder struct {
	project     string
	application string
}

func (b *ProjectAndApplicationFilterPatternBuilder) Build() string {
	return fmt.Sprintf("^%s[\\-\\_]%s\\-v([0-9]{3})$", b.project, b.application)
}
