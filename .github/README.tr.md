# TeaRouter

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/tearouter)](https://goreportcard.com/report/github.com/your-username/tearouter)

**Bubble Tea uygulamaları için GoRouter'dan ilham alan güçlü ve basit bir yönlendirici.**

TeaRouter, `bubbletea` TUI framework'ü ile geliştirilen karmaşık uygulamalarda sayfa (model) yönetimini ve navigasyonu kolaylaştırmak için tasarlanmıştır. Flutter'daki `gorouter` paketinin temel prensiplerini TUI dünyasına taşır.

## Özellikler

- **Yığın Tabanlı Navigasyon**: `Push` ve `Pop` operasyonları ile sayfalar arasında kolayca geçiş yapın.
- **Bildirimsel Yönlendirme**: Rotalarınızı temiz ve okunabilir bir şekilde tanımlayın.
- **Durumu Sıfırlayarak Yönlendirme**: `Go` metodu ile navigasyon geçmişini temizleyerek yeni bir sayfaya gidin.
- **Sayfa Değiştirme**: `Replace` ile mevcut sayfayı yığından çıkarmadan yenisiyle değiştirin.
- **Middleware Desteği**: Rota geçişlerini yakalayarak kimlik doğrulama, loglama gibi ara katman işlemleri ekleyin.

## Kurulum

```bash
go get github.com/your-username/tearouter
```

## Hızlı Başlangıç

Aşağıda, iki sayfa (`home` ve `settings`) arasında geçiş yapan temel bir `tearouter` kullanımı gösterilmiştir.

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/your-username/tearouter"
)

// --- Sayfa Modellerimiz ---

// HomePageModel
type HomePageModel struct{}

func (m HomePageModel) Init() tea.Cmd { return nil }
func (m HomePageModel) View() string {
	return "Ana Sayfa

's' tuşuna basarak ayarlar sayfasına git.
'q' ile çıkış yap."
}
func (m HomePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			// Ayarlar sayfasını yığının üzerine ekle
			return m, tearouter.Redirect(tearouter.Push, "/settings")
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// SettingsPageModel
type SettingsPageModel struct{}

func (m SettingsPageModel) Init() tea.Cmd { return nil }
func (m SettingsPageModel) View() string {
	return "Ayarlar Sayfası

'b' tuşuna basarak geri dön."
}
func (m SettingsPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			// Yığından bir önceki sayfaya dön
			return m, tearouter.Redirect(tearouter.Pop)
		}
	}
	return m, nil
}

// --- Ana Uygulama ---

func main() {
	// Rotaları tanımla
	routes := []tearouter.Route{
		{
			Path:    "/",
			Builder: func() tea.Model { return HomePageModel{} },
		},
		{
			Path:    "/settings",
			Builder: func() tea.Model { return SettingsPageModel{} },
		},
	}

	// Router modelini oluştur
	routerModel := tearouter.Model{
		InitialRoute: "/",
		Routes:       routes,
	}

	p := tea.NewProgram(routerModel)
	if err := p.Start(); err != nil {
		fmt.Printf("Bir hata oluştu: %v", err)
		os.Exit(1)
	}
}
```

## Navigasyon Metotları

Navigasyon, `tearouter.Redirect` komutu ile tetiklenir.

- `tearouter.Go`: Navigasyon yığınını temizler ve belirtilen hedefe yönlendirir. Geri dönmek mümkün değildir. Genellikle login sonrası ana sayfaya yönlendirme gibi durumlar için kullanılır.
  ```go
  return m, tearouter.Redirect(tearouter.Go, "/home")
  ```

- `tearouter.Push`: Yeni bir sayfayı mevcut yığının üzerine ekler. Kullanıcı geri dönebilir.
  ```go
  return m, tearouter.Redirect(tearouter.Push, "/profile")
  ```

- `tearouter.Replace`: Yığındaki mevcut (en üstteki) sayfayı yenisiyle değiştirir. Yığının boyutu değişmez.
  ```go
  return m, tearouter.Redirect(tearouter.Replace, "/profile/edit")
  ```

- `tearouter.Pop`: Yığının en üstündeki sayfayı kaldırır ve bir önceki sayfaya döner. Eğer yığında tek bir sayfa varsa hata döner.
  ```go
  return m, tearouter.Redirect(tearouter.Pop)
  ```

## Middleware Kullanımı

Middleware, her yönlendirme talebini işleyen bir fonksiyondur. Loglama yapabilir veya kullanıcının yetkisi olmayan bir sayfaya gitmesini engelleyebilirsiniz.

Eğer middleware `""` (boş string) dönerse, navigasyon normal şekilde devam eder. Eğer yeni bir path dönerse, kullanıcı o path'e yönlendirilir.

```go
// Örnek: Auth Middleware
var isAuthenticated = false

func authMiddleware(targetPath string) (newPath string) {
    // Login sayfasına veya login sayfasından yapılan yönlendirmelere dokunma
    if targetPath == "/login" || !isAuthenticated {
        return "" // Devam et
    }

    if !isAuthenticated {
        // Eğer kullanıcı giriş yapmamışsa ve korumalı bir sayfaya gitmeye çalışıyorsa
        // onu login sayfasına yönlendir.
        return "/login"
    }

    return "" // Giriş yapmış, devam et
}

func main() {
    routerModel := tearouter.Model{
        // ...
        Middleware: authMiddleware,
    }
    // ...
}
```

## Lisans

Bu proje MIT Lisansı altındadır. Detaylar için `LICENSE` dosyasına bakınız.
