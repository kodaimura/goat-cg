const checkPassword = () => {
	let pw1 = document.getElementById("pw1").value
	let pw2 = document.getElementById("pw2").value

	if (pw1 !== pw2) {
		document.getElementById("error").innerHTML = "Password and Password (confirm) do not match."
		return false
	} else {
		return true
	}
}
