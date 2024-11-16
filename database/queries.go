package database

//nolint:gochecknoglobals
var queries = map[string]string{
	// Create entire database:
	"createDatabase": `-- Turn on foreign key support
PRAGMA foreign_keys = ON;

-- Drop any existing tables
DROP TABLE IF EXISTS issueCreator;
DROP TABLE IF EXISTS issueTag;
DROP TABLE IF EXISTS seriesImprint;
DROP TABLE IF EXISTS seriesGenre;
DROP TABLE IF EXISTS issue;
DROP TABLE IF EXISTS series;

-- Create series table
CREATE TABLE series (
	uuid         TEXT NOT NULL PRIMARY KEY,
	title        TEXT NOT NULL,
	description  TEXT NOT NULL,
	bookCount    INT NOT NULL,
	issueCount   INT NOT NULL,
	volumeCount  INT NOT NULL,
	omnibusCount INT NOT NULL,
	url          TEXT NOT NULL,
	dateUpdated  INT NOT NULL,
	needUpdate   INT NOT NULL
);

-- Create issue table
CREATE TABLE issue (
	uuid            TEXT NOT NULL PRIMARY KEY,
	seriesUUID      INT NOT NULL,
	title           TEXT NOT NULL,
	description     TEXT NOT NULL,
	publisher       TEXT NOT NULL,
	imprint         TEXT NOT NULL,
	issueNumber     TEXT NOT NULL,
	pages           INT NOT NULL,
	publicationDate INT NOT NULL,
	url             TEXT NOT NULL,
	subscription    TEXT NOT NULL,
	toAdd           TEXT,
	FOREIGN KEY (seriesUUID) REFERENCES series(uuid) ON DELETE CASCADE
);

-- Create series genre table
CREATE TABLE seriesGenre (
	uuid  INT NOT NULL,
	genre TEXT NOT NULL,
	PRIMARY KEY (uuid, genre),
	FOREIGN KEY (uuid) REFERENCES series(uuid) ON DELETE CASCADE
);

-- Create series imprint table
CREATE TABLE seriesImprint (
    uuid    INT NOT NULL,
	imprint TEXT NOT NULL,
	PRIMARY KEY (uuid, imprint),
	FOREIGN KEY (uuid) REFERENCES series(uuid) ON DELETE CASCADE
);

-- Create issue tag table
CREATE TABLE issueTag (
    uuid     INT NOT NULL,
	category TEXT NOT NULL,
	name     TEXT NOT NULL,
	PRIMARY KEY (uuid, category, name),
	FOREIGN KEY (uuid) REFERENCES issue(uuid) ON DELETE CASCADE
);

-- Create issue creator table
CREATE TABLE issueCreator (
	uuid        INT NOT NULL,
	type        TEXT NOT NULL,
	name        TEXT NOT NULL,
	displayName TEXT NOT NULL,
	PRIMARY KEY (uuid, type, name, displayName),
	FOREIGN KEY (uuid) REFERENCES issue(uuid) ON DELETE CASCADE
);`,
	// query to verify if database tables exist
	"pingDatabase": `SELECT *
FROM series
LIMIT 1;`,
}
