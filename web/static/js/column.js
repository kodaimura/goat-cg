document.getElementById("modal-del-button").addEventListener("click", (e)=>{
	fetch("", {method: "DELETE"})
	.then(data => {
		window.location = "../columns"
	})
})