package project

type TemplateVars struct {
	ProjectDir    string
	ProjectName   string
	ModuleName    string
	Author        string
	CopyrightYear string
}

type Project struct {
	vars        TemplateVars
	directories map[string]bool
	files       map[string]string
}
