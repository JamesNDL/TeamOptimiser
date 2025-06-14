package optimiser

func DeepCopyTeams(original []Team) []Team {
	copied := make([]Team, len(original))
	for i, team := range original {
		copied[i] = Team{
			TeamId: team.TeamId,
			People: append([]Person{}, team.People...),
		}
	}
	return copied
}
