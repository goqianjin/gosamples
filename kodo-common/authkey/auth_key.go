package authkey

import "github.com/qianjin/kodo-security/kodokey"

type AuthKey struct {
	AK string
	SK string
}

// -------------------
var (
	// dev user

	Dev_Key_general_storage_011 = AuthKey{kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011}
	Dev_Key_general_storage_002 = AuthKey{kodokey.Dev_AK_general_storage_002, kodokey.Dev_SK_general_torage_002}

	// dev admin

	Dev_Key_admin = AuthKey{kodokey.Dev_AK_admin, kodokey.Dev_SK_admin}

	// prod user

	Prod_Key_shenqianjin = AuthKey{kodokey.Prod_AK_shenqianjin, kodokey.Prod_SK_shenqianjin}

	// prod admin

	Prod_Key_admin = AuthKey{kodokey.Prod_AK_admin, kodokey.Prod_SK_admin}
)
