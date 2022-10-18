/*
* (C) 2022 Wenchao Lv
 */

package cmd

const (
	version = "0.1"
	logo    = `
        _  _               _                       
       | |(_)             | |                      
   ___ | | _  _ __    ___ | |__    __ _  _ __  ___ 
  / __|| || || '_ \  / __|| '_ \  / _ || '__|/ _ \
 | (__ | || || |_) | \__ \| | | || (_| || |  |  __/
  \___||_||_|| .__/  |___/|_| |_| \__,_||_|   \___|
             | |                                   
             |_|                                   
			 `
)

func Version() string {
	return version
}

func Logo() string {
	return logo
}
