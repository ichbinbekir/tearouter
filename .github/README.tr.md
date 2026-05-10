# TeaRouter

**Bubble Tea uygulamaları için GoRouter'dan ilham alan güçlü ve basit bir yönlendirici.**

TeaRouter, `bubbletea` TUI framework'ü ile geliştirilen karmaşık uygulamalarda sayfa (model) yönetimini ve navigasyonu kolaylaştırmak için tasarlanmıştır. Flutter'daki `gorouter` paketinin temel prensiplerini TUI dünyasına taşır.

## Özellikler

- **Hiyerarşik Alt Rotalar**: İç içe geçmiş rotalar tanımlayın ve navigasyon yığınını (stack) otomatik oluşturun.
- **Yığın Tabanlı Navigasyon**: `Push` ve `Pop` operasyonları ile sayfalar arasında kolayca geçiş yapın.
- **Bildirimsel Yönlendirme**: Rotalarınızı temiz ve okunabilir bir şekilde tanımlayın.
- **Durumu Sıfırlayarak Yönlendirme**: `Go` metodu ile navigasyon geçmişini temizleyerek yeni bir sayfaya gidin.
- **Sayfa Değiştirme**: `Replace` ile mevcut sayfayı yığından çıkarmadan yenisiyle değiştirin.
- **Middleware Desteği**: Rota geçişlerini yakalayarak kimlik doğrulama, loglama gibi ara katman işlemleri ekleyin.

## Kurulum

```bash
go get github.com/ichbinbekir/tearouter
```

## Hiyerarşik Alt Rotalar (Sub-Routing)

**GoRouter**'dan ilham alan TeaRouter, hiyerarşik rota tanımlarını destekler. `/main/settings/profile` gibi derin bir yola gittiğinizde, TeaRouter otomatik olarak tüm üst modelleri yığına ekler. Bu sayede doğal bir `Pop` davranışı (Profil -> Ayarlar -> Ana Sayfa) sağlanır.

```go
routes := []tearouter.Route{
    {
        Path: "/main",
        Builder: func() tea.Model { return MainModel{} },
        Children: []tearouter.Route{
            {
                Path: "settings", // Göreli yol: /main/settings
                Builder: func() tea.Model { return SettingsModel{} },
                Children: []tearouter.Route{
                    {
                        Path: "profile", // Göreli yol: /main/settings/profile
                        Builder: func() tea.Model { return ProfileModel{} },
                    },
                },
            },
        },
    },
}
```

`tearouter.Redirect(tearouter.Go, "/main/settings/profile")` çağrıldığında:
1. Yığın (stack) temizlenir.
2. `MainModel`, `SettingsModel` ve `ProfileModel` sırasıyla oluşturulup yığına eklenir.
3. Kullanıcı `ProfileModel` ekranını görür.
4. `Pop` yapıldığında kullanıcı doğal olarak `SettingsModel` ekranına geri döner.

## Hızlı Başlangıç

Aşağıda, sayfalar arasında geçiş yapan temel bir `tearouter` kullanımı gösterilmiştir.

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/tearouter"
)

// --- Sayfa Modellerimiz ---

type HomePageModel struct{}
func (m HomePageModel) Init() tea.Cmd { return nil }
func (m HomePageModel) View() string { return "Ana Sayfa\n\nAyarlar için 's', çıkış için 'q' basın." }
func (m HomePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "s": return m, tearouter.Redirect(tearouter.Push, "/settings")
		case "q": return m, tea.Quit
		}
	}
	return m, nil
}

type SettingsPageModel struct{}
func (m SettingsPageModel) Init() tea.Cmd { return nil }
func (m SettingsPageModel) View() string { return "Ayarlar Sayfası\n\nGeri dönmek için 'b' basın." }
func (m SettingsPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "b" {
		return m, tearouter.Redirect(tearouter.Pop)
	}
	return m, nil
}

func main() {
	routes := []tearouter.Route{
		{ Path: "/", Builder: func() tea.Model { return HomePageModel{} } },
		{ Path: "/settings", Builder: func() tea.Model { return SettingsPageModel{} } },
	}

	routerModel := tearouter.Model{
		InitialRoute: "/",
		Routes:       routes,
	}

	if _, err := tea.NewProgram(routerModel).Run(); err != nil {
		fmt.Printf("Hata: %v", err)
		os.Exit(1)
	}
}
```

## Navigasyon Metotları

Navigasyon, `tearouter.Redirect` komutu ile tetiklenir.

- `tearouter.Go`: Hedef yolun tüm hiyerarşisini oluşturur ve yeni yığın olarak ayarlar.
- `tearouter.Push`: Hedef yolun tüm hiyerarşisini mevcut yığının üzerine ekler.
- `tearouter.Replace`: Mevcut yığını hedef yolun tam hiyerarşisiyle değiştirir.
- `tearouter.Pop`: Yığının en üstündeki sayfayı kaldırır ve bir önceki sayfaya döner.

## Middleware Kullanımı

Middleware, kimlik doğrulama gibi işlemler için navigasyon isteklerini yakalamanıza olanak tanır.

```go
func authMiddleware(targetPath string) (newPath string) {
	if !isLoggedIn && targetPath != "/login" {
		return "/login"
	}
	return ""
}
```
