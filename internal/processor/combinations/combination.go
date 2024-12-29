package combinations

type Combination struct {
	GameId      string
	CategoryId  string
	LevelId     *string
	VariableIds []string
	ValueIds    []string
}

func (c *Combination) isValid() bool {
	return len(c.VariableIds) == len(c.ValueIds)
}
