package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Wymagane dane początkowe
const port = "8080"
const author = "Pawel Pastwa"

type WeatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		WindSpeed   float64 `json:"windspeed"`
	} `json:"current_weather"`
}

var cities = map[string]string{
	"Polska, Warszawa":        "52.2297,21.0122",
	"Niemcy, Berlin":          "52.5200,13.4050",
	"Francja, Paryż":          "48.8566,2.3522",
	"Wielka Brytania, Londyn": "51.5074,-0.1278",
}

const tmplHTML = `
<!DOCTYPE html>
<html lang="pl">
<head>
	<meta charset="UTF-8">
	<title>Aktualna Pogoda</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 600px; margin: 40px auto; padding: 20px; background-color: #f4f4f9; }
		.container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
		select, button { padding: 10px; margin-top: 10px; font-size: 16px; }
		button { background-color: #0056b3; color: white; border: none; border-radius: 4px; cursor: pointer; }
		button:hover { background-color: #004494; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Sprawdź aktualną pogodę</h1>
		<form method="POST" action="/">
			<label for="city">Wybierz lokalizację:</label><br>
			<select name="city" id="city">
				{{range $cityName, $coords := .Cities}}
				<option value="{{$cityName}}">{{$cityName}}</option>
				{{end}}
			</select>
			<br>
			<button type="submit">Sprawdź pogodę</button>
		</form>

		{{if .SelectedCity}}
		<hr>
		<h2>Pogoda dla: {{.SelectedCity}}</h2>
		<p><strong>Temperatura:</strong> {{.Temperature}} °C</p>
		<p><strong>Prędkość wiatru:</strong> {{.WindSpeed}} km/h</p>
		{{end}}
	</div>
</body>
</html>
`

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(tmplHTML))
	data := map[string]interface{}{
		"Cities": cities,
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		selectedCity := r.FormValue("city")
		coords := cities[selectedCity]

		if coords != "" {
			var lat, lon float64
			fmt.Sscanf(coords, "%f,%f", &lat, &lon)

			url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true", lat, lon)
			resp, err := http.Get(url)
			if err == nil {
				defer resp.Body.Close()
				var weather WeatherResponse
				json.NewDecoder(resp.Body).Decode(&weather)

				data["SelectedCity"] = selectedCity
				data["Temperature"] = weather.CurrentWeather.Temperature
				data["WindSpeed"] = weather.CurrentWeather.WindSpeed
			}
		}
	}
	tmpl.Execute(w, data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	log.Printf("Data uruchomienia: %s, Autor: %s, Port TCP: %s\n", time.Now().Format("2006-01-02 15:04:05"), author, port)

	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}