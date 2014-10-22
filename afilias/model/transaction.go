package model

type Transaction struct {
	Model
	TLD              string
	Registrar_Ext_ID string
	Registrar_Name   string
	Server_TrID      string
	Command          string
	Object_Type      string
	Object_Name      string
	Transaction_Date string
}
