package autoload

/*
	You can automatically read the .env file when you import it by simply doing the following
		import _ "github.com/pchchv/env/autoload"
*/

import "github.com/pchchv/env"

func init() {
	env.Load()
}
