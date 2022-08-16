document.getElementById("data_type_cls").addEventListener("change", (e)=>{
	switch (e.target.value) {
		case "11":
		case "12":
			document.getElementById("precision").disabled = false
			document.getElementById("scale").disabled = true
			document.getElementById("scale").value = 0
			break;
		case "30":
			document.getElementById("precision").disabled = false
			document.getElementById("scale").disabled = false
			break;
		default:
			document.getElementById("precision").disabled = true
			document.getElementById("scale").disabled = true
			document.getElementById("precision").value = 0
			document.getElementById("scale").value = 0
			break;
	}
})