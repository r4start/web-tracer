var ids_fetcher = function() {
  var request = new XMLHttpRequest();

  request.onload = function() {
    console.log(this.responseText)
  };

  request.open("GET", "/ids", true);
  request.send();
};