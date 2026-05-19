1. main.go

package main
2. 
3. import (
4. &#x09;"encoding/json"
5. &#x09;"fmt"
6. &#x09;"html/template"
7. &#x09;"log"
8. &#x09;"net/http"
9. &#x09;"time"
10. )
11. 
12. // Wymagane dane początkowe
13. const port = "8080"
14. const author = "Pawel Pastwa"
15. 
16. type WeatherResponse struct {
17. &#x09;CurrentWeather struct {
18. &#x09;	Temperature float64 `json:"temperature"`
19. &#x09;	WindSpeed   float64 `json:"windspeed"`
20. &#x09;} `json:"current\_weather"`
21. }
22. 
23. var cities = map\[string]string{
24. &#x09;"Polska, Warszawa":        "52.2297,21.0122",
25. &#x09;"Niemcy, Berlin":          "52.5200,13.4050",
26. &#x09;"Francja, Paryż":          "48.8566,2.3522",
27. &#x09;"Wielka Brytania, Londyn": "51.5074,-0.1278",
28. }
29. 
30. const tmplHTML = `
31. <!DOCTYPE html>
32. <html lang="pl">
33. <head>
34. &#x09;<meta charset="UTF-8">
35. &#x09;<title>Aktualna Pogoda</title>
36. &#x09;<style>
37. &#x09;	body { font-family: Arial, sans-serif; max-width: 600px; margin: 40px auto; padding: 20px; background-color: #f4f4f9; }
38. &#x09;	.container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
39. &#x09;	select, button { padding: 10px; margin-top: 10px; font-size: 16px; }
40. &#x09;	button { background-color: #0056b3; color: white; border: none; border-radius: 4px; cursor: pointer; }
41. &#x09;	button:hover { background-color: #004494; }
42. &#x09;</style>
43. </head>
44. <body>
45. &#x09;<div class="container">
46. &#x09;	<h1>Sprawdź aktualną pogodę</h1>
47. &#x09;	<form method="POST" action="/">
48. &#x09;		<label for="city">Wybierz lokalizację:</label><br>
49. &#x09;		<select name="city" id="city">
50. &#x09;			{{range $cityName, $coords := .Cities}}
51. &#x09;			<option value="{{$cityName}}">{{$cityName}}</option>
52. &#x09;			{{end}}
53. &#x09;		</select>
54. &#x09;		<br>
55. &#x09;		<button type="submit">Sprawdź pogodę</button>
56. &#x09;	</form>
57. 
58. &#x09;	{{if .SelectedCity}}
59. &#x09;	<hr>
60. &#x09;	<h2>Pogoda dla: {{.SelectedCity}}</h2>
61. &#x09;	<p><strong>Temperatura:</strong> {{.Temperature}} °C</p>
62. &#x09;	<p><strong>Prędkość wiatru:</strong> {{.WindSpeed}} km/h</p>
63. &#x09;	{{end}}
64. &#x09;</div>
65. </body>
66. </html>
67. `
68. 
69. func handler(w http.ResponseWriter, r \*http.Request) {
70. &#x09;tmpl := template.Must(template.New("index").Parse(tmplHTML))
71. &#x09;data := map\[string]interface{}{
72. &#x09;	"Cities": cities,
73. &#x09;}
74. 
75. &#x09;if r.Method == http.MethodPost {
76. &#x09;	r.ParseForm()
77. &#x09;	selectedCity := r.FormValue("city")
78. &#x09;	coords := cities\[selectedCity]
79. 
80. &#x09;	if coords != "" {
81. &#x09;		var lat, lon float64
82. &#x09;		fmt.Sscanf(coords, "%f,%f", \&lat, \&lon)
83. 
84. &#x09;		url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f\&longitude=%f\&current\_weather=true", lat, lon)
85. &#x09;		resp, err := http.Get(url)
86. &#x09;		if err == nil {
87. &#x09;			defer resp.Body.Close()
88. &#x09;			var weather WeatherResponse
89. &#x09;			json.NewDecoder(resp.Body).Decode(\&weather)
90. 
91. &#x09;			data\["SelectedCity"] = selectedCity
92. &#x09;			data\["Temperature"] = weather.CurrentWeather.Temperature
93. &#x09;			data\["WindSpeed"] = weather.CurrentWeather.WindSpeed
94. &#x09;		}
95. &#x09;	}
96. &#x09;}
97. &#x09;tmpl.Execute(w, data)
98. }
99. 
100. func healthHandler(w http.ResponseWriter, r \*http.Request) {
101. &#x09;w.WriteHeader(http.StatusOK)
102. &#x09;w.Write(\[]byte("OK"))
103. }
104. 
105. func main() {
106. &#x09;log.Printf("Data uruchomienia: %s, Autor: %s, Port TCP: %s\\n", time.Now().Format("2006-01-02 15:04:05"), author, port)
107. 
108. &#x09;http.HandleFunc("/", handler)
109. &#x09;http.HandleFunc("/health", healthHandler)
110. 
111. &#x09;log.Fatal(http.ListenAndServe(":"+port, nil))
112. }



2\. Dockerfile



FROM golang:1.22-alpine AS builder



WORKDIR /app



RUN go mod init weatherapp



COPY main.go .



RUN CGO\_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .



FROM alpine:3.19



LABEL org.opencontainers.image.authors="Pawel Pastwa"

LABEL org.opencontainers.image.title="Zadanie 1 - Aplikacja Pogodowa"



RUN apk --no-cache add ca-certificates curl



WORKDIR /root/



COPY --from=builder /app/app .



EXPOSE 8080



HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\

&#x20; CMD curl -f http://localhost:8080/health || exit 1



CMD \["./app"]



3\. Komendy



docker build -t mojpogodynka:v1 .

docker run -d -p 8080:8080 --name pogodynka\_kontener mojpogodynka:v1

docker logs pogodynka\_kontener

docker images mojpogodynka:v1

