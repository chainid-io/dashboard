package bolt

import "github.com/chainid-io/dashboard"

// Migrator defines a service to migrate data after a Chain Platform version update.
type Migrator struct {
	UserService            *UserService
	EndpointService        *EndpointService
	ResourceControlService *ResourceControlService
	SettingsService        *SettingsService
	VersionService         *VersionService
	CurrentDBVersion       int
	store                  *Store
}

// NewMigrator creates a new Migrator.
func NewMigrator(store *Store, version int) *Migrator {
	return &Migrator{
		UserService:            store.UserService,
		EndpointService:        store.EndpointService,
		ResourceControlService: store.ResourceControlService,
		SettingsService:        store.SettingsService,
		VersionService:         store.VersionService,
		CurrentDBVersion:       version,
		store:                  store,
	}
}

// Migrate checks the database version and migrate the existing data to the most recent data model.
func (m *Migrator) Migrate() error {

	// Chain Platform < 1.12
	if m.CurrentDBVersion < 1 {
		err := m.updateAdminUserToDBVersion1()
		if err != nil {
			return err
		}
	}

	// Chain Platform 1.12.x
	if m.CurrentDBVersion < 2 {
		err := m.updateResourceControlsToDBVersion2()
		if err != nil {
			return err
		}
		err = m.updateEndpointsToDBVersion2()
		if err != nil {
			return err
		}
	}

	// Chain Platform 1.13.x
	if m.CurrentDBVersion < 3 {
		err := m.updateSettingsToDBVersion3()
		if err != nil {
			return err
		}
	}

	// Chain Platform 1.14.0
	if m.CurrentDBVersion < 4 {
		err := m.updateEndpointsToDBVersion4()
		if err != nil {
			return err
		}
	}

	// https://github.com/chainid-io/dashboard/issues/1235
	if m.CurrentDBVersion < 5 {
		err := m.updateSettingsToVersion5()
		if err != nil {
			return err
		}
	}

	// https://github.com/chainid-io/dashboard/issues/1236
	if m.CurrentDBVersion < 6 {
		err := m.updateSettingsToVersion6()
		if err != nil {
			return err
		}
	}

	// https://github.com/chainid-io/dashboard/issues/1449
	if m.CurrentDBVersion < 7 {
		err := m.updateSettingsToVersion7()
		if err != nil {
			return err
		}
	}

	if m.CurrentDBVersion < 8 {
		err := m.updateEndpointsToVersion8()
		if err != nil {
			return err
		}
	}

	// https: //github.com/chainid-io/dashboard/issues/1396
	if m.CurrentDBVersion < 9 {
		err := m.updateEndpointsToVersion9()
		if err != nil {
			return err
		}
	}

	// https://github.com/chainid-io/dashboard/issues/461
	if m.CurrentDBVersion < 10 {
		err := m.updateEndpointsToVersion10()
		if err != nil {
			return err
		}
	}

	// https://github.com/chainid-io/dashboard/issues/1906
	if m.CurrentDBVersion < 11 {
		err := m.updateEndpointsToVersion11()
		if err != nil {
			return err
		}
	}

	err := m.VersionService.StoreDBVersion(chainid.DBVersion)
	if err != nil {
		return err
	}
	return nil
}