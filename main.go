package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var bellek = make(map[string]interface{})

type OyunDizayn struct {
	Kod    string
	Sonuc  string
	JSkodu template.JS
}

// v8.0 interaktif ve sistem sabitlerini JS tarafına bağlar
func tokenDegeriniCozStr(t Token) string {
	if t.Type == INT {
		return t.Literal
	}
	if t.Type == IDENT {
		switch t.Literal {
		case "kumanda_tus":
			return "sistem_girdi.sonTus"
		case "kumanda_x":
			return "sistem_girdi.fareX"
		case "kumanda_y":
			return "sistem_girdi.fareY"
		case "kumanda_tiklama":
			return "sistem_girdi.fareTiklandi"
		case "sans_konum":
			return "Math.floor(Math.random() * 360)"
		case "sans_boyut":
			return "Math.floor(Math.random() * 60) + 15"
		case "sans_zar":
			return "Math.floor(Math.random() * 6) + 1"
		case "matematik_pi":
			return "3"
		}
		return t.Literal
	}
	return "0"
}

// v6.0 uyumluluğu için Go tarafındaki statik hafıza çözümleyicisi
func tokenDegeriniCozGo(t Token) int {
	if t.Type == INT {
		val, _ := strconv.Atoi(t.Literal)
		return val
	}
	if t.Type == IDENT {
		if val, ok := bellek[t.Literal].(int); ok {
			return val
		}
	}
	return 0
}

func KoduYorumla(tumKod string) (string, string) {
	bellek = make(map[string]interface{})
	var cikisSatirlari []string
	var jsKomutlari []string

	ekranGenislik := 400
	ekranYukseklik := 400
	aktifRenk := "red"

	// Akıllı Karar Mekanizması: Kod interaktif v8.0 mi yoksa klasik v6.0 mı?
	isRealTime := strings.Contains(tumKod, "kumanda_") || strings.Contains(tumKod, "sans_")

	lexer := NewLexer(tumKod)
	var tokens []Token
	for {
		tok := lexer.NextToken()
		if tok.Type == EOF {
			break
		}
		tokens = append(tokens, tok)
	}

	i := 0
	limit := 0

	for i < len(tokens) && limit < 5000 {
		limit++
		tok := tokens[i]

		switch tok.Type {
		case OLSUM:
			i++
			degiskenAdi := tokens[i].Literal
			i++ // =
			i++ // değer
			
			if isRealTime {
				v1Str := tokenDegeriniCozStr(tokens[i])
				if i+1 < len(tokens) && (tokens[i+1].Type == PLUS || tokens[i+1].Type == MINUS || tokens[i+1].Type == ASTERISK || tokens[i+1].Type == SLASH) {
					operator := tokens[i+1].Type
					i += 2
					v2Str := tokenDegeriniCozStr(tokens[i])
					var opSign string
					switch operator {
					case PLUS: opSign = "+"
					case MINUS: opSign = "-"
					case ASTERISK: opSign = "*"
					case SLASH: opSign = "/"
					}
					jsKomutlari = append(jsKomutlari, fmt.Sprintf("%s = %s %s %s;", degiskenAdi, v1Str, opSign, v2Str))
				} else {
					if tokens[i].Type == STRING {
						jsKomutlari = append(jsKomutlari, fmt.Sprintf("%s = \"%s\";", degiskenAdi, tokens[i].Literal))
					} else {
						jsKomutlari = append(jsKomutlari, fmt.Sprintf("%s = %s;", degiskenAdi, v1Str))
					}
				}
			} else {
				v1 := tokenDegeriniCozGo(tokens[i])
				if i+1 < len(tokens) && (tokens[i+1].Type == PLUS || tokens[i+1].Type == MINUS || tokens[i+1].Type == ASTERISK || tokens[i+1].Type == SLASH) {
					operator := tokens[i+1].Type
					i += 2
					v2 := tokenDegeriniCozGo(tokens[i])
					switch operator {
					case PLUS: bellek[degiskenAdi] = v1 + v2
					case MINUS: bellek[degiskenAdi] = v1 - v2
					case ASTERISK: bellek[degiskenAdi] = v1 * v2
					case SLASH:
						if v2 != 0 { bellek[degiskenAdi] = v1 / v2 } else { bellek[degiskenAdi] = 0 }
					}
				} else {
					if tokens[i].Type == STRING {
						bellek[degiskenAdi] = tokens[i].Literal
					} else {
						bellek[degiskenAdi] = v1
					}
				}
			}
			i++

		case YAZDIR:
			i += 2
			icerikToken := tokens[i]
			if isRealTime {
				if icerikToken.Type == IDENT {
					jsKomutlari = append(jsKomutlari, fmt.Sprintf("console.log(%s);", icerikToken.Literal))
				} else {
					jsKomutlari = append(jsKomutlari, fmt.Sprintf("console.log(\"%s\");", icerikToken.Literal))
				}
			} else {
				if icerikToken.Type == IDENT {
					if val, ok := bellek[icerikToken.Literal]; ok {
						cikisSatirlari = append(cikisSatirlari, fmt.Sprintf("%v", val))
					} else {
						cikisSatirlari = append(cikisSatirlari, "Hata: Değişken bulunamadı")
					}
				} else {
					cikisSatirlari = append(cikisSatirlari, icerikToken.Literal)
				}
			}
			i += 2

		case EKRAN:
			i += 2
			ekranGenislik, _ = strconv.Atoi(tokens[i].Literal)
			i += 2
			ekranYukseklik, _ = strconv.Atoi(tokens[i].Literal)
			i += 2

		case RENK:
			i += 2
			aktifRenk = tokens[i].Literal
			i += 2

		case KARE:
			i += 2
			if isRealTime {
				x := tokenDegeriniCozStr(tokens[i]); i += 2
				y := tokenDegeriniCozStr(tokens[i]); i += 2
				b := tokenDegeriniCozStr(tokens[i])
				jsKomutlari = append(jsKomutlari, fmt.Sprintf("ctx.fillStyle = '%s'; ctx.fillRect(%s, %s, %s, %s);", aktifRenk, x, y, b, b))
			} else {
				x := tokenDegeriniCozGo(tokens[i]); i += 2
				y := tokenDegeriniCozGo(tokens[i]); i += 2
				b := tokenDegeriniCozGo(tokens[i])
				jsKomutlari = append(jsKomutlari, fmt.Sprintf("ctx.fillStyle = '%s'; ctx.fillRect(%d, %d, %d, %d);", aktifRenk, x, y, b, b))
			}
			i += 2

		case DAIRE:
			i += 2
			if isRealTime {
				x := tokenDegeriniCozStr(tokens[i]); i += 2
				y := tokenDegeriniCozStr(tokens[i]); i += 2
				r := tokenDegeriniCozStr(tokens[i])
				jsKomutlari = append(jsKomutlari, fmt.Sprintf("ctx.fillStyle = '%s'; ctx.beginPath(); ctx.arc(%s, %s, %s, 0, 2 * Math.PI); ctx.fill();", aktifRenk, x, y, r))
			} else {
				x := tokenDegeriniCozGo(tokens[i]); i += 2
				y := tokenDegeriniCozGo(tokens[i]); i += 2
				r := tokenDegeriniCozGo(tokens[i])
				jsKomutlari = append(jsKomutlari, fmt.Sprintf("ctx.fillStyle = '%s'; ctx.beginPath(); ctx.arc(%d, %d, %d, 0, 2 * Math.PI); ctx.fill();", aktifRenk, x, y, r))
			}
			i += 2

		case EGER:
			i += 2
			if isRealTime {
				v1 := tokenDegeriniCozStr(tokens[i]); i++
				op := tokens[i].Type; i++
				v2 := tokenDegeriniCozStr(tokens[i]); i += 2
				var opSign string
				if op == GT { opSign = ">" } else { opSign = "<" }
				jsKomutlari = append(jsKomutlari, fmt.Sprintf("if (%s %s %s) {", v1, opSign, v2))
			} else {
				v1 := tokenDegeriniCozGo(tokens[i]); i++
				op := tokens[i].Type; i++
				v2 := tokenDegeriniCozGo(tokens[i]); i += 2
				sartSaglandi := false
				if op == GT && v1 > v2 { sartSaglandi = true }
				if op == LT && v1 < v2 { sartSaglandi = true }
				if !sartSaglandi {
					pSkor := 1
					for pSkor > 0 && i < len(tokens) {
						if tokens[i].Type == LBRACE { pSkor++ }
						if tokens[i].Type == RBRACE { pSkor-- }
						i++
					}
				}
			}
			i++

		case DONGU:
			i += 2
			if isRealTime {
				v1 := tokenDegeriniCozStr(tokens[i]); i++
				op := tokens[i].Type; i++
				v2 := tokenDegeriniCozStr(tokens[i]); i += 2
				var opSign string
				if op == GT { opSign = ">" } else { opSign = "<" }
				jsKomutlari = append(jsKomutlari, fmt.Sprintf("while (%s %s %s) {", v1, opSign, v2))
			} else {
				v1 := tokenDegeriniCozGo(tokens[i]); i++
				op := tokens[i].Type; i++
				v2 := tokenDegeriniCozGo(tokens[i]); i += 2
				donguDevam := false
				if op == GT && v1 > v2 { donguDevam = true }
				if op == LT && v1 < v2 { donguDevam = true }
				if !donguDevam {
					pSkor := 1
					for pSkor > 0 && i < len(tokens) {
						if tokens[i].Type == LBRACE { pSkor++ }
						if tokens[i].Type == RBRACE { pSkor-- }
						i++
					}
				}
			}
			i++

		case RBRACE:
			if isRealTime {
				jsKomutlari = append(jsKomutlari, "}")
			} else {
				pSkor := 1
				for j := i - 1; j >= 0; j-- {
					if tokens[j].Type == RBRACE { pSkor++ }
					if tokens[j].Type == LBRACE {
						pSkor--
						if pSkor == 0 && j-5 >= 0 && tokens[j-5].Type == DONGU {
							i = j - 5
							break
						}
					}
				}
			}
			i++

		default:
			i++
		}
	}

	baslangicJS := fmt.Sprintf("canvas.width = %d; canvas.height = %d;\n", ekranGenislik, ekranYukseklik)
	loopBody := strings.Join(jsKomutlari, "\n")
	
	var tamJSkodu string

	if isRealTime {
		tamJSkodu = baslangicJS + fmt.Sprintf(`
			if(!window.sistem_girdi_kayitli) {
				window.sistem_girdi = { sonTus: 200, fareX: 0, fareY: 0, fareTiklandi: 0 };
				window.addEventListener('keydown', (e) => {
					if(e.key === 'ArrowRight' || e.key === 'd') window.sistem_girdi.sonTus += 12;
					if(e.key === 'ArrowLeft' || e.key === 'a') window.sistem_girdi.sonTus -= 12;
					if(['ArrowUp','ArrowDown','ArrowLeft','ArrowRight'].includes(e.key)) e.preventDefault();
				});
				window.addEventListener('mousemove', (e) => {
					const rect = canvas.getBoundingClientRect();
					window.sistem_girdi.fareX = e.clientX - rect.left;
					window.sistem_girdi.fareY = e.clientY - rect.top;
				});
				window.addEventListener('mousedown', () => { window.sistem_girdi.fareTiklandi = 1; });
				window.addEventListener('mouseup', () => { window.sistem_girdi.fareTiklandi = 0; });
				window.sistem_girdi_kayitli = true;
			}

			function render() {
				ctx.clearRect(0, 0, canvas.width, canvas.height);
				try { oyunDongusu(); } catch(err) {}
				requestAnimationFrame(render);
			}

			function oyunDongusu() {
				%s
			}
			requestAnimationFrame(render);
		`, loopBody)
		if len(cikisSatirlari) == 0 {
			cikisSatirlari = append(cikisSatirlari, "[Sistem - Canlı Mod]: .lvn dosyası 60 FPS motoruyla akıcı olarak yürütülüyor.")
		}
	} else {
		tamJSkodu = baslangicJS + fmt.Sprintf("ctx.clearRect(0,0,%d,%d); \n %s", ekranGenislik, ekranYukseklik, loopBody)
		if len(cikisSatirlari) == 0 {
			cikisSatirlari = append(cikisSatirlari, "[Sistem - Statik Mod]: Klasik .lvn komut dizisi başarıyla işlendi.")
		}
	}

	return strings.Join(cikisSatirlari, "\n"), tamJSkodu
}

const htmlSayfasi = `
<!DOCTYPE html>
<html>
<head>
    <title>Levent Studio IDE v8.0</title>
    <link rel="icon" type="image/svg+xml" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><defs><linearGradient id='g' x1='0%' y1='0%' x2='100%' y2='100%'><stop offset='0%' stop-color='%2338bdf8'/><stop offset='100%' stop-color='%23a855f7'/></linearGradient></defs><rect width='100' height='100' rx='22' fill='%230b0f19'/><path d='M35 25 v40 h30 v-8 h-20 v-32 z' fill='url(%23g)'/><circle cx='68' cy='28' r='7' fill='%2334d399'/></svg>">
    <style>
        body { font-family: 'Segoe UI', sans-serif; background: #0f172a; color: #fff; padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; display: flex; gap: 24px; }
        .sol-panel { flex: 1; }
        .sag-panel { width: 440px; background: #1e293b; padding: 20px; border-radius: 12px; border: 2px solid #a855f7; text-align: center; box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3); }
        
        /* v8.0 Kurumsal Dil Logosu Alanı */
        .logo-container { display: flex; align-items: center; justify-content: center; gap: 14px; margin-bottom: 15px; }
        .lvn-logo { width: 48px; height: 48px; filter: drop-shadow(0 0 8px #a855f7); }
        .baslik { font-size: 24px; font-weight: 800; background: linear-gradient(to right, #38bdf8, #a855f7); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
        
        textarea { width: 100%; height: 350px; background: #0b0f19; color: #f8fafc; font-family: 'Courier New', monospace; font-size: 16px; padding: 15px; border: 2px solid #a855f7; border-radius: 12px; box-sizing: border-box; resize: none; line-height: 1.5; }
        button { background: linear-gradient(135deg, #38bdf8, #a855f7); color: #fff; border: none; padding: 14px 20px; font-size: 18px; font-weight: bold; cursor: pointer; width: 100%; margin-top: 10px; border-radius: 12px; transition: 0.2s; box-shadow: 0 4px 6px -1px rgba(168, 85, 247, 0.2); }
        button:hover { transform: scale(1.01); box-shadow: 0 0 20px #a855f7; }
        .output { background: #000; border: 1px solid #1e293b; padding: 15px; margin-top: 15px; min-height: 80px; border-radius: 12px; font-family: monospace; color: #34d399; font-size: 15px; white-space: pre-wrap; }
        canvas { background: #000; border: 4px solid #fff; margin-top: 15px; border-radius: 8px; outline: none; }
        .rehber { text-align: left; background: #0b0f19; padding: 12px; margin-top: 15px; border-radius: 8px; font-size: 13px; color: #cbd5e1; border-left: 4px solid #38bdf8; line-height: 1.5; }
        .dosya-bilgi { font-size: 12px; color: #64748b; font-weight: 600; text-transform: uppercase; margin-bottom: 5px; letter-spacing: 1px; display: block; text-align: left; }
    </style>
</head>
<body>
    <div class="container">
        <div class="sol-panel">
            <div class="logo-container">
                <svg class="lvn-logo" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
                    <defs>
                        <linearGradient id='logoGrad' x1='0%' y1='0%' x2='100%' y2='100%'>
                            <stop offset='0%' stop-color='#38bdf8'/>
                            <stop offset='100%' stop-color='#a855f7'/>
                        </linearGradient>
                    </defs>
                    <rect width='100' height='100' rx='22' fill='#1e293b' stroke='#a855f7' stroke-width='2'/>
                    <path d='M38 28 v42 h26 v-8 h-16 v-34 z' fill='url(#logoGrad)'/>
                    <circle cx='68' cy='32' r='6' fill='#34d399'/>
                </svg>
                <div class="baslik">LEVENT STUDIO IDE v8.0</div>
            </div>
            
            <form method="POST">
                <span class="dosya-bilgi">📁 ÇALIŞMA DOSYASI: kaynak_kod.lvn</span>
                <textarea name="kod" placeholder="# .lvn kodunuzu buraya yazın veya yapıştırın..." required>{{.Kod}}</textarea>
                <button type="submit">LVN MOTORUNU ATEŞLE ▶</button>
            </form>
            <h3>Sistem Konsolu:</h3>
            <div class="output">{{.Sonuc}}</div>
        </div>
        
        <div class="sag-panel">
            <h2>🕹️ CANLI ÇIKTI EKRANI</h2>
            <canvas id="oyunTuvali" width="400" height="400" tabindex="1"></canvas>
            <div class="rehber">
                <b>📂 .lvn Dosya Logosu & Sistemi:</b><br>
                Arayüzün sol üst köşesindeki ikon ve tarayıcı sekmesindeki favicon, dilimizin resmi <b>.lvn</b> logosudur! <br><br>
                <b>💡 Akıllı Hibrit Derleme:</b><br>
                • Eski v6.0 algoritmik testlerinizi aynen çalıştırır.<br>
                • <code>kumanda_tus</code> veya <code>sans_konum</code> yazınca otomatik olarak etkileşimli oyun moduna geçer.
            </div>
        </div>
    </div>

    <script>
        const canvas = document.getElementById('oyunTuvali');
        const ctx = canvas.getContext('2d');
        {{.JSkodu}}
    </script>
</body>
</html>
`

func main() {
	tmpl := template.Must(template.New("full_ide").Parse(htmlSayfasi))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		veri := OyunDizayn{}
		if r.Method == http.MethodPost {
			veri.Kod = r.FormValue("kod")
			sonuc, jskodu := KoduYorumla(veri.Kod)
			veri.Sonuc = sonuc
			veri.JSkodu = template.JS(jskodu)
		}
		tmpl.Execute(w, veri)
	})

	fmt.Println("=================================================")
	fmt.Println("🌟 LEVENT ENGINE v8.0 EVRENSEL SÜRÜM AKTİF!")
	fmt.Println("📂 Özel Dosya Türü: .lvn (Levent Kod Dosyası)")
	fmt.Println("🎨 Entegre Dil Logosu: Aktif (Sekme ve Panel)")
	fmt.Println("👉 Geliştirici Adresi: http://localhost:8080")
	fmt.Println("=================================================")

	http.ListenAndServe(":8080", nil)
}