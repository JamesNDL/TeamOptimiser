package main

import (
	"fmt"
	"main/optimiser"
)

func main() {
	fmt.Println("I am running")
	problemSpace, err := optimiser.CreateProblemSpace(4)

	people := []optimiser.Person{
		optimiser.Person{Name: "Jen", SkillLevel: 10},
		optimiser.Person{Name: "Jack", SkillLevel: 5},
		optimiser.Person{Name: "John", SkillLevel: 1},
		optimiser.Person{Name: "Michael", SkillLevel: 10},
		optimiser.Person{Name: "Zac", SkillLevel: 5},
		optimiser.Person{Name: "James", SkillLevel: 1},
		optimiser.Person{Name: "Grace", SkillLevel: 10},
		optimiser.Person{Name: "Hamish", SkillLevel: 1},
		optimiser.Person{Name: "Chris", SkillLevel: 1},
		optimiser.Person{Name: "David", SkillLevel: 5},
		optimiser.Person{Name: "Jenny", SkillLevel: 1},
		optimiser.Person{Name: "Paul", SkillLevel: 5},
		optimiser.Person{Name: "Wendy", SkillLevel: 10},
		optimiser.Person{Name: "Tim", SkillLevel: 1},
		optimiser.Person{Name: "Jim", SkillLevel: 1},
		optimiser.Person{Name: "Stephen", SkillLevel: 1},
		optimiser.Person{Name: "Ben", SkillLevel: 5},
		optimiser.Person{Name: "Basil", SkillLevel: 1},
	}

	problemSpace.People = people

	constraints := []optimiser.Constraint{
		optimiser.Constraint{Person1: problemSpace.People[0], Person2: problemSpace.People[1], Inclusive: true},
		optimiser.Constraint{Person1: problemSpace.People[5], Person2: problemSpace.People[6], Inclusive: true},
		optimiser.Constraint{Person1: problemSpace.People[1], Person2: problemSpace.People[2], Inclusive: false},
	}

	problemSpace.Constraints = constraints

	if err != nil {
		panic(err)
	}

	optimisationSettings := optimiser.OptimisationSettings{
		Iterations:         10,
		TrialsPerRound:     10,
		SurvivorPercentage: 0.1,
	}

	err = problemSpace.Optimise(optimisationSettings)

	if err != nil {
		fmt.Println(err)
	}

}
