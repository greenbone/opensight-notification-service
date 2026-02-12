// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package originrepository

import "github.com/greenbone/opensight-notification-service/pkg/entities"

const (
	originsTable       = "notification_service.origins"
	deleteOriginsQuery = `DELETE FROM ` + originsTable + ` WHERE service_id = $1`
	createOriginsQuery = `INSERT INTO ` + originsTable + ` (name, class, service_id) VALUES (:name, :class, :service_id)`
	listOriginsQuery   = `SELECT * FROM ` + originsTable + ` ORDER BY service_id, name`
)

type originRow struct {
	Name      string `db:"name"`
	Class     string `db:"class"`
	ServiceID string `db:"service_id"`
}

func toOriginRow(o entities.Origin, serviceID string) originRow {
	return originRow{
		Name:      o.Name,
		Class:     o.Class,
		ServiceID: serviceID,
	}
}

func (r *originRow) toOriginEntity() entities.Origin {
	return entities.Origin(*r)
}
