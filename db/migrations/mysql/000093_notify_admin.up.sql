CREATE TABLE IF NOT EXISTS NotifyAdmin (
    UserId varchar(26) NOT NULL,
    CreateAt bigint(20) DEFAULT NULL,
    RequiredPlan varchar(26) NOT NULL,
    RequiredFeature varchar(26) NOT NULL,
    Trial BOOLEAN NOT NULL,
    PRIMARY KEY (UserId, RequiredFeature, RequiredPlan)
);
