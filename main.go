package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {


	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run main.go <URL>")
		return
	}

	url := os.Args[1]
	fmt.Println("[+] Hedef URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("[-] Bağlantı hatası:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("[-] HTTP Hata Kodu:", resp.Status)
		return
	}

	htmlFile, err := os.Create("output.html")
	if err != nil {
		fmt.Println("[-] Dosya oluşturulamadı:", err)
		return
	}
	defer htmlFile.Close()

	_, err = io.Copy(htmlFile, resp.Body)
	if err != nil {
		fmt.Println("[-] HTML yazılamadı:", err)
		return
	}

	fmt.Println("[+] HTML kaydedildi: output.html")

	htmlBytes, err := os.ReadFile("output.html")
	if err == nil {

		re := regexp.MustCompile(`href="(http[s]?://[^"]+)"`)
		matches := re.FindAllSubmatch(htmlBytes, -1)

		linkFile, _ := os.Create("links.txt")
		defer linkFile.Close()

		for i, m := range matches {
			linkFile.WriteString(fmt.Sprintf("%d. %s\n", i+1, m[1]))
		}

		fmt.Printf("[+] %d adet link links.txt dosyasına kaydedildi\n", len(matches))
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&buf, 90),
	)

	if err != nil {
		fmt.Println("[-] Screenshot alınamadı:", err)
		return
	}

	err = os.WriteFile("screenshot.png", buf, 0644)
	if err != nil {
		fmt.Println("[-] Screenshot kaydedilemedi:", err)
		return
	}

	fmt.Println("[+] Screenshot alındı: screenshot.png")

	fmt.Println("\n[+] Tüm görevler başarıyla tamamlandı.")
}
