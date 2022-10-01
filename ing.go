package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

type IngData struct {
	ID      string
	Nombre  string
	Paginas IngDataPaginas
}

type IngDataPaginas struct {
	Cronograma   string
	Horarios     string
	Notas        string
	Reglamento   string
	Descargas    string
	Novedades    string
	Bibliografia string
}

func getIngUpdates(data *IngData) string {
	client := resty.New()
	client.SetBaseURL("https://www.ing.unlp.edu.ar/catedras/")

	body := make([]string, 0)

	var res *resty.Response
	var err error

	// Call the home to get some cookies
	client.R().Get(fmt.Sprintf("/%s/index.php", data.ID))

	// Cronograma
	res, err = client.R().Get("/template/cronograma.php?id=" + data.ID)
	if err != nil {
		log.Println(err)
	} else {
		raw := res.String()
		if data.Paginas.Cronograma != raw {
			body = append(body, fmt.Sprintf("el [cronograma](<https://www.ing.unlp.edu.ar/catedras/%s/index.php?secc=cronograma>)", data.ID))
		}
		data.Paginas.Cronograma = raw
	}

	// Horarios
	res, err = client.R().Get("/template/horarios.php?id=" + data.ID)
	if err != nil {
		log.Println(err)
	} else {
		raw := res.String()
		if data.Paginas.Horarios != raw {
			body = append(body, fmt.Sprintf("los [horarios](<https://www.ing.unlp.edu.ar/catedras/%s/index.php?secc=horarios>)", data.ID))
		}
		data.Paginas.Horarios = raw
	}

	// Notas
	res, err = client.R().Get(fmt.Sprintf("/%s/list.php?secc=1&id=%s", data.ID, data.ID))
	if err != nil {
		log.Println(err)
	} else {
		raw := res.String()
		if data.Paginas.Notas != raw {
			body = append(body, fmt.Sprintf("la sección de [notas/comisiones](<https://www.ing.unlp.edu.ar/catedras/%s/index.php?secc=notas>)", data.ID))
		}
		data.Paginas.Notas = raw
	}

	// Reglamento
	res, err = client.R().Get("/template/reglamento.php?id=" + data.ID)
	if err != nil {
		log.Println(err)
	} else {
		raw := res.String()
		if data.Paginas.Reglamento != raw {
			body = append(body, fmt.Sprintf("el [reglamento](<https://www.ing.unlp.edu.ar/catedras/%s/index.php?secc=reglamento>)", data.ID))
		}
		data.Paginas.Reglamento = raw
	}

	// Descargas
	res, err = client.R().Get(fmt.Sprintf("/%s/list.php?secc=0&id=%s", data.ID, data.ID))
	if err != nil {
		log.Println(err)
	} else {
		raw := res.String()
		if data.Paginas.Descargas != raw {
			body = append(body, fmt.Sprintf("la sección de [descargas](<https://www.ing.unlp.edu.ar/catedras/%s/index.php?secc=descargas>)", data.ID))
		}
		data.Paginas.Descargas = raw
	}

	// Novedades
	res, err = client.R().Get("/template/novedades.php?id=" + data.ID)
	if err != nil {
		log.Println(err)
	} else {
		raw := res.String()
		if data.Paginas.Novedades != raw {
			body = append(body, fmt.Sprintf("la sección de [novedades](<https://www.ing.unlp.edu.ar/catedras/%s/index.php?secc=novedades>)", data.ID))
		}
		data.Paginas.Novedades = raw
	}

	// Bibliografia
	res, err = client.R().Get("/template/bibliografia.php?id=" + data.ID)
	if err != nil {
		log.Println(err)
	} else {
		raw := res.String()
		if data.Paginas.Bibliografia != raw {
			body = append(body, fmt.Sprintf("la [bibliografia](<https://www.ing.unlp.edu.ar/catedras/%s/index.php?secc=bibliografia>)", data.ID))
		}
		data.Paginas.Bibliografia = raw
	}

	if len(body) >= 2 {
		return "Se actualizó " + strings.Join(body[:len(body)-1], ", ") + " y " + body[len(body)-1]
	} else if len(body) == 1 {
		return "Se actualizó " + body[0]
	}

	return ""
}
