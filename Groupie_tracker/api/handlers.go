package api

import (
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"groupie-tracker/config"
	"groupie-tracker/model"
)

func AppMux() http.Handler {
	mux := http.NewServeMux()                       // мультиплексор HTTP-запросов. Он сопоставляет URL каждого входящего запроса со списком зарегистрированных шаблонов и вызывает обработчик шаблона, который наиболее точно соответствует URL
	fs := http.FileServer(http.Dir("ui"))           // для статических файлов
	mux.Handle("/ui/", http.StripPrefix("/ui", fs)) // fs возвращает обработчик, который обслуживает HTTP-запросы с содержимым файловой системы

	mux.HandleFunc("/", HomePageHandler)
	mux.HandleFunc("/artists/", ArtistPageHandler)
	mux.HandleFunc("/search", SearchHandler)
	return mux
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	searchinput := r.URL.Query().Get("searchinput") // значения в URL, которые идут после searchinput записываются в эту переменную

	artists := new(model.ArtistsList) // инициализировали структуру из пакета model -> ArtistList

	if err := model.GetAllAtritst(&artists.List); err != nil { // при событии search парсим всю инфу и отправляем в созданный artist.List
		ExecuteErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	}
	relations := new(model.RelationList) // инициализировали структуру из пакета model -> RelationList

	if err := model.GetAllRelations(&relations.List); err != nil { // при событии search парсим всю инфу и отправляем в созданный relations.List
		ExecuteErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	}

	artists.AddDatesLocations(relations.List["index"]) // по ключу index, заполняем данные artist -> List -> Artist -> DatesLocations
	relations = nil                                    // очистили realtions, так как далее не используем

	searchResultArtists := new(model.ArtistsList)    // инициализировали структуру из пакета model -> ArtistList
	artists.ChangeKeys()                             // форматируем данные артистов, убираем и меняем символы в DatesLocations
	searchResultArtists.SeachDataList = artists.List // заполняем распарсиными данным из строчки 32

	if errStatus := CheckSearchRequest(r); errStatus != 0 { // проверяем URL запроса, который прилетает к нам на сервер
		ExecuteSearchErrorPage(w, errStatus, *searchResultArtists)
		errorLog.Printf(http.StatusText(errStatus))
		return
	}
	searchResultArtists.SearchInputInArtistsList(artists.List, searchinput) // заполняем данным
	if len(searchResultArtists.List) == 0 {
		ExecuteSearchErrorPage(w, http.StatusNotFound, *searchResultArtists)
		errorLog.Printf(http.StatusText(http.StatusNotFound))
		return
	}

	if tmpl, err := template.ParseFiles(config.IndexTmplPath); err != nil {
		ExecuteErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	} else {
		if err := tmpl.Execute(w, searchResultArtists); err != nil {
			ExecuteErrorPage(w, http.StatusInternalServerError)
			errorLog.Printf(err.Error())
		}
	}
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	if errStatus := CheckIndexRequest(r); errStatus != 0 {
		ExecuteErrorPage(w, errStatus)
		errorLog.Printf(http.StatusText(errStatus))
		return
	}

	artists := new(model.ArtistsList)

	if err := model.GetAllAtritst(&artists.List); err != nil {
		ExecuteErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	}
	relations := new(model.RelationList)
	if err := model.GetAllRelations(&relations.List); err != nil {
		ExecuteErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	}

	artists.AddDatesLocations(relations.List["index"])
	relations = nil
	artists.ChangeKeys()
	artists.SeachDataList = artists.List

	if tmpl, err := template.ParseFiles(config.IndexTmplPath); err != nil {
		ExecuteErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	} else {
		if err := tmpl.Execute(w, artists); err != nil {
			ExecuteErrorPage(w, http.StatusInternalServerError)
			errorLog.Printf(err.Error())
		}
	}
}

func ArtistPageHandler(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if errStatus := CheckArtistRequest(r); errStatus != 0 {
		ExecuteErrorPage(w, errStatus)
		errorLog.Printf(http.StatusText(errStatus))
		return
	}
	artId := strings.TrimPrefix(r.URL.Path, "/artists/")
	artist := new(model.Artist)

	if err := model.GetAtritst(artId, artist); err != nil {
		ExecuteErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	}

	if artist.Id == 0 {
		ExecuteErrorPage(w, http.StatusNotFound)
		errorLog.Printf(http.StatusText(http.StatusNotFound))
		return
	}

	if err := model.GetRelation(artId, artist); err != nil {
		ExecuteErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	}
	artist.ChangeKey()

	if tmpl, err := template.ParseFiles(config.ArtistTmplPath); err != nil {
		ExecuteErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	} else {
		if err := tmpl.Execute(w, artist); err != nil {
			ExecuteErrorPage(w, http.StatusInternalServerError)
			errorLog.Printf(err.Error())
		}
	}
}

func ExecuteErrorPage(w http.ResponseWriter, errStatus int) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tmpl, err := template.ParseFiles(config.ErrorTmplPath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		errorLog.Printf(err.Error())
		return
	}
	w.WriteHeader(errStatus)
	if err := tmpl.Execute(w, http.StatusText(errStatus)); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		errorLog.Printf(err.Error())
		return
	}
}

func ExecuteSearchErrorPage(w http.ResponseWriter, errStatus int, artists model.ArtistsList) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	artists.ErrorText = http.StatusText(errStatus)
	tmpl, err := template.ParseFiles(config.SearchErrorTmplPath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		errorLog.Printf(err.Error())
		return
	}
	w.WriteHeader(errStatus)
	if err := tmpl.Execute(w, artists); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		errorLog.Printf(err.Error())
		return
	}
}
