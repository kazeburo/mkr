package hosts

import (
	"fmt"
	"io"
	"text/template"

	mackerel "github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type hostApp struct {
	client    mackerelclient.Client
	logger    *logger.Logger
	outStream io.Writer
}

type findHostsParam struct {
	verbose bool

	name     string
	service  string
	roles    []string
	statuses []string

	format string
}

func (ha *hostApp) findHosts(param findHostsParam) error {
	hosts, err := ha.client.FindHosts(&mackerel.FindHostsParam{
		Name:     param.name,
		Service:  param.service,
		Roles:    param.roles,
		Statuses: param.statuses,
	})
	if err != nil {
		return err
	}

	switch {
	case param.format != "":
		t, err := template.New("format").Parse(param.format)
		if err != nil {
			return err
		}
		return t.Execute(ha.outStream, hosts)
	case param.verbose:
		return format.PrettyPrintJSON(ha.outStream, hosts)
	default:
		var hostsFormat []*format.Host
		for _, host := range hosts {
			hostsFormat = append(hostsFormat, &format.Host{
				ID:            host.ID,
				Name:          host.Name,
				DisplayName:   host.DisplayName,
				Status:        host.Status,
				RoleFullnames: host.GetRoleFullnames(),
				IsRetired:     host.IsRetired,
				CreatedAt:     format.ISO8601Extended(host.DateFromCreatedAt()),
				IPAddresses:   host.IPAddresses(),
			})
		}
		return format.PrettyPrintJSON(ha.outStream, hostsFormat)
	}
}

type createHostParam struct {
	Name             string
	RoleFullnames    []string
	Status           string
	CustomIdentifier string
}

func (ha *hostApp) createHost(param createHostParam) error {
	hostID, err := ha.client.CreateHost(&mackerel.CreateHostParam{
		Name:             param.Name,
		RoleFullnames:    param.RoleFullnames,
		CustomIdentifier: param.CustomIdentifier,
	})
	ha.dieIf(err)

	ha.log("created", hostID)

	if param.Status != "" {
		err := ha.client.UpdateHostStatus(hostID, param.Status)
		ha.dieIf(err)
		ha.log("updated", fmt.Sprintf("%s %s", hostID, param.Status))
	}
	return nil
}

func (ha *hostApp) log(prefix, message string) {
	if ha.logger != nil {
		ha.logger.Log(prefix, message)
	}
}

func (ha *hostApp) dieIf(err error) {
	if ha.logger != nil {
		ha.logger.DieIf(err)
	}
}
