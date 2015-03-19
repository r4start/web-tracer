var ids_fetcher = function() {
  var request = new XMLHttpRequest();

  request.onload = function() {
    var response_obj = JSON.parse(this.responseText);
    console.log(response_obj.ids);

    $('#debug_id').autocomplete({ source : response_obj.ids.map(String) });
  };

  request.open("GET", "http://" + $('#service_addr').val() + "/ids", true);
  request.send();
};