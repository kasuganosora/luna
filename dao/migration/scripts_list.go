package migration

import migrationScripts "github.com/kabukky/journey/dao/migration/scripts"

// scripts  Migration scripts list
var scripts = []ScriptInterface{
	&migrationScripts.InstallScript{},
}
