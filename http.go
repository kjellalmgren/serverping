package gossm

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func calculateServerUptime(statusAtTime []*statusAtTime) string {
	if len(statusAtTime) == 0 {
		return "unknown"
	}

	var sum float64

	for _, val := range statusAtTime {
		var i float64
		if val.Status {
			i = 1
		} else {
			i = 0
		}
		sum += i
	}

	return fmt.Sprintf("%.2f", sum/float64(len(statusAtTime))*100)
}

func lastStatus(statusAtTime []*statusAtTime) string {
	lastChecked := statusAtTime[len(statusAtTime)-1]
	difference := time.Since(lastChecked.Time)
	status := "OK"
	if !lastChecked.Status {
		status = "ERR"
	}
	return fmt.Sprintf("%s, %.0f seconds ago", status, difference.Seconds())
}

func RunHttp(address string, monitor *Monitor) {
	funcMap := template.FuncMap{
		"calculateServerUptime": calculateServerUptime,
		"lastStatus":            lastStatus,
	}

	t := template.Must(template.New("main").Funcs(funcMap).Parse(`<!DOCTYPE html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	<title>GOSSM - Dashboard</title>
	
    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/css/bootstrap.min.css" integrity="sha384-/Y6pD6FV/Vv2HJnA6t+vslU6fwYXjCFtcEpHbNJ0lyAFsXTsjBbfaDjzALeQsN6M" crossorigin="anonymous">
  </head>
  <body>
	<div class="container">
		<br>
		<center><h1>Dashboard</h1></center>
		<hr>
		<div class="row">
			{{ range $server, $statusAtTime := .}}
			<div class="col-md-4">
				<div class="card" style="margin-top: 5px;">
					<div class="card-body">
						<h4 class="card-title">{{ $server.Name }}</h4>
						<p class="card-text">{{ $server }}<br>tested {{ len $statusAtTime }} times<br>{{ $statusAtTime | lastStatus }}</p>
						<a href="#" class="btn btn-primary">{{ $statusAtTime | calculateServerUptime }}%</a>
					</div>
				</div>
			</div>
			{{ end }}
		</div>
	</div>
    <!-- Optional JavaScript -->
    <!-- jQuery first, then Popper.js, then Bootstrap JS -->
    <script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.11.0/umd/popper.min.js" integrity="sha384-b/U6ypiBEHpOf/4+1nzFpr53nxSS+GLCkfwBdFNTxtclqqenISfwAzpKaMNFNmj4" crossorigin="anonymous"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/js/bootstrap.min.js" integrity="sha384-h0AbiXch4ZDo7tp9hKZ4TsHbi047NrKGLO3SEJAg45jXxnGIfYzk4Si90RDIqNm1" crossorigin="anonymous"></script>
  </body>
</html>`))

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		t.Execute(rw, monitor.serverStatusData.GetServerStatus())
	})

	http.ListenAndServe(address, nil)
}
