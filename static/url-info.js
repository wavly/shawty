document.getElementById("code").addEventListener("submit", function(event) {
  event.preventDefault();

  /** @type string */
  const input = document.getElementById("input").value;
  const errElm = document.getElementById("err");

  if (input.length < 2) {
    errElm.innerText = "2 characters is required for the code"
  } else if (input.length > 8) {
    errElm.innerText = "Only 8 characters are allowed"
  } else {
    window.location.replace(`/stat/${input}`)
  }

  setTimeout(function() {
    errElm.innerText = ""
  }, 2000);

});
