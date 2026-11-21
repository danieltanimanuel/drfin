package renderer

import (
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Renderer struct {
	BasePath string
	cache    map[string]*template.Template
}

func NewRenderer(basePath string) *Renderer {
    return &Renderer{
		BasePath: basePath,
		cache:    make(map[string]*template.Template),
	}
}

const (
	LayoutDate = "2006-01-02"
)

func formatNumber(f *big.Float) string {
    if f == nil {
        return ""
    }
    f64, _ := f.Float64()

    // Format with 2 decimals, dot as decimal separator (temporarily)
    s := fmt.Sprintf("%.2f", f64) // e.g. "11000000.00"

    parts := strings.Split(s, ".")
    integerPart := parts[0] // "11000000"
    sign := ""
    if strings.HasPrefix(integerPart, "-") {
        sign = "-"
        integerPart = integerPart[1:]
    }
    decimalPart := parts[1] // "00"

    // Insert thousands separator '.' from right to left
    var intBuilder strings.Builder
    count := 0
    for i := len(integerPart) - 1; i >= 0; i-- {
        intBuilder.WriteByte(integerPart[i])
        count++
        if count%3 == 0 && i != 0 {
            intBuilder.WriteByte('.')
        }
    }

    // Reverse back
    intWithDots := reverseString(intBuilder.String())

    // Join with ',' as decimal separator
    return sign + intWithDots + "," + decimalPart
}

func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func formatDate(t time.Time) string {
    return t.Format(LayoutDate)
}

func (r *Renderer) Render(w http.ResponseWriter, templatePath string, templateName string, data any) error {
     // Load ALL templates in the folder
     pattern := filepath.Join(r.BasePath, templatePath, "*.html")
     tmpl := template.New("").Funcs(template.FuncMap{
        "formatNumber": formatNumber,
        "formatDate": formatDate,
    })
     tmpl, err := tmpl.ParseGlob(pattern)
     if err != nil {
         return err
     }
 
     // Render the requested template
     return tmpl.ExecuteTemplate(w, templateName, data)
}
