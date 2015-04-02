var IdsFetcher = function() {
  this.fetch_ids = function() {
    var request = new XMLHttpRequest();

    request.onload = this.ids_loaded;

    request.open("GET", "http://" + $('#service_addr').val() + "/ids", true);
    request.send();
  }

  this.ids_loaded = function () {
    setTimeout(this.fetch_ids, 10000);

    if (this.status != 200) {
      return;
    }

    var response_obj = JSON.parse(this.responseText);
    console.log(response_obj.ids);

    $('#debug_id').autocomplete({source: response_obj.ids.map(String)});
  }
};
