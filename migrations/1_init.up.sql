CREATE TABLE `package` (
`id` int unsigned NOT NULL AUTO_INCREMENT,
`name` varchar(2048) NOT NULL DEFAULT '',
`version` varchar(32) NOT NULL DEFAULT '',
`publication_date` datetime NULL,
`title` varchar(2048) NOT NULL DEFAULT '',
`description` varchar(8192) NOT NULL DEFAULT '',
`authors` varchar(2048) NOT NULL DEFAULT '',
`maintainers` varchar(2048) NOT NULL DEFAULT '',
PRIMARY KEY (`id`)
);
