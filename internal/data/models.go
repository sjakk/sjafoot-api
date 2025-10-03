package data

import "database/sql"

type Models struct {
	Users UserModel
	Torcedores TorcedorModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
		Torcedores: TorcedorModel{DB: db},
	}
}
