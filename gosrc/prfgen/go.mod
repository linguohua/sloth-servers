module prfgen

go 1.12

require (
	gconst v0.0.0
	gpubsub v0.0.0
	gscfg v0.0.0
	pokerface v0.0.0
	prserver v0.0.0
)

replace prserver => ../prserver

replace gscfg => ../gscfg

replace pokerface => ../pokerface

replace gconst => ../gconst

replace gpubsub => ../gpubsub
