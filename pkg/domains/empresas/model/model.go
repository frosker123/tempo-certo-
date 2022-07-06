package model

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Agenda struct {
	ID      int64 `json:"ID,omitempty"`
	Empresa `json:"Empresa,omitempty"`
}

type Empresa struct {
	ID       int64  `json:"ID,omitempty"`
	Cnpj     string `json:"Cnpj,omitempty"`
	Fantasia string `json:"fantasia,omitempty"`
	Horario  string `json:"Horario,omitempty"`
}

type UpsertRequest struct {
	Empresa `json:"Empresa,omitempty"`
}

type ListAgendados struct {
	Cnpj     string `json:"cnpj"`
	Fantasia string `json:"fantasia"`
	Inicio   string `json:"Inicio"`
}

type HorarioAgendamento struct {
	ID         int64  `json:"ID,omitempty"`
	Inicio     string `json:"Inicio"`
	Fim        string `json:"Fim"`
	Disponivel bool   `json:"Disponivel"`
}

func NewAgendamento(fantasia, cnpj string, agendado string) (Agenda, error) {

	if len(cnpj) != 14 {
		return Agenda{}, errors.New("cnpj deve ter 14 caracteres")
	}

	agendadoEm, _ := time.Parse("15:04", agendado)
	if agendadoEm.IsZero() {
		agendadoEm = time.Now()
	}

	return Agenda{
		Empresa: Empresa{
			Fantasia: fantasia,
			Cnpj:     cnpj,
			Horario:  agendadoEm.Format("15:04"),
		},
	}, nil
}

func NewListaDisponibilidade() ([]HorarioAgendamento, error) {
	lista, err := os.Open("lista.json")
	if err != nil {
		return nil, err
	}
	defer lista.Close()

	var listaDisponibilidade []HorarioAgendamento
	err = json.NewDecoder(lista).Decode(&listaDisponibilidade)
	if err != nil {
		return nil, err
	}

	return listaDisponibilidade, nil
}
