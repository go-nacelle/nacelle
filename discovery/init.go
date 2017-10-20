package discovery

import (
	"github.com/efritz/reception"
	"github.com/satori/go.uuid"

	"github.com/efritz/nacelle"
)

// TODO - multiple watchers?
// TODO - name should be part of config?

func InitAnnouncer(name string, onDisconnect func(error)) nacelle.ServiceInitializerFunc {
	return func(config nacelle.Config, container *nacelle.ServiceContainer) error {
		client, err := makeClient(config, container)
		if err != nil {
			return err
		}

		service := &reception.Service{
			ID:      uuid.NewV4().String(),
			Name:    name,
			Address: "", // TODO
			Port:    0,  // TODO
		}

		if err := client.Register(service, onDisconnect); err != nil {
			return err
		}

		return container.Set("reception-client", client)
	}
}

func InitHotBackup(name string, onDisconnect func(error)) nacelle.ServiceInitializerFunc {
	return func(config nacelle.Config, container *nacelle.ServiceContainer) error {
		client, err := makeClient(config, container)
		if err != nil {
			return err
		}

		logger := container.GetLogger()
		logger.Info("Starting election")

		elector := reception.NewElector(
			client,
			name,
			reception.WithDisconnectionCallback(onDisconnect),
		)

		if err := elector.Elect(); err != nil {
			return err
		}

		logger.Info("Won election")
		return nil
	}
}

func InitWatcher(name string, f func(*reception.ServiceState)) nacelle.ServiceInitializerFunc {
	return func(config nacelle.Config, container *nacelle.ServiceContainer) error {
		client, err := makeClient(config, container)
		if err != nil {
			return err
		}

		ch, err := client.NewWatcher(name).Start()
		if err != nil {
			return err
		}

		go func() {
			for state := range ch {
				f(state)
			}
		}()

		return nil
	}
}
