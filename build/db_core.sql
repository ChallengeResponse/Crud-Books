/* Create table for books */
DROP TABLE IF EXISTS `books`;
CREATE TABLE `books` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `author` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `publisher` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `publishDate` DATE NOT NULL,
  `rating` tinyint NOT NULL,
  `status` boolean NOT NULL,
  PRIMARY KEY (`id`)
);
