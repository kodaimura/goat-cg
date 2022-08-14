document.getElementById("all").addEventListener("click", (e) => {
	let ls = document.getElementsByName("table_id");

	for (let e of ls) {
		e.checked = true
	}
})


document.getElementById("clear").addEventListener("click", (e) => {
	let ls = document.getElementsByName("table_id");

	for (let e of ls) {
		e.checked = false
	}
})


const getChechedValues = () => {
	let ls = document.getElementsByName("table_id");
	let ret = [];

	for (let x of ls) {
		if (x.checked) {
			ret.push(x.value);
		}
	}
	return ret
}

document.getElementById("cg-goat").addEventListener("click", (e) => {
	let tableids = getChechedValues()
	let dbtype = document.getElementById("db_type").value

	fetch(`./codegen/goat`, {
		method: "POST",
		headers: {"Content-Type": "application/json"},
		body: JSON.stringify({dbtype,tableids})
	})
	.then(response => {
		return response.text()
	})
	.then(filepath => {
		let alink = document.createElement('a');
		alink.download = filepath.substring(5);
		alink.href = filepath;
		alink.click();
		return false;
	})
	.catch(console.error);
})

document.getElementById("cg-ddl").addEventListener("click", (e) => {
	let tableids = getChechedValues()
	let dbtype = document.getElementById("db_type").value

	fetch(`./codegen/ddl`, {
		method: "POST",
		headers: {"Content-Type": "application/json"},
		body: JSON.stringify({dbtype,tableids})
	})
	.then(response => {
		return response.text()
	})
	.then(filepath => {
		let alink = document.createElement('a');
		alink.download = filepath.substring(5);
		alink.href = filepath;
		alink.click();
		return false;
	})
	.catch(console.error);
})