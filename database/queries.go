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
	id           INTEGER PRIMARY KEY,
	title        TEXT NOT NULL,
	description  TEXT NOT NULL,
	bookCount    INT NOT NULL,
	issueCount   INT NOT NULL,
	volumeCount  INT NOT NULL,
	omnibusCount INT NOT NULL,
	url          TEXT NOT NULL UNIQUE,
	uuid         TEXT NOT NULL UNIQUE,
	dateUpdated  INT NOT NULL,
	needUpdate   INT NOT NULL
);

-- Create issue table
CREATE TABLE issue (
	id              INTEGER PRIMARY KEY,
	seriesID        INT NOT NULL,
	title           TEXT NOT NULL,
	description     TEXT NOT NULL,
	publisher       TEXT NOT NULL,
	imprint         TEXT NOT NULL,
	issueNumber     TEXT NOT NULL,
	pages           INT NOT NULL,
	publicationDate INT NOT NULL,
	url             TEXT NOT NULL UNIQUE,
	uuid            TEXT NOT NULL UNIQUE,
	subscription    TEXT NOT NULL,
	toAdd           TEXT
);

-- Create series genre table
CREATE TABLE seriesGenre (
	id    INT NOT NULL,
	genre TEXT NOT NULL,
	PRIMARY KEY (id, genre),
	FOREIGN KEY (id) REFERENCES series(id) ON DELETE CASCADE
);

-- Create series imprint table
CREATE TABLE seriesImprint (
    id      INT NOT NULL,
	imprint TEXT NOT NULL,
	PRIMARY KEY (id, imprint),
	FOREIGN KEY (id) REFERENCES series(id) ON DELETE CASCADE
);

-- Create issue tag table
CREATE TABLE issueTag (
    id       INT NOT NULL,
	category TEXT NOT NULL,
	name     TEXT NOT NULL,
	PRIMARY KEY (id, category, name),
	FOREIGN KEY (id) REFERENCES issue(id) ON DELETE CASCADE
);

-- Create issue creator table
CREATE TABLE issueCreator (
	id          INT NOT NULL,
	type        TEXT NOT NULL,
	name        TEXT NOT NULL,
	displayName TEXT NOT NULL,
	PRIMARY KEY (id, type, name, displayName),
	FOREIGN KEY (id) REFERENCES issue(id) ON DELETE CASCADE
);`,
	// query to verify if database tables exist
	"pingDatabase": `SELECT *
FROM series
LIMIT 1;`,
}
