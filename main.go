/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package TinyKeyValueDB

import (
    
)


// Open(filename string) DB
type TinyKeyValueDB interface{
    Save()
    Get(key string, value interface{}) error
    Put(key string, value interface{}) error
    Delete(key string) error
}

type DBData struct{
	Keys []string
	Values []string
}
