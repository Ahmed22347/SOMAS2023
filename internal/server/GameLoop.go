package server

import (
	"SOMAS2023/internal/common/objects"
	"fmt"
)

func (s *Server) RunGameLoop() {
	// Each agent makes a decision
	for agentId, agent := range s.GetAgentMap() {
		fmt.Printf("Agent %s updating state \n", agentId)
		agent.UpdateGameState(s)
		agent.UpdateAgentInternalState()
		switch agent.DecideAction() {
		case objects.Pedal:
			agent.DecideForce()
		case objects.ChangeBike:
			newBikeId := agent.ChangeBike()
			// Remove the agent from the old bike, if it was on one
			if oldBikeId, ok := s.megaBikeRiders[agentId]; ok {
				s.megaBikes[oldBikeId].RemoveAgent(agentId)
			}
			// Add the agent to the new bike
			s.megaBikes[newBikeId].AddAgent(agent)
			s.megaBikeRiders[agentId] = newBikeId
		default:
			panic("agent decided invalid action")
		}
	}

	for _, bike := range s.GetMegaBikes() {
		bike.Move()
	}

	// Lootbox Distribution
	s.LootboxCheckAndDistributions()

	// Replenish objects
	s.replenishLootBoxes()
	s.replenishMegaBikes()
}

func (s *Server) LootboxCheckAndDistributions() {
	for bikeid, megabike := range s.GetMegaBikes() {
		for lootid, lootbox := range s.GetLootBoxes() {
			if megabike.CheckForCollision(lootbox) {
				// Collision detected
				fmt.Printf("Collision detected between MegaBike %s and LootBox %s \n", bikeid, lootid)
				agents := megabike.GetAgents()
				totAgents := len(agents)

				for _, agent := range agents {
					// this function allows the agent to decide on its allocation parameters
					// these are the parameters that we want to be considered while carrying out
					// the elected protocol for resource allocation
					agent.SetAllocationParameters()

					// in the MVP  the allocation parameters are ignored and
					// the utility share will simply be 1 / the number of agents on the bike
					utilityShare := 1 / totAgents
					lootShare := utilityShare * lootbox.GetTotalResources()
					// Allocate loot based on the calculated utility share
					fmt.Printf("Agent %s allocated %f loot \n", agent.GetID(), lootShare)
					agent.SetEnergyLevel(lootShare)
				}
			}
		}
	}
}

func (s *Server) Start() {
	fmt.Printf("Server initialised with %d agents \n\n", len(s.GetAgentMap()))
	for i := 0; i < s.GetIterations(); i++ {
		fmt.Printf("Game Loop %d running... \n \n", i)
		fmt.Printf("Main game loop running...\n\n")
		s.RunGameLoop()
		fmt.Printf("\nMain game loop finished.\n\n")
		fmt.Printf("Messaging session started...\n\n")
		s.RunMessagingSession()
		fmt.Printf("\nMessaging session completed\n\n")
		fmt.Printf("Game Loop %d completed.\n", i)
	}
}
