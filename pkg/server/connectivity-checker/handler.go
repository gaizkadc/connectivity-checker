/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package connectivity_checker

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}