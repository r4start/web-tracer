var IdsFetcher = function() {
  var obj_ref = this;

  obj_ref.fetch_ids = function() {
    var request = new XMLHttpRequest();

    request.onload = obj_ref.ids_loaded;

    request.open("GET", "http://" + $('#service_addr').val() + "/ids", true);
    request.send();
  };

  obj_ref.ids_loaded = function () {
    setTimeout(obj_ref.fetch_ids, 10000);

    if (this.status != 200) {
      return;
    }

    var response_obj = JSON.parse(obj_ref.responseText);
    console.log(response_obj.ids);

    $('#debug_id').autocomplete({source: response_obj.ids.map(String)});
  };
};
