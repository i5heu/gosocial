package gosocial

func CreateTable() {
	db.Exec("CREATE TABLE `gosocial_comments` ( `ID` int(11) NOT NULL AUTO_INCREMENT, `slug` varchar(500) NOT NULL, `ModRelease` tinyint(4) NOT NULL, `submitTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, `Name` varchar(120) NOT NULL, `Title` varchar(300) NOT NULL, `Text` text NOT NULL, `upvotes` int(10) unsigned NOT NULL, `downvotes` int(10) unsigned NOT NULL, PRIMARY KEY (`ID`), KEY `slug` (`slug`) ) ENGINE=InnoDB DEFAULT CHARSET=latin1")
}
