var id_refresher = null;

var ids_fetcher = function() {
  var request = new XMLHttpRequest();

  request.onload = function() {
    if (this.status != 200) {
      return;
    }

    var response_obj = JSON.parse(this.responseText);
    console.log(response_obj.ids);

    $('#debug_id').autocomplete({ source : response_obj.ids.map(String) });

    stop_fetching();

    setInterval(ids_fetcher, 10000);
  };

  request.open("GET", "http://" + $('#service_addr').val() + "/ids", true);
  request.send();
};

var stop_fetching = function() {
  if (id_refresher === null) {
    return;
  }

  clearInterval(id_refresher);
};
