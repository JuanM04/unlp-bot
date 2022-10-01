package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

type GraficaData struct {
	Inicio       string
	Apuntes      string
	Guias        string
	Animaciones  string
	Parcial      string
	Inicializado bool
}

func getGraficaUpdates(data *GraficaData) string {
	client := resty.New()
	client.SetBaseURL("https://catedras.ing.unlp.edu.ar/grafica/")

	body := make([]string, 0)

	var res *resty.Response
	var err error

	// Inicio
	res, err = client.R().Get("/")
	if err != nil {
		log.Println(err)
	} else {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res.String()))
		if err != nil {
			log.Println(err)
		} else {
			content, err := doc.Find(".main-inner").Html()
			if err != nil {
				log.Println(err)
			} else {
				if data.Inicio != content {
					body = append(body, "la página de [inicio](<https://catedras.ing.unlp.edu.ar/grafica/>)")
				}
				data.Inicio = content
			}
		}
	}

	// Apuntes
	// Loguearse una vez para que funcione
	res, err = client.R().SetFormData(map[string]string{"acpwd-pass": "GPI2022"}).Post("/biblioteca/apuntes-de-catedra")
	if err != nil {
		log.Println(err)
	} else {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res.String()))
		if err != nil {
			log.Println(err)
		} else {
			content, err := doc.Find(".main-inner").Html()
			if err != nil {
				log.Println(err)
			} else {
				if data.Apuntes != content {
					body = append(body, "los [apuntes de cátedra](<https://catedras.ing.unlp.edu.ar/grafica/biblioteca/apuntes-de-catedra>)")
				}
				data.Apuntes = content
			}
		}
	}

	// Guías
	res, err = client.R().Get("/biblioteca/guias-de-trabajo")
	if err != nil {
		log.Println(err)
	} else {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res.String()))
		if err != nil {
			log.Println(err)
		} else {
			content, err := doc.Find(".main-inner").Html()
			if err != nil {
				log.Println(err)
			} else {
				if data.Guias != content {
					body = append(body, "las [guías de trabajo](<https://catedras.ing.unlp.edu.ar/grafica/biblioteca/guias-de-trabajo>)")
				}
				data.Guias = content
			}
		}
	}

	// Animaciones
	res, err = client.R().Get("/biblioteca/animaciones")
	if err != nil {
		log.Println(err)
	} else {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res.String()))
		if err != nil {
			log.Println(err)
		} else {
			content, err := doc.Find(".main-inner").Html()
			if err != nil {
				log.Println(err)
			} else {
				if data.Animaciones != content {
					body = append(body, "las [animaciones y videos](<https://catedras.ing.unlp.edu.ar/grafica/biblioteca/animaciones>)")
				}
				data.Animaciones = content
			}
		}
	}

	// Parcial
	res, err = client.R().Get("/parcial")
	if err != nil {
		log.Println(err)
	} else {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res.String()))
		if err != nil {
			log.Println(err)
		} else {
			content, err := doc.Find(".main-inner").Html()
			if err != nil {
				log.Println(err)
			} else {
				if data.Parcial != content {
					body = append(body, "la página del [parcial](<https://catedras.ing.unlp.edu.ar/grafica/parcial>)")
				}
				data.Parcial = content
			}
		}
	}

	if data.Inicializado {
		if len(body) >= 2 {
			return "Se actualizó " + strings.Join(body[:len(body)-1], ", ") + " y " + body[len(body)-1]
		} else if len(body) == 1 {
			return "Se actualizó " + body[0]
		}
	} else {
		data.Inicializado = true
	}

	return ""
}
