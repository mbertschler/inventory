var scanner
function startCodeScanner() {
	scanner = new Instascan.Scanner({
		video: document.getElementById("scanVideo"),
		mirror: false,
	})
	scanner.addListener("scan", function (content) {
		guiapi("scanCode", content)
	})
	Instascan.Camera.getCameras().then(function (cameras) {
		if (cameras.length > 0) {
			scanner.start(cameras[cameras.length-1])
		} else {
			window.alert("No cameras found.")
			console.error("No cameras found.")
		}
	}).catch(function (e) {
		window.alert("Can't get cameras. Check console")
		console.error(e)
	})
}

function stopCodeScanner() {
	scanner.stop().then(function () {
		console.log("scanner stopped")
	}).catch(function (e) {
		window.alert("Can't stop scanner. Check console")
		console.error(e)
	})
}

function sendForm(action, selector) {
	var elements = $(selector)
	var data = {}
	for (var i = 0; i < elements.length; i++) {
		data[elements[i].name] = elements[i].value
	}
	guiapi(action, data)
}

function sendInput(action, event) {
	guiapi(action, event.target.value)
}

// ================= GUI API =================
function guiapi(name, args) {
	var req = {
		Actions: [{
			Name: name,
			Args: args,
		}]
	}
	$.ajax({
		method: "POST",
		url: "/guiapi/",
		data: JSON.stringify(req),
		success: function (data) {
			var ret = JSON.parse(data)
			handleResponse(ret)
		},
		error: function (error) {
			console.error("error:", error)
		},
	})
}

var callableFunctions = {
	"redirect": redirect,
	"startScan": startCodeScanner,
	"stopScan": stopCodeScanner,
}

function redirect(path) {
	window.location = path
}

function handleResponse(resp) {
	for (var i =0; i < resp.Results.length; i++) {
		var r = resp.Results[i]
		if (r.Error){
			console.error("[" + r.Error.Code + "]", r.Error.Message, r.Error)
			window.alert("guiapi error, check console")
			continue
		}
		if (r.HTML) {
			for (var j = 0; j < r.HTML.length; j++) {
				var update = r.HTML[j]
				if (update.Operation == 1) {
					$(update.Selector).html(update.Content)
				} else {
					console.warn("update type not implemented :(", update)
				}
			}
		}
		if (r.JS) {
			for (var j=0; j< r.JS.length; j++) {
				var call = r.JS[j]
				var func = callableFunctions[call.Name]
				if (func) {
					func(call.Arguments)
				} else {
					console.warn("function call not implemented :(", call)
				}
			}
		}
	}
}
