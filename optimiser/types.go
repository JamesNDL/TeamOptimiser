package optimiser

import (
	"errors"
	"fmt"
)

type Person struct {
	Name       string
	SkillLevel int
}

type Team struct {
	TeamId int
	People []Person
}

func (t *Team) CalculateTotalSkill() (int, error) {
	var total int

	if len(t.People) == 0 {
		return 0, errors.New("calculate total skill: team list has no people in it")
	}

	for _, person := range t.People {
		total += person.SkillLevel
	}

	return total, nil
}

// should 2 people be paired up or not
type Constraint struct {
	Person1   Person
	Person2   Person
	Inclusive bool
}

type ProblemSpace struct {
	People      []Person
	Teams       []Team
	Constraints []Constraint
}

func CreateProblemSpace(numberOfTeams int) (ProblemSpace, error) {

	if numberOfTeams <= 1 {
		return ProblemSpace{}, errors.New(fmt.Sprintf("numberOfTeams needs to be greater than 1!  Got: %d", numberOfTeams))
	}

	prob := ProblemSpace{}
	for i := 0; i < numberOfTeams; i++ {
		prob.Teams = append(prob.Teams, Team{
			TeamId: i,
			People: []Person{},
		})
	}

	return prob, nil
}

type Trial struct {
	People []Person
	Teams  []Team
	Cost   int
}
