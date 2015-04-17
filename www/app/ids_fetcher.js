var ids_fetcher = function() {
  var obj_ref = this;

  obj_ref.timeout_obj = null;

  obj_ref.fetch_ids = function() {
    var request = new XMLHttpRequest();

    request.onload = obj_ref.ids_loaded;

    request.open("GET", "http://" + $('#service_addr').val() + "/ids", true);
    request.send();
  };

  obj_ref.ids_loaded = function () {
    // If response took a lot of time and previous call had already come, then we should remove it.
    if (obj_ref.timeout_obj != null) {
      clearTimeout(obj_ref.timeout_obj);
    }

    obj_ref.timeout_obj = setTimeout(obj_ref.fetch_ids, 10000);

    if (this.status != 200) {
      return;
    }

    var response_obj = JSON.parse(this.responseText);
    console.log(response_obj.ids);

    $('#debug_id').autocomplete({source: response_obj.ids.map(String)});
  };
};
