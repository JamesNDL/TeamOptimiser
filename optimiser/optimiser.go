package optimiser

import (
	"errors"
	"fmt"
	"main/math_functions"
	"math/rand/v2"
	"sort"
)

func (t *Trial) CalculateCost() error {
	// the cost is the "range" i.e. the difference between the best and the worst team
	var costs []int
	var maxSkill int
	var minSkill int
	var largestTeam int
	var smallestTeam int

	var teamSizes = make([]int, len(t.Teams))

	for i, team := range t.Teams {
		teamCost, err := team.CalculateTotalSkill()
		if err != nil {
			return err
		}
		costs = append(costs, teamCost)
		teamSizes[i] = len(team.People)
	}

	for i, skillTotal := range costs {
		if i == 0 {
			maxSkill = skillTotal
			minSkill = skillTotal
			continue
		}

		if skillTotal > maxSkill {
			maxSkill = skillTotal
		}

		if skillTotal < minSkill {
			minSkill = skillTotal
		}

	}

	for i, teamSize := range teamSizes {
		if i == 0 {
			largestTeam = teamSize
			smallestTeam = teamSize
			continue
		}

		if teamSize > largestTeam {
			largestTeam = teamSize
		}

		if teamSize < smallestTeam {
			smallestTeam = teamSize
		}
	}

	t.Cost = (maxSkill - minSkill) + (largestTeam - smallestTeam)

	return nil
}

func (t *Trial) Clone() Trial {
	newTeams := DeepCopyTeams(t.Teams)
	newPeople := append([]Person{}, t.People...)
	return Trial{
		People: newPeople,
		Teams:  newTeams,
		Cost:   t.Cost,
	}
}

func (t *Trial) Assign() error {
	shuffled, err := math_functions.RandomSampleWithoutReplacement(t.People, len(t.People))
	if err != nil {
		return err
	}

	// Clear all teams before assigning
	for i := range t.Teams {
		t.Teams[i].People = nil
	}

	team := 0
	for _, person := range shuffled {
		t.Teams[team].People = append(t.Teams[team].People, person)
		team++
		if team >= len(t.Teams) {
			team = 0
		}
	}

	return nil
}

func (t *Trial) MixTeams() error {
	indexes := make([]int, len(t.Teams))
	for i := range indexes {
		indexes[i] = i
	}

	teamIdx, err := math_functions.RandomSampleWithoutReplacement(indexes, 2)
	if err != nil {
		return err
	}

	team1 := &t.Teams[teamIdx[0]]
	team2 := &t.Teams[teamIdx[1]]

	if len(team1.People) == 0 || len(team2.People) == 0 {
		return nil // skip if no one to swap
	}

	i1 := rand.IntN(len(team1.People))
	i2 := rand.IntN(len(team2.People))

	// ✅ Swap directly, no remove/append needed
	team1.People[i1], team2.People[i2] = team2.People[i2], team1.People[i1]

	return nil
}

type Round struct {
	Trials    []Trial
	Iteration int
}

func (r *Round) SortTrials() {
	// find the lowest cost
	sort.Slice(r.Trials, func(i, j int) bool {
		return r.Trials[i].Cost < r.Trials[j].Cost
	})
}

func (r *Round) GenerateInitialTrials(trialsPerRound int, people []Person, teams []Team) error {
	if trialsPerRound <= 0 {
		return errors.New("trials per round cannot be 0 or below")
	}

	if len(people) <= 1 {
		return errors.New("must have more than 1 person for an optimisation")
	}

	if len(teams) <= 1 {
		return errors.New("must have more than 1 team for an optimisation")
	}

	for i := 0; i < trialsPerRound; i++ { // todo: pre-allocate the size of the slice to
		r.Trials = append(r.Trials,
			Trial{
				People: people,
				Teams:  DeepCopyTeams(teams),
				Cost:   0,
			},
		)
	}

	for i := 0; i < trialsPerRound; i++ {
		err := r.Trials[i].Assign()
		if err != nil {
			return err
		}

		err = r.Trials[i].CalculateCost()
		if err != nil {
			return err
		}
	}
	return nil
}

type OptimisationSettings struct {
	Iterations         int
	TrialsPerRound     int
	SurvivorPercentage float64
}

func (s *OptimisationSettings) NumberOfSurvivors() int {
	return int(float64(s.TrialsPerRound) * s.SurvivorPercentage)
}

func (p *ProblemSpace) Optimise(settings OptimisationSettings) error {

	if len(p.People) == 0 {
		return errors.New("There are no people in the problem space to optmise")
	}

	var rounds []Round //todo: pre-allocate the number of rounds because we know this number

	people := make([]Person, len(p.People))
	copy(people, p.People)

	round0 := Round{
		Iteration: 0,
	}
	err := round0.GenerateInitialTrials(settings.TrialsPerRound, p.People, p.Teams)
	fmt.Println("Initial Trails Generated.")

	if err != nil {
		panic(err)
	}

	round0.SortTrials()
	fmt.Println(fmt.Sprintf("Initial Cost %d", round0.Trials[0].Cost))

	rounds = append(rounds, round0)

	for iteration := 1; iteration < settings.Iterations; iteration++ {
		round := Round{
			Iteration: iteration,
		}

		previousRound := rounds[iteration-1]
		precedingTrials := previousRound.Trials[:settings.NumberOfSurvivors()]
		trialIdx := 0

		for i := 0; i < settings.TrialsPerRound; i++ { //todo: turn this into a go-routine
			newTrial := precedingTrials[trialIdx].Clone()

			if i > 0 {
				err = newTrial.MixTeams()

				if err != nil {
					return err
				}
			}

			err = newTrial.CalculateCost()

			if err != nil {
				return err
			}
			round.Trials = append(round.Trials, newTrial)

			trialIdx++

			if trialIdx >= len(precedingTrials) {
				trialIdx = 0
			}

		}

		round.SortTrials()
		rounds = append(rounds, round)
		fmt.Println(fmt.Sprintf("Iteration: %d, Cost: %d", iteration, round.Trials[0].Cost))

	}

	lastRound := rounds[len(rounds)-1]
	bestTrial := lastRound.Trials[0]
	fmt.Println("============================================")
	fmt.Println("Best Cost: ", bestTrial.Cost)
	for i, team := range bestTrial.Teams {
		fmt.Println("Team ", i)
		skill, _ := team.CalculateTotalSkill()
		fmt.Println("Total Skill ", skill)
		fmt.Println("Team Members: ")
		for j, person := range team.People {
			fmt.Println(fmt.Sprintf("	%d. name: %s skill: %d", j, person.Name, person.SkillLevel))
		}
		fmt.Println(" ")
	}
	return nil
}
