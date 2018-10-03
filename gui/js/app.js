console.log("app.js loaded");

// ================= GUI API =================
function callGuiAPI(name, args) {
	var req = {
		Actions: [{
			Name: name,
			Args: args,
		}]
	}
	$.ajax({
		method: "POST",
		url: "/gui/",
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
	"setURL": setURL,
}

function setURL(args) {
	history.pushState(args[0], args[1], args[2])
}

function handleResponse(resp) {
	for (var i=0; i< resp.Results.length; i++) {
		var r = resp.Results[i]
		if (r.HTML) {
			for (var j=0; j< r.HTML.length; j++) {
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
