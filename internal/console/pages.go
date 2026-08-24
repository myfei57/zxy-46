package console

import (
	"encoding/json"
	"io"
	"net/http"
)

const pageShellTop = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>`

const pageShellBottom = `</title>
<style>
body{font-family:sans-serif;margin:2rem;background:#f4f6f8;color:#222}
nav a{margin-right:1rem;color:#0b5cab;text-decoration:none}
table{border-collapse:collapse;width:100%;background:#fff}
th,td{border:1px solid #ccc;padding:0.5rem;text-align:left}
th{background:#e8eef4}
h1{color:#0b3d66}
</style>
</head>
<body>
<nav>
<a href="/zones">分区</a>
<a href="/plants">冷站</a>
<a href="/vav">VAV</a>
<a href="/alarms">告警</a>
</nav>
`

const indexPageHTML = `<h1>BldgHVAC 智能楼宇暖通空调群控平台</h1>
<p>请从上方导航进入分区、冷站、VAV 与告警控制台。</p>
`

const zonesPageHTML = `<h1>分区控制台</h1>
<table id="rows"><thead><tr><th>分区</th><th>楼栋</th><th>设定</th><th>供水</th><th>状态</th><th>舒适度</th></tr></thead><tbody></tbody></table>
<script>
fetch('/api/zones').then(function(r){return r.json()}).then(function(d){
  var rows='';
  d.zones.forEach(function(z){
    rows+='<tr><td>'+z.name+'</td><td>'+z.building+'</td><td>'+z.setpoint+'</td><td>'+z.supply+'</td><td>'+z.label+'</td><td>'+z.score+'</td></tr>';
  });
  document.getElementById('rows').tBodies[0].innerHTML=rows;
});
</script>
`

const plantsPageHTML = `<h1>冷站控制台</h1>
<table id="rows"><thead><tr><th>机组</th><th>状态</th><th>导泵</th><th>主机</th><th>负荷</th><th>主机</th></tr></thead><tbody></tbody></table>
<script>
fetch('/api/plants').then(function(r){return r.json()}).then(function(d){
  var rows='';
  d.plants.forEach(function(p){
    rows+='<tr><td>'+p.name+'</td><td>'+p.state+'</td><td>'+(p.pump?'运行':'停止')+'</td><td>'+(p.compressor?'运行':'停止')+'</td><td>'+p.loadKw+'</td><td>'+(p.lead?'主':'备')+'</td></tr>';
  });
  document.getElementById('rows').tBodies[0].innerHTML=rows;
  document.getElementById('baseline').textContent=d.baseline;
});
</script>
<p>供水温度基准：<span id="baseline">-</span></p>
`

const vavPageHTML = `<h1>VAV 末端控制台</h1>
<table id="rows"><thead><tr><th>末端</th><th>分区</th><th>风阀</th><th>再热</th><th>风量</th><th>通风模式</th></tr></thead><tbody></tbody></table>
<script>
fetch('/api/vav').then(function(r){return r.json()}).then(function(d){
  var rows='';
  d.vav.forEach(function(v){
    rows+='<tr><td>'+v.name+'</td><td>'+v.zone+'</td><td>'+v.damper+'</td><td>'+(v.reheat?'开':'关')+'</td><td>'+v.flow+'</td><td>'+v.mode+'</td></tr>';
  });
  document.getElementById('rows').tBodies[0].innerHTML=rows;
  document.getElementById('baseline').textContent=d.baseline;
});
</script>
<p>共享供水温度基线：<span id="baseline">-</span></p>
`

const alarmsPageHTML = `<h1>告警台</h1>
<table id="rows"><thead><tr><th>对象</th><th>消息</th></tr></thead><tbody></tbody></table>
<script>
fetch('/api/alarms').then(function(r){return r.json()}).then(function(d){
  var rows='';
  d.alarms.forEach(function(a){
    rows+='<tr><td>'+a.zone+'</td><td>'+a.message+'</td></tr>';
  });
  document.getElementById('rows').tBodies[0].innerHTML=rows;
});
</script>
`

func renderPage(w http.ResponseWriter, title, content string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, pageShellTop+title+pageShellBottom+content)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
