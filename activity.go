package objectmapper

import (
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/core/data/resolve"
)

func init() {
	_ = activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})
var resolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{
	".":        &resolve.ScopeResolver{},
	"env":      &resolve.EnvResolver{},
	"property": &property.Resolver{},
	"loop":     &resolve.LoopResolver{},
})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Debugf("Setting: %v", s)

	act := &Activity{} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (bool, error) {

	var err error

	input := &Input{}
	err = ctx.GetInputObject(input)
	ctx.Logger().Debugf("Input: %v", input)
	if err != nil {
		return true, err
	}

	mapperFactory := mapper.NewFactory(resolver)

	var inVarMapper mapper.Mapper
	inVarMapper, err = mapperFactory.NewMapper(input.InVar)
	if err != nil {
		return true, err
	}

	var inVarValue interface{}
	inVarValue, err = inVarMapper.Apply(ctx.ActivityHost().Scope())
	if err != nil {
		return true, err
	}
	ctx.Logger().Debugf("outValue: %v", inVarValue)

	output := &Output{OutVar: inVarValue}
	ctx.Logger().Debugf("Output: %v", output)

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
