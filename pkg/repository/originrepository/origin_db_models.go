// SPDX-FileCopyrightText: 2026 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package originrepository

import "github.com/greenbone/opensight-notification-service/pkg/entities"

const (
	originsTable       = "notification_service.origins"
	deleteOriginsQuery = `DELETE FROM ` + originsTable + ` WHERE namespace = $1`
	createOriginsQuery = `INSERT INTO ` + originsTable + ` (name, class, namespace) VALUES (:name, :class, :namespace)`
	listOriginsQuery   = `SELECT * FROM ` + originsTable + ` ORDER BY namespace, name`
)

type originRow struct {
	Name      string `db:"name"`
	Class     string `db:"class"`
	Namespace string `db:"namespace"`
}

func toOrignRow(o entities.Origin, namespace string) originRow {
	return originRow{
		Name:      o.Name,
		Class:     o.Class,
		Namespace: namespace,
	}
}

func (r *originRow) toOriginEntity() entities.Origin {
	return entities.Origin(*r)
}
