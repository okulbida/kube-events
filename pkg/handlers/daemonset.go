package handlers

import (
	"fmt"
	"github.com/anodot/anodot-common/pkg/common"
	"github.com/anodot/anodot-common/pkg/events"
	v1 "k8s.io/api/apps/v1"
)

type DaemonsetHandler struct {
	UserEventConfiguration
}

func (d *DaemonsetHandler) SupportedEvent() string {
	return "daemonset"
}

func (d *DaemonsetHandler) EventData(event Event) ([]events.Event, error) {
	allEvents := make([]events.Event, 0)

	switch event.EventType {
	case "update":
		if event.New == nil || event.Old == nil {
			//TODO better error message
			return nil, fmt.Errorf("unable to retrieve deployment information")
		}

		newDep := event.New.(*v1.DaemonSet)
		oldDeployment := event.Old.(*v1.DaemonSet)

		deploymentName := newDep.Name

		//image changed
		for _, newC := range newDep.Spec.Template.Spec.Containers {
			for _, oldC := range oldDeployment.Spec.Template.Spec.Containers {
				if newC.Name == oldC.Name {
					if newC.Image != oldC.Image {
						res := events.Event{
							Title:       fmt.Sprintf("'%s' daemonset image changed", deploymentName),
							Description: fmt.Sprintf("%s image daemonset from '%s' to '%s'", deploymentName, oldC.Image, newC.Image),
							Category:    d.Category,
							Source:      d.Source,
							Properties: []events.EventProperties{
								{Key: "daemonset", Value: deploymentName},
								{Key: "namespace", Value: newDep.Namespace},
								{Key: "container", Value: newC.Name}},
							StartDate: common.AnodotTimestamp{Time: event.EventTime},
						}
						allEvents = append(allEvents, res)
					}
				}
				//TODO env check ?
				/*if !reflect.DeepEqual(c.Env, oldC.Env) {
				}*/
			}
		}

	case "delete":

	}

	return allEvents, nil
}
