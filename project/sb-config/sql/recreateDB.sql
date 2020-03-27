use sports_betting_db;
drop table if exists Bet;
drop table if exists Odd;
drop table if exists Line;
drop table if exists GameEvent;
drop table if exists Game2Team;
drop table if exists Game;
drop table if exists Player2Team;
drop table if exists Player;
drop table if exists Team;
drop table if exists League;
drop table if exists Sport;
drop table if exists Site2User;
drop table if exists `Site`;
drop table if exists User;

create table Sport(
	ID bigint unsigned auto_increment primary key,
	`Name` varchar(255) not null,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table League(
	ID bigint unsigned auto_increment primary key,
	`Name` varchar(255) not null,
	Country varchar(255) not null,
	SportID bigint unsigned not null,
	foreign key (SportID) references Sport(ID) on update CASCADE on delete RESTRICT,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table Player(
	ID bigint unsigned auto_increment primary key,
	`Name` varchar(255) not null,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table Team(
	ID bigint unsigned auto_increment primary key,
	`Name` varchar(255) not null,
	LeagueID bigint unsigned not null,
	foreign key (LeagueID) references League(ID) on update CASCADE on delete RESTRICT,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table Player2Team(
	ID bigint unsigned auto_increment primary key,
	PlayerID bigint unsigned not null,
	TeamID bigint unsigned not null,
	foreign key (PlayerID) references Player(ID) on update CASCADE on delete RESTRICT,
	foreign key (TeamID) references Team(ID) on update CASCADE on delete RESTRICT
);

create table Game (
	ID bigint unsigned auto_increment primary key,
	`Name` varchar(255) not null,
	SportID bigint unsigned not null,
	foreign key (SportID) references Sport(ID) on update CASCADE on delete RESTRICT,
	LeagueID bigint unsigned not null,
	foreign key (LeagueID) references League(ID) on update CASCADE on delete RESTRICT,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table Game2Team(
	ID bigint unsigned auto_increment primary key,
	GameID bigint unsigned not null,
	foreign key (GameID) references Game(ID) on update CASCADE on delete RESTRICT,
	TeamID bigint unsigned not null,
	foreign key (TeamID) references Team(ID) on update CASCADE on delete RESTRICT
);

create table GameEvent (
	ID bigint unsigned auto_increment primary key,
	Date timestamp not null,
	EventType int not null,
	IsFrozen bit default 0,
	RelatedGameID bigint unsigned not null,
	foreign key (RelatedGameID) references Game(ID) on update CASCADE on delete RESTRICT,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table Line(
	ID bigint unsigned auto_increment primary key,
	OddFormat int not null,
	LineType int not null,
	`Description` varchar(255) not null,
	EventID bigint unsigned not null,
	foreign key (EventID) references GameEvent(ID) on update CASCADE on delete RESTRICT,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table Odd(
	ID bigint unsigned auto_increment primary key,
	Source varchar(255) not null,
	`Values` varchar(50) not null, -- 1.2:5.6:6.7
	LineID bigint unsigned not null,
	foreign key (LineID) references Line(ID) on update CASCADE on delete RESTRICT,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table Site(
	ID bigint unsigned auto_increment primary key,
	`Name` varchar(255) not null,
	Country varchar(50) not null,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table User(
	ID bigint unsigned auto_increment primary key,
	`Name` varchar(255) not null,
	`Password` varchar(255) not null,
	`Roles` int not null,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);

create table Site2User(
	ID bigint unsigned auto_increment primary key,
	SiteID bigint unsigned not null,
	foreign key (SiteID) references Site(ID) on update CASCADE on delete RESTRICT,
	UserID bigint unsigned not null,
	foreign key (UserID) references User(ID) on update CASCADE on delete RESTRICT
);

create table Bet(
	ID bigint unsigned auto_increment primary key,
	IsPayed bit default 0,
	BetOnIndex unsigned int not null,
	LineID bigint unsigned not null,
	foreign key (LineID) references Line(ID) on update CASCADE on delete RESTRICT,
	UserID bigint unsigned not null,
	foreign key (UserID) references User(ID) on update CASCADE on delete RESTRICT,
	`Value` decimal(15, 3) not null,

	CDate timestamp default NOW(),
	MDate timestamp default NOW(),
	IsDeleted bit default 0
);
