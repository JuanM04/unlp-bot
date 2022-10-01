package main

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type InfoData struct {
	ID            int
	Nombre        string
	Inicializado  bool
	UltimoMensaje CarteleraInfoMensaje
}

type CarteleraInfoResponse struct {
	Total    int                    `json:"total"`
	From     int                    `json:"from"`
	Count    int                    `json:"count"`
	Mensajes []CarteleraInfoMensaje `json:"mensajes"`
}

type CarteleraInfoMensaje struct {
	Materia        string `json:"materia"`
	Titulo         string `json:"titulo"`
	Cuerpo         string `json:"cuerpo"`
	Fecha          string `json:"fecha"`
	Autor          string `json:"autor"`
	IsAnulado      bool   `json:"is_anulado"`
	FechaAnulacion string `json:"fecha_anulacion"`
	Adjuntos       []struct {
		Nombre     string `json:"nombre"`
		PublicPath string `json:"public_path"`
	} `json:"adjuntos"`
}

func getInfoUpdates(data *InfoData) []CarteleraInfoMensaje {
	client := resty.New()

	resp, err := client.R().SetResult(CarteleraInfoResponse{}).Get(fmt.Sprintf("http://gestiondeaulas.info.unlp.edu.ar/cartelera/data/0/5?idMateria=%d", data.ID))
	if err != nil {
		log.Println(err)
	} else {
		result := resp.Result().(*CarteleraInfoResponse)

		if len(result.Mensajes) > 0 {
			if data.Inicializado {
				updates := make([]CarteleraInfoMensaje, 0)
				for _, mensaje := range result.Mensajes {
					if mensaje.Fecha != data.UltimoMensaje.Fecha {
						updates = append(updates, mensaje)
					}
				}
				data.UltimoMensaje = result.Mensajes[0]
				return updates
			} else {
				data.Inicializado = true
				data.UltimoMensaje = result.Mensajes[0]
			}

		}
	}

	return nil
}
