# Sprawozdanie - Zadanie (CI/CD Pipeline)

**Autor:** Paweł Pastwa

## Konfiguracja etapów zadania

Łańcuch GitHub Actions został skonfigurowany tak, aby zautomatyzować proces budowy, skanowania i publikacji obrazu.

1. **Wsparcie dla wielu architektur:**
   Wykorzystano akcję `docker/setup-qemu-action` oraz silnik `Docker Buildx`, co pozwoliło na dodanie flagi `platforms: linux/amd64,linux/arm64` w kroku docelowego budowania.

2. **Zarządzanie Cache:**
   Utworzono dedykowane repozytorium `weather-cache` na DockerHub. Skonfigurowano eksport typu `registry` w trybie `mode=max` (eksportowane są wszystkie warstwy, również te z pośrednich etapów wieloetapowego pliku Dockerfile). Przyspiesza to znacznie kolejne przebiegi pipeline'u.

3. **Test CVE:**
   Wybrano skaner **Trivy** w formie gotowej akcji GitHub (`aquasecurity/trivy-action`). Jest on najlepszym wyborem, ponieważ integruje się bezpośrednio z przepływem GH Actions. Obraz najpierw budowany jest do lokalnego demona, skanowany pod kątem błędów `CRITICAL,HIGH` z ustawioną flagą `exit-code: 1`. Jeśli pojawią się podatności, pipeline zostaje przerwany na tym etapie i obraz multi-arch nie trafia do GHCR.

## Przyjęty sposób tagowania obrazów i danych cache (Uzasadnienie)

Zastosowano następującą strategię tagowania, opartą na rekomendacjach "Docker Best Practices" dotyczących środowisk CI/CD:

* **Dla obrazu docelowego (GHCR):** Obrazy są tagowane podwójnie:
    1. Z użyciem krótkiego hasha z gita (np. `sha-f2a8b9c`).
    2. Z tagiem `latest` (przypisanym do najnowszego udanego wdrożenia z gałęzi głównej).
    *Uzasadnienie:* Używanie commit SHA zapewnia niezmienność (immutability) i pozwala na natychmiastową identyfikację, z jakiej wersji kodu powstał konkretny obraz, co ułatwia ewentualny "rollback" w przypadku awarii. Tag `latest` stanowi z kolei wygodę dla deweloperów chcących pobrać najnowszą stabilną wersję. (*Źródło: [Docker Tagging Best Practices](https://docs.docker.com/develop/dev-best-practices/)*).
* **Dla danych cache (DockerHub):** Dane są wgrywane do oddzielnego repozytorium z pojedynczym, stałym tagiem `:buildcache`.
    *Uzasadnienie:* Stały tag zapewnia, że GitHub Actions zawsze odpytuje to samo zaktualizowane miejsce, a nadpisywanie starych wpisów cache zapobiega niekontrolowanemu rozrostowi zajmowanego miejsca w rejestrze. Tryb `max` gwarantuje wgranie do tagu kompletnego drzewa warstw.