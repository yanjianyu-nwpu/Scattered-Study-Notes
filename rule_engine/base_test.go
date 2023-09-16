package rule_engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type MyFact struct {
	IntAttribute     int64
	StringAttribute  string
	BooleanAttribute bool
	FloatAttribute   float64
	TimeAttribute    time.Time
	WhatToSay        string
}

func (mf *MyFact) GetWhatToSay(sentence string) string {
	return fmt.Sprintf("Let say \"%s\"", sentence)
}
func TestBase(t *testing.T) {
	drls := `
	rule CheckValues "Check the default values" salience 10 {
		when 
			MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
		then
			MF.WhatToSay = MF.GetWhatToSay("Hello Grule");
			Retract("CheckValues");
	}
	`
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	bs := pkg.NewBytesResource([]byte(drls))
	err := ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", bs)
	if err != nil {
		panic(err)
	}

	myFact := &MyFact{
		IntAttribute:     123,
		StringAttribute:  "Some string value",
		BooleanAttribute: true,
		FloatAttribute:   1.234,
		TimeAttribute:    time.Now(),
	}

	dataCtx := ast.NewDataContext()
	err = dataCtx.Add("MF", myFact)
	if err != nil {
		panic(err)
	}
	engine := engine.NewGruleEngine()
	err = engine.Execute(dataCtx, knowledgeLibrary.GetKnowledgeBase("TutorialRules", "0.0.1"))
	if err != nil {
		panic(err)
	}

	t.Log(myFact.WhatToSay)
}
