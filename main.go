package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func distribute(roles []string, groupsByRole map[string]Stack, members []Member) {
	for _, role := range roles {
		var stack Stack
		groupsByRole[role] = stack
	}

	for _, member := range members {
		stack := groupsByRole[member.role]
		stack.Push(member)
		groupsByRole[member.role] = stack
	}

	for role, group := range groupsByRole {
		rand.Shuffle(len(group), func(i, j int) { group[i], group[j] = group[j], group[i] })
		groupsByRole[role] = group
	}
}

func assignTeam(groupsByRole map[string]Stack, dynamics map[string]int, teams [][]Member) [][]Member {
	isEnd := false
	for {
		if isEnd {
			break
		}

		team := make([]Member, 0)

		for role, group := range groupsByRole {
			if group.IsEmpty() {
				isEnd = true
				break
			}

			for i := 0; i < dynamics[role]; i++ {
				member, isNotEmpty := group.Pop()
				if isNotEmpty {
					team = append(team, member)
				}
			}

			groupsByRole[role] = group
		}

		if !isEnd {
			teams = append(teams, team)
		}
	}
	return teams
}

func printTeamDetails(teams [][]Member) {
	count := 1
	for _, team := range teams {
		memberDetails := fmt.Sprintf("%-20s | %-20s | %-10s\n", "Name", "Role", "Profiency")
		for _, member := range team {
			memberDetail := fmt.Sprintf("%-20s | %-20s | %-10d\n", member.name, member.role, member.profiency)
			memberDetails += memberDetail
		}
		fmt.Printf("Team %d\n", count)
		fmt.Print(memberDetails)
		fmt.Println()
		count += 1
	}
}

func setup() ([]string, map[string]Stack, [][]Member, map[string]int, []Member) {
	roles := []string{}
	members := []Member{}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter how many roles: ")
	
	numRoles, _ := reader.ReadString('\n')
	numOfRoles, err := strconv.Atoi(removeNextLine(numRoles))
	if err != nil {
		log.Fatalln(err)
	}

	if numOfRoles <= 0 {
		log.Fatalln("Number of roles cannot be negative or zero.")
	}
	
	for i:=0; i<numOfRoles; i++ {
		fmt.Print("Enter role: ")
		role, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		roles = append(roles, strings.Replace(removeNextLine(role), " ", "_", -1))
	}

	groupsByRole := make(map[string]Stack, len(roles))
	teams := make([][]Member, 0)
	dynamics := make(map[string]int)

	for _, role := range roles {
		format := fmt.Sprintf("Enter maximum number of people required for %s: ", role)
		fmt.Print(format)
		numPeople, _ := reader.ReadString('\n')
		numOfPeople, err := strconv.Atoi(removeNextLine(numPeople))
		if err != nil {
			log.Fatalln(err)
		}
		dynamics[role] = numOfPeople
	}

	for {
		fmt.Print("Enter name: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		if name == "\n" {
			break
		}

		fmt.Print("Enter role: ")
		role, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		if role == "\n" {
			break
		}

		fmt.Print("Enter profiency: ")
		profiency, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		if profiency == "\n" {
			break
		}
		numProfiency, err := strconv.Atoi(removeNextLine(profiency))
		if err != nil {
			log.Fatalln(err)
		}

		members = append(
			members, 
			Member{
				name: removeNextLine(name),
				role: removeNextLine(role),
				profiency: numProfiency,
			},
		)
		fmt.Println(members)

	}

	return roles, groupsByRole, teams, dynamics, members
}

func main() {
	rand.Seed(time.Now().UnixNano())
	roles, groupsByRole, teams, dynamics, members := setup()
	
	distribute(roles, groupsByRole, members)
	teams = assignTeam(groupsByRole, dynamics, teams)

	printTeamDetails(teams)
}
