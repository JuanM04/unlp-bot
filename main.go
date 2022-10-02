package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func getART() *time.Location {
	return time.FixedZone("America/Argentina/Buenos_Aires", int((-3 * time.Hour).Seconds()))
}

func main() {
	godotenv.Load()

	catedrasDeIngenieria := []IngData{
		{ID: "F0301", Nombre: "Matemática A"},
		{ID: "F0302", Nombre: "Matemática B"},
		{ID: "F0303", Nombre: "Física I"},
		{ID: "U1901", Nombre: "Química para Ingeniería"},
	}
	catedraDeGrafica := GraficaData{}
	catedrasDeInformatica := []InfoData{
		{ID: 105, Nombre: "Programación I"},
		{ID: 107, Nombre: "Programación I (Redictado)"},
		{ID: 106, Nombre: "Programación II"},
		{ID: 108, Nombre: "Programación II (Redictado)"},
	}

	log.Println("Initializing data...")

	for i := range catedrasDeIngenieria {
		getIngUpdates(&catedrasDeIngenieria[i])
	}
	getGraficaUpdates(&catedraDeGrafica)
	for i := range catedrasDeInformatica {
		getInfoUpdates(&catedrasDeInformatica[i])
	}

	scheduler := gocron.NewScheduler(getART())

	scheduler.Cron("*/15 * * * *").Do(func() {
		log.Println("Checking for updates...")
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
	})

	log.Println("Started UNLP bot")

	scheduler.StartBlocking()
}
