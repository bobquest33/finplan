package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"os"
	"strings"
	"time"
)

type timeframe int

const (
	monthly timeframe = iota + 1
	annual
)

type MonthTime struct {
    time.Time
}

const mtLayout = "Jan 06"

func (mt *MonthTime) UnmarshalJSON(b []byte) (err error) {
    if b[0] == '"' && b[len(b)-1] == '"' {
        b = b[1 : len(b)-1]
    }
    mt.Time, err = time.Parse(mtLayout, string(b))
    return
}

func (mt *MonthTime) MarshalJSON() ([]byte, error) {
    return []byte(mt.Time.Format(mtLayout)), nil
}

var nilTime = (time.Time{}).UnixNano()
func (mt *MonthTime) IsSet() bool {
    return mt.UnixNano() != nilTime
}

func parseMonth(month string) MonthTime {
	t, err := time.Parse(mtLayout, month)
	panicOn(err)
	return MonthTime{t}
}

func monthString(month MonthTime) string {
	return month.Time.Format(mtLayout)
}


func main() {

	revexpFiles := flag.String("revexps", "", "list of comma-separated filenames containing revexps")
		ccFiles := flag.String("ccs", "", "list of comma-separated filenames containing ccs")
	flag.Parse()

	allRevexps := make(revexps)
	allCCs := make(ccs)

	var revexps revexps
	for _, revexpFile := range strings.Split(*revexpFiles, ",") {
		f, err := os.Open(revexpFile)
		panicOn(err)
		err = json.NewDecoder(f).Decode(&revexps)
		panicOn(err)
		for k, v := range revexps {
			allRevexps[k] = v
		}
	}
	
	ccs := make([]*cc, 0)
	for _, ccFile := range strings.Split(*ccFiles, ",") {
		f, err := os.Open(ccFile)
		panicOn(err)		
		err = json.NewDecoder(f).Decode(&ccs)
		for _, cc := range ccs {
			allCCs[cc.Name] = cc
		}
	}

	render(project(allRevexps, allCCs, 36))

}

func render(months []month) {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Projections</title>
		<style>
			th, td {
				padding: 10px;
			}
		</style>
	</head>
	<body>
		<table>
			<thead>
				<tr>
					<th rowspan="2">month 
					<th rowspan="2">netrev 
					{{range (index . 0).CCs}}
						<th colspan="3">{{.CC.Name}}
					{{end}}					
					<th  rowspan="2">total ri 
				</tr>
				<tr>
					{{range (index . 0).CCs}}
							<th>payment 
							<th>balance 
							<th>ri
						
					{{end}}	
				</tr>	
			</thead>
			{{range .}}
				<tr>
					<td>{{ .Format }}
					<td>{{ .Netrev }}
					{{range .CCs}}
						<td>{{ .Payment }}
						<td>{{ .Balance }}
						<td>{{ .RI }}
					{{end}}
					<td>{{ .RI }}				
				</tr>
			{{end}}
	</body>
</html>`
	t, err := template.New("webpage").Parse(tpl)
	panicOn(err)

	err = t.Execute(os.Stdout, months)
	panicOn(err)
}

// normalize sets time to the 20th of the current month
func normalize(t time.Time) MonthTime {
	return MonthTime{time.Date(t.Year(), t.Month(), 20, 0, 0, 0, 0, time.UTC)}
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
