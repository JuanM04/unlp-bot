package main

import (
	"fmt"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	catedrasDeIngenieria := []IngData{
		{ID: "F0301", Nombre: "Matemática A"},
		{ID: "F0302", Nombre: "Matemática B"},
		{ID: "F0303", Nombre: "Física I"},
		{ID: "U1901", Nombre: "Química para Ingeniería"},
	}
	catedrasDeInformatica := []InfoData{
		{ID: 105, Nombre: "Programación I"},
		{ID: 107, Nombre: "Programación I (Redictado)"},
		{ID: 106, Nombre: "Programación II"},
		{ID: 108, Nombre: "Programación II (Redictado)"},
	}
	catedraDeGrafica := GraficaData{}

	ingUpdates := make([]string, 0)
	for i := range catedrasDeIngenieria {
		message := getIngUpdates(&catedrasDeIngenieria[i])
		if message != "" {
			ingUpdates = append(ingUpdates, fmt.Sprintf("**%s** - %s", catedrasDeIngenieria[i].Nombre, message))
		}
	}

	graficaUpdate := getGraficaUpdates(&catedraDeGrafica)
	if graficaUpdate != "" {
		ingUpdates = append(ingUpdates, fmt.Sprintf("**Gráfica para Ingeniería** - %s", graficaUpdate))
	}

	if len(ingUpdates) > 0 {
		invokeWebhook(WebhookPayload{
			Content: strings.Join(ingUpdates, "\n"),
		})
	}

	mdConverter := md.NewConverter("https://gestiondeaulas.info.unlp.edu.ar/cartelera/", true, nil)
	for i := range catedrasDeInformatica {
		updates := getInfoUpdates(&catedrasDeInformatica[i])
		if len(updates) > 0 {
			embeds := make([]WebhookEmbed, 0)
			for _, update := range updates {
				cuerpo, _ := mdConverter.ConvertString(update.Cuerpo)
				adjuntos := make([]WebhookEmbedField, 0)

				for _, adjunto := range update.Adjuntos {
					adjuntos = append(adjuntos, WebhookEmbedField{
						Name:  adjunto.Nombre,
						Value: adjunto.PublicPath,
					})
				}

				embeds = append(embeds, WebhookEmbed{
					Title:       update.Titulo,
					Type:        "rich",
					Description: cuerpo,
					Color:       15158332,
					Author: WebhookEmbedAuthor{
						Name: strings.Trim(update.Autor, " "),
					},
					Fields: adjuntos,
				})
			}
			invokeWebhook(WebhookPayload{
				Content: fmt.Sprintf("Novedades de **%s** - [ver todas →](<https://gestiondeaulas.info.unlp.edu.ar/cartelera/#form[materia]=%d>)", catedrasDeInformatica[i].Nombre, catedrasDeInformatica[i].ID),
				Embeds:  embeds,
			})
		}
	}
}
