package handler

import (
	"dmaas/internal/app/dmaas/context"
	"dmaas/internal/app/dmaas/dto"
	"fmt"
)

type SourceHandler struct {
	Context *context.ApplicationContext
}

func (h SourceHandler) HandleSources(messages chan dto.SourceMessage) {
	for message := range messages {
		fmt.Println(fmt.Sprintf("handled %T message", message))
		switch message.Action {
		case dto.ImportDatabaseAction:
			fmt.Println(fmt.Sprintf("importing database with Title=%v", message.Source.Title))
			h.Context.SourceUseCase.ImportDatabase(*message.Source)
			break
		case dto.RemoveDatabaseAction:
			fmt.Println(fmt.Sprintf("removing database with Title=%v", message.Source.Title))
			h.Context.SourceUseCase.DeleteDatabase(*message.Source)
			break
		}
	}
}
