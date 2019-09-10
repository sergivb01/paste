package utils

var (
	CreatePastesTable = `CREATE TABLE IF NOT EXISTS pastes(
		id varchar(16) unique primary key,
		title varchar(25),
		content text,
		created timestamptz,
		expires timestamptz
	);`
)
