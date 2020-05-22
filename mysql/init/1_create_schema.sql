# Install Mroonga, CJK-ready fulltext engine
INSTALL SONAME 'ha_mroonga';
DROP FUNCTION IF EXISTS last_insert_grn_id;
CREATE FUNCTION last_insert_grn_id RETURNS INTEGER SONAME 'ha_mroonga.so';
DROP FUNCTION IF EXISTS mroonga_snippet;
CREATE FUNCTION mroonga_snippet RETURNS STRING SONAME 'ha_mroonga.so';
DROP FUNCTION IF EXISTS mroonga_command;
CREATE FUNCTION mroonga_command RETURNS STRING SONAME 'ha_mroonga.so';
DROP FUNCTION IF EXISTS mroonga_escape;
CREATE FUNCTION mroonga_escape RETURNS STRING SONAME 'ha_mroonga.so';

CREATE TABLE IF NOT EXISTS `message` (
    `id` CHAR(36) NOT NULL,
    `channel_id` CHAR(36) NOT NULL,
    `created_at` DATETIME(6) NOT NULL,
    `text` TEXT CHARSET `utf8mb4` COLLATE `utf8mb4_bin` NOT NULL,
    FULLTEXT KEY `full_text` (`text`)
) ENGINE=Mroonga DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `seen_channel` (
    id CHAR(36) PRIMARY KEY,
    last_processed_message DATETIME(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
